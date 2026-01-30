package repository

import (
	"linuxdo-review/models"

	"gorm.io/gorm"
)

// UserRepository 用户仓库
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 创建用户
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// FindByID 根据ID查找用户
func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail 根据邮箱查找用户
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByLinuxDoID 根据Linux.do ID查找用户
func (r *UserRepository) FindByLinuxDoID(linuxDoID string) (*models.User, error) {
	var user models.User
	err := r.db.Where("linux_do_id = ?", linuxDoID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// UpdateRole 更新用户角色
func (r *UserRepository) UpdateRole(id uint, role models.UserRole) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("role", role).Error
}

// List 获取用户列表(分页)
func (r *UserRepository) List(offset, limit int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	if err := r.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// CountByRole 统计指定角色的用户数量
func (r *UserRepository) CountByRole(role models.UserRole) (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("role = ?", role).Count(&count).Error
	return count, err
}

// CountAll 统计所有用户数量
func (r *UserRepository) CountAll() (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Count(&count).Error
	return count, err
}

// CountTodayNew 统计今日新增用户数量
func (r *UserRepository) CountTodayNew() (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("DATE(created_at) = DATE('now')").Count(&count).Error
	return count, err
}

// HasAdmin 检查是否存在管理员账号
func (r *UserRepository) HasAdmin() (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("role = ?", models.RoleAdmin).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
