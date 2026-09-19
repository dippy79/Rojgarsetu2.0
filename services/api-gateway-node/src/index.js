const express = require('express');
const { createProxyMiddleware } = require('http-proxy-middleware');
const rateLimit = require('express-rate-limit');
const { RedisStore } = require('rate-limit-redis');
const Redis = require('redis');
const csrf = require('csurf');
const cookieParser = require('cookie-parser');
const cors = require('cors');
require('dotenv').config();

const app = express();

// 1. BASIC MIDDLEWARE
app.use((req, res, next) => {
  console.log(`[API Gateway] Incoming: ${req.method} ${req.url}`);
  next();
});

app.use(cookieParser());
app.use(cors({
  origin: function(origin, callback) {
    const allowed = process.env.ALLOWED_ORIGINS
      ? process.env.ALLOWED_ORIGINS.split(',')
      : [];

    // FAIL-CLOSED: If no allowed origins configured, reject all
    if (allowed.length === 0) {
      console.error('FATAL: ALLOWED_ORIGINS not configured. Failing closed for security.');
      return callback(new Error('CORS Policy: ALLOWED_ORIGINS missing'));
    }

    if (!origin || allowed.includes(origin)) {
      callback(null, true);
    } else {
      console.warn(`[CORS] Blocked request from unauthorized origin: ${origin}`);
      callback(new Error('Not allowed by CORS'));
    }
  },
  credentials: true,
  methods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS', 'PATCH'],
  allowedHeaders: ['Content-Type', 'Authorization', 'X-Requested-With']
}));

// 2. PROXY CONFIGURATION (CONSOLIDATED)
const BACKEND_TARGET = process.env.BACKEND_SERVICE_URL || 'http://localhost:8084';
const AUTH_TARGET = process.env.AUTH_SERVICE_URL || 'http://auth-service:8081';
const AI_TARGET = process.env.AI_ENGINE_URL || 'http://ai-engine:8000';
const CRAWLER_TARGET = process.env.CRAWLER_SERVICE_URL || 'http://crawler:8080';

console.log(`[API Gateway] Routing initialized.`);

const proxyOptions = {
  changeOrigin: true,
  on: {
    proxyReq: (proxyReq, req) => {
      console.log(`[API Gateway] PROXY -> ${req.method} ${req.url}`);
      if (req.headers.authorization) proxyReq.setHeader('Authorization', req.headers.authorization);
      else if (req.cookies?.access_token) proxyReq.setHeader('Authorization', `Bearer ${req.cookies.access_token}`);

      if (req.body && Object.keys(req.body).length > 0) {
        const bodyData = JSON.stringify(req.body);
        proxyReq.setHeader('Content-Type', 'application/json');
        proxyReq.setHeader('Content-Length', Buffer.byteLength(bodyData));
        proxyReq.write(bodyData);
      }
    },
    error: (err, req, res) => {
      console.error(`[API Gateway] PROXY ERROR (${req.url}):`, err.message);
      if (!res.headersSent) res.status(502).json({ error: 'Upstream unavailable' });
    }
  }
};

// 3. SPECIAL ROUTES (CSRF BYPASS)
// Definitively route ALL auth traffic to the Go Backend
// Java Auth service is being decommissioned for unified Go Auth.
const authProxy = createProxyMiddleware({
  target: BACKEND_TARGET,
  pathRewrite: {
    '^/api/v1/auth': '/api/v1/auth',
    '^/api/auth': '/api/v1/auth'
  },
  ...proxyOptions
});

app.use('/api/v1/auth', authProxy);
app.use('/api/auth', authProxy);

// 4. RATE LIMITING & CSRF (Applied to other routes)
app.use(express.json({ limit: '1mb' }));

const generalLimiter = rateLimit({
  windowMs: 60000,
  max: 100,
  message: { error: 'Too many requests' }
});

const csrfProtection = csrf({ cookie: { httpOnly: true, sameSite: 'strict' } });
app.get('/api/csrf-token', csrfProtection, (req, res) => {
  res.json({ csrfToken: req.csrfToken() });
});

// 5. REMAINING PROXIES
// Use pathFilter to avoid stripping the prefix
app.use(createProxyMiddleware({
  pathFilter: '/api/v1',
  target: BACKEND_TARGET,
  ...proxyOptions
}));

app.use(createProxyMiddleware({
  pathFilter: '/api/crawler',
  target: CRAWLER_TARGET,
  ...proxyOptions
}));

app.use('/api/ai', createProxyMiddleware({
  target: AI_TARGET,
  pathRewrite: { '^/api/ai': '' },
  ...proxyOptions
}));

// 6. HEALTH & FALLBACK
app.get('/health', (req, res) => res.json({ status: 'UP' }));

app.use('/api', createProxyMiddleware({
  target: BACKEND_TARGET,
  ...proxyOptions
}));
app.use((req, res) => {
  console.warn(`[API Gateway] 404 Not Found: ${req.url}`);
  res.status(404).json({ error: 'Route not found' });
});

const PORT = process.env.PORT || 3000;
app.listen(PORT, '0.0.0.0', () => {
  console.log(`[API Gateway] Listening on port ${PORT}`);
});
