/**
 * RojgarSetu 2.0 - User Routes
 * Profile, applications, settings
 */

const express = require('express');
const { body, validationResult } = require('express-validator');
const { asyncHandler, AppError } = require('../utils/errorHandler');
const db = require('../database');
const redis = require('../utils/redis');

const router = express.Router();

// Middleware to check authentication
const requireAuth = asyncHandler(async (req, res, next) => {
  // Token validation would happen here via middleware
  if (!req.user) {
    throw new AppError('Unauthorized', 401, 'UNAUTHORIZED');
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

// Get user profile
router.get('/profile', requireAuth, asyncHandler(async (req, res, next) => {
  const userId = req.user.userId;

  const result = await db.query(
    `SELECT id, name, email, phone, resume_url, skills, created_at
     FROM users WHERE id = $1`,
    [userId]
  );

  if (result.rows.length === 0) {
    throw new AppError('User not found', 404, 'USER_NOT_FOUND');
  }

  const user = result.rows[0];
  res.json({
    status: 'success',
    data: {
      id: user.id,
      name: user.name,
      email: user.email,
      phone: user.phone,
      resumeUrl: user.resume_url,
      skills: user.skills,
      createdAt: user.created_at
    },
    error: null
  });
}));

// Update profile
router.put('/profile', requireAuth, [
  body('name').optional().trim().notEmpty(),
  body('phone').optional().trim(),
  body('skills').optional().isArray(),
  body('resumeUrl').optional().isURL()
], validateRequest, asyncHandler(async (req, res, next) => {
  const userId = req.user.userId;
  const { name, phone, skills, resumeUrl } = req.body;

  const updates = [];
  const params = [];
  let paramIndex = 1;

  if (name) {
    updates.push(`name = $${paramIndex++}`);
    params.push(name);
  }
  if (phone) {
    updates.push(`phone = $${paramIndex++}`);
    params.push(phone);
  }
  if (skills) {
    updates.push(`skills = $${paramIndex++}`);
    params.push(skills);
  }
  if (resumeUrl) {
    updates.push(`resume_url = $${paramIndex++}`);
    params.push(resumeUrl);
  }

  if (updates.length === 0) {
    throw new AppError('No fields to update', 400, 'NO_FIELDS');
  }

  params.push(userId);
  const result = await db.query(
    `UPDATE users SET ${updates.join(', ')} WHERE id = $${paramIndex}
     RETURNING id, name, email, phone, resume_url, skills`,
    params
  );

  // Invalidate cache
  await redis.del(`user:${userId}`);

  const user = result.rows[0];
  res.json({
    status: 'success',
    data: {
      id: user.id,
      name: user.name,
      email: user.email,
      phone: user.phone,
      resumeUrl: user.resume_url,
      skills: user.skills
    },
    error: null
  });
}));

// Apply for job
router.post('/apply-job', requireAuth, [
  body('jobId').notEmpty().withMessage('Job ID is required')
], validateRequest, asyncHandler(async (req, res, next) => {
  const userId = req.user.userId;
  const { jobId } = req.body;

  // Check if job exists
  const jobResult = await db.query('SELECT id FROM jobs WHERE id = $1 AND is_active = true', [jobId]);
  if (jobResult.rows.length === 0) {
    throw new AppError('Job not found or not active', 404, 'JOB_NOT_FOUND');
  }

  // Check if already applied
  const existingApp = await db.query(
    'SELECT id FROM applications WHERE user_id = $1 AND job_id = $2',
    [userId, jobId]
  );

  if (existingApp.rows.length > 0) {
    throw new AppError('Already applied to this job', 409, 'ALREADY_APPLIED');
  }

  // Create application
  const result = await db.query(
    `INSERT INTO applications (user_id, job_id, status) 
     VALUES ($1, $2, 'applied') 
     RETURNING id, user_id, job_id, status, applied_at`,
    [userId, jobId]
  );

  // Create notification
  await db.query(
    `INSERT INTO notifications (user_id, title, message, type)
     VALUES ($1, $2, $3, $4)`,
    [userId, 'Job Application Submitted', 'Your application has been submitted successfully', 'application_update']
  );

  const app = result.rows[0];
  res.status(201).json({
    status: 'success',
    data: {
      id: app.id,
      jobId: app.job_id,
      status: app.status,
      appliedAt: app.applied_at
    },
    error: null
  });
}));

// Get user applications
router.get('/applications', requireAuth, asyncHandler(async (req, res, next) => {
  const userId = req.user.userId;

  const result = await db.query(
    `SELECT a.id, a.status, a.applied_at, a.updated_at,
            j.id as job_id, j.title as job_title, j.location as job_location,
            c.name as company_name
     FROM applications a
     JOIN jobs j ON a.job_id = j.id
     LEFT JOIN companies c ON j.company_id = c.id
     WHERE a.user_id = $1
     ORDER BY a.applied_at DESC`,
    [userId]
  );

  const applications = result.rows.map(a => ({
    id: a.id,
    status: a.status,
    appliedAt: a.applied_at,
    updatedAt: a.updated_at,
    job: {
      id: a.job_id,
      title: a.job_title,
      location: a.job_location,
      company: a.company_name
    }
  }));

  res.json({
    status: 'success',
    data: { applications },
    error: null
  });
}));

// Get user notifications
router.get('/notifications', requireAuth, asyncHandler(async (req, res, next) => {
  const userId = req.user.userId;
  const { unreadOnly } = req.query;

  let query = `SELECT id, title, message, type, is_read, data, created_at 
               FROM notifications WHERE user_id = $1`;
  
  if (unreadOnly === 'true') {
    query += ' AND is_read = false';
  }
  
  query += ' ORDER BY created_at DESC LIMIT 50';

  const result = await db.query(query, [userId]);

  const notifications = result.rows.map(n => ({
    id: n.id,
    title: n.title,
    message: n.message,
    type: n.type,
    isRead: n.is_read,
    data: n.data,
    createdAt: n.created_at
  }));

  res.json({
    status: 'success',
    data: { notifications },
    error: null
  });
}));

// Mark notification as read
router.put('/notifications/:id/read', requireAuth, asyncHandler(async (req, res, next) => {
  const userId = req.user.userId;
  const { id } = req.params;

  await db.query(
    'UPDATE notifications SET is_read = true WHERE id = $1 AND user_id = $2',
    [id, userId]
  );

  res.json({
    status: 'success',
    data: { message: 'Notification marked as read' },
    error: null
  });
}));

module.exports = router;

