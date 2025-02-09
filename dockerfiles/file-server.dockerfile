FROM golang:1.23.6-alpine AS builder

WORKDIR /app

COPY go.* /app/

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o file-server ./cmd/main.go

FROM alpine:latest

RUN apk add --no-cache bash

WORKDIR /app

RUN mkdir -p /app/uploads

RUN chmod -R 755 /app/uploads

COPY --from=builder /app/file-server .

EXPOSE 8080

CMD ["./file-server"]
