/**
 * RojgarSetu 2.0 - Recommendations Routes
 * AI-powered job and course recommendations
 */

const express = require('express');
const { asyncHandler, AppError } = require('../utils/errorHandler');
const db = require('../database');
const redis = require('../utils/redis');
const axios = require('axios');
const config = require('../config');

const router = express.Router();

// Get personalized job recommendations
router.get('/jobs', asyncHandler(async (req, res, next) => {
  const userId = req.user?.userId;
  
  if (!userId) {
    return res.json({
      status: 'success',
      data: { jobs: [], message: 'Login for personalized recommendations' },
      error: null
    });
  }

  // Check cache
  const cacheKey = `job-recommendations:${userId}`;
  const cached = await redis.get(cacheKey);
  if (cached) {
    return res.json(JSON.parse(cached));
  }

  try {
    // Call AI engine service
    const aiResponse = await axios.post(`${config.services.aiEngine}/recommend/jobs`, {
      userId
    }, {
      timeout: 5000
    });

    const response = {
      status: 'success',
      data: aiResponse.data,
      error: null
    };

    await redis.setex(cacheKey, 900, JSON.stringify(response));
    res.json(response);
  } catch (err) {
    // Fallback: Get jobs based on user skills
    const userResult = await db.query('SELECT skills FROM users WHERE id = $1', [userId]);
    const userSkills = userResult.rows[0]?.skills || [];

    if (userSkills.length === 0) {
      return res.json({
        status: 'success',
        data: { jobs: [] },
        error: null
      });
    }

    const jobsResult = await db.query(
      `SELECT j.id, j.title, j.location, j.job_type, j.salary_min, j.salary_max,
              c.name as company_name
       FROM jobs j
       LEFT JOIN companies c ON j.company_id = c.id
       WHERE j.is_active = true
       ORDER BY j.posted_at DESC
       LIMIT 10`
    );

    const response = {
      status: 'success',
      data: { jobs: jobsResult.rows, source: 'fallback' },
      error: null
    };

    await redis.setex(cacheKey, 300, JSON.stringify(response));
    res.json(response);
  }
}));

// Get personalized course recommendations
router.get('/courses', asyncHandler(async (req, res, next) => {
  const userId = req.user?.userId;
  
  if (!userId) {
    return res.json({
      status: 'success',
      data: { courses: [] },
      error: null
    });
  }

  const cacheKey = `course-recommendations:${userId}`;
  const cached = await redis.get(cacheKey);
  if (cached) {
    return res.json(JSON.parse(cached));
  }

  try {
    const aiResponse = await axios.post(`${config.services.aiEngine}/recommend/courses`, {
      userId
    }, {
      timeout: 5000
    });

    const response = {
      status: 'success',
      data: aiResponse.data,
      error: null
    };

    await redis.setex(cacheKey, 900, JSON.stringify(response));
    res.json(response);
  } catch (err) {
    // Fallback
    const result = await db.query(
      `SELECT * FROM courses ORDER BY created_at DESC LIMIT 10`
    );

    const response = {
      status: 'success',
      data: { courses: result.rows, source: 'fallback' },
      error: null
    };

    await redis.setex(cacheKey, 300, JSON.stringify(response));
    res.json(response);
  }
}));

// Refresh recommendations (trigger AI recalculation)
router.post('/refresh', asyncHandler(async (req, res, next) => {
  const userId = req.user?.userId;
  
  if (!userId) {
    throw new AppError('Authentication required', 401, 'UNAUTHORIZED');
  }

  // Invalidate cache
  await redis.del(`job-recommendations:${userId}`);
  await redis.del(`course-recommendations:${userId}`);

  res.json({
    status: 'success',
    data: { message: 'Recommendations cache cleared, will be regenerated on next request' },
    error: null
  });
}));

module.exports = router;

