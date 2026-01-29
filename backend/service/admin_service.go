package service

import (
	"errors"

	"linuxdo-review/dto"
	"linuxdo-review/models"
	"linuxdo-review/repository"
)

// AdminService 管理后台服务
type AdminService struct {
	userRepo   *repository.UserRepository
	postRepo   *repository.PostRepository
	voteRepo   *repository.VoteRepository
	configRepo *repository.ConfigRepository
}

// NewAdminService 创建管理后台服务
func NewAdminService(
	userRepo *repository.UserRepository,
	postRepo *repository.PostRepository,
	voteRepo *repository.VoteRepository,
	configRepo *repository.ConfigRepository,
) *AdminService {
	return &AdminService{
		userRepo:   userRepo,
		postRepo:   postRepo,
		voteRepo:   voteRepo,
		configRepo: configRepo,
	}
}

// ListUsers 获取用户列表(分页)
func (s *AdminService) ListUsers(page, pageSize int) ([]*models.User, int64, error) {
	offset := (page - 1) * pageSize
	return s.userRepo.List(offset, pageSize)
}

// GetUserByID 根据ID获取用户
func (s *AdminService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

// UpdateUserRole 更新用户角色
func (s *AdminService) UpdateUserRole(id uint, role models.UserRole) error {
	// 验证用户是否存在
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 验证角色值是否有效
	if role < models.RoleNormal || role > models.RoleAdmin {
		return errors.New("无效的角色值")
	}

	// 不能修改自己的角色(防止管理员降级自己)
	// 这个检查可以在Handler层根据当前用户ID进行

	// 更新角色
	user.Role = role
	return s.userRepo.UpdateRole(id, role)
}

// GetConfigs 获取所有配置
func (s *AdminService) GetConfigs() ([]*models.SystemConfig, error) {
	return s.configRepo.GetAll()
}

// GetConfig 获取单个配置
func (s *AdminService) GetConfig(key string) (string, error) {
	return s.configRepo.Get(key)
}

// UpdateConfig 更新配置
func (s *AdminService) UpdateConfig(key, value, description string) error {
	if key == "" {
		return errors.New("配置键不能为空")
	}
	return s.configRepo.Set(key, value, description)
}

// BatchUpdateConfigs 批量更新配置
func (s *AdminService) BatchUpdateConfigs(configs []dto.UpdateConfigRequest) error {
	for _, cfg := range configs {
		if err := s.UpdateConfig(cfg.Key, cfg.Value, ""); err != nil {
			return err
		}
	}
	return nil
}

// GetStats 获取统计数据
func (s *AdminService) GetStats() (*dto.StatsResponse, error) {
	stats := &dto.StatsResponse{}

	// 用户统计
	totalUsers, err := s.userRepo.CountAll()
	if err != nil {
		return nil, err
	}
	stats.TotalUsers = totalUsers

	certifiedCount, _ := s.userRepo.CountByRole(models.RoleCertified)
	adminCount, _ := s.userRepo.CountByRole(models.RoleAdmin)
	stats.CertifiedUsers = certifiedCount + adminCount

	stats.TodayNewUsers, _ = s.userRepo.CountTodayNew()

	// 帖子统计
	stats.TotalPosts, _ = s.postRepo.CountAll()

	pendingCount, _ := s.postRepo.CountByStatus(models.StatusPending)
	stats.PendingPosts = pendingCount

	firstReviewCount, _ := s.postRepo.CountByStatus(models.StatusFirstReview)
	stats.FirstReviewPosts = firstReviewCount

	secondReviewCount, _ := s.postRepo.CountByStatus(models.StatusSecondReview)
	stats.SecondReviewPosts = secondReviewCount

	stats.ApprovedPosts, _ = s.postRepo.CountByStatus(models.StatusApproved)
	stats.RejectedPosts, _ = s.postRepo.CountByStatus(models.StatusRejected)
	stats.TodayNewPosts, _ = s.postRepo.CountTodayNew()
	stats.TodayApproved, _ = s.postRepo.CountTodayApproved()

	// 投票统计
	stats.TotalVotes, _ = s.voteRepo.CountAll()

	return stats, nil
}

// DeleteUser 删除用户(软删除或禁用)
func (s *AdminService) DeleteUser(id uint) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 不能删除管理员账户
	if user.IsAdmin() {
		return errors.New("不能删除管理员账户")
	}

	// 这里可以实现软删除或禁用逻辑
	// 目前先返回未实现错误
	return errors.New("删除用户功能暂未实现")
}
