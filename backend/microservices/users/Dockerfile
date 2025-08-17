# syntax=docker/dockerfile:experimental

# --- Build stage ---

FROM go:1.24.5-alpine AS builder

RUN apk add upx=5.0.0 --no-cache

WORKDIR /app

RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=bind,src=backend/go.mod,target=go.mod \
  --mount=type=bind,src=backend/go.sum,target=go.sum \
  go mod download

ENV CGO_ENABLED=0

RUN --mount=type=cache,target=/go/pkg/mod  \
  --mount=type=bind,src=backend,target=. \
  go build -ldflags="-s -w" -o /usr/local/bin/users-microservice ./microservices/users/cmd

RUN upx --best --no-backup /usr/local/bin/users-microservice

# --- Packaging stage ---

FROM scratch AS packager

USER default

WORKDIR /

COPY --from=builder /usr/local/bin/users-microservice /usr/local/bin/users-microservice

EXPOSE 4000/tcp

ENTRYPOINT [ "users-microservice" ]
