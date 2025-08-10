FROM golang:1.24.5-alpine AS builder

WORKDIR /app

COPY nimbus-lib ./nimbus-lib
COPY notification-srv ./notification-srv

WORKDIR /app/notification-srv

RUN go mod download

RUN go build -o notification-srv ./cmd/server/main.go


FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/notification-srv .

CMD ["./notification-srv"]