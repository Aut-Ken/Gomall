package router

import (
	"github.com/gin-gonic/gin"
)

// APIVersion API版本信息
type APIVersion struct {
	Version   string
	Prefix    string
	Deprecated bool
}

// API 版本常量
const (
	APIVersionV1 = "v1"
)

// VersionRegistry API版本注册表
type VersionRegistry struct {
	versions map[string]*APIVersion
}

// NewVersionRegistry 创建版本注册表
func NewVersionRegistry() *VersionRegistry {
	return &VersionRegistry{
		versions: make(map[string]*APIVersion),
	}
}

// Register 注册API版本
func (r *VersionRegistry) Register(version, prefix string) {
	r.versions[version] = &APIVersion{
		Version: version,
		Prefix:  prefix,
	}
}

// Get 获取API版本信息
func (r *VersionRegistry) Get(version string) *APIVersion {
	return r.versions[version]
}

// GetPrefix 获取版本前缀
func (r *VersionRegistry) GetPrefix(version string) string {
	if v, ok := r.versions[version]; ok {
		return v.Prefix
	}
	return "/api/" + version
}

// SetupVersionedRoutes 设置带版本的路由
func SetupVersionedRoutes(r *gin.Engine) {
	versionRegistry := NewVersionRegistry()
	versionRegistry.Register(APIVersionV1, "/api/v1")

	// 默认版本路由组（指向最新版本v1）
	defaultGroup := r.Group("/api")
	{
		setupDefaultRoutes(defaultGroup, versionRegistry)
	}

	// V1版本路由组
	v1Group := r.Group("/api/v1")
	{
		setupV1Routes(v1Group)
	}
}

// setupDefaultRoutes 设置默认路由（带版本重定向提示）
func setupDefaultRoutes(group *gin.RouterGroup, versions *VersionRegistry) {
	group.GET("/version", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"current_version": APIVersionV1,
			"available_versions": []string{APIVersionV1},
			"message": "建议使用 /api/v1 前缀访问API",
		})
	})
}

// setupV1Routes 设置V1版本路由
func setupV1Routes(group *gin.RouterGroup) {
	// 健康检查（各版本独立）
	group.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"version": APIVersionV1,
		})
	})

	// 注意：实际的用户模块路由在 router.go 中定义
	_ = group.Group("/user")
}

// VersionedHandler 版本化处理器包装器
func VersionedHandler(handler gin.HandlerFunc, version string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("api_version", version)
		handler(c)
	}
}

// DeprecationMiddleware 过时API中间件
func DeprecationMiddleware(sunsetDate string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Sunset", sunsetDate)
		c.Header("Deprecation", "true")

		// 可选：添加提醒头
		c.Header("X-API-Deprecated", "true")
		c.Header("X-API-Sunset-Date", sunsetDate)

		c.Next()
	}
}
