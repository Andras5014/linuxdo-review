package service

import (
	"errors"
	"strconv"

	"linuxdo-review/models"
	"linuxdo-review/repository"

	"gorm.io/gorm"
)

// ReviewService 审核服务
type ReviewService struct {
	postRepo     *repository.PostRepository
	userRepo     *repository.UserRepository
	configRepo   *repository.ConfigRepository
	emailService *EmailService
}

// NewReviewService 创建审核服务
func NewReviewService(
	postRepo *repository.PostRepository,
	userRepo *repository.UserRepository,
	configRepo *repository.ConfigRepository,
	emailService *EmailService,
) *ReviewService {
	return &ReviewService{
		postRepo:     postRepo,
		userRepo:     userRepo,
		configRepo:   configRepo,
		emailService: emailService,
	}
}

// CheckAndPromote 检查帖子是否满足进入二级审核的条件
func (s *ReviewService) CheckAndPromote(postID uint) error {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return err
	}

	if post.Status != models.StatusFirstReview {
		return nil // 不在一级审核阶段,无需处理
	}

	// 获取配置
	minVotes, approvalRate := s.getReviewConfig()

	// 检查票数是否达到阈值
	totalVotes := post.TotalVotes()
	if totalVotes < minVotes {
		return nil // 票数不足
	}

	// 计算赞率
	rate := post.ApprovalRate()

	if rate >= float64(approvalRate) {
		// 进入二级审核
		return s.postRepo.UpdateStatus(postID, models.StatusSecondReview)
	} else {
		// 拒绝
		return s.postRepo.UpdateStatus(postID, models.StatusRejected)
	}
}

// Approve 通过审核(认证用户提交邀请码)
func (s *ReviewService) Approve(postID uint, reviewerID uint, inviteCode string) error {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("帖子不存在")
		}
		return err
	}

	if post.Status != models.StatusSecondReview {
		return errors.New("帖子不在二级审核阶段")
	}

	// TODO: 可以对邀请码进行加密存储
	if err := s.postRepo.Approve(postID, reviewerID, inviteCode); err != nil {
		return err
	}

	// 发送邮件通知申请者
	if s.emailService != nil && post.User != nil && post.User.Email != "" {
		_ = s.emailService.SendInviteCode(post.User.Email, post.User.Username, inviteCode)
	}

	return nil
}

// ApproveWithNotification 通过审核并发送邮件通知(完整流程)
func (s *ReviewService) ApproveWithNotification(postID uint, reviewerID uint, inviteCode string) error {
	// 先获取完整的帖子信息(包含用户)
	post, err := s.postRepo.FindByIDWithReviewer(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("帖子不存在")
		}
		return err
	}

	if post.Status != models.StatusSecondReview {
		return errors.New("帖子不在二级审核阶段")
	}

	// 执行审核通过
	if err := s.postRepo.Approve(postID, reviewerID, inviteCode); err != nil {
		return err
	}

	// 发送邮件通知申请者
	if s.emailService != nil && post.User != nil && post.User.Email != "" {
		_ = s.emailService.SendInviteCode(post.User.Email, post.User.Username, inviteCode)
	}

	return nil
}

// Reject 拒绝申请
func (s *ReviewService) Reject(postID uint, reason string) error {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("帖子不存在")
		}
		return err
	}

	if post.Status != models.StatusFirstReview && post.Status != models.StatusSecondReview {
		return errors.New("帖子状态不允许拒绝")
	}

	if err := s.postRepo.Reject(postID, reason); err != nil {
		return err
	}

	// 发送拒绝通知邮件
	if s.emailService != nil && post.User != nil && post.User.Email != "" {
		_ = s.emailService.SendRejectionNotification(post.User.Email, post.User.Username, post.Title, reason)
	}

	return nil
}

// RejectWithNotification 拒绝申请并发送邮件通知(完整流程)
func (s *ReviewService) RejectWithNotification(postID uint, reason string) error {
	// 先获取完整的帖子信息(包含用户)
	post, err := s.postRepo.FindByIDWithReviewer(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("帖子不存在")
		}
		return err
	}

	if post.Status != models.StatusFirstReview && post.Status != models.StatusSecondReview {
		return errors.New("帖子状态不允许拒绝")
	}

	if err := s.postRepo.Reject(postID, reason); err != nil {
		return err
	}

	// 发送拒绝通知邮件
	if s.emailService != nil && post.User != nil && post.User.Email != "" {
		_ = s.emailService.SendRejectionNotification(post.User.Email, post.User.Username, post.Title, reason)
	}

	return nil
}

// GetUserByID 获取用户信息
func (s *ReviewService) GetUserByID(userID uint) (*models.User, error) {
	if s.userRepo == nil {
		return nil, errors.New("用户仓库未初始化")
	}
	return s.userRepo.FindByID(userID)
}

// GetPostByID 获取帖子(供审核用)
func (s *ReviewService) GetPostByID(postID uint) (*models.Post, error) {
	return s.postRepo.FindByIDWithReviewer(postID)
}

// GetNextForReview 获取下一个待审核的帖子并锁定
func (s *ReviewService) GetNextForReview(userID uint, skipIDs []uint) (*models.Post, error) {
	post, err := s.postRepo.GetNextForReview(userID, skipIDs)
	if err != nil {
		return nil, err
	}

	// 锁定帖子
	if err := s.postRepo.LockPost(post.ID, userID); err != nil {
		return nil, errors.New("帖子已被其他审核员锁定")
	}

	return post, nil
}

// SkipPost 跳过当前帖子（解锁）
func (s *ReviewService) SkipPost(postID uint, userID uint) error {
	return s.postRepo.UnlockPost(postID, userID)
}

// CheckLockAndApprove 检查锁定状态并通过审核
func (s *ReviewService) CheckLockAndApprove(postID uint, reviewerID uint, inviteCode string) error {
	// 检查帖子是否被其他用户锁定
	locked, err := s.postRepo.IsPostLocked(postID, reviewerID)
	if err != nil {
		return err
	}
	if locked {
		return errors.New("帖子已被其他审核员锁定，请刷新后重试")
	}

	return s.ApproveWithNotification(postID, reviewerID, inviteCode)
}

// CheckLockAndReject 检查锁定状态并拒绝
func (s *ReviewService) CheckLockAndReject(postID uint, userID uint, reason string) error {
	// 检查帖子是否被其他用户锁定
	locked, err := s.postRepo.IsPostLocked(postID, userID)
	if err != nil {
		return err
	}
	if locked {
		return errors.New("帖子已被其他审核员锁定，请刷新后重试")
	}

	return s.RejectWithNotification(postID, reason)
}

// GetReviewCount 获取待二级审核的数量
func (s *ReviewService) GetReviewCount() (int64, error) {
	return s.postRepo.CountForSecondReview()
}

// getReviewConfig 获取审核配置
func (s *ReviewService) getReviewConfig() (minVotes int, approvalRate int) {
	// 默认值
	minVotes = 10
	approvalRate = 70

	// 从数据库获取配置
	if val, err := s.configRepo.Get("min_votes"); err == nil {
		if v, err := strconv.Atoi(val); err == nil {
			minVotes = v
		}
	}

	if val, err := s.configRepo.Get("approval_rate"); err == nil {
		if v, err := strconv.Atoi(val); err == nil {
			approvalRate = v
		}
	}

	return
}
