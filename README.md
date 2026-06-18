# BikeStoreGolang

A microservices e-commerce backend for a bicycle store, built in Go with gRPC service communication, NATS event streaming, PostgreSQL + Redis persistence, and a React frontend — designed as a learning project covering clean architecture, asynchronous messaging, and full-stack integration.

![Go](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go&logoColor=white)
![gRPC](https://img.shields.io/badge/gRPC-protobuf-4285F4?logo=google&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-7-DC382D?logo=redis&logoColor=white)
![NATS](https://img.shields.io/badge/NATS-2-199bdb?logo=natsdotio&logoColor=white)
![React](https://img.shields.io/badge/React-18-61DAFB?logo=react&logoColor=white)
![Vite](https://img.shields.io/badge/Vite-5-646CFF?logo=vite&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-compose-2496ED?logo=docker&logoColor=white)

---

## Overview

BikeStoreGolang is a full-stack microservice application for managing a bicycle store. It covers the complete shopping flow: user authentication, product browsing, order placement, and payment processing. Each domain is isolated in its own service, services communicate over gRPC, and asynchronous side-effects (order events, payment notifications) flow over NATS.

The project applies a layered architecture inside each service: `domain` (entities and repository interfaces), `usecase` (business logic), `repository` (Postgres + Redis implementations), and `delivery` (gRPC handlers and NATS subscribers/publishers).

---

## Features

- **Auth service** — user registration and login with JWT, session management via Redis, bcrypt password hashing
- **Product service** — CRUD for bicycle catalog; products typed as `road`, `mountain`, `hybrid`, or `electric`; stock management use case; Redis caching of product reads
- **Order service** — create and cancel orders (with order-item line entries), cache layer via Redis, NATS publisher for order events; integrates with product service via local module replace
- **Payment service** — payment processing use case with Redis-based distributed lock (`lock_repo`), webhook handler for async payment callbacks, NATS publisher for payment events
- **API gateway** — single HTTP entry point (Gin), JWT auth middleware, Prometheus metrics (`http_requests_total`, `http_request_duration_seconds`, `grpc_client_connections_total`), structured Logrus logging, Swagger spec
- **React frontend** — Vite build; API clients for auth, products, and orders; route-based navigation

---

## Architecture

```
Browser
  │
  ▼
api-gateway  (Gin HTTP :8080)
  │  JWT middleware · Prometheus /metrics · Swagger /api/swagger.yaml
  │
  ├──gRPC──► auth-service    (PostgreSQL + Redis)
  ├──gRPC──► product-service (PostgreSQL + Redis cache)
  ├──gRPC──► order-service   (PostgreSQL + Redis cache + NATS pub)
  └──gRPC──► payment-service (PostgreSQL + Redis lock + NATS pub + HTTP webhook)

NATS ◄──── order-service, payment-service (async events)
```

---

## Tech Stack

| Layer | Technologies |
|---|---|
| Language | Go 1.23 |
| HTTP framework | Gin |
| Service communication | gRPC + Protocol Buffers |
| Messaging | NATS |
| Primary database | PostgreSQL |
| Cache / locks | Redis |
| Observability | Prometheus client, Logrus |
| Frontend | React 18, Vite, Axios |
| Containers | Docker, Docker Compose |
| API docs | Swagger (YAML) |

---

## Project Structure

```
BikeStoreGolang/
├── api-gateway/
│   ├── cmd/main.go                  # Gin router, Prometheus middleware, route registration
│   ├── internal/
│   │   ├── auth/                    # JWT parsing and auth middleware
│   │   ├── handlers/                # auth, product, order HTTP handlers
│   │   ├── service/                 # service layer wrapping gRPC clients
│   │   └── client/                  # gRPC client constructors
│   ├── proto/                        # .proto definitions (auth, product, order)
│   ├── api/swagger.yaml              # OpenAPI/Swagger spec
│   └── configs/config.yaml
├── services/
│   ├── auth-service/
│   │   ├── internal/domain/         # User entity + repository interface
│   │   ├── internal/usecase/        # auth_usecase, session_usecase
│   │   ├── internal/repository/
│   │   │   ├── postgres/            # user_repo
│   │   │   └── redis/               # token_repo
│   │   └── internal/delivery/
│   │       ├── grpc/                # gRPC handler + server setup
│   │       └── nats/                # publisher + subscriber
│   ├── order-service/
│   │   ├── internal/domain/         # Order, OrderItem entities + repo interface
│   │   ├── internal/usecase/        # order_usecase, payment_usecase
│   │   ├── internal/repository/
│   │   │   ├── postgres/            # order_repo
│   │   │   └── redis/               # cache_repo
│   │   └── internal/delivery/
│   │       ├── grpc/
│   │       └── nats/
│   ├── payment-service/
│   │   ├── internal/domain/         # Payment, Transaction entities
│   │   ├── internal/usecase/        # payment_usecase, webhook_usecase
│   │   ├── internal/repository/
│   │   │   ├── postgres/            # payment_repo
│   │   │   └── redis/               # lock_repo (distributed lock)
│   │   └── internal/delivery/
│   │       ├── grpc/
│   │       ├── http/                # webhook handler
│   │       └── nats/
│   └── product-service/
│       ├── internal/domain/         # Product entity + repository interface
│       ├── internal/usecase/        # product_usecase, stock_usecase
│       ├── internal/repository/
│       │   ├── postgres/            # product_repo
│       │   └── redis/               # cache_repo
│       └── internal/delivery/
│           ├── grpc/
│           └── nats/
└── frontend/
    ├── src/
    │   ├── api/                     # auth.js, orders.js, products.js
    │   ├── App.jsx / routes.jsx
    │   └── utils/                   # api helper, auth helper
    └── vite.config.js
```

---

## Getting Started

### Prerequisites

- Go 1.23+
- PostgreSQL (one database or separate per service)
- Redis
- NATS server
- Docker + Docker Compose (for the full stack)
- Node.js 20+ (for the frontend)

### Configuration

Each service uses a `configs/config.yaml` and an optional `.env` file. Copy and fill in the values before running.

### Run individual services

```bash
# Start the auth service
cd services/auth-service && go run cmd/main.go

# Start the product service
cd services/product-service && go run cmd/main.go

# Start the order service
cd services/order-service && go run cmd/main.go

# Start the payment service
cd services/payment-service && go run cmd/main.go

# Start the API gateway
cd api-gateway && go run cmd/main.go
```

### Run the frontend

```bash
cd frontend
npm install
npm run dev
# Vite dev server at http://localhost:5173
```

### Run everything with Docker Compose

```bash
docker compose up --build
```

---

## API Endpoints

| Method | Path | Description |
|---|---|---|
| `POST` | `/login` | Authenticate user, receive JWT |
| `POST` | `/register` | Create new account |
| `GET` | `/me` | Get current user profile |
| `POST` | `/logout` | Invalidate session |
| `POST` | `/refresh-token` | Issue new access token |
| `GET` | `/products` | List products |
| `POST` | `/products/search` | Search products |
| `POST` | `/products` | Create product |
| `GET` | `/products/:id` | Get product by ID |
| `PUT` | `/products/:id` | Update product |
| `DELETE` | `/products/:id` | Delete product |
| `POST` | `/products/:id/stock` | Adjust stock quantity |
| `POST` | `/orders` | Place an order |
| `GET` | `/orders/:id` | Get order by ID |
| `GET` | `/orders/user/:user_id` | List orders for a user |
| `POST` | `/orders/:id/cancel` | Cancel order |
| `POST` | `/orders/:id/approve` | Approve order |
| `GET` | `/metrics` | Prometheus metrics |
| `GET` | `/health` | Health check |

---

Adil Ormanov — [GitHub](https://github.com/Adilforest)
