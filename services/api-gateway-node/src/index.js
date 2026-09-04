const express = require('express');
const { createProxyMiddleware } = require('http-proxy-middleware');
const rateLimit = require('express-rate-limit');
const RedisStore = require('rate-limit-redis');
const Redis = require('redis');
const csrf = require('csurf');
const cookieParser = require('cookie-parser');
require('dotenv').config();

const app = express();
app.use(express.json({ limit: '1mb' }));
app.use(express.urlencoded({ limit: '1mb', extended: true }));
app.use(cookieParser());

const cors = require('cors');
app.use(cors({
  origin: function(origin, callback) {
    const allowed = process.env.ALLOWED_ORIGINS
      ? process.env.ALLOWED_ORIGINS.split(',')
      : [
          'http://localhost:3000',
          'http://localhost:3001',
          'http://localhost:3002',
          'http://localhost:8080',
          'http://localhost:80',
          'http://localhost',
          'http://127.0.0.1:3000',
          'http://127.0.0.1:3001',
          'http://127.0.0.1:3002',
          'http://127.0.0.1:8080'
        ];
    if (!origin || allowed.includes(origin)) {
      callback(null, true);
    } else {
      callback(new Error('Not allowed by CORS'));
    }
  },
  credentials: true,
  methods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS', 'PATCH'],
  allowedHeaders: ['Content-Type', 'Authorization', 'X-Requested-With']
}));

// Redis client for rate limiting with graceful degradation
let redisClient = null;
let redisAvailable = false;

if (process.env.REDIS_URL || process.env.REDIS_HOST) {
  try {
    const redisUrl = process.env.REDIS_URL || `redis://${process.env.REDIS_HOST || 'localhost'}:${process.env.REDIS_PORT || 6379}`;
    redisClient = Redis.createClient({
      url: redisUrl,
      socket: {
        reconnectStrategy: (retries) => {
          if (retries > 3) {
            console.error('[API Gateway] Redis reconnection failed after 3 attempts');
            return new Error('Redis reconnection failed');
          }
          return Math.min(retries * 100, 3000);
        }
      }
    });

    redisClient.on('error', (err) => {
      console.error('[API Gateway] Redis client error:', err.message);
      redisAvailable = false;
    });

    redisClient.on('connect', () => {
      console.log('[API Gateway] Redis client connected');
      redisAvailable = true;
    });

    redisClient.on('disconnect', () => {
      console.warn('[API Gateway] Redis client disconnected');
      redisAvailable = false;
    });

    (async () => {
      try {
        await redisClient.connect();
        redisAvailable = true;
      } catch (err) {
        console.error('[API Gateway] Failed to connect to Redis, falling back to in-memory rate limiting:', err.message);
        redisAvailable = false;
      }
    })();
  } catch (err) {
    console.error('[API Gateway] Redis initialization failed, falling back to in-memory rate limiting:', err.message);
    redisAvailable = false;
  }
} else {
  console.log('[API Gateway] Redis URL not configured, using in-memory rate limiting');
}

// Rate limiting with Redis backend or fallback to in-memory
const createRateLimiter = (windowMs, max, message) => {
  const config = {
    windowMs,
    max,
    message,
    standardHeaders: true,
    legacyHeaders: false,
  };

  if (redisAvailable && redisClient) {
    config.store = new RedisStore({
      client: redisClient,
      prefix: 'rate_limit:',
    });
    console.log('[API Gateway] Using Redis-backed rate limiting');
  } else {
    console.log('[API Gateway] Using in-memory rate limiting');
  }

  return rateLimit(config);
};

const generalLimiter = createRateLimiter(60000, 100, { error: 'Too many requests, please try again later' });
const loginLimiter = createRateLimiter(60000, 5, { error: 'Too many login attempts, please try again later' });
const applyLimiter = createRateLimiter(60000, 10, { error: 'Too many application attempts, please try again later' });

// Apply rate limiters before proxy middleware
app.use('/auth/login', loginLimiter);
app.use('/api/', generalLimiter);

// Apply specific limiters for job application endpoints
app.use('/api/v1/*/apply', applyLimiter);
app.use('/api/v1/gov-jobs', applyLimiter);
app.use('/api/v1/priv-jobs', applyLimiter);

// CSRF Protection (applies only to state-changing methods)
const csrfSecret = process.env.CSRF_SECRET || 'default-csrf-secret-change-in-production';
const csrfProtection = csrf({
  cookie: {
    httpOnly: true,
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'strict',
    maxAge: 3600000, // 1 hour
  },
  secret: csrfSecret,
  ignoreMethods: ['GET', 'HEAD', 'OPTIONS'],
});

