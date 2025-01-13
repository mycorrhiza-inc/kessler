import getConfig from "next/config";

// This url is for hitting the kessler api from the client
//

export type RuntimeEnvConfig = {
  public_api_url?: string;
  public_posthog_key?: string;
  public_posthog_host?: string;
  internal_api_url?: string;
  deployment_env?: string;
  flags?: {
    enable_all_features?: boolean;
  };
};

export const runtimeConfig: RuntimeEnvConfig = {
  public_api_url: process.env.PUBLIC_KESSLER_API_URL,
  internal_api_url: process.env.INTERNAL_KESSLER_API_URL,
  public_posthog_key: process.env.PUBLIC_POSTHOG_KEY,
  public_posthog_host: process.env.PUBLIC_POSTHOG_HOST,
  deployment_env: process.env.REACT_APP_ENV || "production",
  flags: {
    enable_all_features: true,
  },
};

export const internalAPIURL = runtimeConfig.internal_api_url;
