import { envSchema } from './env';

// Validate server-side environment variables
const parsed = envSchema.safeParse(process.env);
if (!parsed.success) {
  console.error('❌ Server env validation error:', parsed.error.format());
  throw new Error('Invalid server environment variables');
}

export const serverEnv = parsed.data;
