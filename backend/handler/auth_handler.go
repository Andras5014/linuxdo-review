package handler

import (
	"linuxdo-review/config"
	"linuxdo-review/dto"
	"linuxdo-review/middleware"
	"linuxdo-review/pkg/response"
	"linuxdo-review/service"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *service.AuthService
	cfg         *config.Config
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService *service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cfg:         cfg,
	}
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user, err := h.authService.Register(&req)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, dto.ToUserResponse(user))
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	result, err := h.authService.Login(&req)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, result)
}

// Me 获取当前用户信息
func (h *AuthHandler) Me(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		response.Error(c, "获取用户信息失败")
		return
	}

	response.Success(c, dto.ToUserResponse(user))
}

// OAuthLinuxDo Linux.do OAuth跳转
func (h *AuthHandler) OAuthLinuxDo(c *gin.Context) {
	authURL, state, err := h.authService.GetOAuthURL()
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	// 返回OAuth授权URL和state
	// 前端可以选择直接跳转或者存储state后跳转
	response.Success(c, &dto.OAuthURLResponse{
		URL:   authURL,
		State: state,
	})
}

// OAuthLinuxDoRedirect Linux.do OAuth直接重定向
func (h *AuthHandler) OAuthLinuxDoRedirect(c *gin.Context) {
	authURL, _, err := h.authService.GetOAuthURL()
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	// 直接重定向到OAuth授权页面
	c.Redirect(302, authURL)
}

// OAuthLinuxDoCallback Linux.do OAuth回调
func (h *AuthHandler) OAuthLinuxDoCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")

	// 检查是否有错误
	if errorParam != "" {
		errorDesc := c.Query("error_description")
		response.Error(c, "OAuth授权失败: "+errorDesc)
		return
	}

	if code == "" {
		response.BadRequest(c, "缺少授权码")
		return
	}

	if state == "" {
		response.BadRequest(c, "缺少state参数")
		return
	}

	// 处理OAuth回调
	result, err := h.authService.HandleOAuthCallback(code, state)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	// 如果配置了前端回调URL,重定向到前端页面
	// 否则直接返回JSON
	frontendURL := c.Query("redirect_uri")
	if frontendURL != "" {
		// 重定向到前端,带上token
		c.Redirect(302, frontendURL+"?token="+result.Token)
		return
	}

	response.Success(c, result)
}

// GetSystemStatus 获取系统状态（是否已初始化）
func (h *AuthHandler) GetSystemStatus(c *gin.Context) {
	initialized, err := h.authService.IsSystemInitialized()
	if err != nil {
		response.Error(c, "获取系统状态失败")
		return
	}

	response.Success(c, &dto.SystemStatusResponse{
		Initialized: initialized,
	})
}

// SetupAdmin 初始化管理员账号
func (h *AuthHandler) SetupAdmin(c *gin.Context) {
	var req dto.SetupAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user, err := h.authService.SetupAdmin(&req)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, dto.ToUserResponse(user))
}