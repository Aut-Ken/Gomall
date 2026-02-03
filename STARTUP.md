# GoMall 启动指南

## 环境要求

| 软件 | 版本要求 | 用途 |
|------|----------|------|
| Golang | 1.20+ | 开发语言 |
| MySQL | 8.0+ | 数据库 |
| Redis | 7.0+ | 缓存（可选） |
| RabbitMQ | 3.12+ | 消息队列（可选） |

---

## 第一步：修改配置

编辑 `conf/config.yaml`，修改数据库密码：

```yaml
database:
  host: "localhost"
  port: 3306
  username: "root"
  password: "你的MySQL密码"  # ← 修改这里
  name: "Gomall"
```

---

## 第二步：初始化数据库

在命令行执行：

```bash
# Windows
mysql -u root -p < deploy\mysql\init.sql

# Linux/Mac
mysql -u root -p < deploy/mysql/init.sql
```

**或使用 MySQL 客户端执行：**
```sql
source E:/Gomall/deploy/mysql/init.sql
```

---

## 第三步：安装依赖

```bash
go mod download
```

---

## 第四步：启动服务

```bash
# 方式一：直接运行
go run main.go -config conf/config.yaml

# 方式二：使用 Makefile（推荐）
make run
```

启动成功后看到以下输出：
```
正在初始化配置...
配置初始化成功
正在连接数据库...
数据库连接成功
服务启动成功，访问地址: http://localhost:8080
健康检查: http://localhost:8080/health
```

---

## 验证服务

打开浏览器访问：

| 地址 | 说明 |
|------|------|
| http://localhost:8080/health | 健康检查 |
| http://localhost:8080/api/product | 商品列表（无需登录） |

---

## 测试账号

数据库已预置以下测试账号（密码均为：`123456`）：

| 用户名 | 邮箱 | 用途 |
|--------|------|------|
| admin | admin@Gomall.com | 管理员 |
| testuser | test@Gomall.com | 普通用户 |

**密码加密后的值**：`$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi`

---

## Docker 启动（推荐）

如果你的环境有 Docker，可以使用一键启动：

```bash
# 启动所有服务（MySQL、Redis、RabbitMQ、GoMall）
docker-compose up -d

# 查看日志
docker-compose logs -f app

# 停止服务
docker-compose down
```

Docker 启动后自动初始化数据库，访问地址：http://localhost:8080

---

## API 测试示例

### 1. 用户注册
```bash
curl -X POST http://localhost:8080/api/user/register \
  -H "Content-Type: application/json" \
  -d '{"username":"newuser","password":"123456","email":"new@test.com","phone":"15000000000"}'
```

### 2. 用户登录
```bash
curl -X POST http://localhost:8080/api/user/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}'
```

**返回的 `token` 后续请求需要放在 Header 中：**
```
Authorization: Bearer <你的token>
```

### 3. 获取商品列表
```bash
curl http://localhost:8080/api/product
```

### 4. 创建订单（需登录）
```bash
curl -X POST http://localhost:8080/api/order \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <你的token>" \
  -d '{"product_id":1,"quantity":1}'
```

---

## 常见问题

### Q: 数据库连接失败？
A: 检查 `conf/config.yaml` 中的数据库密码是否正确，确保 MySQL 服务已启动。

### Q: Redis 连接失败？
A: 不影响基础功能，程序会自动降级到数据库方案。

### Q: 端口被占用？
A: 修改 `conf/config.yaml` 中的 `app.port` 为其他端口。

### Q: 秒杀功能如何使用？
A: 先调用 `POST /api/seckill` 预热库存，再调用秒杀接口。

---

## 项目结构

```
gomall/
├── conf/config.yaml          # 配置文件
├── deploy/
│   ├── docker-compose.yml    # Docker部署
│   └── mysql/init.sql        # 数据库初始化
├── internal/                  # 业务代码
│   ├── api/                  # HTTP处理器
│   ├── service/              # 业务逻辑
│   ├── repository/           # 数据访问
│   ├── model/                # 数据模型
│   ├── redis/                # Redis操作
│   └── rabbitmq/             # 消息队列
├── pkg/                       # 公共工具
├── main.go                    # 入口文件
└── Makefile                   # 构建脚本
```
