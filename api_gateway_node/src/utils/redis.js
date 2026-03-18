/**
 * RojgarSetu 2.0 - Redis Client
 * Caching and session management
 */

const Redis = require('ioredis');
const config = require('../config');
const { logger } = require('./logger');

const redis = new Redis({
  host: config.redis.host,
  port: config.redis.port,
  password: config.redis.password,
  db: config.redis.db,
  retryStrategy: (times) => {
    const delay = Math.min(times * 50, 2000);
    return delay;
  },
  maxRetriesPerRequest: 3
});

redis.on('connect', () => {
  logger.info('Redis connected');
});

redis.on('error', (err) => {
  logger.error('Redis error', { error: err.message });
});

// Cache helpers
const cacheGet = async (key) => {
  try {
    const data = await redis.get(key);
    return data ? JSON.parse(data) : null;
  } catch (err) {
    logger.error('Redis get error', { key, error: err.message });
    return null;
  }
};

const cacheSet = async (key, value, expiration = 3600) => {
  try {
    await redis.setex(key, expiration, JSON.stringify(value));
    return true;
  } catch (err) {
    logger.error('Redis set error', { key, error: err.message });
    return false;
  }
};

const cacheDel = async (key) => {
  try {
    await redis.del(key);
    return true;
  } catch (err) {
    logger.error('Redis delete error', { key, error: err.message });
    return false;
  }
};

// Token blacklist
const addToBlacklist = async (token, expiresIn) => {
  try {
    await redis.setex(`blacklist:${token}`, expiresIn, '1');
    return true;
  } catch (err) {
    logger.error('Redis blacklist error', { error: err.message });
    return false;
  }
};

const isBlacklisted = async (token) => {
  try {
    const result = await redis.get(`blacklist:${token}`);
    return result === '1';
  } catch (err) {
    return false;
  }
};

// Rate limiting helper
const checkRateLimit = async (key, limit, window) => {
  try {
    const current = await redis.incr(key);
    if (current === 1) {
      await redis.expire(key, window);
    }
    return {
      allowed: current <= limit,
      remaining: Math.max(0, limit - current),
      reset: await redis.ttl(key)
    };
  } catch (err) {
    logger.error('Redis rate limit error', { key, error: err.message });
    return { allowed: true, remaining: limit, reset: window };
  }
};

module.exports = {
  redis,
  get: cacheGet,
  set: cacheSet,
  setex: cacheSet,
  del: cacheDel,
  addToBlacklist,
  isBlacklisted,
  checkRateLimit
};

