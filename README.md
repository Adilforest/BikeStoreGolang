# BikeStoreGolang 🚴‍♂️

**Online велосипедный магазин** на микросервисной архитектуре (Go, gRPC, React).

## 📌 Основной стек

- **Бэкенд**:
  - Микросервисы на Go (Clean Architecture)
  - gRPC для межсервисного взаимодействия
  - NATS (Message Queue)
  - PostgreSQL + Redis
- **Фронтенд**: React + Vite
- **Инфраструктура**: Docker, Prometheus/Grafana

## 🏗️ Структура проекта

BikeStoreGolang/
├── api-gateway/ # API Gateway (HTTP -> gRPC)
├── services/ # Микросервисы (auth, order, payment, product)
├── frontend/ # React-приложение
├── deployments/ # Docker, k8s, мониторинг
└── proto/ # gRPC-контракты

## 🚀 Запуск (production)

```bash
docker-compose -f deployments/docker-compose.prod.yml up --build
🔍 Мониторинг
Grafana: http://localhost:3000 (логин: admin, пароль: admin)

Prometheus: http://localhost:9090
```
