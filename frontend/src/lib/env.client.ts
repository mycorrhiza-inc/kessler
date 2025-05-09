import { z } from 'zod';
import { envSchema } from './env';

// Pick only public vars
const clientSchema = envSchema.pick({
  NEXT_PUBLIC_API_BASE: true,
  NEXT_PUBLIC_FEATURE_FLAG: true,
});

declare global {
  interface Window {
    __ENV__?: Record<string, any>;
  }
}

const raw = window.__ENV__;
const result = clientSchema.safeParse(raw);
if (!result.success) {
  console.error('❌ Client env validation error:', result.error.format());
  throw new Error('Invalid client environment variables');
}

export const clientEnv = result.data;
