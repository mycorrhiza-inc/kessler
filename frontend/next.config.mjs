/** @type {import('next').NextConfig} */
const nextConfig = {
  // transpilePackages: ["@clerk/clerk-react", "@saas-ui/clerk"],
  env: {
    NEXT_PUBLIC_BASE_URL: process.env.NEXT_PUBLIC_BASE_URL || 'http://app.kessler.xyz',
  },
  async rewrites() {
    return [
      {
        source: '/:path*',
        destination: `https://app.kessler.xyz/:path*`,
      },
    ];
  },
};

export default nextConfig;
