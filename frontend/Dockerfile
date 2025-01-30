FROM node:23.6.1-alpine3.20 as BUILD_IMAGE
WORKDIR /app
COPY ./package.json ./   
COPY ./package-lock.json ./     
# install dependencies
RUN npm install --force
COPY ./tsconfig.json ./
COPY . .
# build with a dummy .env for build time
# RUN touch .env.local 
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
