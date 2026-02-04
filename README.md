# GoMall 高并发分布式电商秒杀系统

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.23-blue.svg)](https://go.dev/)
[![Gin Framework](https://img.shields.io/badge/Gin-1.9-green.svg)](https://gin-gonic.com/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

</div>

## 项目简介

GoMall 是一个从**单体架构逐步演进到微服务架构**的电商实战项目，专注于解决高并发场景下的秒杀难题。项目完整实现了电商核心业务模块，并重点攻克了**超卖、少卖、高并发**等技术痛点。

### 核心特性

| 特性 | 实现方案 | 效果 |
|------|----------|------|
| **高性能** | Gin + Redis 缓存 + Lua 原子脚本 | 支撑万级 QPS 秒杀 |
| **高可用** | RabbitMQ 异步削峰 + 优雅关闭 | 系统稳定运行 |
| **可观测** | OpenTelemetry + Prometheus + Zap | 全链路监控 |
| **可扩展** | 模块化设计 + 服务注册发现 | 支持微服务拆分 |

---

## 技术栈

### 核心开发
| 技术 | 版本 | 用途 |
|------|------|------|
| Go | 1.23+ | 后端开发语言 |
| Gin | v1.9.1 | HTTP Web 框架 |
| GORM | v1.25.5 | MySQL ORM |
| gRPC | v1.60.1 | 微服务通信 |
| Viper | v1.18.2 | 配置管理 |

### 中间件
| 技术 | 版本 | 用途 |
|------|------|------|
| MySQL | 8.0+ | 持久化存储 |
| Redis | 7.0+ | 缓存、分布式锁、计数器 |
| RabbitMQ | 3.12+ | 流量削峰、异步解耦 |

### 运维监控
| 技术 | 用途 |
|------|------|
| OpenTelemetry + Jaeger | 分布式链路追踪 |
| Prometheus | 指标监控 |
| Swagger | API 文档 |
| Uber Zap | 结构化日志 |
| Docker | 容器化部署 |

---

## 架构设计

```
┌─────────────────────────────────────────────────────────────────┐
│                         用户请求                                  │
└─────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                      API Gateway (8080)                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐   │
│  │ JWT 认证    │  │ 限流中间件  │  │ Prometheus 指标采集   │   │
│  └─────────────┘  └─────────────┘  └─────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                        Gin HTTP Server                           │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                    Controller 层                          │   │
│  │  用户 │ 商品 │ 订单 │ 购物车 │ 秒杀 │ 健康检查            │   │
│  └──────────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                    Service 层                             │   │
│  │  业务逻辑处理 + Redis/Lua 库存扣减 + RabbitMQ 消息投递   │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                                 │
              ┌──────────────────┼──────────────────┐
              ▼                  ▼                  ▼
       ┌──────────┐      ┌──────────┐      ┌──────────┐
       │  MySQL   │      │  Redis   │      │RabbitMQ  │
       │  (持久化) │      │  (缓存)  │      │ (异步)   │
       └──────────┘      └──────────┘      └──────────┘
```

### 秒杀核心流程

```
用户请求 → Redis 预检查 → Lua 原子扣减 → MQ 异步下单 → 数据库落库
   │            │             │            │           │
   │            │             │            │           │
   └────────────┴─────────────┴────────────┴───────────┘
                        │
                        ▼
                  返回"排队中"
```

---

## 目录结构

```
gomall/
├── cmd/                      # 程序入口
│   └── main.go               # 主程序入口
├── conf/                     # 配置文件
│   ├── config.yaml           # 默认配置
│   ├── config-dev.yaml      # 开发环境
│   └── config-prod.yaml      # 生产环境
├── deploy/                   # 部署配置
│   ├── docker-compose.yml   # Docker Compose
│   └── mysql/init.sql       # 数据库初始化
├── docs/                    # Swagger 文档
├── internal/                # 内部业务代码
│   ├── api/                 # HTTP Handlers
│   │   ├── handler.go       # 用户/商品/订单处理器
│   │   ├── cart_handler.go  # 购物车
│   │   ├── seckill_handler.go # 秒杀
│   │   └── health_check.go  # 健康检查
│   ├── config/              # 配置加载
│   ├── database/            # MySQL 连接
│   ├── gateway/             # API 网关
│   ├── grpc/                # gRPC 服务
│   ├── logger/              # Zap 日志
│   ├── metrics/             # Prometheus 指标
│   ├── middleware/          # 中间件
│   │   ├── auth.go          # JWT 认证
│   │   ├── ratelimit.go     # 限流
│   │   └── metrics.go       # 指标
│   ├── model/               # 数据模型
│   ├── rabbitmq/            # 消息队列
│   ├── redis/               # Redis 客户端
│   ├── registry/            # 服务注册发现
│   ├── repository/          # 数据访问层
│   ├── router/              # 路由配置
│   ├── service/             # 业务逻辑
│   │   ├── service.go       # 基础服务
│   │   └── seckill.go       # 秒杀核心
│   └── tracing/             # 链路追踪
├── pkg/                     # 公共工具
│   ├── jwt/                 # JWT 工具
│   └── password/            # 密码加密
├── scripts/                 # 运维脚本
├── ARCHITECTURE.md          # 架构文档
├── Dockerfile
├── Makefile
└── README.md
```

---

## 快速开始

### 环境要求

- Go 1.23+
- MySQL 8.0
- Redis 7.0+
- RabbitMQ 3.12+ (可选)

### 1. 克隆项目

```bash
git clone https://github.com/yourusername/gomall.git
cd gomall
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 配置数据库

编辑 `conf/config-dev.yaml`：

```yaml
database:
  host: "localhost"
  port: 3306
  username: "root"
  password: "your_password"
  name: "gomall"

redis:
  host: "localhost"
  port: 6379
  password: ""

rabbitmq:
  host: "localhost"
  port: 5672
  username: "guest"
  password: "guest"
```

### 4. 初始化数据库

```bash
mysql -u root -p < deploy/mysql/init.sql
```

### 5. 启动服务

```bash
# 开发模式 (默认使用 config-dev.yaml)
make run

# 生产模式 (使用 config-prod.yaml)
make run-prod

# 或直接运行
go run main.go -env dev
```

### 6. 访问服务

| 服务 | 地址 |
|------|------|
| API 服务 | http://localhost:8080 |
| 健康检查 | http://localhost:8080/health |
| 就绪检查 | http://localhost:8080/ready |
| Prometheus 指标 | http://localhost:8080/metrics |
| Swagger 文档 | http://localhost:8080/swagger/index.html |

---

## API 文档

### 用户模块

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/user/register` | 用户注册 |
| POST | `/api/user/login` | 用户登录 |
| GET | `/api/user/profile` | 获取用户信息 |

### 商品模块

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/product` | 商品列表 |
| GET | `/api/product/:id` | 商品详情 |
| POST | `/api/product` | 创建商品 (需登录) |
| PUT | `/api/product/:id` | 更新商品 (需登录) |
| DELETE | `/api/product/:id` | 删除商品 (需登录) |

### 订单模块

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/order` | 创建订单 (需登录) |
| GET | `/api/order` | 订单列表 (需登录) |
| GET | `/api/order/:order_no` | 订单详情 (需登录) |
| POST | `/api/order/:order_no/pay` | 支付订单 (需登录) |
| POST | `/api/order/:order_no/cancel` | 取消订单 (需登录) |

### 购物车模块

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/cart` | 添加商品 (需登录) |
| GET | `/api/cart` | 购物车列表 (需登录) |
| PUT | `/api/cart` | 更新数量 (需登录) |
| DELETE | `/api/cart` | 删除商品 (需登录) |
| DELETE | `/api/cart/clear` | 清空购物车 (需登录) |

### 秒杀模块

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/seckill` | 秒杀接口 (需登录) |
| POST | `/api/seckill/init` | 初始化库存 (需管理员) |

---

## 核心功能实现

### 1. 高并发秒杀 (Redis + Lua)

```go
// 使用 Lua 脚本原子扣减库存，防止超卖
func decrStockWithLua(ctx context.Context, productID uint, quantity int) (int, error) {
    script := redis.NewScript(`
        local stock = redis.call('GET', KEYS[1])
        if stock == false then return -1 end
        stock = tonumber(stock)
        if stock < tonumber(ARGV[1]) then return -1 end
        redis.call('DECRBY', KEYS[1], ARGV[1])
        return stock - ARGV[1]
    `)
    return script.Run(ctx, redis.Client, []string{key}, quantity).Int()
}
```

**秒杀流程：**
1. 预加载库存到 Redis
2. 用户请求检查 Redis 库存
3. Lua 脚本原子扣减
4. 发送消息到 RabbitMQ
5. 异步消费者创建订单

### 2. JWT 认证

- bcrypt 密码加密
- JWT Token 生成与验证
- 中间件拦截认证

### 3. 限流策略

| 场景 | QPS | 突发 | 实现方式 |
|------|-----|------|----------|
| 全局 | 1000 | 2000 | 本地限流 |
| API | 100 | 200 | 本地限流 |
| 秒杀 | 5 | 10 | Redis 分布式 |
| 登录 | 10 | 20 | 本地限流 |

### 4. 链路追踪

集成 OpenTelemetry，通过 OTLP gRPC 导出到 Jaeger：

```yaml
tracing:
  enabled: true
  service_name: "gomall-service"
  jaeger_endpoint: "localhost:4317"
```

### 5. 优雅关闭

支持 SIGHUP 信号热重载配置：

```bash
kill -HUP <pid>
```

---

## Docker 部署

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
| GoMall | 8080 |
| MySQL | 3306 |
| Redis | 6379 |
| RabbitMQ | 5672 / 15672 |
| Jaeger UI | 16686 |
| Jaeger OTLP | 4317 |

---

## Makefile 命令

```bash
make deps          # 下载依赖
make build         # 编译项目
make run           # 开发模式运行
make run-prod      # 生产模式运行
make stop          # 停止服务
make clean         # 清理构建文件
make test          # 运行测试
make docker-build  # 构建 Docker 镜像
make docker-run    # 启动 Docker 服务
make docker-stop   # 停止 Docker 服务
make logs          # 查看日志
make backup        # 数据库备份
make swag          # 生成 Swagger 文档
make help          # 显示帮助
```

---

## 微服务架构

### 部署模式

**单体模式 (默认)：**
- 单一进程，所有功能集成
- 部署简单，适合中小流量

**微服务模式：**

| 服务 | 端口 | 职责 |
|------|------|------|
| API Gateway | 8080 | 请求路由、认证 |
| User Service | 8081 | 用户管理 |
| Product Service | 8082 | 商品管理 |
| Order Service | 8083 | 订单管理 |
| Stock Service | 8084 | 库存管理 |

### 启动微服务

```bash
# Docker Compose 方式
docker-compose -f deploy/docker-compose-microservices.yml up -d

# 独立启动
./app -service=user -port=8081
./app -service=product -port=8082
./app -service=order -port=8083
./app -service=stock -port=8084
./app -gateway -port=8080
```

---

## 项目进度

| 阶段 | 状态 | 内容 |
|------|------|------|
| Phase 1 | ✅ | 单体架构基础、用户/商品/订单模块 |
| Phase 2 | ✅ | 高并发秒杀、Redis + Lua + RabbitMQ |
| Phase 3 | ✅ | 购物车模块 |
| Phase 4 | ✅ | 可观测性、健康检查、限流、监控 |
| Phase 5 | ✅ | 微服务架构、服务注册发现 |

---

## 贡献指南

1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

---

## 许可证

本项目采用 MIT License 开源。

---

## 致谢

感谢所有开源贡献者！
