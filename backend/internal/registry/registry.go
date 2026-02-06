package registry

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ServiceInfo 服务信息
type ServiceInfo struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Host      string            `json:"host"`
	Port      int               `json:"port"`
	Tags      []string          `json:"tags"`
	Metadata  map[string]string `json:"metadata"`
	CreatedAt time.Time        `json:"created_at"`
}

// Registry 服务注册中心接口
type Registry interface {
	Register(ctx context.Context, service *ServiceInfo) error
	Deregister(ctx context.Context, serviceID string) error
	Discover(ctx context.Context, serviceName string) ([]*ServiceInfo, error)
	Close() error
}

// InMemoryRegistry 内存注册中心（单机模式）
type InMemoryRegistry struct {
	services map[string][]*ServiceInfo
	mu       sync.RWMutex
}

// NewInMemoryRegistry 创建内存注册中心
func NewInMemoryRegistry() *InMemoryRegistry {
	return &InMemoryRegistry{
		services: make(map[string][]*ServiceInfo),
	}
}

// Register 服务注册
func (r *InMemoryRegistry) Register(ctx context.Context, service *ServiceInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	service.CreatedAt = time.Now()
	r.services[service.Name] = append(r.services[service.Name], service)

	return nil
}

// Deregister 服务注销
func (r *InMemoryRegistry) Deregister(ctx context.Context, serviceID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for name, list := range r.services {
		for i, s := range list {
			if s.ID == serviceID {
				r.services[name] = append(list[:i], list[i+1:]...)
				return nil
			}
		}
	}
	return nil
}

// Discover 服务发现
func (r *InMemoryRegistry) Discover(ctx context.Context, serviceName string) ([]*ServiceInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.services[serviceName], nil
}

// Close 关闭
func (r *InMemoryRegistry) Close() error {
	return nil
}

// ConsulRegistry Consul/Redis 注册中心
type ConsulRegistry struct {
	client   Registry
	prefix   string
	services map[string]string // serviceName -> addr
}

// NewConsulRegistry 创建 Consul 注册中心
func NewConsulRegistry(client Registry) *ConsulRegistry {
	return &ConsulRegistry{
		client:   client,
		prefix:   "gomall:",
		services: make(map[string]string),
	}
}

// Register 服务注册
func (r *ConsulRegistry) Register(ctx context.Context, service *ServiceInfo) error {
	service.CreatedAt = time.Now()

	// 保存服务地址
	r.services[service.Name] = fmt.Sprintf("%s:%d", service.Host, service.Port)

	// 同时注册到后端存储
	return r.client.Register(ctx, service)
}

// Deregister 服务注销
func (r *ConsulRegistry) Deregister(ctx context.Context, serviceID string) error {
	return r.client.Deregister(ctx, serviceID)
}

// Discover 服务发现
func (r *ConsulRegistry) Discover(ctx context.Context, serviceName string) ([]*ServiceInfo, error) {
	return r.client.Discover(ctx, serviceName)
}

// GetServiceAddr 获取服务地址（简化方法）
func (r *ConsulRegistry) GetServiceAddr(serviceName string) string {
	return r.services[serviceName]
}

// Close 关闭
func (r *ConsulRegistry) Close() error {
	return r.client.Close()
}

// ServiceDiscovery 服务发现客户端
type ServiceDiscovery struct {
	registry Registry
}

// NewServiceDiscovery 创建服务发现客户端
func NewServiceDiscovery(registry Registry) *ServiceDiscovery {
	return &ServiceDiscovery{
		registry: registry,
	}
}

// GetService 获取服务实例（带负载均衡）
func (sd *ServiceDiscovery) GetService(serviceName string) (*ServiceInfo, error) {
	services, err := sd.registry.Discover(context.Background(), serviceName)
	if err != nil {
		return nil, err
	}

	if len(services) == 0 {
		return nil, fmt.Errorf("服务 %s 未找到", serviceName)
	}

	// 简单轮询负载均衡
	return services[0], nil
}

// GetAllServices 获取所有服务实例
func (sd *ServiceDiscovery) GetAllServices(serviceName string) ([]*ServiceInfo, error) {
	return sd.registry.Discover(context.Background(), serviceName)
}

// NewRegistry 创建注册中心
func NewRegistry(typ string) Registry {
	switch typ {
	case "memory":
		return NewInMemoryRegistry()
	default:
		return NewInMemoryRegistry()
	}
}
