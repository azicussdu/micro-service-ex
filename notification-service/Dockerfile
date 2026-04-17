FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/notification-service ./cmd

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /bin/notification-service /usr/local/bin/notification-service

CMD ["notification-service"]
