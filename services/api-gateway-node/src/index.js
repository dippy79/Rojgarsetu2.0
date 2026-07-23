const express = require('express');
const { createProxyMiddleware } = require('http-proxy-middleware');
const cors = require('cors');
require('dotenv').config();

const app = express();
const PORT = process.env.PORT || 8080;

// ── CORS ────────────────────────────────────────────────────────────────
app.use(cors({
  origin: process.env.ALLOWED_ORIGIN || '*',
  methods: ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'OPTIONS'],
  allowedHeaders: ['Content-Type', 'Authorization'],
}));

// ── Body parsers ────────────────────────────────────────────────────────
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

// ── Health check ────────────────────────────────────────────────────────
app.get('/health', (_req, res) => {
  res.status(200).json({
    status: 'ok',
    service: 'rojgarsetu-api-gateway',
    timestamp: new Date().toISOString(),
  });
});

// ── Metrics passthrough ─────────────────────────────────────────────────
app.get('/metrics', (_req, res) => {
  res.status(200).send('# API Gateway metrics endpoint\n');
});

// ── Backend proxy (Go) ──────────────────────────────────────────────────
const BACKEND_TARGET = process.env.BACKEND_URL || 'http://backend:8080';
app.use('/api', createProxyMiddleware({
  target: BACKEND_TARGET,
  changeOrigin: true,
  pathRewrite: { '^/api': '' },
  on: {
    proxyReq: (proxyReq, req) => {
      if (req.headers.authorization) {
        proxyReq.setHeader('Authorization', req.headers.authorization);
      }
    },
    proxyRes: (proxyRes) => {
      delete proxyRes.headers['x-powered-by'];
    },
    error: (err, _req, res) => {
      console.error('[API Gateway] Backend proxy error:', err.message);
      if (res.writeHead) {
        res.writeHead(502, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'Backend service unavailable' }));
      }
    },
  },
}));

// ── Auth service proxy (Java/Spring Boot) ───────────────────────────────
const AUTH_TARGET = process.env.AUTH_URL || 'http://auth-service:8081';
app.use('/auth', createProxyMiddleware({
  target: AUTH_TARGET,
  changeOrigin: true,
  pathRewrite: { '^/auth': '' },
  on: {
    error: (err, _req, res) => {
      console.error('[API Gateway] Auth proxy error:', err.message);
      if (res.writeHead) {
        res.writeHead(502, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'Auth service unavailable' }));
      }
    },
  },
}));

// ── AI engine proxy (Python) ────────────────────────────────────────────
const AI_TARGET = process.env.AI_URL || 'http://ai-engine:8000';
app.use('/ai', createProxyMiddleware({
  target: AI_TARGET,
  changeOrigin: true,
  pathRewrite: { '^/ai': '' },
  on: {
    error: (err, _req, res) => {
      console.error('[API Gateway] AI engine proxy error:', err.message);
      if (res.writeHead) {
        res.writeHead(502, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'AI engine service unavailable' }));
      }
    },
  },
}));

// ── 404 handler ─────────────────────────────────────────────────────────
app.use((_req, res) => {
  res.status(404).json({ error: 'Not found' });
});

// ── Global error handler ────────────────────────────────────────────────
app.use((err, _req, res, _next) => {
  console.error('[API Gateway] Unhandled error:', err);
  res.status(500).json({ error: 'Internal server error' });
});

// ── Start server ────────────────────────────────────────────────────────
app.listen(PORT, '0.0.0.0', () => {
  console.log(`[API Gateway] Running on port ${PORT}`);
  console.log(`[API Gateway] Backend proxy → ${BACKEND_TARGET}`);
  console.log(`[API Gateway] Auth proxy    → ${AUTH_TARGET}`);
  console.log(`[API Gateway] AI engine proxy→ ${AI_TARGET}`);
});

module.exports = app;

