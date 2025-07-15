FROM golang:1.24.5-alpine AS builder

WORKDIR /app

COPY weather-upd-subscription-srv/go.mod weather-upd-subscription-srv/go.sum ./
COPY nimbus-proto ./nimbus-proto
COPY weather-upd-subscription-srv ./weather-upd-subscription-srv

WORKDIR /app/weather-upd-subscription-srv

RUN go mod download

RUN go build -o weather-upd-subscription-srv ./cmd/server/main.go

EXPOSE 8080

CMD ["./weather-upd-subscription-srv"]
