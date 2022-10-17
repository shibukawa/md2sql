/** @type {import('next').NextConfig} */
const nextConfig = {
  basePath: '/md2sql',
  reactStrictMode: true,
  swcMinify: true,
  images: {
    unoptimized: true
  }
}

module.exports = nextConfig
