# BIG UPDATE: Changed the root image and source of the configuration files for deployment
FROM node:22.14-alpine3.20 AS build_image
WORKDIR /app
COPY ./package.json ./  
COPY ./package-lock.json ./    
COPY ./postcss.config.js ./
ENV NODE_ENV=development
# install dependencies
RUN  npm install --force
COPY . .
# build
# RUN npm run build
FROM node:22.14-alpine3.20
WORKDIR /app
# copy from build image
COPY --from=build_image /app/package.json ./package.json
COPY --from=build_image /app/next.config.js ./next.config.js
COPY --from=build_image /app/node_modules ./node_modules
COPY --from=build_image /app/public ./public
COPY --from=build_image /app/postcss.config.js ./
COPY ./tsconfig.json ./

EXPOSE 300 
