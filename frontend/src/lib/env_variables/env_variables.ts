import assert from "assert";
import { env } from "next-runtime-env";
import { z } from "zod";

// Zod schema for runtime environment configuration
export const RuntimeEnvConfigSchema = z.object({
  public_api_url: z.string().nonempty(),
  internal_api_url: z.string().nonempty(),
  public_posthog_key: z.string().nonempty(),
  public_posthog_host: z.string().nonempty(),
  deployment_env: z.string().nonempty(),
  version_hash: z.string().nonempty(),
});

// Type derived from Zod schema
export type RuntimeEnvConfig = z.infer<typeof RuntimeEnvConfigSchema>;

// Empty default config for client fallback
export const emptyRuntimeConfig: RuntimeEnvConfig = {
  public_api_url: "",
  internal_api_url: "",
  public_posthog_key: "",
  public_posthog_host: "",
  deployment_env: "",
  version_hash: "",
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

export const getContextualAPIURL = () => {
  if (window == undefined) {
    const result = env("INTERNAL_KESSLER_API_URL");
    assert(result != undefined, "INTERNAL API URL IS UNDEFINED");
    return result;
  }
  const result = env("PUBLIC_KESSLER_API_URL");
  assert(result != undefined, "PUBLIC API URL IS UNDEFINED");
  return result;
};

export function getUniversalEnvConfig(): RuntimeEnvConfig {
  const rawConfig = {
    public_api_url: removeBackslash(env("PUBLIC_KESSLER_API_URL")),
    internal_api_url: removeBackslash(env("INTERNAL_KESSLER_API_URL")),
    public_posthog_key: env("PUBLIC_POSTHOG_KEY"),
    public_posthog_host: env("PUBLIC_POSTHOG_HOST"),
    deployment_env: env("REACT_APP_ENV") ?? "production",
    version_hash: env("VERSION_HASH") ?? "unknown",
    flags: {
      // next-runtime-env returns strings for env vars, so compare string
      enable_all_features: env("ENABLE_ALL_FEATURES") === "true",
    },
  };

  return RuntimeEnvConfigSchema.parse(rawConfig);
}

/**
 * Client-specific alias for getting runtime env config.
 */
export function getClientRuntimeEnv(): RuntimeEnvConfig {
  return getUniversalEnvConfig();
}
