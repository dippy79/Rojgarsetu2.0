const express = require('express');
const { createProxyMiddleware } = require('http-proxy-middleware');
const rateLimit = require('express-rate-limit');
require('dotenv').config();

const app = express();
app.use(express.json({ limit: '1mb' }));
app.use(express.urlencoded({ limit: '1mb', extended: true }));

const cors = require('cors');
app.use(cors({
  origin: function(origin, callback) {
    const allowed = [
      'http://localhost:3000',
      'http://localhost:3001',
      'http://localhost:3002',
      'http://localhost:80',
      'http://localhost',
      'http://127.0.0.1:3000',
      'http://127.0.0.1:3001',
      'http://127.0.0.1:3002'
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

// Rate limiting
const generalLimiter = rateLimit({
  windowMs: 60000,
  max: 100,
  message: { error: 'Too many requests, please try again later' }
});

const loginLimiter = rateLimit({
  windowMs: 60000,
  max: 5,
  message: { error: 'Too many login attempts, please try again later' }
});

const applyLimiter = rateLimit({
  windowMs: 60000,
  max: 10,
  message: { error: 'Too many application attempts, please try again later' }
});

// Apply rate limiters before proxy middleware
app.use('/auth/login', loginLimiter);
app.use('/api/', generalLimiter);

// Apply specific limiters for job application endpoints
app.use('/api/v1/gov-jobs', applyLimiter);
app.use('/api/v1/priv-jobs', applyLimiter);

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

// Dynamic target URLs from env with local fallbacks
const BACKEND_SERVICE_URL = process.env.BACKEND_SERVICE_URL || 'http://localhost:8080';
const CRAWLER_SERVICE_URL = process.env.CRAWLER_SERVICE_URL || 'http://localhost:8081';

// Proxy rules for Go Backend
app.use('/api/v1', createProxyMiddleware({
  target: BACKEND_SERVICE_URL,
  changeOrigin: true,
  pathRewrite: { '^/api/v1': '/api/v1' },
}));

// Proxy rules for Crawler Service
app.use('/api/crawler', createProxyMiddleware({
  target: CRAWLER_SERVICE_URL,
  changeOrigin: true,
}));

app.get('/health', (req, res) => {
  res.status(200).json({ status: 'UP', gateway: 'API Gateway Online' });
});

// Proxy helper for JSON request bodies
function fixRequestBody(proxyReq, req) {
  if (!req.body || Object.keys(req.body).length === 0) {
    return null;
  }
  const bodyData = JSON.stringify(req.body);
  proxyReq.setHeader('Content-Type', 'application/json');
  proxyReq.setHeader('Content-Length', Buffer.byteLength(bodyData));
  return bodyData;
}

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

const BACKEND_TARGET = process.env.BACKEND_SERVICE_URL || 'http://backend:8083';

const AI_TARGET = process.env.AI_ENGINE_URL || process.env.AI_URL || 'http://ai-engine:8000';

app.use('/api/jobs/recommendations/me', createProxyMiddleware({
  target: AI_TARGET,
  changeOrigin: true,
  pathRewrite: { '^/api/jobs/recommendations/me': '/recommend/jobs' },
  ...proxyOptions,
}));

app.use('/api', createProxyMiddleware({
  target: BACKEND_TARGET,
  pathRewrite: { '^/api': '' },
  ...proxyOptions,
  filter: (pathname) => !pathname.startsWith('/api/jobs/recommendations/me'),
}));

const AUTH_TARGET = process.env.AUTH_SERVICE_URL || process.env.AUTH_URL || 'http://auth-service:8081';
app.use('/auth', createProxyMiddleware({ target: AUTH_TARGET, pathRewrite: { '^/auth': '' }, ...proxyOptions }));

app.use('/ai', createProxyMiddleware({ target: AI_TARGET, pathRewrite: { '^/ai': '' }, ...proxyOptions }));

app.use((_req, res) => res.status(404).json({ error: 'Route not found' }));
app.use((err, _req, res, _next) => { 
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
  server.close(() => process.exit(0)); 
  setTimeout(() => process.exit(1), 10000).unref(); 
};
process.on('SIGTERM', () => shutdown('SIGTERM'));
process.on('SIGINT', () => shutdown('SIGINT'));

module.exports = app;