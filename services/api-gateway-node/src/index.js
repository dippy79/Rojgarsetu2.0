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

    if (allowed.length === 0) {
      console.error('FATAL: ALLOWED_ORIGINS not configured. Failing closed.');
      return callback(new Error('CORS Policy: ALLOWED_ORIGINS missing'));
    }

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

// 2. PROXY CONFIGURATION
const BACKEND_TARGET = process.env.BACKEND_SERVICE_URL || 'http://backend:8083';
const AI_TARGET = process.env.AI_ENGINE_URL || 'http://ai-engine:8000';
const CRAWLER_TARGET = process.env.CRAWLER_SERVICE_URL || 'http://crawler:8080';

const AI_SECRET_KEY = (process.env.AI_ENGINE_API_KEY || '').trim();

const proxyOptions = {
  changeOrigin: true,
  on: {
    proxyReq: (proxyReq, req) => {
      if (req.originalUrl.startsWith('/api/ai')) {
        proxyReq.setHeader('X-AI-Secret-Key', AI_SECRET_KEY);
      }
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

// 3. AUTH ROUTES (EXPLICIT FORWARDING)
app.use('/api/v1/auth', createProxyMiddleware({
  target: BACKEND_TARGET,
  pathRewrite: { '^/api/v1/auth': '/api/v1/auth' }, // Ensure path remains /api/v1/auth
  ...proxyOptions
}));

app.use('/api/auth', createProxyMiddleware({
  target: BACKEND_TARGET,
  pathRewrite: { '^/api/auth': '/api/v1/auth' },
  ...proxyOptions
}));

// 4. RATE LIMITING & CSRF
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
app.use('/api/v1', generalLimiter, createProxyMiddleware({ target: BACKEND_TARGET, ...proxyOptions }));
app.use('/api/crawler', createProxyMiddleware({ target: CRAWLER_TARGET, ...proxyOptions }));
app.use('/api/ai', createProxyMiddleware({ target: AI_TARGET, pathRewrite: { '^/api/ai': '' }, ...proxyOptions }));

// 6. HEALTH & FALLBACK
app.get('/health', (req, res) => res.json({ status: 'UP' }));

app.use('/api', createProxyMiddleware({
  target: BACKEND_TARGET,
  ...proxyOptions
}));

const PORT = process.env.PORT || 3000;
app.listen(PORT, '0.0.0.0', () => {
  console.log(`[API Gateway] Listening on port ${PORT}`);
});
