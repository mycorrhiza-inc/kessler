/** @type {import('next').NextConfig} */
const nextConfig = {
  publicRuntimeConfig: {
    // Will be available on both server and client
    KESSLER_API_URL: process.env.KESSLER_API_URL,
  },
  serverRuntimeConfig: {
    // Will only be available on the server side
    INTERNAL_KESSLER_API_URL: process.env.INTERNAL_KESSLER_API_URL,
  },
};

module.exports = nextConfig;
