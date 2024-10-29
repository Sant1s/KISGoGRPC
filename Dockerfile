FROM golang:latest AS builder
WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o main cmd/service/main.go

FROM ubuntu:latest

COPY --from=builder /app/main /usr/local/bin/main
COPY config.yaml /usr/local/bin/config.yaml
ENV CONFIG_PATH=/usr/local/bin/config.yaml

CMD ["main"]
