FROM golang:1.24.4-alpine AS builder

WORKDIR /app

RUN apk add --no-cache \
    gcc \
    bash \
    musl-dev \
    libc-dev \
    make \
    git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o lock-bot ./cmd/lock-bot/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/lock-bot .
CMD ["./lock-bot"]
