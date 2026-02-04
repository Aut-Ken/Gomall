# 🛒 GoMall - 高并发分布式电商秒杀系统

> 一个基于 Golang + Gin + GORM + MySQL + Redis + RabbitMQ + gRPC 构建的分布式电商平台。
> 本项目旨在解决高并发场景下的"超卖"、"少卖"问题，并实践微服务架构拆分与治理。

## 📖 项目简介 (Introduction)

**GoMall** 是一个从单体架构逐步演进到微服务架构的电商实战项目。项目涵盖了电商核心业务模块（用户、商品、订单、库存），并重点攻克**秒杀高并发**场景下的技术难点。

**核心目标：**
- **高并发：** 通过 Redis 缓存、Lua 脚本、消息队列削峰，支撑万级 QPS 秒杀。
- **高可用：** 结合 Docker 容器化部署，保障系统稳定性。
- **分布式：** 实践 gRPC 微服务通信、分布式链路追踪（预留）。

---

## 🛠 技术栈 (Tech Stack)

### 核心开发
| 技术 | 用途 |
|------|------|
| Golang 1.20+ | 后端开发语言 |
| Gin | 高性能 HTTP Web 框架 |
| GORM | MySQL 数据库操作 |
| gRPC + Protobuf | 微服务通信 |
| Viper | 配置管理 |

### 中间件 & 存储
| 技术 | 用途 |
|------|------|
| MySQL 8.0 | 持久化存储 |
| Redis 7.0 | 缓存、分布式锁、计数器 |
| RabbitMQ | 流量削峰、异步解耦 |

### 运维 & 监控
| 技术 | 用途 |
|------|------|
| Docker | 容器化部署 |
| Docker Compose | 本地开发环境 |
| OpenTelemetry | 链路追踪标准 |
| Jaeger (OTLP gRPC) | 分布式追踪系统 |
| golang.org/x/time | 本地限流 |

---

## 📂 目录结构 (Directory Structure)

```
gomall/
├── cmd/                    # 程序入口
│   └── main.go             # 主程序入口
├── conf/                   # 配置文件
│   └── config.yaml         # 应用配置
├── deploy/                 # 部署配置
│   ├── docker-compose.yml  # Docker Compose 配置
│   └── mysql/
│       └── init.sql        # 数据库初始化脚本
├── internal/               # 内部业务代码
│   ├── api/                # HTTP Handlers (Controllers)
│   │   ├── handler.go       # 用户、商品、订单处理器
│   │   ├── cart_handler.go  # 购物车处理器
│   │   └── seckill_handler.go  # 秒杀处理器
│   ├── config/             # 配置加载
│   │   └── config.go
│   ├── database/           # 数据库连接
│   │   └── database.go
│   ├── grpc/               # gRPC 服务
│   │   ├── grpc.go         # gRPC 服务实现
│   │   └── proto/          # Protobuf 定义
│   ├── middleware/         # 中间件
│   │   ├── auth.go         # JWT 认证中间件
│   │   └── ratelimit.go    # 限流中间件
│   ├── tracing/            # 链路追踪 (OpenTelemetry/Jaeger)
│   │   └── tracing.go
│   ├── model/              # 数据模型 (GORM)
│   │   └── model.go
│   ├── rabbitmq/           # RabbitMQ 消息队列
│   │   └── rabbitmq.go
│   ├── redis/              # Redis 缓存
│   │   └── redis.go
│   ├── repository/         # 数据访问层
│   │   └── repository.go
│   ├── router/             # 路由配置
│   │   └── router.go
│   └── service/            # 业务逻辑层
│       ├── service.go      # 用户、商品、订单、购物车服务
│       └── seckill.go      # 秒杀服务
├── pkg/                    # 公共工具库
│   ├── jwt/                # JWT 工具
│   │   └── jwt.go
│   └── password/           # 密码加密
│       └── password.go
├── .gitignore
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## 🚀 快速开始

### 1. 环境要求

- Golang 1.20+
- MySQL 8.0
- Redis 7.0+
- RabbitMQ 3.12+ (可选)

### 2. 安装依赖

```bash
# 下载依赖
go mod download

# 或者使用 Makefile
make deps
```

### 3. 配置数据库

编辑 `conf/config.yaml`，修改数据库连接信息：

```yaml
database:
  host: "localhost"
  port: 3306
  username: "root"
  password: "your_password"
  name: "gomall"
```

### 4. 初始化数据库

```bash
# 执行 SQL 脚本
mysql -u root -p gomall < deploy/mysql/init.sql
```

### 5. 启动服务

```bash
# 方式一：直接运行
go run main.go -config conf/config.yaml

# 方式二：使用 Makefile
make run

