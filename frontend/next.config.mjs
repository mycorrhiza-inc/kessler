/** @type {import('next').NextConfig} */
// import * as Sentry from "@sentry/nextjs";
// const SENTRY_DSN = process.env.SENTRY_DSN || process.env.NEXT_PUBLIC_SENTRY_DSN;
// Sentry.init({
//   dsn:
//     SENTRY_DSN ||
//     "https://f7a73c4f60af41e08f8b11a62c734d5d@glitchtip.mycor.io/1",
// });
// Its a shitty idea to leave this in git, but it should be fine for testing
import { withSentryConfig } from "@sentry/nextjs";

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

export default withSentryConfig(nextConfig, {
  // For all available options, see:
  // https://github.com/getsentry/sentry-webpack-plugin#options

  org: "mycorrhiza",
  project: "kessler",
  sentryUrl: "https://glitchtip.mycor.io/",

  // Only print logs for uploading source maps in CI
  silent: !process.env.CI,

  // For all available options, see:
  // https://docs.sentry.io/platforms/javascript/guides/nextjs/manual-setup/

  // Upload a larger set of source maps for prettier stack traces (increases build time)
  widenClientFileUpload: true,

  // Route browser requests to Sentry through a Next.js rewrite to circumvent ad-blockers.
  // This can increase your server load as well as your hosting bill.
  // Note: Check that the configured route will not match with your Next.js middleware, otherwise reporting of client-
  // side errors will fail.
  tunnelRoute: "/monitoring",

  // Hides source maps from generated client bundles
  hideSourceMaps: true,

  // Automatically tree-shake Sentry logger statements to reduce bundle size
  disableLogger: true,

  // Enables automatic instrumentation of Vercel Cron Monitors. (Does not yet work with App Router route handlers.)
  // See the following for more information:
  // https://docs.sentry.io/product/crons/
  // https://vercel.com/docs/cron-jobs
  automaticVercelMonitors: true,
});
