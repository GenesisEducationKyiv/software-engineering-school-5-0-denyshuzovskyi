FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o nimbus-notify ./cmd/api-server/main.go

EXPOSE 8080

CMD ["./nimbus-notify"]
