FROM golang:1.24.5-alpine AS builder

WORKDIR /app

COPY notification-srv/go.mod notification-srv/go.sum ./
COPY nimbus-proto ./nimbus-proto
COPY notification-srv ./notification-srv

WORKDIR /app/notification-srv

RUN go mod download

RUN go build -o notification-srv ./cmd/server/main.go

EXPOSE 8080

CMD ["./notification-srv"]
