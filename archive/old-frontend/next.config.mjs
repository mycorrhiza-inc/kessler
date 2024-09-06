/** @type {import('next').NextConfig} */
const nextConfig = {
  // transpilePackages: ["@clerk/clerk-react", "@saas-ui/clerk"],
  experimental: {
    // You may not need this, it's just to support moduleResolution: 'node16'
    extensionAlias: {
      ".js": [".tsx", ".ts", ".jsx", ".js"],
    },
    turbo: {
      resolveAlias: {
        // Turbopack does not support standard ESM import paths yet
        "./Sample.js": "./app/Sample.tsx",
        /**
         * Critical: prevents " ⨯ ./node_modules/canvas/build/Release/canvas.node
         * Module parse failed: Unexpected character '�' (1:0)" error
         */
        canvas: "./empty-module.ts",
      },
    },
  },
  /**
   * Critical: prevents ''import', and 'export' cannot be used outside of module code" error
   * See https://github.com/vercel/next.js/pull/66817
   */
  swcMinify: false,
  env: {
    NEXT_PUBLIC_BASE_URL:
      process.env.NEXT_PUBLIC_BASE_URL || "http://app.kessler.xyz",
  },
  async rewrites() {
    return [
      {
        source: "/:path*",
        destination: `https://app.kessler.xyz/:path*`,
      },
    ];
  },
};
