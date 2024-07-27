# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS build_base
WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o ./user_auth ./cmd/main.go

RUN go build -o ./user_auth ./cmd/main.go

FROM alpine:3.18
WORKDIR /app

COPY --from=build_base /app/user_auth ./user_auth

CMD ["/app/user_auth"]