/**
 * RojgarSetu 2.0 - Database Connection
 * PostgreSQL with connection pooling
 */

const { Pool } = require('pg');
const config = require('../config');
const { logger } = require('../utils/logger');

const pool = new Pool({
  host: config.db.host,
  port: config.db.port,
  database: config.db.database,
  user: config.db.user,
  password: config.db.password,
  max: config.db.max,
  idleTimeoutMillis: config.db.idleTimeoutMillis,
  connectionTimeoutMillis: config.db.connectionTimeoutMillis
});

// Test connection
pool.on('connect', () => {
  logger.info('New database connection established');
});

pool.on('error', (err) => {
  logger.error('Unexpected database error', { error: err.message });
});

// Query helper
const query = async (text, params) => {
  const start = Date.now();
  const res = await pool.query(text, params);
  const duration = Date.now() - start;
  
  logger.debug('Executed query', { text, duration, rows: res.rowCount });
  return res;
};

// Transaction helper
const transaction = async (callback) => {
  const client = await pool.connect();
  try {
    await client.query('BEGIN');
    const result = await callback(client);
    await client.query('COMMIT');
    return result;
  } catch (error) {
    await client.query('ROLLBACK');
    throw error;
  } finally {
    client.release();
  }
};

// Health check
const healthCheck = async () => {
  try {
    const result = await pool.query('SELECT 1');
    return result.rowCount === 1;
  } catch (error) {
    logger.error('Database health check failed', { error: error.message });
    return false;
  }
};

module.exports = {
  query,
  pool,
  transaction,
  healthCheck
};

