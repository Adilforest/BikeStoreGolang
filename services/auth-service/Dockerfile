# Build stage
FROM golang:1.23

WORKDIR /app

# Копируем go.mod и go.sum для скачивания зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Компилируем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service ./cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Копируем бинарник из builder
COPY --from=builder /auth-service .
# Копируем конфиги
COPY configs ./configs

# Указываем порт (должен совпадать с вашим gRPC-портом)
EXPOSE 50051

# Команда для запуска
CMD ["./auth-service"]