package middleware

import (
	"linuxdo-review/models"
	"linuxdo-review/pkg/response"

	"github.com/gin-gonic/gin"
)

// RequireCertified 要求认证用户(role >= 1)
func RequireCertified() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := GetUserRole(c)
		if role < int(models.RoleCertified) {
			response.Forbidden(c, "需要认证用户权限")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireAdmin 要求管理员(role == 2)
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := GetUserRole(c)
		if role != int(models.RoleAdmin) {
			response.Forbidden(c, "需要管理员权限")
			c.Abort()
			return
		}
		c.Next()
	}
}
