# syntax=docker/dockerfile:1

FROM golang:1.21-alpine as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o posts cmd/posts/main.go

FROM alpine:3.16

WORKDIR /app

RUN apk update && \
    apk add --no-cache \
    openssh-keygen \
    openssl

COPY --from=builder /build/posts ./
COPY ./configs/posts/config.json ./config.json

CMD ./posts -c config.json

EXPOSE 8001