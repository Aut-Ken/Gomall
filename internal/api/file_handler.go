package api

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gomall/internal/config"
	"gomall/internal/middleware"
	"gomall/internal/response"

	"github.com/gin-gonic/gin"
)

// FileHandler 文件处理层
type FileHandler struct {
}

// NewFileHandler 创建文件处理器
func NewFileHandler() *FileHandler {
	return &FileHandler{}
}

// UploadRequest 上传请求
type UploadRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

// UploadResponse 上传响应
type UploadResponse struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

// Upload 单文件上传
// @Summary 单文件上传
// @Description 上传图片文件，支持 jpg, jpeg, png, gif
// @Tags 文件上传
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "上传文件"
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/upload [post]
func (h *FileHandler) Upload(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "需要登录")
		return
	}

	// 读取配置文件中的上传配置
	uploadConfig := config.GetApp().Sub("upload")
	if uploadConfig == nil {
		// 使用默认值
		uploadConfig = config.Config.Sub("app")
	}

	// 获取上传配置
	maxSize := uploadConfig.GetInt64("max_size") // 单位：MB，默认5MB
	if maxSize <= 0 {
		maxSize = 5
	}

	allowedTypes := uploadConfig.GetStringSlice("allowed_types")
	if len(allowedTypes) == 0 {
		allowedTypes = []string{"jpg", "jpeg", "png", "gif"}
	}

	uploadPath := uploadConfig.GetString("path")
	if uploadPath == "" {
		uploadPath = "./uploads"
	}

	domain := uploadConfig.GetString("domain")
	if domain == "" {
		domain = "http://localhost:8080"
	}

	// 读取文件
	file, err := c.FormFile("file")
	if err != nil {
		response.FailWithMsg(c, response.CodeUploadFileEmpty, "请选择要上传的文件")
		return
	}

	// 检查文件大小
	if file.Size > maxSize*1024*1024 {
		response.FailWithMsg(c, response.CodeUploadFileTooLarge, fmt.Sprintf("文件大小不能超过 %dMB", maxSize))
		return
	}

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	ext = strings.TrimPrefix(ext, ".")
	if !contains(allowedTypes, ext) {
		response.FailWithMsg(c, response.CodeUploadFileTypeError, "不支持的文件类型，仅支持: "+strings.Join(allowedTypes, ", "))
		return
	}

	// 生成文件名
	filename := fmt.Sprintf("%d%s%s", time.Now().UnixNano(), randomString(8), ext)
	dir := filepath.Join(uploadPath, time.Now().Format("2006-01-02"))

	// 创建目录
	if err := os.MkdirAll(dir, 0755); err != nil {
		response.FailWithMsg(c, response.CodeUploadSaveFailed, "创建上传目录失败")
		return
	}

	// 保存文件
	dst := filepath.Join(dir, filename)
	src, err := file.Open()
	if err != nil {
		response.FailWithMsg(c, response.CodeUploadSaveFailed, "打开文件失败")
		return
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		response.FailWithMsg(c, response.CodeUploadSaveFailed, "创建文件失败")
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		response.FailWithMsg(c, response.CodeUploadSaveFailed, "保存文件失败")
		return
	}

	// 返回文件访问URL
	fileURL := fmt.Sprintf("%s/uploads/%s/%s", domain, time.Now().Format("2006-01-02"), filename)

	response.OkWithData(c, UploadResponse{
		URL:      fileURL,
		Filename: file.Filename,
		Size:     file.Size,
	})
}

// UploadMulti 多文件上传
// @Summary 多文件上传
// @Description 一次上传多个图片文件
// @Tags 文件上传
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "上传文件"
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/upload/multi [post]
func (h *FileHandler) UploadMulti(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "需要登录")
		return
	}

	// 读取配置文件中的上传配置
	uploadConfig := config.GetApp().Sub("upload")
	if uploadConfig == nil {
		uploadConfig = config.Config.Sub("app")
	}

	maxSize := uploadConfig.GetInt64("max_size")
	if maxSize <= 0 {
		maxSize = 5
	}

	allowedTypes := uploadConfig.GetStringSlice("allowed_types")
	if len(allowedTypes) == 0 {
		allowedTypes = []string{"jpg", "jpeg", "png", "gif"}
	}

	uploadPath := uploadConfig.GetString("path")
	if uploadPath == "" {
		uploadPath = "./uploads"
	}

	domain := uploadConfig.GetString("domain")
	if domain == "" {
		domain = "http://localhost:8080"
	}

	// 读取多个文件
	form, err := c.MultipartForm()
	if err != nil {
		response.FailWithMsg(c, response.CodeUploadFileEmpty, "请选择要上传的文件")
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		response.FailWithMsg(c, response.CodeUploadFileEmpty, "请选择要上传的文件")
		return
	}

	// 限制文件数量
	if len(files) > 10 {
		response.FailWithMsg(c, response.CodeUploadTooManyFiles, "单次最多上传10个文件")
		return
	}

	var results []UploadResponse
	for _, file := range files {
		// 检查文件大小
		if file.Size > maxSize*1024*1024 {
			continue // 跳过过大文件
		}

		// 检查文件类型
		ext := strings.ToLower(filepath.Ext(file.Filename))
		ext = strings.TrimPrefix(ext, ".")
		if !contains(allowedTypes, ext) {
			continue // 跳过不支持的文件类型
		}

		// 生成文件名
		filename := fmt.Sprintf("%d%s%s", time.Now().UnixNano(), randomString(8), ext)
		dir := filepath.Join(uploadPath, time.Now().Format("2006-01-02"))

		// 创建目录
		os.MkdirAll(dir, 0755)

		// 保存文件
		dst := filepath.Join(dir, filename)
		src, err := file.Open()
		if err != nil {
			continue
		}
		defer src.Close()

		out, err := os.Create(dst)
		if err != nil {
			continue
		}
		defer out.Close()

		io.Copy(out, src)

		// 返回文件访问URL
		fileURL := fmt.Sprintf("%s/uploads/%s/%s", domain, time.Now().Format("2006-01-02"), filename)

		results = append(results, UploadResponse{
			URL:      fileURL,
			Filename: file.Filename,
			Size:     file.Size,
		})
	}

	response.OkWithData(c, results)
}

// contains 检查字符串是否在切片中
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// randomString 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	time.Sleep(time.Nanosecond)
	return string(result)
}

// SetupStatic 配置静态文件服务
func SetupStatic(r *gin.Engine) {
	uploadConfig := config.GetApp().Sub("upload")
	if uploadConfig == nil {
		uploadConfig = config.Config.Sub("app")
	}

	uploadPath := uploadConfig.GetString("path")
	if uploadPath == "" {
		uploadPath = "./uploads"
	}

	r.Static("/uploads", uploadPath)
}
