/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'export',
  images: {
    unoptimized: true,
    remotePatterns: [
      { protocol: 'https', hostname: 'i.ytimg.com' },
      { protocol: 'https', hostname: 'img.youtube.com' },
      { protocol: 'https', hostname: '**' },
    ],
  },
  // Disabling optimizePackageImports for stability in CI
  experimental: {
    // optimizePackageImports: ['lucide-react'],
  },
  webpack: (config) => {
    // Optimize memory usage for webpack in CI
    config.cache = false;
    return config;
  },
};
module.exports = nextConfig;
