import { z } from "zod";

const envSchema = z.object({
  public_api_url: z.string().min(1),
  internal_api_url: z.string().min(1),
  public_posthog_key: z.string().min(1),
  public_posthog_host: z.string().min(1),
});

export type EnvConfig = z.infer<typeof envSchema>;

const defaults: EnvConfig = {
  public_api_url: "http://localhost",
  internal_api_url: "http://backend-server:4041",
  public_posthog_host: "REPLACE THIS",
  public_posthog_key: "REPLACE THIS",
};

function getEnvConfig(): EnvConfig {
  // Server-side
  if (typeof window === "undefined") {
    return {
      ...defaults,
      internal_api_url: process.env.INTERNAL_KESSLER_API_URL || defaults.internal_api_url,
    };
  }

  // Client-side
  try {
    const config = {
      ...defaults,
      public_api_url: process.env.NEXT_PUBLIC_KESSLER_API_URL?.replace(/\/$/, "") || defaults.public_api_url,
      public_posthog_key: process.env.NEXT_PUBLIC_POSTHOG_KEY || defaults.public_posthog_key,
      public_posthog_host: process.env.NEXT_PUBLIC_POSTHOG_HOST || defaults.public_posthog_host,
    };

    return envSchema.parse(config);
  } catch (error) {
    console.warn("Environment validation failed, using defaults:", error);
    return defaults;
  }
}

export function getAPIURL(): string {
  const config = getEnvConfig();
  return typeof window === "undefined"
    ? config.internal_api_url
    : config.public_api_url;
}

export { getEnvConfig };