// CSRF token endpoint for frontend (must be before CSRF protection application)
app.get('/api/csrf-token', csrfProtection, (req, res) => {
  res.json({ csrfToken: req.csrfToken() });
});

// Apply CSRF protection to state-changing API routes only
app.use('/api', csrfProtection);

// Input validation middleware
app.use((req, res, next) => {
  // Sanitize query parameters
  if (req.query) {
    for (const key in req.query) {
      if (typeof req.query[key] === 'string') {
        req.query[key] = req.query[key].trim();
        // Basic XSS prevention
        req.query[key] = req.query[key]
          .replace(/</g, '&lt;')
          .replace(/>/g, '&gt;')
          .replace(/"/g, '&quot;')
          .replace(/'/g, '&#x27;');
      }
    }
  }
  next();
});

const PORT = process.env.PORT || 3002;

// Consolidated Target Configuration
const BACKEND_TARGET = process.env.BACKEND_SERVICE_URL || 'http://backend:8083';
const AUTH_TARGET = process.env.AUTH_SERVICE_URL || process.env.AUTH_URL || 'http://auth-service:8081';
const AI_TARGET = process.env.AI_ENGINE_URL || process.env.AI_URL || 'http://ai-engine:8000';
const CRAWLER_TARGET = process.env.CRAWLER_SERVICE_URL || 'http://crawler:8080';

app.get('/health', (req, res) => {
  res.status(200).json({ status: 'UP', gateway: 'API Gateway Online' });
});

// Proxy options with retry
const proxyOptions = {
  changeOrigin: true,
  timeout: 30000,
  proxyTimeout: 30000,
  onProxyReq: (proxyReq, req) => {
    if (req.headers.authorization) proxyReq.setHeader('Authorization', req.headers.authorization);
    proxyReq.setHeader('X-Forwarded-For', req.ip);
    const bodyData = fixRequestBody(proxyReq, req);
    if (bodyData) proxyReq.write(bodyData);
  },
  onProxyRes: (proxyRes) => { delete proxyRes.headers['x-powered-by']; },
  onError: (err, _req, res) => {
    console.error(`[API Gateway] Proxy error:`, err.message);
    if (!res.headersSent) res.status(502).json({ error: 'Upstream service unavailable' });
  },
};

// Route Definitions
app.use('/api/v1', createProxyMiddleware({
  target: BACKEND_TARGET,
  ...proxyOptions,
}));

app.use('/api/crawler', createProxyMiddleware({
  target: CRAWLER_TARGET,
  ...proxyOptions,
}));

app.use('/api/jobs/recommendations/me', createProxyMiddleware({
  target: AI_TARGET,
  pathRewrite: { '^/api/jobs/recommendations/me': '/recommend/jobs' },
  ...proxyOptions,
}));

app.use('/auth', createProxyMiddleware({ target: AUTH_TARGET, pathRewrite: { '^/auth': '' }, ...proxyOptions }));
app.use('/ai', createProxyMiddleware({ target: AI_TARGET, pathRewrite: { '^/ai': '' }, ...proxyOptions }));
app.use('/api', createProxyMiddleware({ target: BACKEND_TARGET, pathRewrite: { '^/api': '' }, ...proxyOptions }));

app.use((_req, res) => res.status(404).json({ error: 'Route not found' }));

// General error handler with CSRF error handling
app.use((err, _req, res, _next) => {
  if (err.code === 'EBADCSRFTOKEN') {
    return res.status(403).json({ error: 'Invalid CSRF token' });
  }
  console.error('[API Gateway] Unhandled error:', err);
  if (!res.headersSent) res.status(500).json({ error: 'Internal server error' });
});

const server = app.listen(PORT, '::', () => {
  console.log(`[API Gateway] Running on http://[::]:${PORT}`);
  console.log(`[API Gateway] Backend proxy → ${BACKEND_TARGET}`);
  console.log(`[API Gateway] Auth proxy    → ${AUTH_TARGET}`);
  console.log(`[API Gateway] AI engine proxy→ ${AI_TARGET}`);
});

const shutdown = (signal) => {
  console.log(`[API Gateway] Received ${signal}. Shutting down...`);
  
  // Close Redis connection if available
  if (redisClient && redisAvailable) {
    redisClient.quit()
      .then(() => console.log('[API Gateway] Redis connection closed'))
      .catch((err) => console.error('[API Gateway] Error closing Redis connection:', err.message))
      .finally(() => {
        server.close(() => process.exit(0));
        setTimeout(() => process.exit(1), 10000).unref();
      });
  } else {
    server.close(() => process.exit(0));
    setTimeout(() => process.exit(1), 10000).unref();
  }
};
process.on('SIGTERM', () => shutdown('SIGTERM'));
process.on('SIGINT', () => shutdown('SIGINT'));

module.exports = app;