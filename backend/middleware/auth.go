package middleware

import (
	"linuxdo-review/models"
	"linuxdo-review/pkg/response"

	"github.com/gin-gonic/gin"
)

// RequireCertified 要求认证用户(trust_level >= 3 或管理员)用于二级审核
func RequireCertified() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := GetUserRole(c)
		trustLevel := GetTrustLevel(c)

		// 管理员或信任等级>=3可以进行二级审核
		if role == int(models.RoleAdmin) || trustLevel >= 3 {
			c.Next()
			return
		}

		response.Forbidden(c, "需要 Linux.do 信任等级3及以上或管理员权限")
		c.Abort()
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

// RequireLinuxDoBinding 要求绑定LinuxDo账号
func RequireLinuxDoBinding() gin.HandlerFunc {
	return func(c *gin.Context) {
		linuxDoID := GetLinuxDoID(c)
		if linuxDoID == "" {
			response.Forbidden(c, "请先绑定 Linux.do 账号")
			c.Abort()
			return
		}
		c.Next()
	}
}
