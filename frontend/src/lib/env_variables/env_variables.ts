import getConfig from "next/config";
import { z } from "zod";

// Zod schema for runtime environment configuration
export const RuntimeEnvConfigSchema = z.object({
  public_api_url: z.string().optional(),
  internal_api_url: z.string().optional(),
  public_posthog_key: z.string().optional(),
  public_posthog_host: z.string().optional(),
  deployment_env: z.string().optional(),
  version_hash: z.string().optional(),
  flags: z
    .object({
      enable_all_features: z.boolean().optional(),
    })
    .optional(),
});

// Type derived from Zod schema
export type RuntimeEnvConfig = z.infer<typeof RuntimeEnvConfigSchema>;

// Script element ID for embedding config in HTML
export const envScriptId = "env-config";

// Helper to trim trailing slash
const removeBackslash = (val: string | undefined): string => {
  if (!val) return "";
  if (val.endsWith("/")) {
    return val.slice(0, -1);
  }
  return val;
};

// Raw runtime configuration (may contain undefined values)
const rawRuntimeConfig = {
  public_api_url: removeBackslash(process.env.PUBLIC_KESSLER_API_URL),
  internal_api_url: removeBackslash(process.env.INTERNAL_KESSLER_API_URL),
  public_posthog_key: process.env.PUBLIC_POSTHOG_KEY,
  public_posthog_host: process.env.PUBLIC_POSTHOG_HOST,
  deployment_env: process.env.REACT_APP_ENV || "production",
  version_hash: process.env.VERSION_HASH || "unknown",
  flags: {
    enable_all_features: true,
  },
};

// Validated runtimeConfig for server usage
export const runtimeConfig: RuntimeEnvConfig =
  RuntimeEnvConfigSchema.parse(rawRuntimeConfig);

// Empty default config for client fallback
export const emptyRuntimeConfig: RuntimeEnvConfig =
  RuntimeEnvConfigSchema.parse({
    public_api_url: "",
    internal_api_url: "",
    public_posthog_key: "",
    public_posthog_host: "",
    deployment_env: "",
    version_hash: "",
    flags: { enable_all_features: true },
  });

// Universal getter for runtime config, works on client and server
export function getUniversalEnvConfig(): RuntimeEnvConfig {
  if (typeof window === "undefined") {
    // Server-side: return validated config directly
    return runtimeConfig;
  }
  // Client-side: read from embedded <script> tag
  const script = window.document.getElementById(
    envScriptId,
  ) as HTMLScriptElement;
  const raw = script ? JSON.parse(script.innerText) : {};
  return RuntimeEnvConfigSchema.parse(raw);
}
