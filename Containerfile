FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN set -eux; \
  apk update; \
  apk upgrade; \
  apk add --no-cache make musl-dev gcc sqlite; \
  make

ENTRYPOINT ["/app/bin/contrived"]
