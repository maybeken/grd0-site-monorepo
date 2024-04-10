FROM oven/bun:latest as builder
WORKDIR /app
COPY . .
WORKDIR /app/vue-web
RUN bun install
RUN bun run build