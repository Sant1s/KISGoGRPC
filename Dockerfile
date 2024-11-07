FROM golang:latest AS builder
WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o main main.go

FROM ubuntu:latest

WORKDIR /app
COPY --from=builder /app/main /app/main
COPY ./config.yaml /app/config.yaml
ENV CONFIG_PATH=./config.yaml

CMD ["/app/main"]
