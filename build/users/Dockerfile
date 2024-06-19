# syntax=docker/dockerfile:1

FROM golang:1.22-alpine as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o users cmd/users/main.go

FROM alpine:3.16

WORKDIR /app

RUN apk update && \
    apk add --no-cache \
    openssh-keygen \
    openssl

COPY --from=builder /build/users ./
COPY ./configs/users/config.json ./config.json

RUN ssh-keygen -t rsa -P "" -b 4096 -m PEM -f rs256.rsa
RUN openssl rsa -in rs256.rsa -pubout -outform PEM -out rs256.rsa.pub

CMD ./users -c config.json

EXPOSE 8000