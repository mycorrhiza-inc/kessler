FROM node:23.6.1-alpine3.20 AS frontend-builder
# Not sure why this is here, but I dont think useradd is a cmd availible on alpine
# RUN useradd -ms /bin/sh -u 1001 app
# USER app
WORKDIR /app

COPY package.json package-lock.json ./   
RUN npm install --force
COPY ./tsconfig.json ./
COPY ./build_time_variables.env ./.env
COPY . .
# Im, sorry - nic
ENV NEXT_PUBLIC_KESSLER_API_URL=https://api.kessler.xyz
ENV NEXT_PUBLIC_NIGHTLY_KESSLER_API_URL=https://nightly-api.kessler.xyz
ENV NEXT_PUBLIC_SUPABASE_URL=https://kpvkpczxcclxslabfzeu.supabase.co
ENV NEXT_PUBLIC_SUPABASE_ANON_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImtwdmtwY3p4Y2NseHNsYWJmemV1Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3MjUwNTQzOTIsImV4cCI6MjA0MDYzMDM5Mn0.9kR-oYUM5SqjmteQbPE1w-ABX8-0sSGldGAGsegCHfs
ENV NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY=pk_live_Y2xlcmsua2Vzc2xlci54eXok

RUN npm run build

FROM node:23.6.1-alpine3.20
WORKDIR /app
# copy from build image
COPY --from=frontend-builder /app/package.json ./package.json
COPY --from=frontend-builder /app/node_modules ./node_modules
COPY --from=frontend-builder /app/.next ./.next
COPY --from=frontend-builder /app/public ./public
# COPY --from=frontend-builder /app/tsconfig.json ./
COPY ./tailwind.config.ts ./
COPY ./postcss.config.js ./

EXPOSE 3000
