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
ENV NEXT_PUBLIC_SUPABASE_URL=https://kpvkpczxcclxslabfzeu.supabase.co
ENV NEXT_PUBLIC_SUPABASE_ANON_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImtwdmtwY3p4Y2NseHNsYWJmemV1Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3MjUwNTQzOTIsImV4cCI6MjA0MDYzMDM5Mn0.9kR-oYUM5SqjmteQbPE1w-ABX8-0sSGldGAGsegCHfs
ENV NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY=pk_test_YWNlLXdhbGxhYnktODQuY2xlcmsuYWNjb3VudHMuZGV2JA
ENV NEXT_PUBLIC_CLERK_FRONTEND_API=ace-wallaby-84.clerk.accounts.dev

RUN --mount=type=cache,id=pnpm,target=/root/.local/share/pnpm/store pnpm fetch --frozen-lockfile
COPY package.json ./
RUN --mount=type=cache,id=pnpm,target=/root/.local/share/pnpm/store pnpm install --frozen-lockfile
COPY . .

CMD ["pnpm", "dev"]


