# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /build
ENV GOPROXY=direct
ADD . .
ENV GOCACHE=/root/.cache/go-build
RUN apk add git
RUN --mount=type=cache,target="/root/.cache/go-build" cd api; go build -o /build/grd0_blog_api

# Final stage
FROM alpine:latest
WORKDIR /opt/app
COPY --from=builder /build/grd0_blog_api .
CMD ["./grd0_blog_api"]
