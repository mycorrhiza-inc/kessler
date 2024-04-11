FROM node:21-alpine

WORKDIR /app

COPY package.json ./
COPY postcss.config.js ./
COPY tailwind.config.js ./
COPY node_modules ./
COPY .env.local ./

RUN npm install --save-dev
