
# BikeStoreGolang

**BikeStoreGolang** is a microservice-based application for managing an online bicycle store. It is built using Clean Architecture principles, API Gateway, and several microservices: authentication, product management, order handling, and payments.

## 📦 Project Includes

- **API Gateway** (in Go)
- **Auth Service**
- **Product Service**
- **Order Service**
- **Payment Service**

## ⚙️ Technologies

- **Backend:** Go, gRPC, NATS, MongoDB, Redis
- **DevOps:** Docker, Docker Compose

## 🚀 How to Run Locally

1. Make sure Docker and Docker Compose are installed.
2. Clone the repository:
   ```bash
   git clone https://github.com/Adilforest/BikeStoreGolang.git
   cd BikeStoreGolang
   ```
3. Start the application:
   ```bash
   docker-compose up --build
   ```


## 🔌 gRPC Endpoints (Main)

Each microservice uses gRPC. The API Gateway translates HTTP requests to gRPC calls:

- **Auth Service**
  - `Login`
  - `Register`
  - `Activate`
  - `ForgotPassword`
  - `ResetPassword`
  - `RefreshToken`
  - `GetMe`
  - `Logout`

- **Product Service**
  - `ListProducts`
  - `SearchProducts`
  - `CreateProduct`
  - `GetProduct`
  - `UpdateProduct`
  - `DeleteProduct`
  - `ChangeProductStock`

- **Order Service**
  - `CreateOrder`
  - `GetOrder`
  - `ListOrdersByUser`
  - `CancelOrder`
  - `ApproveOrder`

## 📌 Implemented HTTP Routes

### 🔐 Auth Routes

```http
POST   /login               -> Login
POST   /register            -> Register
GET    /activate            -> Activate
POST   /forgot-password     -> ForgotPassword
POST   /reset-password      -> ResetPassword
POST   /refresh-token       -> RefreshToken
GET    /me                  -> GetMe
POST   /logout              -> Logout
```

### 📦 Product Routes

```http
GET    /products            -> ListProducts
GET    /products/search     -> SearchProducts
POST   /products            -> CreateProduct
GET    /products/:id        -> GetProduct
PUT    /products/:id        -> UpdateProduct
DELETE /products/:id        -> DeleteProduct
POST   /products/:id/stock  -> ChangeProductStock
```

### 📦 Order Routes

```http
POST   /orders              -> CreateOrder
GET    /orders/:id          -> GetOrder
GET    /orders/user/:id     -> ListOrdersByUser
POST   /orders/:id/cancel   -> CancelOrder
POST   /orders/:id/approve  -> ApproveOrder
```

## ✅ Implemented Features

- User registration and authentication with JWT
- Product CRUD operations, search, and stock control
- Order creation, cancellation, and approval workflows
- gRPC-based service communication with HTTP gateway
- Dockerized infrastructure for rapid deployment

---

