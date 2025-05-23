import assert from "assert";
import { env } from "next-runtime-env";
import { z } from "zod";

// Zod schema for runtime environment configuration
export const RuntimeEnvConfigSchema = z.object({
  public_api_url: z.string().nonempty(),
  internal_api_url: z.string().nonempty(),
  public_posthog_key: z.string().nonempty(),
  public_posthog_host: z.string().nonempty(),
  // deployment_env: z.string().nonempty(),
  // version_hash: z.string().nonempty(),
});

// Type derived from Zod schema
export type RuntimeEnvConfig = z.infer<typeof RuntimeEnvConfigSchema>;

// Empty default config for client fallback
export const emptyRuntimeConfig: RuntimeEnvConfig = {
  public_api_url: "",
  internal_api_url: "",
  public_posthog_key: "",
  public_posthog_host: "",
  // deployment_env: "",
  // version_hash: "",
};

// Helper to trim trailing slash
const removeBackslash = (val: string | undefined): string => {
  if (!val) return "";
  return val.endsWith("/") ? val.slice(0, -1) : val;
};

/**
 * Universal getter for runtime config, works on server and client.
 * Uses next-runtime-env under the hood and validates via Zod schema.
 */

export const getContextualAPIURL = (): String => {
  if (typeof window === "undefined") {
    const result = getUniversalEnvConfig().internal_api_url;
    return result;
  }
  const result = getUniversalEnvConfig().public_api_url;
  return result;
};

export function getUniversalEnvConfig(): RuntimeEnvConfig {
  if (typeof window === "undefined") {
    const serverConfig: RuntimeEnvConfig = {
      public_api_url: process.env.NEXT_PUBLIC_KESSLER_API_URL || "unknownurl",
      internal_api_url:
        process.env.INTERNAL_KESSLER_API_URL || "http://backend-go:4041",
      public_posthog_host: "REPLACE THIS",
      public_posthog_key: "REPLACE THIS",
    };
    return serverConfig;
  }
  const defaultConfig: RuntimeEnvConfig = {
    public_api_url: "http://localhost",
    internal_api_url: "http://backend-go:4041",
    public_posthog_host: "REPLACE THIS",
    public_posthog_key: "REPLACE THIS",
    // deployment_env: "dev",
  };
  try {
    const rawConfig = {
      public_api_url: removeBackslash(env("NEXT_PUBLIC_KESSLER_API_URL")),
      internal_api_url: "http://backend-go:4041",
      public_posthog_key: env("NEXT_PUBLIC_POSTHOG_KEY") || "unknown",
      public_posthog_host: env("NEXT_PUBLIC_POSTHOG_HOST") || "unknown",
      // version_hash: env("VERSION_HASH") ?? "unknown",
    };
    return RuntimeEnvConfigSchema.parse(rawConfig);
  } catch (error) {
    console.log(
      "FAILED TO GENERATE UNIVERSAL ENVIORNMENT VARIABLES, FALLING BACK TO DEFAULTS.",
    );
    console.log(error);
    return defaultConfig;
  }
}

/**
 * Client-specific alias for getting runtime env config.
 */
export function getClientRuntimeEnv(): RuntimeEnvConfig {
  return getUniversalEnvConfig();
}
