FROM golang:1.23 AS builder

WORKDIR /app

COPY services/product-service/go.mod services/product-service/go.sum ./
COPY services/product-service/vendor ./vendor
COPY services/product-service .

RUN CGO_ENABLED=0 GOOS=linux go build -o product-service ./cmd/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/product-service .
COPY services/product-service/.env .env
EXPOSE 50052
CMD ["./product-service"]