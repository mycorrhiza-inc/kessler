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
RUN pnpm build




FROM base
WORKDIR /app
COPY --from=deps /app/node_modules /app/node_modules
COPY --from=build /app/dist /app/dist
ENV NODE_ENV production
CMD ["node", "./dist/index.js"]



# FROM node:23.6.1-alpine3.20 AS frontend-builder
# # Not sure why this is here, but I dont think useradd is a cmd availible on alpine
# # RUN useradd -ms /bin/sh -u 1001 app
# # USER app
# WORKDIR /app
#
# COPY package.json package-lock.json ./   
# RUN NODE_ENV=development npm install --force
# COPY ./tsconfig.json ./
# COPY . .
# # Im, sorry - nic
# ENV NEXT_PUBLIC_KESSLER_API_URL=https://api.kessler.xyz
# ENV NEXT_PUBLIC_NIGHTLY_KESSLER_API_URL=https://nightly-api.kessler.xyz
# ENV NEXT_PUBLIC_SUPABASE_URL=https://kpvkpczxcclxslabfzeu.supabase.co
# ENV NEXT_PUBLIC_SUPABASE_ANON_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImtwdmtwY3p4Y2NseHNsYWJmemV1Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3MjUwNTQzOTIsImV4cCI6MjA0MDYzMDM5Mn0.9kR-oYUM5SqjmteQbPE1w-ABX8-0sSGldGAGsegCHfs
# ENV NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY=pk_test_YWNlLXdhbGxhYnktODQuY2xlcmsuYWNjb3VudHMuZGV2JA
# ENV NEXT_PUBLIC_CLERK_FRONTEND_API=ace-wallaby-84.clerk.accounts.dev
#
#
# RUN npm run build
#
# FROM node:23.6.1-alpine3.20
# WORKDIR /app
# # copy from build image
# COPY --from=frontend-builder /app/package.json ./package.json
# COPY --from=frontend-builder /app/node_modules ./node_modules
# COPY --from=frontend-builder /app/.next ./.next
# COPY --from=frontend-builder /app/public ./public
# # COPY --from=frontend-builder /app/tsconfig.json ./
# # COPY ./tailwind.config.ts ./
# COPY ./postcss.config.js ./
#
# EXPOSE 3000
