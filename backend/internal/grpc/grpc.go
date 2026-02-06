package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"gomall/backend/internal/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server gRPC服务器
type Server struct {
	grpcServer *grpc.Server
}

// NewServer 创建gRPC服务器
func NewServer() *Server {
	return &Server{}
}

// Start 启动gRPC服务器
func (s *Server) Start() error {
	grpcConfig := config.GetGRPCConfig()
	port := grpcConfig.GetInt("port")

	// 创建TCP监听
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("监听端口失败: %w", err)
	}

	// 创建gRPC服务器
	s.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(logInterceptor),
	)

	log.Printf("gRPC服务启动中，端口: %d", port)
	return s.grpcServer.Serve(lis)
}

// Stop 停止gRPC服务器
func (s *Server) Stop() {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
}

// logInterceptor 日志拦截器
func logInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("gRPC调用: %s", info.FullMethod)

	resp, err := handler(ctx, req)
	if err != nil {
		log.Printf("gRPC调用失败: %v", err)
	}

	return resp, err
}

// ErrorToStatus 将错误转换为gRPC状态
func ErrorToStatus(err error) error {
	if err == nil {
		return nil
	}
	return status.Errorf(codes.Internal, err.Error())
}

// ValidateProductRequest 验证商品请求参数
func ValidateProductRequest(req interface{}) error {
	// 这里可以添加参数验证逻辑
	return nil
}
