FROM golang:1.20-alpine3.17 as builder

RUN apk update && \
    apk add --no-cache file build-base

WORKDIR /go/src/app

COPY . .
RUN mkdir /app

RUN go mod download
RUN CGO_ENABLED=1 GOPROXY=direct go build -o /app/birdhouse -mod=mod  ./cmd/main.go

# deploy-stage
FROM alpine:latest

RUN apk update && \
    apk add --no-cache curl bash

RUN mkdir  /key

WORKDIR /app

COPY --from=builder /app ./


RUN apk update && \
    apk add --no-cache curl bash

HEALTHCHECK --interval=30s --start-period=1m --timeout=30s --retries=3 \
    CMD curl --silent --fail --fail-early http://127.0.0.1:$PORT/docs || exit 1
EXPOSE $PORT
ENTRYPOINT ["/app/birdhouse"]
