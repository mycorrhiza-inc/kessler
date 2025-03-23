FROM node:23.6.1-alpine3.20 as frontend-builder
RUN useradd -ms /bin/sh -u 1001 app
USER app
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
COPY --from=BUILD_IMAGE /app/package.json ./package.json
COPY --from=BUILD_IMAGE /app/node_modules ./node_modules
COPY --from=BUILD_IMAGE /app/.next ./.next
COPY --from=BUILD_IMAGE /app/public ./public
COPY --from=BUILD_IMAGE /app/tsconfig.json ./
COPY ./tailwind.config.ts ./
COPY ./postcss.config.js ./

EXPOSE 3000
