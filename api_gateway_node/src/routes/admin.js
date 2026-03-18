/**
 * RojgarSetu 2.0 - Admin Routes
 * Admin-only endpoints for managing jobs, courses, users
 */

const express = require('express');
const { body, validationResult } = require('express-validator');
const { asyncHandler, AppError } = require('../utils/errorHandler');
const db = require('../database');
const redis = require('../utils/redis');

const router = express.Router();

// Middleware to check admin role
const requireAdmin = asyncHandler(async (req, res, next) => {
  if (!req.user || req.user.role !== 'admin') {
    throw new AppError('Admin access required', 403, 'ADMIN_REQUIRED');
  }
  next();
});

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

// ===== Jobs Management =====

// Get all jobs (including inactive)
router.get('/jobs', requireAdmin, asyncHandler(async (req, res, next) => {
  const { page = 1, limit = 20, source, active } = req.query;
  const offset = (page - 1) * limit;

  let whereClause = '';
  const params = [];
  let paramIndex = 1;

  if (source) {
    whereClause += ` WHERE source = $${paramIndex++}`;
    params.push(source);
  }

  if (active !== undefined) {
    const isActive = active === 'true';
    whereClause += whereClause ? ` AND is_active = $${paramIndex++}` : ` WHERE is_active = $${paramIndex++}`;
    params.push(isActive);
  }

  const countResult = await db.query(`SELECT COUNT(*) FROM jobs ${whereClause}`, params);
  const total = parseInt(countResult.rows[0].count);

  const result = await db.query(
    `SELECT j.*, c.name as company_name
     FROM jobs j
     LEFT JOIN companies c ON j.company_id = c.id
     ${whereClause}
     ORDER BY j.created_at DESC
     LIMIT $${paramIndex} OFFSET $${paramIndex + 1}`,
    [...params, parseInt(limit), offset]
  );

  res.json({
    status: 'success',
    data: {
      jobs: result.rows,
      pagination: { page: parseInt(page), limit: parseInt(limit), total, totalPages: Math.ceil(total / limit) }
    },
    error: null
  });
}));

// Create job
router.post('/jobs', requireAdmin, [
  body('title').notEmpty(),
  body('companyId').optional().isUUID(),
  body('source').notEmpty(),
  body('location').optional(),
  body('jobType').optional(),
  body('salaryMin').optional().isInt(),
  body('salaryMax').optional().isInt(),
  body('eligibility').optional(),
  body('description').optional(),
  body('applicationUrl').optional().isURL()
], validateRequest, asyncHandler(async (req, res, next) => {
  const { title, companyId, source, location, jobType, salaryMin, salaryMax, eligibility, description, applicationUrl } = req.body;

  const result = await db.query(
    `INSERT INTO jobs (title, company_id, source, location, job_type, salary_min, salary_max, eligibility, description, application_url)
     VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
     RETURNING *`,
    [title, companyId, source, location, jobType, salaryMin, salaryMax, eligibility, description, applicationUrl]
  );

  // Invalidate jobs cache
  const keys = await redis.redis.keys('jobs:*');
  if (keys.length > 0) {
    await redis.redis.del(...keys);
  }

  res.status(201).json({
    status: 'success',
    data: result.rows[0],
    error: null
  });
}));

// Update job
router.put('/jobs/:id', requireAdmin, asyncHandler(async (req, res, next) => {
  const { id } = req.params;
  const updates = req.body;

  const setClauses = Object.keys(updates).map((key, idx) => {
    const column = key.replace(/([A-Z])/g, '_$1').toLowerCase();
    return `${column} = $${idx + 1}`;
  }).join(', ');

  const values = [...Object.values(updates), id];

  const result = await db.query(
    `UPDATE jobs SET ${setClauses} WHERE id = $${Object.keys(updates).length + 1} RETURNING *`,
    values
  );

  if (result.rows.length === 0) {
    throw new AppError('Job not found', 404, 'JOB_NOT_FOUND');
  }

  // Invalidate cache
  await redis.del(`job:${id}`);
  const keys = await redis.redis.keys('jobs:*');
  if (keys.length > 0) {
    await redis.redis.del(...keys);
  }

  res.json({
    status: 'success',
    data: result.rows[0],
    error: null
  });
}));

// Delete job
router.delete('/jobs/:id', requireAdmin, asyncHandler(async (req, res, next) => {
  const { id } = req.params;

  await db.query('DELETE FROM jobs WHERE id = $1', [id]);

  // Invalidate cache
  await redis.del(`job:${id}`);

  res.json({
    status: 'success',
    data: { message: 'Job deleted successfully' },
    error: null
  });
}));

// ===== Courses Management =====

router.get('/courses', requireAdmin, asyncHandler(async (req, res, next) => {
  const result = await db.query('SELECT * FROM courses ORDER BY created_at DESC');
  res.json({ status: 'success', data: { courses: result.rows }, error: null });
}));

router.post('/courses', requireAdmin, [
  body('title').notEmpty(),
  body('provider').notEmpty(),
  body('url').isURL()
], validateRequest, asyncHandler(async (req, res, next) => {
  const { title, provider, url, skills, duration, level, isFree, price } = req.body;

  const result = await db.query(
    `INSERT INTO courses (title, provider, url, skills, duration, level, is_free, price)
     VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`,
    [title, provider, url, skills, duration, level, isFree ?? true, price]
  );

  res.status(201).json({ status: 'success', data: result.rows[0], error: null });
}));

// ===== Users Management =====

router.get('/users', requireAdmin, asyncHandler(async (req, res, next) => {
  const { page = 1, limit = 20 } = req.query;
  const offset = (page - 1) * limit;

  const countResult = await db.query('SELECT COUNT(*) FROM users');
  const total = parseInt(countResult.rows[0].count);

  const result = await db.query(
    'SELECT id, name, email, skills, created_at FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2',
    [parseInt(limit), offset]
  );

  res.json({
    status: 'success',
    data: {
      users: result.rows,
      pagination: { page: parseInt(page), limit: parseInt(limit), total, totalPages: Math.ceil(total / limit) }
    },
    error: null
  });
}));

router.delete('/users/:id', requireAdmin, asyncHandler(async (req, res, next) => {
  const { id } = req.params;
  await db.query('DELETE FROM users WHERE id = $1', [id]);
  res.json({ status: 'success', data: { message: 'User deleted' }, error: null });
}));

// ===== Statistics =====

router.get('/stats', requireAdmin, asyncHandler(async (req, res, next) => {
  const [userCount, jobCount, applicationCount, courseCount] = await Promise.all([
    db.query('SELECT COUNT(*) FROM users'),
    db.query('SELECT COUNT(*) FROM jobs WHERE is_active = true'),
    db.query('SELECT COUNT(*) FROM applications'),
    db.query('SELECT COUNT(*) FROM courses')
  ]);

  res.json({
    status: 'success',
    data: {
      users: parseInt(userCount.rows[0].count),
      activeJobs: parseInt(jobCount.rows[0].count),
      applications: parseInt(applicationCount.rows[0].count),
      courses: parseInt(courseCount.rows[0].count)
    },
    error: null
  });
}));

module.exports = router;

