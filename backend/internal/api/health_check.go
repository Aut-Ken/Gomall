package api

import (
	"net/http"

	db "gomall/backend/internal/database"
	"gomall/backend/internal/rabbitmq"
	rds "gomall/backend/internal/redis"

	"github.com/gin-gonic/gin"
)

// HealthCheck 健康检查处理器
type HealthCheck struct{}

// NewHealthCheck 创建健康检查处理器
func NewHealthCheck() *HealthCheck {
	return &HealthCheck{}
}

// Health 健康检查端点
// @Summary 健康检查
// @Description 检查应用和依赖服务的健康状态
// @Tags 系统
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h *HealthCheck) Health(c *gin.Context) {
	status := "healthy"
	components := make(map[string]string)

	// 检查数据库
	if err := db.Ping(); err != nil {
		components["database"] = "unhealthy: " + err.Error()
		status = "unhealthy"
	} else {
		components["database"] = "healthy"
	}

	// 检查Redis
	if err := rds.Ping(); err != nil {
		components["redis"] = "unhealthy: " + err.Error()
		status = "unhealthy"
	} else {
		components["redis"] = "healthy"
	}

	// 检查RabbitMQ
	if err := rabbitmq.Ping(); err != nil {
		// RabbitMQ故障不影响整体健康状态（可降级）
		components["rabbitmq"] = "degraded: " + err.Error()
		if status == "healthy" {
			status = "degraded"
		}
	} else {
		components["rabbitmq"] = "healthy"
	}

	code := http.StatusOK
	if status == "unhealthy" {
		code = http.StatusServiceUnavailable
	}

	c.JSON(code, gin.H{
		"status":     status,
		"components": components,
	})
}

// Ready 就绪检查端点
// @Summary 就绪检查
// @Description 检查应用是否准备好接收流量（用于K8s就绪探针）
// @Tags 系统
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /ready [get]
func (h *HealthCheck) Ready(c *gin.Context) {
	// 检查数据库是否可用（就绪必要条件）
	if err := db.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not_ready",
			"reason": "database not available",
		})
		return
	}

	// Redis可选（部分功能可降级）
	redisStatus := "ok"
	if err := rds.Ping(); err != nil {
		redisStatus = "degraded"
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"components": gin.H{
			"database": "ok",
			"redis":    redisStatus,
		},
	})
}
