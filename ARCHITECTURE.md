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

- **Language**: Go 1.23+
- **Web Framework**: Gin (High performance HTTP web framework)
- **Database ORM**: Gorm (MySQL interaction)
- **Cache & kv Store**: Redis (Used for caching, distributed locks, and inventory counters)
- **Message Queue**: RabbitMQ (Used for traffic peaking/shaving and asynchronous decoupling)
- **Configuration**: Viper (Support for YAML/JSON configs, environment variables, hot reload)
- **Tracing**: OpenTelemetry + Jaeger via OTLP gRPC (Distributed tracing)
- **Rate Limiting**: golang.org/x/time/rate (Local) + Redis (Distributed)
- **Metrics**: Prometheus (HTTP request counts, latency histograms, business metrics)
- **Logging**: Uber Zap (Structured JSON logging)
- **API Documentation**: Swagger/OpenAPI 3.0
- **Validation**: go-playground/validator (Request validation)
- **Payment**: WeChat Pay API (Sandbox environment)

## 3. Layered Design

The application follows a strict layered architecture:

### 3.1 Interface Layer (`internal/api` & `internal/router`)
- **Router**: Defines HTTP routes using Gin.
- **Handlers**: Handles HTTP requests, parameter validation, and response formatting.
- **Middleware**: Handles cross-cutting concerns:
  - JWT Authentication (`internal/middleware/auth.go`)
  - Rate Limiting (`internal/middleware/ratelimit.go`)
  - Structured Logging (`internal/middleware/logger.go`)
  - Prometheus Metrics (`internal/middleware/metrics.go`)
  - Error Handling (`internal/middleware/error_handler.go`)
  - Parameter Validation (`internal/middleware/validator.go`)

**Handlers Include:**
- `handler.go` - User, Product, Order handlers
- `cart_handler.go` - Shopping cart
- `seckill_handler.go` - Seckill operations
- `auth_handler.go` - JWT refresh, password change, logout
- `file_handler.go` - File upload (single/multi)
- `wechat_pay_handler.go` - WeChat Pay integration
- `health_check.go` - Health checks

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

### 4.6 Unified Response & Error Codes

**Standard Response Format:**
```json
{
  "code": 0,
  "message": "success",
  "data": {...},
  "trace_id": "xxx"
}
```

**Error Code System:**
| Code Range | Module |
|------------|--------|
| 0 | Success |
| 400-500 | System errors (400=BadRequest, 401=Unauthorized, 403=Forbidden, 404=NotFound, 500=ServerError) |
| 10001-10099 | User module |
| 20001-20099 | Product module |
| 30001-30099 | Order module |
| 40001-40099 | Payment module |
| 50001-50099 | Cart module |
| 60001-60099 | Seckill module |
| 70001-70099 | File upload module |

### 4.7 Parameter Validation

- Uses `go-playground/validator` for request validation
- Supports struct tags: `binding:"required,min=3,max=50,email"`
- Custom error messages in Chinese/English
- Validation for JSON, Form, and Query parameters

### 4.8 WeChat Pay Integration (Sandbox)

**Flow:**
1. Frontend calls unified-order API
2. Backend calls WeChat API with signed request
3. Returns payment QR code URL
4. User scans QR code to pay
5. WeChat sends async notification to callback URL
6. Backend verifies signature and updates order status
7. Frontend polls order status

**Features:**
- Unified Order API
- Payment notification handling
- Order query
- Order close
- Refund support
- MD5 signature verification

**Files:**
- `internal/service/wechat_pay.go` - Payment service logic
- `internal/api/wechat_pay_handler.go` - API handlers

### 4.9 Distributed Tracing (OpenTelemetry + Jaeger via OTLP gRPC)
- **OpenTelemetry Integration**: Standardized tracing API
- **OTLP gRPC Exporter**: Sends traces to Jaeger via gRPC protocol (port 4317)
- **Key Features**:
  - Request latency tracking
  - Error tracking and correlation
  - Custom attributes (UserID, ProductID, OrderNo)
  - Trace context propagation

