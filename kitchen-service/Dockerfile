FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/kitchen-service ./cmd

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /bin/kitchen-service /usr/local/bin/kitchen-service

CMD ["kitchen-service"]
