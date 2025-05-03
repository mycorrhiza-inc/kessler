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
# build with a dummy .env for build time
# RUN touch .env
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
