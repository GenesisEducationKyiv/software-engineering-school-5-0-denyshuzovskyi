FROM golang:1.24.5-alpine AS builder

WORKDIR /app

COPY nimbus-proto ./nimbus-proto
COPY weather-upd-subscription-srv ./weather-upd-subscription-srv

WORKDIR /app/weather-upd-subscription-srv

RUN go mod download

RUN go build -o weather-upd-subscription-srv ./cmd/server/main.go


FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/weather-upd-subscription-srv .

CMD ["./weather-upd-subscription-srv"]