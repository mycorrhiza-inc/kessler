/** @type {import('next').NextConfig} */
const nextConfig = {
  transpilePackages: ["@clerk/clerk-react", "@saas-ui/clerk"],
  env: {
    NEXT_PUBLIC_BASE_URL: process.env.NEXT_PUBLIC_BASE_URL || 'http://localhost:3000',
  },
};

export default nextConfig;
