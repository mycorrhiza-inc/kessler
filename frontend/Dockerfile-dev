# BIG UPDATE: Changed the root image and source of the configuration files for deployment
FROM node:22.4.1-alpine3.20 as build_image
WORKDIR /app
COPY ./package.json ./  
COPY ./package-lock.json ./    
# install dependencies
RUN npm install --force
COPY . .
# build
# RUN npm run build
FROM node:22.4.1-alpine3.20
WORKDIR /app
# copy from build image
COPY --from=build_image /app/package.json ./package.json
COPY --from=build_image /app/node_modules ./node_modules
COPY --from=build_image /app/public ./public
COPY ./tsconfig.json ./
COPY ./tailwind.config.ts ./
COPY ./postcss.config.js ./

EXPOSE 3000
