FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api-gateway ./cmd/api-gateway

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /bin/api-gateway /usr/local/bin/api-gateway

EXPOSE 8080

CMD ["api-gateway"]
