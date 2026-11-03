FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY relay_server.go .
RUN go build -o relay relay_server.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/relay .
EXPOSE 7000 7001 7002
CMD ["./relay"]