# 方式三：Docker 部署
make docker-build
make docker-run
```

### 6. 访问服务

- **服务地址**: http://localhost:8080
- **健康检查**: http://localhost:8080/health

---

## 📡 API 文档

### 用户模块

| 方法 | 路径 | 说明 | 参数 |
|------|------|------|------|
| POST | /api/user/register | 用户注册 | username, password, email, phone |
| POST | /api/user/login | 用户登录 | username, password |
| GET | /api/user/profile | 获取个人信息 | Authorization Header |

### 商品模块

| 方法 | 路径 | 说明 | 参数 |
|------|------|------|------|
| GET | /api/product | 商品列表 | page, page_size, category |
| GET | /api/product/:id | 商品详情 | - |
| POST | /api/product | 创建商品 | name, price, stock... (需登录) |
| PUT | /api/product/:id | 更新商品 | name, price, stock... (需登录) |
| DELETE | /api/product/:id | 删除商品 | - (需登录) |

### 订单模块

| 方法 | 路径 | 说明 | 参数 |
|------|------|------|------|
| POST | /api/order | 创建订单 | product_id, quantity (需登录) |
| GET | /api/order | 订单列表 | page, page_size (需登录) |
| GET | /api/order/:order_no | 订单详情 | - (需登录) |
| POST | /api/order/:order_no/pay | 支付订单 | - (需登录) |
| POST | /api/order/:order_no/cancel | 取消订单 | - (需登录) |

### 购物车模块

| 方法 | 路径 | 说明 | 参数 |
|------|------|------|------|
| POST | /api/cart | 添加商品到购物车 | product_id, quantity (需登录) |
| GET | /api/cart | 获取购物车列表 | - (需登录) |
| PUT | /api/cart | 更新购物车商品数量 | product_id, quantity (需登录) |
| DELETE | /api/cart | 删除购物车商品 | product_id (需登录) |
| DELETE | /api/cart/clear | 清空购物车 | - (需登录) |

### 秒杀模块

| 方法 | 路径 | 说明 | 参数 |
|------|------|------|------|
| POST | /api/seckill | 秒杀接口 | product_id (需登录) |
| POST | /api/seckill/init | 初始化秒杀库存 | product_id, stock (需管理员) |

---

## 🔧 Makefile 命令

```bash
make deps        # 下载依赖
make build       # 编译项目
make run         # 运行项目
make stop        # 停止服务
make clean       # 清理构建文件
make test        # 运行测试
make docker-build  # 构建Docker镜像
make docker-run    # 启动Docker服务
make docker-stop   # 停止Docker服务
make logs         # 查看日志
make help         # 显示帮助信息
```

---

## 🔐 核心功能实现

### 1. 用户认证 (JWT)
- 用户注册密码使用 bcrypt 加密
- 登录后生成 JWT Token
- 中间件验证 Token 有效性

### 2. 高并发秒杀 (Redis + Lua)
```
流程：
1. 秒杀开始前预加载库存到 Redis
2. 用户请求先检查 Redis 库存（内存级别，快速）
3. 使用 Lua 脚本原子扣减库存（防止超卖）
4. 扣减成功发送消息到 RabbitMQ 异步创建订单
5. 订单超时未支付自动取消（延迟队列）
```

### 3. 库存安全
- **数据库层面**：事务 + 悲观锁
- **Redis 层面**：Lua 脚本原子操作
- **消息队列**：异步削峰，流量控制

### 4. 限流与熔断
- **本地限流**：基于 golang.org/x/time/rate 实现 IP 级别限流
- **分布式限流**：基于 Redis 实现滑动窗口算法，支持多实例共享
- **限流配置**：
  - 全局：1000 QPS，突发 2000
  - API：100 QPS，突发 200
  - 秒杀：5 QPS，突发 10
  - 登录：10 QPS，突发 20

### 5. 链路追踪 (OpenTelemetry + Jaeger)
- **集成 OpenTelemetry**：标准化的链路追踪方案
- **OTLP gRPC 导出器**：通过 gRPC 协议将追踪数据发送到 Jaeger
- **Jaeger 可视化**：支持请求链路、延迟分析、错误追踪
- **自定义属性**：支持 UserID、ProductID、OrderNo 等业务标签

---

## 🐳 Docker 部署

### 使用 Docker Compose

```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f app

# 停止服务
docker-compose down
```

### 服务端口

| 服务 | 端口 |
|------|------|
| GoMall App | 8080 |
| MySQL | 3306 |
| Redis | 6379 |
| RabbitMQ | 5672 (AMQP), 15672 (管理) |
| Jaeger UI | 16686 |
| Jaeger OTLP gRPC | 4317 |
| Jaeger OTLP HTTP | 4318 |

---

## 📝 开发计划

- [x] Phase 1: 单体架构基础
  - [x] 数据库表结构设计
  - [x] 用户模块 (注册、JWT 登录、鉴权)
  - [x] 商品模块 (CRUD、列表展示)
  - [x] 基础下单流程
- [x] Phase 2: 高并发秒杀核心
  - [x] Redis 缓存预热
  - [x] Lua 脚本实现库存原子扣减
  - [x] RabbitMQ 异步创建订单
  - [x] 解决超卖问题
- [x] Phase 3: 购物车模块
  - [x] 添加商品到购物车
  - [x] 获取/更新/删除购物车商品
  - [x] 清空购物车
- [x] Phase 4: 稳定性与部署
  - [x] 接入 Jaeger/OpenTelemetry 链路追踪
  - [x] Docker Compose 一键部署
  - [x] 限流中间件 (IP + Redis 分布式限流)
- [ ] Phase 5: 微服务拆分 (待实现)
  - [ ] 拆分为 User/Product/Order/Stock 独立服务
  - [ ] 引入 gRPC 进行服务间通信
  - [ ] 使用 Consul 进行服务注册与发现

---

## 📄 许可证

MIT License

---

## 🙏 致谢

本项目参考了多个优秀的开源项目，致敬所有开源贡献者！
