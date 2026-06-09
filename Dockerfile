# syntax=docker/dockerfile:1

# ---- Build stage ----
# Runs on the builder's native platform and cross-compiles for the target
# platform, so a single buildx invocation can produce both amd64 and arm64.
FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src

# Cache module downloads.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o /out/server ./cmd/server

# ---- Runtime stage ----
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /out/server /app/server
COPY web/ /app/web/

# Self-signed TLS server; listens on 8443 by default (override with PORT).
EXPOSE 8443

CMD ["./server"]
