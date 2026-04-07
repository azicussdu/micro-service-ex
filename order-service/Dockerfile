FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/order-service ./cmd/order-service

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /bin/order-service /usr/local/bin/order-service

EXPOSE 8082

CMD ["order-service"]
