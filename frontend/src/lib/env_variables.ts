import getConfig from "next/config";

// This url is for hitting the kessler api from the client
//

export type RuntimeEnvConfig = {
  public_api_url?: string;
  internal_api_url?: string;
  deployment_env?: string;
  flags?: {
    enable_all_features?: boolean;
  };
};

export const runtimeEnvConfig: RuntimeEnvConfig = {
  public_api_url: process.env.PUBLIC_KESSLER_API_URL,
  internal_api_url: process.env.INTERNAL_KESSLER_API_URL,
  deployment_env: process.env.REACT_APP_ENV || "production",
  flags: {
    enable_all_features: true,
  },
};
