FROM node:24-alpine AS base
FROM base AS deps
RUN corepack enable
WORKDIR /app
COPY pnpm-lock.yaml ./
RUN --mount=type=cache,id=pnpm,target=/root/.local/share/pnpm/store pnpm fetch --frozen-lockfile
COPY package.json ./
RUN --mount=type=cache,id=pnpm,target=/root/.local/share/pnpm/store pnpm install --frozen-lockfile --prod



FROM base AS build
RUN corepack enable
WORKDIR /app
COPY pnpm-lock.yaml ./

ENV NEXT_PUBLIC_KESSLER_API_URL=https://api.kessler.xyz
ENV NEXT_PUBLIC_NIGHTLY_KESSLER_API_URL=https://nightly-api.kessler.xyz
ENV NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY=pk_test_YWNlLXdhbGxhYnktODQuY2xlcmsuYWNjb3VudHMuZGV2JA
ENV NEXT_PUBLIC_CLERK_FRONTEND_API=ace-wallaby-84.clerk.accounts.dev

RUN --mount=type=cache,id=pnpm,target=/root/.local/share/pnpm/store pnpm fetch --frozen-lockfile
COPY package.json ./
RUN --mount=type=cache,id=pnpm,target=/root/.local/share/pnpm/store pnpm install --frozen-lockfile
COPY . .
RUN pnpm build




FROM base
WORKDIR /app
COPY --from=deps /app/node_modules /app/node_modules
COPY --from=build /app/dist /app/dist
ENV NODE_ENV production
CMD ["node", "./dist/index.js"]


