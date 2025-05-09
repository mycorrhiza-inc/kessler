# Frontend Environment Refactor

## Overview

We need a clean, type-safe, and secure way to share environment variables between server and client code in the frontend. The current hack pulls in `process.env` at runtime and hydrates it via scripts. Below are architectural ideas for refactoring this setup.

---

## Goals

- **Type Safety:** Use schemas (e.g. Zod) to validate/parse all env vars at startup.
- **Separation of Concerns:** Clearly distinguish server-only vs. client-safe variables.
- **Single Source of Truth:** Define env names/types once; re-use across client/server.
- **Secure Exposure:** Only expose a vetted subset (e.g. prefixed with `NEXT_PUBLIC_`) to the browser.
- **SSR & CSR Compatible:** Support Next.js getServerSideProps and client components alike.
- **Minimal Bundle Impact:** Don’t bloat the client bundle with unused vars.

---

## 1. Central Type Definitions & Schema

Create `src/lib/env.ts`:

```typescript
import { z } from 'zod';

// 1. Declare full schema (both server & client entries)
const envSchema = z.object({
  DATABASE_URL: z.string().url(),
  NEXT_PUBLIC_API_BASE: z.string().url(),
  NEXT_PUBLIC_FEATURE_FLAG: z.boolean().optional(),
  // ...add more
});

type FullEnv = z.infer<typeof envSchema>;

export { envSchema, FullEnv };
```

**Benefits:**
- Schema enforces required keys & types.
- Single source of truth for var names.

---

## 2. Separate Server & Client Entrypoints

- `src/lib/env.server.ts`
- `src/lib/env.client.ts`

```typescript
// env.server.ts
import { envSchema } from './env';

const _parsed = envSchema.safeParse(process.env);
if (!_parsed.success) {
  console.error(_parsed.error.format());
  throw new Error('Invalid server env vars');
}
export const serverEnv = _parsed.data;
```

```typescript
// env.client.ts
import { z } from 'zod';
import { envSchema } from './env';

// Build-time generated JSON stub:
// window.__ENV__ = { NEXT_PUBLIC_API_BASE: "...", NEXT_PUBLIC_FEATURE_FLAG: true };

const clientSchema = envSchema.pick({
  NEXT_PUBLIC_API_BASE: true,
  NEXT_PUBLIC_FEATURE_FLAG: true,
});

const raw = (window as any).__ENV__;
const result = clientSchema.safeParse(raw);
if (!result.success) {
  console.error('Client env parse error', result.error.format());
  throw new Error('Invalid client env');
}
export const clientEnv = result.data;
```

**Load Flow:**
1. **Server build/runtime:** `process.env` & Zod parsing.
2. **Client bundle:** Inlined small JSON of only public vars.

---

## 3. Build-Time JSON Generation

Add a script (`scripts/gen-env-json.ts`) that:
1. Reads `process.env`.
2. Filters by `NEXT_PUBLIC_` prefix.
3. Writes `public/env.json` under `frontend/public`.

Then in your HTML root template (e.g. `_document.tsx`), inject:

```html
<script>
  window.__ENV__ = %ENV_JSON%;
</script>
```

#### Example `gen-env-json.ts`
```ts
import fs from 'fs';

const PREFIX = 'NEXT_PUBLIC_';
const obj = Object.entries(process.env)
  .filter(([k]) => k.startsWith(PREFIX))
  .reduce<Record<string,string>>((acc, [k,v]) => { acc[k]=v!; return acc; }, {});

fs.writeFileSync(
  './frontend/public/env.json',
  JSON.stringify(obj, null, 2)
);
```

Hook into `package.json` build step:
```json
"scripts": {
  "prebuild": "ts-node scripts/gen-env-json.ts",
  "build": "next build"
}
```

---

## 4. Runtime Injection in Next.js

In `pages/_document.tsx` (or in App router `layout.tsx`):

```tsx
import fs from 'fs';
import path from 'path';

export default function Document() {
  const envJson = fs.readFileSync(
    path.resolve('./public/env.json'),
    'utf-8'
  );

  return (
    <Html>
      <Head />
      <body>
        <script
          dangerouslySetInnerHTML={{ __html: `window.__ENV__ = ${envJson}` }}
        />
        <Main />
        <NextScript />
      </body>
    </Html>
  );
}
```

---

## 5. Client Hook & Provider

```ts
// useEnv.ts
import { clientEnv } from './env.client';
export function useEnv() {
  return clientEnv;
}
```

Use in React:
```tsx
import { useEnv } from 'src/lib/useEnv';

export default function Page() {
  const { NEXT_PUBLIC_API_BASE } = useEnv();
  // ...
}
```

---

## 6. Alternative Patterns

- **Next.js `publicRuntimeConfig`:** Define in `next.config.js`, accessible via `getConfig()`.
- **`import.meta.env` (Vite/RSC):** Leverage Vite's built-in env system when migrating.
- **Edge Functions:** Ensure JSON injection works under the edge runtime.

---

## 7. CI & Validation

- Add a CI step to run `ts-node scripts/gen-env-json.ts` and catch errors.
- Lint against missing keys.

---

## Summary

This approach:
- Centralizes schema & types.
- Enforces validation at build/runtime.
- Cleanly separates server-only vs. client-safe vars.
- Keeps client bundle minimal.

Feel free to iterate on this, or discuss alternatives!