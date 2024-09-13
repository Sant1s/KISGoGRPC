FROM golang:latest AS builder
WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o main cmd/moderation/main.go

FROM ubuntu:latest

COPY --from=builder /app/main /usr/local/bin/main
COPY .env /usr/local/bin/.env
ENV CONFIG_PATH=/usr/local/bin/.env

CMD ["main"]
