const nextJest = require('next/jest');
const path = require('node:path');

const createJestConfig = nextJest({
  dir: path.resolve(__dirname),
});

/** @type {import('jest').Config} */
const customJestConfig = {
  testEnvironment: 'jsdom',
  setupFilesAfterEnv: ['<rootDir>/jest.setup.js'],
};

module.exports = createJestConfig(customJestConfig);