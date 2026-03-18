/**
 * RojgarSetu 2.0 - Auth Routes
 * Register, Login, Logout, Refresh Token
 */

const express = require('express');
const { body, validationResult } = require('express-validator');
const jwt = require('jsonwebtoken');
const bcrypt = require('bcryptjs');
const { asyncHandler, AppError } = require('../utils/errorHandler');
const { logger } = require('../utils/logger');
const config = require('../config');
const db = require('../database');

const router = express.Router();

// Validation middleware
const validateRequest = (req, res, next) => {
  const errors = validationResult(req);
  if (!errors.isEmpty()) {
    return res.status(400).json({
      status: 'error',
      data: null,
      error: {
        code: 400,
        message: 'Validation failed',
        details: errors.array()
      }
    });
  }
  next();
};

// Register
router.post('/register', [
  body('name').trim().notEmpty().withMessage('Name is required'),
  body('email').isEmail().normalizeEmail().withMessage('Valid email is required'),
  body('password').isLength({ min: 8 }).withMessage('Password must be at least 8 characters')
], validateRequest, asyncHandler(async (req, res, next) => {
  const { name, email, password } = req.body;

  // Check if user exists
  const existingUser = await db.query(
    'SELECT id FROM users WHERE email = $1',
    [email]
  );

  if (existingUser.rows.length > 0) {
    throw new AppError('Email already registered', 409, 'EMAIL_EXISTS');
  }

  // Hash password
  const salt = await bcrypt.genSalt(10);
  const passwordHash = await bcrypt.hash(password, salt);

  // Create user
  const result = await db.query(
    `INSERT INTO users (name, email, password_hash) 
     VALUES ($1, $2, $3) 
     RETURNING id, name, email, created_at`,
    [name, email, passwordHash]
  );

  const user = result.rows[0];

  // Generate tokens
  const accessToken = jwt.sign(
    { userId: user.id, email: user.email },
    config.jwt.secret,
    { expiresIn: config.jwt.expiresIn }
  );

  const refreshToken = jwt.sign(
    { userId: user.id },
    config.jwt.secret,
    { expiresIn: config.jwt.refreshExpiresIn }
  );

  logger.info('User registered', { userId: user.id, email });

  res.status(201).json({
    status: 'success',
    data: {
      user: {
        id: user.id,
        name: user.name,
        email: user.email,
        createdAt: user.created_at
      },
      accessToken,
      refreshToken
    },
    error: null
  });
}));

// Login
router.post('/login', [
  body('email').isEmail().normalizeEmail().withMessage('Valid email is required'),
  body('password').notEmpty().withMessage('Password is required')
], validateRequest, asyncHandler(async (req, res, next) => {
  const { email, password } = req.body;

  // Find user
  const result = await db.query(
    'SELECT id, name, email, password_hash FROM users WHERE email = $1',
    [email]
  );

  if (result.rows.length === 0) {
    throw new AppError('Invalid credentials', 401, 'INVALID_CREDENTIALS');
  }

  const user = result.rows[0];

  // Check password
  const isMatch = await bcrypt.compare(password, user.password_hash);
  if (!isMatch) {
    throw new AppError('Invalid credentials', 401, 'INVALID_CREDENTIALS');
  }

  // Generate tokens
  const accessToken = jwt.sign(
    { userId: user.id, email: user.email },
    config.jwt.secret,
    { expiresIn: config.jwt.expiresIn }
  );

  const refreshToken = jwt.sign(
    { userId: user.id },
    config.jwt.secret,
    { expiresIn: config.jwt.refreshExpiresIn }
  );

  logger.info('User logged in', { userId: user.id });

  res.json({
    status: 'success',
    data: {
      user: {
        id: user.id,
        name: user.name,
        email: user.email
      },
      accessToken,
      refreshToken
    },
    error: null
  });
}));

// Refresh token
router.post('/refresh', [
  body('refreshToken').notEmpty().withMessage('Refresh token is required')
], validateRequest, asyncHandler(async (req, res, next) => {
  const { refreshToken } = req.body;

  try {
    const decoded = jwt.verify(refreshToken, config.jwt.secret);
    
    // Generate new access token
    const accessToken = jwt.sign(
      { userId: decoded.userId },
      config.jwt.secret,
      { expiresIn: config.jwt.expiresIn }
    );

    res.json({
      status: 'success',
      data: { accessToken },
      error: null
    });
  } catch (err) {
    throw new AppError('Invalid refresh token', 401, 'INVALID_TOKEN');
  }
}));

// Logout (client-side token removal, but we can blacklist if needed)
router.post('/logout', asyncHandler(async (req, res, next) => {
  // In production, you might want to blacklist the token
  logger.info('User logged out');

  res.json({
    status: 'success',
    data: { message: 'Logged out successfully' },
    error: null
  });
}));

module.exports = router;

