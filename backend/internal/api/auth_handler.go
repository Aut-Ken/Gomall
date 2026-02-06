package api

import (
	"gomall/backend/internal/repository"
	"gomall/backend/pkg/jwt"
	"gomall/backend/pkg/password"

	"github.com/gin-gonic/gin"
	"gomall/backend/internal/response"
)

// AuthHandler 认证接口处理层
type AuthHandler struct {
	userRepo *repository.UserRepository
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		userRepo: repository.NewUserRepository(),
	}
}

// RefreshTokenRequest 刷新Token请求结构
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshToken 刷新Token
// @Summary 刷新Token
// @Description 使用refresh_token刷新访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param req body RefreshTokenRequest true "刷新令牌"
// @Success 200 {object} response.Response
// @Router /api/auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, response.CodeBadRequest, "refresh_token不能为空")
		return
	}

	jwtUtil := jwt.NewJWT()
	claims, err := jwtUtil.ParseToken(req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, "无效的刷新令牌")
		return
	}

	// 生成新的Token
	newToken, err := jwtUtil.GenerateToken(claims.UserID, claims.Username, claims.Email)
	if err != nil {
		response.ServerError(c, "Token生成失败")
		return
	}

	response.OkWithData(c, gin.H{"token": newToken})
}

// ChangePasswordRequest 修改密码请求结构
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=20"`
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 用户修改自己的密码
// @Tags 认证
// @Accept json
// @Produce json
// @Param req body ChangePasswordRequest true "密码信息"
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if req.OldPassword == req.NewPassword {
		response.FailWithMsg(c, response.CodeBadRequest, "新密码不能与旧密码相同")
		return
	}

	// 获取用户原始信息（包含密码）
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	// 验证旧密码
	if !password.CheckPassword(req.OldPassword, user.Password) {
		response.FailWithMsg(c, response.CodeUserPasswordError, "旧密码错误")
		return
	}

	// 修改密码
	newHashedPassword, err := password.HashPassword(req.NewPassword)
	if err != nil {
		response.ServerError(c, "密码加密失败")
		return
	}

	user.Password = newHashedPassword
	if err := h.userRepo.Update(user); err != nil {
		response.ServerError(c, "密码修改失败")
		return
	}

	response.Ok(c)
}

// Logout 退出登录
// @Summary 退出登录
// @Description 用户退出登录（客户端清除Token即可，服务端无需处理）
// @Tags 认证
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	response.Ok(c)
}
