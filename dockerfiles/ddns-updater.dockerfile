FROM golang:1.23.6-alpine AS builder

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ddns-updater ./cmd/main.go

FROM alpine:latest

RUN apk add --no-cache bash

WORKDIR /app

COPY --from=builder /app/ddns-updater .

CMD ["./ddns-updater"]
