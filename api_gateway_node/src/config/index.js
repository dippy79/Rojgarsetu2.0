/**
 * RojgarSetu 2.0 - Configuration
 */

module.exports = {
  nodeEnv: process.env.NODE_ENV || 'development',
  port: process.env.PORT || 3000,
  
  // Database
  db: {
    host: process.env.DB_HOST || 'localhost',
    port: parseInt(process.env.DB_PORT) || 5432,
    database: process.env.DB_NAME || 'rojgarsetu',
    user: process.env.DB_USER || (() => { throw new Error('DB_USER required'); })(),
    password: process.env.DB_PASSWORD || (() => { throw new Error('DB_PASSWORD required'); })(),
    max: 20,
    idleTimeoutMillis: 30000,
    connectionTimeoutMillis: 2000
  },

  // Redis
  redis: {
    host: process.env.REDIS_HOST || 'localhost',
    port: parseInt(process.env.REDIS_PORT) || 6379,
    password: process.env.REDIS_PASSWORD || undefined,
    db: 0
  },

  // JWT
  jwt: {
    secret: process.env.JWT_SECRET || (() => { throw new Error('JWT_SECRET required'); })(),
    expiresIn: process.env.JWT_EXPIRES_IN || '24h',
    refreshExpiresIn: process.env.JWT_REFRESH_EXPIRES_IN || '7d'
  },

  // Rate limiting
  rateLimit: {
    window: parseInt(process.env.RATE_LIMIT_WINDOW) || 15,
    max: parseInt(process.env.RATE_LIMIT_MAX) || 100
  },

  // CORS
  cors: {
    origin: process.env.CORS_ORIGIN || '*',
    methods: ['GET', 'POST', 'PUT', 'DELETE', 'PATCH'],
    allowedHeaders: ['Content-Type', 'Authorization'],
    credentials: true
  },

  // External services
  services: {
    authService: process.env.AUTH_SERVICE_URL || 'http://localhost:8081',
    crawlerService: process.env.CRAWLER_SERVICE_URL || 'http://localhost:8082',
    aiEngine: process.env.AI_ENGINE_URL || 'http://localhost:8000'
  },

  // Logging
  logging: {
    level: process.env.LOG_LEVEL || 'info'
  }
};
