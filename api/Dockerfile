# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /build
ENV GOPROXY=https://goproxy.io,direct
ADD . .
RUN cd api; go build -o ../grd0_blog_api

# Final stage
FROM alpine:latest
WORKDIR /opt/app
COPY --from=builder /build/grd0_blog_api .
CMD ["./grd0_blog_api"]
