Could you redo all the code in these files to use this library instead of the existing implementation. All variables in the system are safe to expose publicly

kessler/frontend/src/lib/env_variables
   │ │ │ │ 󰛦 env_variables.ts          
   │ │ │ │  env_variables_hydration_scr
   │ │ │ └  env_variables_root_script.ts

Could you use the same API for exposing them on the frontend that was previously used to minimize complexity. Please keep the same existing apis and functionality like getClientRuntimeEnv and getUniversalEnvConfig. And also keep zod around to help debug issues around enviornment variables not being properly set.

# Making env public 🛠

In some cases you might have control over the naming of the environment variables. (Or you simply don't want to prefix them with `NEXT_PUBLIC_`.) In this case you can use the `makeEnvPublic` utility function to make them public.

## Example

```ts
// next.config.js

const { makeEnvPublic } = require('next-runtime-env');

// Given that `FOO` is declared as a regular env var, not a public one. This
// will make it public and available as `NEXT_PUBLIC_FOO`.
makeEnvPublic('FOO');

// Or you can make multiple env vars public at once.
makeEnvPublic(['BAR', 'BAZ']);
```

> You can also use the experimental instrumentation hook introduced in Next.js 13. See the `with-app-router` example for more details.


# Exposing custom environment variables 🛠

- [Exposing custom environment variables 🛠](#exposing-custom-environment-variables-)
  - [Using the script approach (recommend)](#using-the-script-approach-recommend)
    - [Example](#example)
  - [Using the context approach](#using-the-context-approach)
    - [Example](#example-1)

## Using the script approach (recommend)

You might not only want to expose environment variables that are prefixed with `NEXT_PUBLIC_`. In this case you can use the `EnvScript` to expose custom environment variables to the browser.

### Example

```tsx
// app/layout.tsx
// This is as of Next.js 14, but you could also use other dynamic functions
import { unstable_noStore as noStore } from 'next/cache';
import { EnvScript } from 'next-runtime-env';

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  noStore(); // Opt into dynamic rendering

  // This value will be evaluated at runtime
  return (
    <html lang="en">
      <head>
        <EnvScript
          env={{
            NEXT_PUBLIC_: process.env.NEXT_PUBLIC_FOO,
            BAR: process.env.BAR,
            BAZ: process.env.BAZ,
            notAnEnvVar: 'not-an-env-var',
          }}
        />
      </head>
      <body>
        {children}
      </body>
    </html>
  );
}
```

## Using the context approach

You might not only want to expose environment variables that are prefixed with `NEXT_PUBLIC_`. In this case you can use the `EnvProvider` to expose custom environment variables to the context.

### Example

```tsx
// app/layout.tsx
// This is as of Next.js 14, but you could also use other dynamic functions
import { unstable_noStore as noStore } from 'next/cache';
import { EnvProvider } from 'next-runtime-env';

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  noStore(); // Opt into dynamic rendering

  // This value will be evaluated at runtime
  return (
    <html lang="en">
      <body>
        <EnvProvider
          env={{
            NEXT_PUBLIC_: process.env.NEXT_PUBLIC_FOO,
            BAR: process.env.BAR,
            BAZ: process.env.BAZ,
            notAnEnvVar: 'not-an-env-var',
          }}
        >
          {children}
        </EnvProvider>
      </body>
    </html>
  );
}
```

## Libary Docs
# 🌐 Next.js Runtime Environment Configuration

**Effortlessly populate your environment at runtime, not just at build time, with `next-runtime-env`.**

🌟 **Highlights:**
- **Isomorphic Design:** Works seamlessly on both server and browser, and even in middleware.
- **Next.js 13 & 14 Ready:** Fully compatible with the latest Next.js features.
- **`.env` Friendly:** Use `.env` files during development, just like standard Next.js.

### 🤔 Why `next-runtime-env`?

In the modern software development landscape, the "[Build once, deploy many][build-once-deploy-many-link]" philosophy is key. This principle, essential for easy deployment and testability, is a [cornerstone of continuous delivery][fundamental-principle-link] and is embraced by the [twelve-factor methodology][twelve-factor-link]. However, front-end development, particularly with Next.js, often lacks support for this - requiring separate builds for different environments. `next-runtime-env` is our solution to bridge this gap in Next.js.

### 📦 Introducing `next-runtime-env`

`next-runtime-env` dynamically injects environment variables into your Next.js application at runtime. This approach adheres to the "build once, deploy many" principle, allowing the same build to be used across various environments without rebuilds.

### 🤝 Compatibility Notes

- **Next.js 14:** Use `next-runtime-env@3.x` for optimal caching support.
- **Next.js 13:** Opt for [`next-runtime-env@2.x`][app-router-branch-link], tailored for the App Router.
- **Next.js 12/13 Page Router:** Stick with [`next-runtime-env@1.x`][pages-router-branch-link].

### 🔖 Version Guide

- **1.x:** Next.js 12/13 Page Router
- **2.x:** Next.js 13 App Router
- **3.x:** Next.js 14 with advanced caching

### 🚀 Getting Started

In your `app/layout.tsx`, add:

```js
// app/layout.tsx
import { PublicEnvScript } from 'next-runtime-env';

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <head>
        <PublicEnvScript />
      </head>
      <body>
        {children}
      </body>
    </html>
  );
}
```

The `PublicEnvScript` component automatically exposes all environment variables prefixed with `NEXT_PUBLIC_` to the browser. For custom variable exposure, refer to [EXPOSING_CUSTOM_ENV.md](docs/EXPOSING_CUSTOM_ENV.md).

### 🧑‍💻 Usage

Access your environment variables easily:

```tsx
// app/client-page.tsx
'use client';
import { env } from 'next-runtime-env';

export default function SomePage() {
  const NEXT_PUBLIC_FOO = env('NEXT_PUBLIC_FOO');
  return <main>NEXT_PUBLIC_FOO: {NEXT_PUBLIC_FOO}</main>;
}
```

### 🛠 Utilities

Need to expose non-prefixed environment variables to the browser? Check out [MAKING_ENV_PUBLIC.md](docs/MAKING_ENV_PUBLIC.md).

