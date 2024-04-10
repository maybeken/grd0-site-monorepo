FROM oven/bun:latest as builder
WORKDIR /app
COPY . .
WORKDIR /app/vue-web
RUN bun install
RUN bun run build-only
COPY --from=builder /app/vue-web/dist dist/