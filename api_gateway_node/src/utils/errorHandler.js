/**
 * RojgarSetu 2.0 - Error Handler
 * Zero 404/500 errors - structured error responses
 */

const { logger } = require('./logger');

/**
 * Standard error response format
 */
const errorResponse = (res, statusCode, message, code = 'ERROR') => {
  return res.status(statusCode).json({
    status: 'error',
    data: null,
    error: {
      code: statusCode,
      message: message,
      code: code
    }
  });
};

/**
 * Success response format
 */
const successResponse = (res, data, statusCode = 200) => {
  return res.status(statusCode).json({
    status: 'success',
    data: data,
    error: null
  });
};

/**
 * Custom error class
 */
class AppError extends Error {
  constructor(message, statusCode, errorCode = 'ERROR') {
    super(message);
    this.statusCode = statusCode;
    this.errorCode = errorCode;
    this.isOperational = true;
    Error.captureStackTrace(this, this.constructor);
  }
}

/**
 * Not found handler - 404
 */
const notFoundHandler = (req, res, next) => {
  const error = new AppError(`Route ${req.originalUrl} not found`, 404, 'NOT_FOUND');
  next(error);
};

/**
 * Global error handler
 */
const errorHandler = (err, req, res, next) => {
  // Log the error
  logger.error('Error occurred:', {
    message: err.message,
    stack: err.stack,
    url: req.originalUrl,
    method: req.method,
    body: req.body,
    params: req.params,
    query: req.query
  });

  // Default error
  let statusCode = err.statusCode || 500;
  let message = err.message || 'Internal server error';
  let errorCode = err.errorCode || 'INTERNAL_ERROR';

  // Handle specific error types
  if (err.name === 'ValidationError') {
    statusCode = 400;
    message = 'Validation failed';
    errorCode = 'VALIDATION_ERROR';
  } else if (err.name === 'JsonWebTokenError') {
    statusCode = 401;
    message = 'Invalid token';
    errorCode = 'INVALID_TOKEN';
  } else if (err.name === 'TokenExpiredError') {
    statusCode = 401;
    message = 'Token expired';
    errorCode = 'TOKEN_EXPIRED';
  } else if (err.code === '23505') { // PostgreSQL unique violation
    statusCode = 409;
    message = 'Resource already exists';
    errorCode = 'DUPLICATE_RESOURCE';
  } else if (err.code === '23503') { // PostgreSQL foreign key violation
    statusCode = 400;
    message = 'Invalid reference';
    errorCode = 'INVALID_REFERENCE';
  }

  // Don't leak error details in production
  if (process.env.NODE_ENV === 'production' && statusCode === 500) {
    message = 'Internal server error';
  }

  return res.status(statusCode).json({
    status: 'error',
    data: null,
    error: {
      code: statusCode,
      message: message,
      code: errorCode
    }
  });
};

/**
 * Async handler wrapper - catches async errors
 */
const asyncHandler = (fn) => (req, res, next) => {
  Promise.resolve(fn(req, res, next)).catch(next);
};

module.exports = {
  AppError,
  errorHandler,
  notFoundHandler,
  errorResponse,
  successResponse,
  asyncHandler
};

