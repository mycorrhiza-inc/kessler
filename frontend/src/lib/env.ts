import { z } from 'zod';

/**
 * Full environment schema, including server-only and public variables.
 */
export const envSchema = z.object({
  DATABASE_URL: z.string().url(),
  NEXT_PUBLIC_API_BASE: z.string().url(),
  NEXT_PUBLIC_FEATURE_FLAG: z.boolean().optional(),
  // Add more environment variables here
});

type FullEnv = z.infer<typeof envSchema>;

export type { FullEnv };
