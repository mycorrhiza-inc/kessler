/** @type {import('next').NextConfig} */
const { makeEnvPublic } = require("next-runtime-env");

// Make all relevant environment variables public

const nextConfig = {
  // No need for publicRuntimeConfig/serverRuntimeConfig anymore
  // All vars are exposed via next-runtime-env at runtime
  experimental: {
    runtime: "experimental-edge", // or other if using edge functions
  },
};

module.exports = nextConfig;

makeEnvPublic([
  "NEXT_PUBLIC_KESSLER_API_URL",
  "INTERNAL_KESSLER_API_URL",
  "NEXT_PUBLIC_POSTHOG_KEY",
  "NEXT_PUBLIC_POSTHOG_HOST",
  "REACT_APP_ENV",
  "VERSION_HASH",
  "ENABLE_ALL_FEATURES",
]);
