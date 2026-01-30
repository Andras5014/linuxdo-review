package main

import (
	"fmt"
	"log"
	"os"

	"linuxdo-review/config"
	"linuxdo-review/database"
	"linuxdo-review/handler"
	"linuxdo-review/repository"
	"linuxdo-review/router"
	"linuxdo-review/service"
)

func main() {
	// 配置文件路径：优先使用环境变量，其次当前目录，最后父目录
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		if _, err := os.Stat("config.yaml"); err == nil {
			configPath = "config.yaml"
		} else {
			configPath = "../config.yaml"
		}
	}

	// 加载配置
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库
	if err := database.Init(cfg.Database.Path); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	db := database.GetDB()

	// 初始化Repository层
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	voteRepo := repository.NewVoteRepository(db)
	configRepo := repository.NewConfigRepository(db)

	// 初始化默认配置
	if err := configRepo.InitDefaults(); err != nil {
		log.Printf("初始化默认配置失败: %v", err)
	}

	// 初始化Service层
	emailService := service.NewEmailService(cfg)
	authService := service.NewAuthService(userRepo, cfg, emailService)
	postService := service.NewPostService(postRepo, voteRepo, configRepo, userRepo, cfg)
	reviewService := service.NewReviewService(postRepo, userRepo, configRepo, emailService)
	adminService := service.NewAdminService(userRepo, postRepo, voteRepo, configRepo)

	// 初始化Handler层
	authHandler := handler.NewAuthHandler(authService, cfg)
	postHandler := handler.NewPostHandler(postService, reviewService)
	reviewHandler := handler.NewReviewHandler(reviewService, postService)
	adminHandler := handler.NewAdminHandler(adminService)

	// 设置路由
	r := router.SetupRouter(cfg, authHandler, postHandler, reviewHandler, adminHandler)

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("服务器启动在 http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
