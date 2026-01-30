package router

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"linuxdo-review/config"
	"linuxdo-review/handler"
	"linuxdo-review/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由

func SetupRouter(
	cfg *config.Config,
	authHandler *handler.AuthHandler,
	postHandler *handler.PostHandler,
	reviewHandler *handler.ReviewHandler,
	adminHandler *handler.AdminHandler,
) *gin.Engine {
	// 设置运行模式
	gin.SetMode(cfg.Server.Mode)

	r := gin.Default()

	// 全局中间件
	r.Use(middleware.CORS())

	// API路由组
	api := r.Group("/api")
	{
		// 系统状态（公开接口）
		api.GET("/system/status", authHandler.GetSystemStatus)
		api.POST("/system/setup", authHandler.SetupAdmin)

		// 认证相关(无需登录)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/oauth/linuxdo", authHandler.OAuthLinuxDo)                  // 获取OAuth URL(API方式)
			auth.GET("/oauth/linuxdo/redirect", authHandler.OAuthLinuxDoRedirect) // 直接重定向到OAuth页面
			auth.GET("/oauth/linuxdo/callback", authHandler.OAuthLinuxDoCallback)

			// 需要登录
			auth.GET("/me", middleware.JWTAuth(cfg), authHandler.Me)
		}

		// 帖子相关
		posts := api.Group("/posts")
		{
			// 公开接口(可选登录,用于显示投票状态)
			posts.GET("", middleware.OptionalJWTAuth(cfg), postHandler.List)
			posts.GET("/:id", middleware.OptionalJWTAuth(cfg), postHandler.Get)

			// 需要登录
			posts.POST("", middleware.JWTAuth(cfg), postHandler.Create)
			posts.POST("/:id/vote", middleware.JWTAuth(cfg), postHandler.Vote)

			// 认证用户专属(二级审核列表)
			posts.GET("/review", middleware.JWTAuth(cfg), middleware.RequireCertified(), postHandler.ListForReview)
		}

		// 用户相关(需要登录)
		user := api.Group("/user", middleware.JWTAuth(cfg))
		{
			user.GET("/posts", postHandler.MyPosts)
			user.GET("/profile", authHandler.Me)
			user.PUT("/profile", authHandler.UpdateProfile)
			user.GET("/bindlinuxdo", authHandler.GetBindLinuxDoURL)
			user.POST("/unbindlinuxdo", authHandler.UnbindLinuxDo)
			user.POST("/bindmail", authHandler.BindEmail)       // LinuxDO用户绑定邮箱
			user.POST("/email/code", authHandler.SendEmailCode) // 发送邮箱验证码
			user.POST("/email/change", authHandler.ChangeEmail) // 修改邮箱
			user.PUT("/avatar", authHandler.UpdateAvatar)       // 更新头像
		}

		// 审核相关(认证用户专属)
		review := api.Group("/review", middleware.JWTAuth(cfg), middleware.RequireCertified())
		{
			review.GET("/next", reviewHandler.GetNext)         // 获取下一个待审核的帖子
			review.POST("/:id/skip", reviewHandler.Skip)       // 跳过当前帖子
			review.POST("/:id/approve", reviewHandler.Approve) // 通过审核
			review.POST("/:id/reject", reviewHandler.Reject)   // 拒绝申请
		}

		// 管理后台(管理员专属)
		admin := api.Group("/admin", middleware.JWTAuth(cfg), middleware.RequireAdmin())
		{
			// 用户管理
			admin.GET("/users", adminHandler.ListUsers)
			admin.GET("/users/:id", adminHandler.GetUser)
			admin.PUT("/users/:id", adminHandler.UpdateUserRole)

			// 配置管理
			admin.GET("/configs", adminHandler.GetConfigs)
			admin.PUT("/configs", adminHandler.UpdateConfig)
			admin.PUT("/configs/batch", adminHandler.BatchUpdateConfigs)

			// 数据统计
			admin.GET("/stats", adminHandler.GetStats)
		}
	}

	// 静态文件服务 (生产环境)
	staticPath := os.Getenv("STATIC_PATH")
	if staticPath == "" {
		staticPath = "./static"
	}

	// 检查静态文件目录是否存在
	if _, err := os.Stat(staticPath); err == nil {
		// 服务静态文件
		r.Static("/assets", filepath.Join(staticPath, "assets"))

		// 处理 index.html 和 SPA 路由
		r.NoRoute(func(c *gin.Context) {
			// 如果是 API 请求，返回 404
			if strings.HasPrefix(c.Request.URL.Path, "/api") {
				c.JSON(http.StatusNotFound, gin.H{"error": "API not found"})
				return
			}

			// 其他请求返回 index.html (SPA 路由)
			c.File(filepath.Join(staticPath, "index.html"))
		})
	}

	return r
}
