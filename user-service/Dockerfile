FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/user-service ./cmd/user-service

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /bin/user-service /usr/local/bin/user-service

EXPOSE 8081

CMD ["user-service"]
