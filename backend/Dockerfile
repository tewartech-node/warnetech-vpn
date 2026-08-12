# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o vpn-server main.go

# Final stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates iptables

WORKDIR /app

COPY --from=builder /app/vpn-server .
COPY config/server.json ./config/server.json

EXPOSE 443/udp
EXPOSE 443/tcp
EXPOSE 51820/udp

ENTRYPOINT ["./vpn-server", "--config=/app/config/server.json"]
