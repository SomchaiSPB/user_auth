# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS build_base
WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o ./user_auth ./cmd/worker/main.go

FROM alpine:3.18
WORKDIR /app

COPY --from=build_base /app/user_auth ./user_auth

CMD ["/app/user_auth"]