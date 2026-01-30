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

// OAuthLinuxDoCallback Linux.do OAuth回调（支持登录和绑定两种模式）
func (h *AuthHandler) OAuthLinuxDoCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")

	// 获取前端基础地址
	frontendBaseURL := h.cfg.Server.FrontendURL
	if frontendBaseURL == "" {
		scheme := "http"
		if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
			scheme = "https"
		}
		frontendBaseURL = scheme + "://" + c.Request.Host
	}

	// 默认登录回调地址
	frontendCallbackURL := frontendBaseURL + "/oauth/callback"

	// 检查是否有错误
	if errorParam != "" {
		errorDesc := c.Query("error_description")
		c.Redirect(302, frontendCallbackURL+"?error="+errorParam+"&error_description="+errorDesc)
		return
	}

	if code == "" {
		c.Redirect(302, frontendCallbackURL+"?error=missing_code&error_description=缺少授权码")
		return
	}

	if state == "" {
		c.Redirect(302, frontendCallbackURL+"?error=missing_state&error_description=缺少state参数")
		return
	}

	// 处理OAuth回调
	result, err := h.authService.HandleOAuthCallback(code, state)
	if err != nil {
		// 如果是绑定模式失败，重定向到绑定回调页面
		c.Redirect(302, frontendCallbackURL+"?error=oauth_failed&error_description="+err.Error())
		return
	}

	// 根据模式重定向到不同的前端页面
	if result.IsBindMode {
		// 绑定模式：直接重定向到个人信息页面，带上绑定成功标识
		profileURL := frontendBaseURL + "/profile?bindSuccess=true"
		c.Redirect(302, profileURL)
	} else {
		// 登录模式：重定向到登录回调页面
		c.Redirect(302, frontendCallbackURL+"?token="+result.LoginResponse.Token)
	}
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

// GetBindLinuxDoURL 获取绑定LinuxDO的OAuth URL
func (h *AuthHandler) GetBindLinuxDoURL(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	// 获取用户的 token（从 Authorization header 中提取）
	authHeader := c.GetHeader("Authorization")
	token := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	}

	authURL, state, err := h.authService.GetBindLinuxDoURL(userID, token)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, &dto.OAuthURLResponse{
		URL:   authURL,
		State: state,
	})
}

// UnbindLinuxDo 解绑LinuxDO
func (h *AuthHandler) UnbindLinuxDo(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	user, err := h.authService.UnbindLinuxDo(userID)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, dto.ToUserResponse(user))
}

// BindEmail 绑定邮箱（LinuxDO用户专用）
func (h *AuthHandler) BindEmail(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	var req dto.BindEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user, err := h.authService.BindEmail(userID, &req)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, dto.ToUserResponse(user))
}

// UpdateProfile 更新用户资料
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user, err := h.authService.UpdateProfile(userID, &req)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, dto.ToUserResponse(user))
}

// SendEmailCode 发送邮箱验证码
func (h *AuthHandler) SendEmailCode(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	var req dto.SendEmailCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	err := h.authService.SendEmailVerificationCode(userID, req.Email)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.SuccessMessage(c, "验证码已发送")
}

// ChangeEmail 修改邮箱
func (h *AuthHandler) ChangeEmail(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	var req dto.ChangeEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user, err := h.authService.ChangeEmail(userID, &req)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, dto.ToUserResponse(user))
}

// UpdateAvatar 更新头像
func (h *AuthHandler) UpdateAvatar(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	var req dto.UpdateAvatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user, err := h.authService.UpdateAvatar(userID, &req)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, dto.ToUserResponse(user))
}
