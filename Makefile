.PHONY: all build run stop clean test proto deps

# GoMall Makefile

# 编译输出目录
BUILD_DIR := bin

# 编译选项
CGO_ENABLED := 0
GOOS := linux
GOARCH := amd64
LDFLAGS := -s -w

# 默认目标
all: deps build

# 下载依赖
deps:
	@echo "下载依赖..."
	go mod download
	@echo "依赖下载完成"

# 编译项目
build:
	@echo "编译项目..."
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/main .
	@echo "编译完成: $(BUILD_DIR)/main"

# 运行项目
run: deps
	@echo "启动服务..."
	go run main.go -config conf/config.yaml

# 停止服务
stop:
	@echo "停止服务..."
	@if [ -f $(BUILD_DIR)/main.pid ]; then \
		kill $$(cat $(BUILD_DIR)/main.pid) 2>/dev/null; \
		rm $(BUILD_DIR)/main.pid; \
	fi
	@echo "服务已停止"

# 清理构建文件
clean:
	@echo "清理构建文件..."
	rm -rf $(BUILD_DIR)
	rm -f *.test
	@echo "清理完成"

# 运行测试
test:
	@echo "运行测试..."
	go test -v ./...
	@echo "测试完成"

# 生成proto文件
proto:
	@echo "生成proto文件..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/*.proto
	@echo "proto文件生成完成"

# Docker 构建
docker-build:
	@echo "构建Docker镜像..."
	docker build -t gomall:latest .
	@echo "Docker镜像构建完成"

# Docker 运行
docker-run:
	@echo "启动Docker服务..."
	docker-compose up -d
	@echo "Docker服务启动完成"

# Docker 停止
docker-stop:
	@echo "停止Docker服务..."
	docker-compose down
	@echo "Docker服务已停止"

# 查看日志
logs:
	docker-compose logs -f app

# 数据库迁移
migrate:
	@echo "执行数据库迁移..."
	go run -tags mysql cmd/migrate/main.go
	@echo "迁移完成"

# 代码检查
lint:
	@echo "代码检查..."
	golangci-lint run ./...
	@echo "检查完成"

# 帮助信息
help:
	@echo "GoMall Makefile 命令:"
	@echo ""
	@echo "  make deps        - 下载依赖"
	@echo "  make build       - 编译项目"
	@echo "  make run         - 运行项目"
	@echo "  make stop        - 停止服务"
	@echo "  make clean       - 清理构建文件"
	@echo "  make test        - 运行测试"
	@echo "  make proto       - 生成proto文件"
	@echo "  make docker-build  - 构建Docker镜像"
	@echo "  make docker-run    - 启动Docker服务"
	@echo "  make docker-stop   - 停止Docker服务"
	@echo "  make logs        - 查看日志"
	@echo "  make migrate     - 数据库迁移"
	@echo "  make lint        - 代码检查"
	@echo "  make help        - 显示帮助信息"
