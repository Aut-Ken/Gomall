# 基础镜像
FROM golang:1.23-alpine AS builder

# 设置工作目录
WORKDIR /app

# 设置Go代理
ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 运行镜像
FROM alpine:latest AS runner

# 设置时区
RUN apk --no-cache add tzdata

# 创建非root用户
RUN addgroup -g 1000 appgroup && \
    adduser -u 1000 -G appgroup -s /bin/sh -D appuser

# 设置工作目录
WORKDIR /app

# 复制编译好的应用
COPY --from=builder /app/main .
COPY --from=builder /app/conf ./conf

# 更改文件所有者
RUN chown -R appuser:appgroup /app

# 切换用户
USER appuser

# 暴露端口
EXPOSE 8080

# 启动命令
ENTRYPOINT ["./main"]
CMD ["-config", "conf/config-dev.yaml"]
