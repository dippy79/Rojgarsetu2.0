const path = require('path');

/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'export',
  // Silencing workspace root warning
  experimental: {
    // optimizePackageImports: ['lucide-react'],
  },
  turbopack: {
    root: path.resolve(__dirname, '..'),
  },
  allowedDevOrigins: ['169.254.148.199'],
  images: {
    unoptimized: true,
    remotePatterns: [
      { protocol: 'https', hostname: 'i.ytimg.com' },
      { protocol: 'https', hostname: 'img.youtube.com' },
      { protocol: 'https', hostname: '**' },
    ],
  },
  webpack: (config) => {
    config.cache = false;
    return config;
  },
};
module.exports = nextConfig;
