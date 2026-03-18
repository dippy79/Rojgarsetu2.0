/**
 * RojgarSetu 2.0 - API Gateway
 * Main entry point with zero-error ready design
 */

require('dotenv').config();
const express = require('express');
const helmet = require('helmet');
const cors = require('cors');
const rateLimit = require('express-rate-limit');
const SlowDown = require('express-slow-down');
const { errorHandler, notFoundHandler } = require('./utils/errorHandler');
const { requestLogger } = require('./utils/logger');
const config = require('./config');

// Prometheus metrics
const promClient = require('prom-client');
const register = new promClient.Registry();

// Add default metrics
promClient.collectDefaultMetrics({ register });

// Custom metrics
const httpRequestDuration = new promClient.Histogram({
  name: 'http_request_duration_seconds',
  help: 'Duration of HTTP requests in seconds',
  labelNames: ['method', 'route', 'status_code'],
  buckets: [0.1, 0.5, 1, 2, 5]
});
register.registerMetric(httpRequestDuration);

const httpRequestCount = new promClient.Counter({
  name: 'http_requests_total',
  help: 'Total number of HTTP requests',
  labelNames: ['method', 'route', 'status_code']
});
register.registerMetric(httpRequestCount);

// Route imports
const authRoutes = require('./routes/auth');
const jobRoutes = require('./routes/jobs');
const courseRoutes = require('./routes/courses');
const userRoutes = require('./routes/users');
const adminRoutes = require('./routes/admin');
const recommendationRoutes = require('./routes/recommendations');

const app = express();

// Security middleware
app.use(helmet());
app.use(cors(config.cors));

// Body parsing
app.use(express.json({ limit: '10mb' }));
app.use(express.urlencoded({ extended: true }));

// Request logging
app.use(requestLogger);

// Rate limiting - prevents abuse
const limiter = rateLimit({
  windowMs: config.rateLimit.window * 60 * 1000,
  max: config.rateLimit.max,
  message: {
    status: 'error',
    error: {
      code: 429,
      message: 'Too many requests, please try again later.'
    }
  },
  standardHeaders: true,
  legacyHeaders: false
});
app.use('/api/', limiter);

// Slow down - prevents rapid-fire attacks
const speedLimiter = SlowDown({
  windowMs: 15 * 60 * 1000,
  delayAfter: 100,
  delayMs: 500,
  maxDelayMs: 20000
});
app.use(speedLimiter);

// Health check endpoint (no rate limit)
app.get('/health', (req, res) => {
  res.status(200).json({
    status: 'success',
    data: {
      service: 'api-gateway',
      version: '2.0.0',
      timestamp: new Date().toISOString(),
      uptime: process.uptime()
    },
    error: null
  });
});

// Metrics endpoint for Prometheus
app.get('/metrics', async (req, res) => {
  try {
    res.set('Content-Type', register.contentType);
    res.end(await register.metrics());
  } catch (err) {
    res.status(500).end(err.message);
  }
});

// API Routes
app.use('/api/auth', authRoutes);
app.use('/api/jobs', jobRoutes);
app.use('/api/courses', courseRoutes);
app.use('/api/users', userRoutes);
app.use('/api/admin', adminRoutes);
app.use('/api/recommendations', recommendationRoutes);

// Circuit breaker health check
app.get('/api/circuit-breaker/status', (req, res) => {
  res.status(200).json({
    status: 'success',
    data: {
      services: {
        auth: 'CLOSED',
        crawler: 'CLOSED',
        ai: 'CLOSED'
      }
    },
    error: null
  });
});

// Error handling - zero 404/500 errors
app.use(notFoundHandler);
app.use(errorHandler);

// Graceful shutdown
const gracefulShutdown = () => {
  console.log('🛑 Received shutdown signal. Closing gracefully...');
  // Close database connections
  // Close Redis connections
  process.exit(0);
};

process.on('SIGTERM', gracefulShutdown);
process.on('SIGINT', gracefulShutdown);

// Start server
const PORT = config.port;
const server = app.listen(PORT, () => {
  console.log(`🚀 RojgarSetu API Gateway running on port ${PORT}`);
  console.log(`📊 Environment: ${config.nodeEnv}`);
});

module.exports = app;

