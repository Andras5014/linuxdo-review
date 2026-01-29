package database

import (
	"os"
	"path/filepath"

	"linuxdo-review/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Init 初始化数据库连接
func Init(dbPath string) error {
	// 确保数据库目录存在
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	// 自动迁移数据库表结构
	if err = autoMigrate(); err != nil {
		return err
	}

	return nil
}

// autoMigrate 自动迁移数据库表
func autoMigrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Vote{},
		&models.SystemConfig{},
	)
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
