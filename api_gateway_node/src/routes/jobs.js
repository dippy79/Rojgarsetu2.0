/**
 * RojgarSetu 2.0 - Jobs Routes
 * Job search, details, applications
 */

const express = require('express');
const { body, query, validationResult } = require('express-validator');
const { asyncHandler, AppError } = require('../utils/errorHandler');
const { logger } = require('../utils/logger');
const db = require('../database');
const redis = require('../utils/redis');

const router = express.Router();

// Validation middleware
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

// Get all jobs with pagination and filters
router.get('/', [
  query('page').optional().isInt({ min: 1 }).withMessage('Page must be a positive integer'),
  query('limit').optional().isInt({ min: 1, max: 100 }).withMessage('Limit must be between 1 and 100'),
  query('location').optional().trim().escape(),
  query('source').optional().trim().escape(),
  query('jobType').optional().trim().escape()
], validateRequest, asyncHandler(async (req, res, next) => {
  const page = parseInt(req.query.page) || 1;
  const limit = parseInt(req.query.limit) || 20;
  const offset = (page - 1) * limit;
  const { location, source, jobType, search } = req.query;

  // Build query dynamically
  let whereClause = 'WHERE j.is_active = true';
  const params = [];
  let paramIndex = 1;

  if (location) {
    whereClause += ` AND j.location ILIKE $${paramIndex}`;
    params.push(`%${location}%`);
    paramIndex++;
  }

  if (source) {
    whereClause += ` AND j.source = $${paramIndex}`;
    params.push(source);
    paramIndex++;
  }

  if (jobType) {
    whereClause += ` AND j.job_type = $${paramIndex}`;
    params.push(jobType);
    paramIndex++;
  }

  if (search) {
    whereClause += ` AND (j.title ILIKE $${paramIndex} OR j.description ILIKE $${paramIndex})`;
    params.push(`%${search}%`);
    paramIndex++;
  }

  // Check cache first
  const cacheKey = `jobs:${page}:${limit}:${JSON.stringify(req.query)}`;
  const cached = await redis.get(cacheKey);
  if (cached) {
    logger.debug('Jobs cache hit', { key: cacheKey });
    return res.json(JSON.parse(cached));
  }

  // Get total count
  const countResult = await db.query(
    `SELECT COUNT(*) FROM jobs j ${whereClause}`,
    params
  );
  const total = parseInt(countResult.rows[0].count);

  // Get jobs
  const jobsResult = await db.query(
    `SELECT j.id, j.title, j.location, j.job_type, j.salary_min, j.salary_max, 
            j.eligibility, j.description, j.application_url, j.posted_at, j.source,
            c.id as company_id, c.name as company_name, c.logo_url as company_logo
     FROM jobs j
     LEFT JOIN companies c ON j.company_id = c.id
     ${whereClause}
     ORDER BY j.posted_at DESC
     LIMIT $${paramIndex} OFFSET $${paramIndex + 1}`,
    [...params, limit, offset]
  );

  const jobs = jobsResult.rows.map(job => ({
    id: job.id,
    title: job.title,
    location: job.location,
    jobType: job.job_type,
    salaryMin: job.salary_min,
    salaryMax: job.salary_max,
    eligibility: job.eligibility,
    description: job.description,
    applicationUrl: job.application_url,
    postedAt: job.posted_at,
    source: job.source,
    company: {
      id: job.company_id,
      name: job.company_name,
      logo: job.company_logo
    }
  }));

  const response = {
    status: 'success',
    data: {
      jobs,
      pagination: {
        page,
        limit,
        total,
        totalPages: Math.ceil(total / limit)
      }
    },
    error: null
  };

  // Cache for 5 minutes
  await redis.setex(cacheKey, 300, JSON.stringify(response));

  res.json(response);
}));

// Get single job by ID
router.get('/:id', asyncHandler(async (req, res, next) => {
  const { id } = req.params;

  // Check cache
  const cacheKey = `job:${id}`;
  const cached = await redis.get(cacheKey);
  if (cached) {
    return res.json(JSON.parse(cached));
  }

  const result = await db.query(
    `SELECT j.id, j.title, j.location, j.job_type, j.salary_min, j.salary_max, 
            j.eligibility, j.description, j.application_url, j.posted_at, j.source,
            c.id as company_id, c.name as company_name, c.logo_url as company_logo,
            c.website as company_website
     FROM jobs j
     LEFT JOIN companies c ON j.company_id = c.id
     WHERE j.id = $1 AND j.is_active = true`,
    [id]
  );

  if (result.rows.length === 0) {
    throw new AppError('Job not found', 404, 'JOB_NOT_FOUND');
  }

  const job = result.rows[0];
  const response = {
    status: 'success',
    data: {
      id: job.id,
      title: job.title,
      location: job.location,
      jobType: job.job_type,
      salaryMin: job.salary_min,
      salaryMax: job.salary_max,
      eligibility: job.eligibility,
      description: job.description,
      applicationUrl: job.application_url,
      postedAt: job.posted_at,
      source: job.source,
      company: {
        id: job.company_id,
        name: job.company_name,
        logo: job.company_logo,
        website: job.company_website
      }
    },
    error: null
  };

  // Cache for 10 minutes
  await redis.setex(cacheKey, 600, JSON.stringify(response));

  res.json(response);
}));

// Search jobs (advanced search)
router.post('/search', [
  body('query').optional().trim().escape(),
  body('filters').optional().isObject(),
  body('page').optional().isInt({ min: 1 }),
  body('limit').optional().isInt({ min: 1, max: 100 })
], validateRequest, asyncHandler(async (req, res, next) => {
  const { query: searchQuery, filters = {}, page = 1, limit = 20 } = req.body;
  const offset = (page - 1) * limit;

  // Build advanced search query
  // This would integrate with Elasticsearch in production
  // For now, use PostgreSQL full-text search

  const response = {
    status: 'success',
    data: {
      jobs: [],
      pagination: { page, limit, total: 0, totalPages: 0 }
    },
    error: null
  };

  res.json(response);
}));

// Get recommended jobs for user
router.get('/recommendations/me', asyncHandler(async (req, res, next) => {
  // Get user from auth middleware
  const userId = req.user?.userId;

  if (!userId) {
    return res.json({
      status: 'success',
      data: { jobs: [] },
      error: null
    });
  }

  // Get cached recommendations
  const cacheKey = `recommendations:${userId}`;
  const cached = await redis.get(cacheKey);
  if (cached) {
    return res.json(JSON.parse(cached));
  }

  // Get from AI service
  const response = {
    status: 'success',
    data: { jobs: [] },
    error: null
  };

  // Cache for 15 minutes
  await redis.setex(cacheKey, 900, JSON.stringify(response));

  res.json(response);
}));

module.exports = router;

