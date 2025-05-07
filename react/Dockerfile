# syntax=docker/dockerfile:1
FROM node:20-alpine3.21 AS base
WORKDIR /app

FROM base AS deps
COPY package.json pnpm-lock.yaml ./
RUN corepack enable pnpm && pnpm i --frozen-lockfile

FROM base AS build
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN corepack enable pnpm && pnpm run build
 
FROM base AS final
COPY --from=build /app/dist ./dist
RUN npm i -g serve
COPY .env .
EXPOSE 3000
CMD [ "serve", "-s", "dist" ]