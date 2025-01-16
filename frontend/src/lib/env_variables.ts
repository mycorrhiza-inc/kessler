import getConfig from "next/config";

// This url is for hitting the kessler api from the client
//

// All of these clients end up as json on the frontend and can be accessed by anyone looking at the html. Dont put sensitive data here
export type RuntimeEnvConfig = {
  public_api_url?: string;
  public_posthog_key?: string;
  public_posthog_host?: string;
  internal_api_url?: string;
  deployment_env?: string;
  version_hash?: string;
  flags?: {
    enable_all_features?: boolean;
  };
};

const removeBackslash = (val: string | undefined): string => {
  if (!val) return "";
  if (val.endsWith("/")) {
    return val.slice(0, -1);
  }
  return val;
};

export const runtimeConfig: RuntimeEnvConfig = {
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

export const emptyRuntimeConfig: RuntimeEnvConfig = {
  public_api_url: "",
  internal_api_url: "",
  public_posthog_key: "",
  public_posthog_host: "",
  deployment_env: "",
  version_hash: "",
  flags: {
    enable_all_features: true,
  },
};

export const internalAPIURL = runtimeConfig.internal_api_url;
