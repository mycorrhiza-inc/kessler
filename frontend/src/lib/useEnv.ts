import { clientEnv } from './env.client';

/**
 * Hook to access client environment variables.
 */
export function useEnv() {
  return clientEnv;
}
