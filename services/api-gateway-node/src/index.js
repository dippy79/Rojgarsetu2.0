const express = require('express');
const { createProxyMiddleware } = require('http-proxy-middleware');
require('dotenv').config();

const app = express();
const PORT = process.env.PORT || 8000;

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

app.listen(PORT, () => {
  console.log(`API Gateway running on port ${PORT}`);
});

// Proxy options with retry
const proxyOptions = {
  changeOrigin: true,
  timeout: 30000,
  proxyTimeout: 30000,
  onProxyReq: (proxyReq, req) => {
    fixRequestBody(proxyReq, req);
    if (req.headers.authorization) proxyReq.setHeader('Authorization', req.headers.authorization);
    proxyReq.setHeader('X-Forwarded-For', req.ip);
  },
  onProxyRes: (proxyRes) => { delete proxyRes.headers['x-powered-by']; },
  onError: (err, _req, res) => {
    console.error(`[API Gateway] Proxy error:`, err.message);
    if (!res.headersSent) res.status(502).json({ error: 'Upstream service unavailable' });
  },
};

const BACKEND_TARGET = process.env.BACKEND_SERVICE_URL || 'http://backend:8083';
app.use('/api', createProxyMiddleware({ target: BACKEND_TARGET, pathRewrite: { '^/api': '' }, ...proxyOptions }));

const AUTH_TARGET = process.env.AUTH_URL || 'http://auth-service:8081';
app.use('/auth', createProxyMiddleware({ target: AUTH_TARGET, pathRewrite: { '^/auth': '' }, ...proxyOptions }));

const AI_TARGET = process.env.AI_URL || 'http://ai-engine:8000';
app.use('/ai', createProxyMiddleware({ target: AI_TARGET, pathRewrite: { '^/ai': '' }, ...proxyOptions }));

app.use(express.json({ limit: '10mb' }));
app.use(express.urlencoded({ extended: true, limit: '10mb' }));
app.use((_req, res) => res.status(404).json({ error: 'Route not found' }));
app.use((err, _req, res, _next) => { console.error('[API Gateway] Unhandled error:', err); if (!res.headersSent) res.status(500).json({ error: 'Internal server error' }); });

const server = app.listen(PORT, '::', () => {
  console.log(`[API Gateway] Running on http://[::]:${PORT}`);
  console.log(`[API Gateway] Backend proxy → ${BACKEND_TARGET}`);
  console.log(`[API Gateway] Auth proxy    → ${AUTH_TARGET}`);
  console.log(`[API Gateway] AI engine proxy→ ${AI_TARGET}`);
});

const shutdown = (signal) => { console.log(`[API Gateway] Received ${signal}. Shutting down...`); server.close(() => process.exit(0)); setTimeout(() => process.exit(1), 10000).unref(); };
process.on('SIGTERM', () => shutdown('SIGTERM'));
process.on('SIGINT', () => shutdown('SIGINT'));
module.exports = app;