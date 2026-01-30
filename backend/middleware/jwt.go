package middleware

import (
	"strings"

	"linuxdo-review/config"
	"linuxdo-review/pkg/jwt"
	"linuxdo-review/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	// ContextUserIDKey 上下文中用户ID的key
	ContextUserIDKey = "user_id"
	// ContextUserEmailKey 上下文中用户邮箱的key
	ContextUserEmailKey = "user_email"
	// ContextUsernameKey 上下文中用户名的key
	ContextUsernameKey = "username"
	// ContextUserRoleKey 上下文中用户角色的key
	ContextUserRoleKey = "user_role"
	// ContextTrustLevelKey 上下文中信任等级的key
	ContextTrustLevelKey = "trust_level"
	// ContextLinuxDoIDKey 上下文中LinuxDo ID的key
	ContextLinuxDoIDKey = "linuxdo_id"
)

// JWTAuth JWT认证中间件
func JWTAuth(cfg *config.Config) gin.HandlerFunc {
	jwtManager := jwt.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpireHours)

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		// Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Unauthorized(c, "无效的Authorization格式")
			c.Abort()
			return
		}

		claims, err := jwtManager.ParseToken(parts[1])
		if err != nil {
			response.Unauthorized(c, err.Error())
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextUserEmailKey, claims.Email)
		c.Set(ContextUsernameKey, claims.Username)
		c.Set(ContextUserRoleKey, claims.Role)
		c.Set(ContextTrustLevelKey, claims.TrustLevel)
		c.Set(ContextLinuxDoIDKey, claims.LinuxDoID)

		c.Next()
	}
}

// OptionalJWTAuth 可选的JWT认证中间件(不强制要求登录)
func OptionalJWTAuth(cfg *config.Config) gin.HandlerFunc {
	jwtManager := jwt.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpireHours)

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			claims, err := jwtManager.ParseToken(parts[1])
			if err == nil {
				c.Set(ContextUserIDKey, claims.UserID)
				c.Set(ContextUserEmailKey, claims.Email)
				c.Set(ContextUsernameKey, claims.Username)
				c.Set(ContextUserRoleKey, claims.Role)
				c.Set(ContextTrustLevelKey, claims.TrustLevel)
				c.Set(ContextLinuxDoIDKey, claims.LinuxDoID)
			}
		}

		c.Next()
	}
}

// GetUserID 从上下文获取用户ID
func GetUserID(c *gin.Context) uint {
	userID, exists := c.Get(ContextUserIDKey)
	if !exists {
		return 0
	}
	return userID.(uint)
}

// GetUserRole 从上下文获取用户角色
func GetUserRole(c *gin.Context) int {
	role, exists := c.Get(ContextUserRoleKey)
	if !exists {
		return 0
	}
	return role.(int)
}

// GetTrustLevel 从上下文获取信任等级
func GetTrustLevel(c *gin.Context) int {
	level, exists := c.Get(ContextTrustLevelKey)
	if !exists {
		return 0
	}
	return level.(int)
}

// GetLinuxDoID 从上下文获取LinuxDo ID
func GetLinuxDoID(c *gin.Context) string {
	id, exists := c.Get(ContextLinuxDoIDKey)
	if !exists {
		return ""
	}
	return id.(string)
}
