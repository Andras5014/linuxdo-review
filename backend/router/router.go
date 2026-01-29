package router

import (
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
			auth.GET("/oauth/linuxdo", authHandler.OAuthLinuxDo)               // 获取OAuth URL(API方式)
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
		}

		// 审核相关(认证用户专属)
		review := api.Group("/posts", middleware.JWTAuth(cfg), middleware.RequireCertified())
		{
			review.POST("/:id/approve", reviewHandler.Approve)
			review.POST("/:id/reject", reviewHandler.Reject)
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

	return r
}