### 4.10 Microservices Architecture

**Supported Deployment Modes:**

1. **Monolithic Mode (Default)**
   - Single binary, all features in one process
   - Simplest deployment
   - Suitable for small to medium traffic

2. **Microservices Mode**
   - API Gateway routes requests to dedicated services
   - Independent scaling per service
   - Service discovery and registration
   - Suitable for high scalability needs

**Service Components:**

| Service | Port | Responsibility |
|---------|------|----------------|
| API Gateway | 8080 | Request routing, authentication |
| User Service | 8081 | User registration, login |
| Product Service | 8082 | Product CRUD |
| Order Service | 8083 | Order management |
| Stock Service | 8084 | Inventory, seckill operations |

**Service Discovery:**
- In-memory registry (single instance)
- Redis-based registry (distributed)
- Health checks and automatic deregistration

### 4.8 Observability Stack

**Health Checks:**
- `/health` - Liveness probe, checks all dependencies
- `/ready` - Readiness probe, checks if app is ready for traffic

**Prometheus Metrics:**
- HTTP request count/latency histograms
- Business metrics: orders, seckill attempts, user logins
- Infrastructure metrics: DB connections, Redis ping, RabbitMQ messages

**Structured Logging:**
- Uber Zap integration
- JSON format for production, console for development
- Request/response logging middleware

**Configuration Management:**
- Viper with YAML files
- Environment variable override (GOMALL_* prefix)
- Hot reload via SIGHUP signal

## 5. Deployment

- **Single Binary**: The entire application compiles into a single binary (`gomall.exe` / `main`).
- **Containerization**: `Dockerfile` provided for containerized deployment.
- **Orchestration**: `docker-compose.yml` orchestrates the App, MySQL, Redis, and RabbitMQ.

## 7. File Upload

**Features:**
- Single file upload (`POST /api/upload`)
- Multi file upload (`POST /api/upload/multi`)
- Supported formats: jpg, jpeg, png, gif
- Configurable file size limit
- Static file serving (`/uploads`)

**Configuration:**
```yaml
upload:
  path: "./uploads"  # Storage path
  max_size: 10       # Max file size in MB
  allowed_types:     # Allowed file types
    - jpg
    - jpeg
    - png
    - gif
```

## 6. Directory Map

| Path | Purpose |
|------|---------|
| `cmd/` | Main application entry point (`main.go`). |
| `internal/api/` | HTTP Handlers - User, Product, Order, Cart, Seckill, Auth, File, WeChatPay, HealthCheck. |
| `internal/service/` | Business Logic - User, Product, Order, Cart, Seckill, WeChatPay services. |
| `internal/repository/` | DB and Cache interactions. |
| `internal/model/` | Data entities - User, Product, Order, Cart, Stock models. |
| `internal/middleware/` | JWT Auth, Admin Auth, Rate Limiting, Logging, Metrics, Error Handling, Validation. |
| `internal/response/` | Unified response format and error codes. |
| `internal/router/` | Gin route definitions + Swagger integration. |
| `internal/rabbitmq/` | Message queue producer/consumer. |
| `internal/redis/` | Redis client and Lua scripts. |
| `internal/tracing/` | OpenTelemetry + Jaeger tracing integration. |
| `internal/grpc/` | gRPC service definitions (future use). |
| `internal/logger/` | Uber Zap structured logging. |
| `internal/metrics/` | Prometheus metrics definitions. |
| `internal/registry/` | Service discovery and registration. |
| `internal/gateway/` | API Gateway for microservices. |
| `internal/config/` | Configuration loading with Viper. |
| `internal/database/` | MySQL connection with GORM. |
| `conf/` | Configuration files (`config.yaml`, `config-dev.yaml`, `config-prod.yaml`). |
| `docs/` | Swagger API documentation. |
| `scripts/` | Database backup scripts (Linux/Mac + Windows). |
| `deploy/` | Docker deployment configurations. |
| `pkg/` | Shared utilities (JWT, Password hashing). |
