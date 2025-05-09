import { clientEnv } from "./env_variables/env.client";

/**
 * Hook to access client environment variables.
 */
export function useEnv() {
  return clientEnv;
}
