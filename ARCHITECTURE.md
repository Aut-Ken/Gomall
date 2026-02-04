# Gomall Architecture Overview

## 1. System Architecture

GoMall is designed as a **Modular Monolith** application aimed at handling high-concurrency e-commerce scenarios (specifically Seckill/Flash Sales). While it is currently a single deployable unit, its internal structure is layered and modular, preparing it for potential future microservices splitting.

### High-Level Diagram

```mermaid
graph TD
    User[User Client] -->|HTTP/REST| LB[Load Balancer]
    LB -->|Port 8080| App[Gomall Application]

    subgraph "Gomall Application (Monolith)"
        API[API Layer (Gin)] --> Service[Service Layer]
        Service --> Repo[Repository Layer]

        Service -->|Async| MQ_Pub[RabbitMQ Producer]
        MQ_Sub[RabbitMQ Consumer] -->|Background| Service

        subgraph "Cross-Cutting"
            Auth[JWT Auth]
            Rate[Rate Limiter]
            Tracing[OpenTelemetry]
        end
    end

    subgraph "Infrastructure"
        Repo -->|ORM| MySQL[(MySQL 8.0)]
        Repo -->|Cache/Lock| Redis[(Redis 7.0)]
        MQ_Pub --> RabbitMQ[(RabbitMQ)]
        RabbitMQ --> MQ_Sub
        Tracing -.->|Export| Jaeger[(Jaeger)]
    end
```

## 2. Technical Stack

- **Language**: Go 1.20+
- **Web Framework**: Gin (High performance HTTP web framework)
- **Database ORM**: Gorm (MySQL interaction)
- **Cache & kv Store**: Redis (Used for caching, distributed locks, and inventory counters)
- **Message Queue**: RabbitMQ (Used for traffic peaking/shaving and asynchronous decoupling)
- **Configuration**: Viper (Support for YAML/JSON configs)
- **Tracing**: OpenTelemetry + Jaeger via OTLP gRPC (Distributed tracing)
- **Rate Limiting**: golang.org/x/time/rate (Local) + Redis (Distributed)

## 3. Layered Design

The application follows a strict layered architecture:

### 3.1 Interface Layer (`internal/api` & `internal/router`)
- **Router**: Defines HTTP routes using Gin.
- **Handlers**: Handles HTTP requests, parameter validation, and response formatting.
- **Middleware**: Handles cross-cutting concerns:
  - JWT Authentication (`internal/middleware/auth.go`)
  - Rate Limiting (`internal/middleware/ratelimit.go`)

### 3.2 Service Layer (`internal/service`)
- Contains the core business logic.
- **Seckill Logic**:
  - Uses Redis Lua scripts for atomic inventory deduction.
  - Pushes successful purchase requests to RabbitMQ.
- **Order Processing**:
  - Background goroutines consume messages from RabbitMQ to create orders in the database asynchronously.

### 3.3 Data Access Layer (`internal/repository` & `internal/model`)
- **Model**: Defines the database schema structs (Gorm models).
- **Repository**: Encapsulates all direct database and cache operations.

## 4. Key Workflows

### 4.1 High-Concurrency Seckill Flow
To prevent "overselling" and ensure system stability during traffic spikes:

1.  **Request**: User sends a seckill request.
2.  **Pre-Check (Redis)**: System checks stock in Redis (in-memory speed).
3.  **Atomic Deduct**: Lua script atomically decrements stock in Redis.
4.  **Async Queue**: If deduction succeeds, a message is sent to RabbitMQ.
5.  **Response**: Immediate response to user ("Queued").
6.  **Processing**: Background worker consumes message and creates the DB order.
7.  **Consistency**: Database transaction ensures final data consistency.

### 4.2 Standard Shopping Flow
1.  **Browse**: Read-heavy operations (cached in Redis).
2.  **Order**: Database transaction ensures atomicity of Order Creation + Inventory Deduction.

### 4.3 Shopping Cart Flow
1.  **Add**: User adds product to cart (Redis hash storage for performance).
2.  **Manage**: Update quantity, remove items, or clear cart.
3.  **Checkout**: Convert cart items to order (batch processing).

### 4.4 Admin Operations
- **Product Management**: Create, Update, Delete products (AdminAuth required).
- **Stock Initialization**: Pre-load inventory to Redis before seckill events (AdminAuth required).

### 4.5 Rate Limiting
Two types of rate limiting are implemented:

**Local Rate Limiting (IP-based):**
- Uses `golang.org/x/time/rate` for in-memory rate limiting
- Suitable for single-instance deployments
- Configurable per endpoint (Global, API, Seckill, Login)

**Distributed Rate Limiting (Redis-based):**
- Uses Redis sorted sets with Lua scripts for atomic operations
- Supports multi-instance deployments
- Sliding window algorithm for accurate rate limiting

### 4.6 Distributed Tracing (OpenTelemetry + Jaeger via OTLP gRPC)
- **OpenTelemetry Integration**: Standardized tracing API
- **OTLP gRPC Exporter**: Sends traces to Jaeger via gRPC protocol (port 4317)
- **Key Features**:
  - Request latency tracking
  - Error tracking and correlation
  - Custom attributes (UserID, ProductID, OrderNo)
  - Trace context propagation

## 5. Deployment

- **Single Binary**: The entire application compiles into a single binary (`gomall.exe` / `main`).
- **Containerization**: `Dockerfile` provided for containerized deployment.
- **Orchestration**: `docker-compose.yml` orchestrates the App, MySQL, Redis, and RabbitMQ.

## 6. Directory Map

| Path | Purpose |
|------|---------|
| `cmd/` | Main application entry point (`main.go`). |
| `internal/api/` | HTTP Handlers (Controllers) - User, Product, Order, Cart, Seckill. |
| `internal/service/` | Business Logic - User, Product, Order, Cart, Seckill services. |
| `internal/repository/` | DB and Cache interactions. |
| `internal/model/` | Data entities - User, Product, Order, Cart, Stock models. |
| `internal/middleware/` | JWT Authentication, Admin authorization, Rate Limiting. |
| `internal/router/` | Gin route definitions. |
| `internal/rabbitmq/` | Message queue producer/consumer. |
| `internal/redis/` | Redis client and Lua scripts. |
| `internal/tracing/` | OpenTelemetry + Jaeger tracing integration. |
| `internal/grpc/` | gRPC service definitions (future use). |
| `conf/` | Configuration files (`config.yaml`). |
| `pkg/` | Shared utilities (JWT, Password hashing). |
