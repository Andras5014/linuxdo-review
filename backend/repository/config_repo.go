package repository

import (
	"linuxdo-review/models"

	"gorm.io/gorm"
)

// ConfigRepository 配置仓库
type ConfigRepository struct {
	db *gorm.DB
}

// NewConfigRepository 创建配置仓库
func NewConfigRepository(db *gorm.DB) *ConfigRepository {
	return &ConfigRepository{db: db}
}

// Get 获取配置值
func (r *ConfigRepository) Get(key string) (string, error) {
	var config models.SystemConfig
	err := r.db.Where("key = ?", key).First(&config).Error
	if err != nil {
		return "", err
	}
	return config.Value, nil
}

// Set 设置配置值
func (r *ConfigRepository) Set(key, value, description string) error {
	config := models.SystemConfig{
		Key:         key,
		Value:       value,
		Description: description,
	}
	return r.db.Save(&config).Error
}

// GetAll 获取所有配置
func (r *ConfigRepository) GetAll() ([]*models.SystemConfig, error) {
	var configs []*models.SystemConfig
	err := r.db.Find(&configs).Error
	return configs, err
}

// Delete 删除配置
func (r *ConfigRepository) Delete(key string) error {
	return r.db.Where("key = ?", key).Delete(&models.SystemConfig{}).Error
}

// InitDefaults 初始化默认配置
func (r *ConfigRepository) InitDefaults() error {
	defaults := []models.SystemConfig{
		{Key: "min_votes", Value: "10", Description: "进入二级审核的最小票数"},
		{Key: "approval_rate", Value: "70", Description: "赞率阈值(百分比)"},
	}

	for _, config := range defaults {
		// 如果配置不存在则创建
		var existing models.SystemConfig
		result := r.db.Where("key = ?", config.Key).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			if err := r.db.Create(&config).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
