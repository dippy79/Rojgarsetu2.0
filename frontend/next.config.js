/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'export',
  // Silencing workspace root warning
  experimental: {
    // optimizePackageImports: ['lucide-react'],
  },
  // In Next.js 15+, turbopack root can be set here or inferred.
  // The warning suggests setting it directly.
  // @ts-ignore
  turbopack: {
    root: '..',
  },
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
