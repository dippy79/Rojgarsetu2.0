 /**
 * RojgarSetu 2.0 - Courses Routes
 * Course listing and recommendations
 */

const express = require('express');
const { query, validationResult } = require('express-validator');
const { asyncHandler, AppError } = require('../utils/errorHandler');
const db = require('../database');
const redis = require('../utils/redis');

const router = express.Router();

const validateRequest = (req, res, next) => {
  const errors = validationResult(req);
  if (!errors.isEmpty()) {
    return res.status(400).json({
      status: 'error',
      data: null,
      error: { code: 400, message: 'Validation failed', details: errors.array() }
    });
  }
  next();
};

// Get all courses
router.get('/', [
  query('page').optional().isInt({ min: 1 }),
  query('limit').optional().isInt({ min: 1, max: 100 }),
  query('provider').optional().trim().escape(),
  query('level').optional().trim().escape(),
  query('free').optional().isBoolean()
], validateRequest, asyncHandler(async (req, res, next) => {
  const page = parseInt(req.query.page) || 1;
  const limit = parseInt(req.query.limit) || 20;
  const offset = (page - 1) * limit;
  const { provider, level, free } = req.query;

  let whereClause = '';
  const params = [];
  let paramIndex = 1;

  if (provider) {
    whereClause += ` WHERE provider = $${paramIndex}`;
    params.push(provider);
    paramIndex++;
  }

  if (level) {
    whereClause += whereClause ? ` AND level = $${paramIndex}` : ` WHERE level = $${paramIndex}`;
    params.push(level);
    paramIndex++;
  }

  if (free !== undefined) {
    const isFree = free === 'true';
    whereClause += whereClause ? ` AND is_free = $${paramIndex}` : ` WHERE is_free = $${paramIndex}`;
    params.push(isFree);
    paramIndex++;
  }

  // Check cache
  const cacheKey = `courses:${page}:${limit}:${JSON.stringify(req.query)}`;
  const cached = await redis.get(cacheKey);
  if (cached) {
    return res.json(JSON.parse(cached));
  }

  // Get total count
  const countResult = await db.query(`SELECT COUNT(*) FROM courses ${whereClause}`);
  const total = parseInt(countResult.rows[0].count);

  // Get courses
  const result = await db.query(
    `SELECT id, title, provider, url, skills, duration, level, thumbnail_url, is_free, price
     FROM courses ${whereClause}
     ORDER BY created_at DESC
     LIMIT $${paramIndex} OFFSET $${paramIndex + 1}`,
    [...params, limit, offset]
  );

  const courses = result.rows.map(c => ({
    id: c.id,
    title: c.title,
    provider: c.provider,
    url: c.url,
    skills: c.skills,
    duration: c.duration,
    level: c.level,
    thumbnailUrl: c.thumbnail_url,
    isFree: c.is_free,
    price: c.price
  }));

  const response = {
    status: 'success',
    data: {
      courses,
      pagination: { page, limit, total, totalPages: Math.ceil(total / limit) }
    },
    error: null
  };

  await redis.setex(cacheKey, 300, JSON.stringify(response));
  res.json(response);
}));

// Get single course
router.get('/:id', asyncHandler(async (req, res, next) => {
  const { id } = req.params;

  const cacheKey = `course:${id}`;
  const cached = await redis.get(cacheKey);
  if (cached) {
    return res.json(JSON.parse(cached));
  }

  const result = await db.query(
    `SELECT id, title, provider, url, skills, duration, level, thumbnail_url, is_free, price, created_at
     FROM courses WHERE id = $1`,
    [id]
  );

  if (result.rows.length === 0) {
    throw new AppError('Course not found', 404, 'COURSE_NOT_FOUND');
  }

  const c = result.rows[0];
  const response = {
    status: 'success',
    data: {
      id: c.id,
      title: c.title,
      provider: c.provider,
      url: c.url,
      skills: c.skills,
      duration: c.duration,
      level: c.level,
      thumbnailUrl: c.thumbnail_url,
      isFree: c.is_free,
      price: c.price,
      createdAt: c.created_at
    },
    error: null
  };

  await redis.setex(cacheKey, 600, JSON.stringify(response));
  res.json(response);
}));

// Get recommended courses based on user skills
router.get('/recommendations/me', asyncHandler(async (req, res, next) => {
  const userId = req.user?.userId;
  const cacheKey = `course-recommendations:${userId}`;
  
  const cached = await redis.get(cacheKey);
  if (cached) {
    return res.json(JSON.parse(cached));
  }

  // Get courses matching user skills
  const result = await db.query(
    `SELECT id, title, provider, url, skills, duration, level, thumbnail_url, is_free, price
     FROM courses 
     ORDER BY created_at DESC
     LIMIT 10`
  );

  const courses = result.rows.map(c => ({
    id: c.id,
    title: c.title,
    provider: c.provider,
    url: c.url,
    skills: c.skills,
    duration: c.duration,
    level: c.level,
    thumbnailUrl: c.thumbnail_url,
    isFree: c.is_free,
    price: c.price
  }));

  const response = {
    status: 'success',
    data: { courses },
    error: null
  };

  await redis.setex(cacheKey, 900, JSON.stringify(response));
  res.json(response);
}));

module.exports = router;

