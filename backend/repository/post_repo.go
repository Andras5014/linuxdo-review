package repository

import (
	"time"

	"linuxdo-review/models"

	"gorm.io/gorm"
)

// 锁定超时时间（分钟）
const LockTimeout = 5

// PostRepository 帖子仓库
type PostRepository struct {
	db *gorm.DB
}

// NewPostRepository 创建帖子仓库
func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

// Create 创建帖子
func (r *PostRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

// FindByID 根据ID查找帖子
func (r *PostRepository) FindByID(id uint) (*models.Post, error) {
	var post models.Post
	err := r.db.Preload("User").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// FindByIDWithReviewer 根据ID查找帖子(包含审核者)
func (r *PostRepository) FindByIDWithReviewer(id uint) (*models.Post, error) {
	var post models.Post
	err := r.db.Preload("User").Preload("Reviewer").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// Update 更新帖子
func (r *PostRepository) Update(post *models.Post) error {
	return r.db.Save(post).Error
}

// ListByStatus 根据状态获取帖子列表(分页)
func (r *PostRepository) ListByStatus(status models.PostStatus, offset, limit int) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var total int64

	query := r.db.Model(&models.Post{}).Where("status = ?", status)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("User").Order("created_at DESC").Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// ListByUserID 根据用户ID获取帖子列表
func (r *PostRepository) ListByUserID(userID uint, offset, limit int) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var total int64

	query := r.db.Model(&models.Post{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("User").Order("created_at DESC").Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// ListForFirstReview 获取一级审核列表(社区投票中)
func (r *PostRepository) ListForFirstReview(offset, limit int) ([]*models.Post, int64, error) {
	return r.ListByStatus(models.StatusFirstReview, offset, limit)
}

// ListForSecondReview 获取二级审核列表(等待提交邀请码)
func (r *PostRepository) ListForSecondReview(offset, limit int) ([]*models.Post, int64, error) {
	return r.ListByStatus(models.StatusSecondReview, offset, limit)
}

// UpdateVotes 更新帖子票数
func (r *PostRepository) UpdateVotes(postID uint, upVotes, downVotes int) error {
	return r.db.Model(&models.Post{}).Where("id = ?", postID).
		Updates(map[string]interface{}{
			"up_votes":   upVotes,
			"down_votes": downVotes,
		}).Error
}

// UpdateStatus 更新帖子状态
func (r *PostRepository) UpdateStatus(postID uint, status models.PostStatus) error {
	return r.db.Model(&models.Post{}).Where("id = ?", postID).Update("status", status).Error
}

// Approve 通过审核(更新状态、审核者、邀请码、审核时间)
func (r *PostRepository) Approve(postID uint, reviewerID uint, inviteCode string) error {
	return r.db.Model(&models.Post{}).Where("id = ?", postID).
		Updates(map[string]interface{}{
			"status":      models.StatusApproved,
			"reviewer_id": reviewerID,
			"invite_code": inviteCode,
			"reviewed_at": gorm.Expr("datetime('now')"),
		}).Error
}

// Reject 拒绝申请(更新状态和拒绝原因)
func (r *PostRepository) Reject(postID uint, reason string) error {
	updates := map[string]interface{}{
		"status":      models.StatusRejected,
		"reviewed_at": gorm.Expr("datetime('now')"),
	}
	if reason != "" {
		updates["reject_reason"] = reason
	}
	return r.db.Model(&models.Post{}).Where("id = ?", postID).Updates(updates).Error
}

// PromoteToSecondReview 提升到二级审核
func (r *PostRepository) PromoteToSecondReview(postID uint) error {
	return r.db.Model(&models.Post{}).Where("id = ?", postID).
		Update("status", models.StatusSecondReview).Error
}

// CountByStatus 统计指定状态的帖子数量
func (r *PostRepository) CountByStatus(status models.PostStatus) (int64, error) {
	var count int64
	err := r.db.Model(&models.Post{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

// CountAll 统计所有帖子数量
func (r *PostRepository) CountAll() (int64, error) {
	var count int64
	err := r.db.Model(&models.Post{}).Count(&count).Error
	return count, err
}

// CountTodayNew 统计今日新增帖子数量
func (r *PostRepository) CountTodayNew() (int64, error) {
	var count int64
	err := r.db.Model(&models.Post{}).Where("DATE(created_at) = DATE('now')").Count(&count).Error
	return count, err
}

// CountTodayApproved 统计今日通过的帖子数量
func (r *PostRepository) CountTodayApproved() (int64, error) {
	var count int64
	err := r.db.Model(&models.Post{}).
		Where("status = ? AND DATE(updated_at) = DATE('now')", models.StatusApproved).
		Count(&count).Error
	return count, err
}

// ListWithFilter 根据可选状态过滤获取帖子列表(分页)
// status 为 nil 时返回所有状态的帖子
// status 为 -1 时返回投票中和待二级审核的帖子
func (r *PostRepository) ListWithFilter(status *int, offset, limit int) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var total int64

	query := r.db.Model(&models.Post{})

	if status != nil {
		if *status == -1 {
			// 特殊值 -1: 返回投票中和待二级审核的帖子(申请列表默认视图)
			query = query.Where("status IN ?", []models.PostStatus{models.StatusFirstReview, models.StatusSecondReview})
		} else {
			query = query.Where("status = ?", *status)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("User").Order("created_at DESC").Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// HasApprovedPost 检查用户是否有已通过的帖子
func (r *PostRepository) HasApprovedPost(userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Post{}).
		Where("user_id = ? AND status = ?", userID, models.StatusApproved).
		Count(&count).Error
	return count > 0, err
}

// HasPendingOrVotingPost 检查用户是否有投票中或待二级审核的帖子
func (r *PostRepository) HasPendingOrVotingPost(userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Post{}).
		Where("user_id = ? AND status IN ?", userID, []models.PostStatus{
			models.StatusPending,
			models.StatusFirstReview,
			models.StatusSecondReview,
		}).
		Count(&count).Error
	return count > 0, err
}

// GetNextForReview 获取下一个待二级审核的帖子（排除已锁定或过期锁定的）
func (r *PostRepository) GetNextForReview(userID uint, skipIDs []uint) (*models.Post, error) {
	var post models.Post
	now := time.Now()
	lockExpiry := now.Add(-LockTimeout * time.Minute)

	query := r.db.Preload("User").
		Where("status = ?", models.StatusSecondReview).
		Where("(locked_by IS NULL OR locked_by = ? OR locked_at < ?)", userID, lockExpiry)

	// 排除已跳过的帖子
	if len(skipIDs) > 0 {
		query = query.Where("id NOT IN ?", skipIDs)
	}

	err := query.Order("created_at ASC").First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// LockPost 锁定帖子（防止并发操作）
func (r *PostRepository) LockPost(postID uint, userID uint) error {
	now := time.Now()
	lockExpiry := now.Add(-LockTimeout * time.Minute)

	// 只有未锁定或锁定过期或被同一用户锁定的帖子才能被锁定
	result := r.db.Model(&models.Post{}).
		Where("id = ? AND status = ?", postID, models.StatusSecondReview).
		Where("(locked_by IS NULL OR locked_by = ? OR locked_at < ?)", userID, lockExpiry).
		Updates(map[string]interface{}{
			"locked_by": userID,
			"locked_at": now,
		})

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

// UnlockPost 解锁帖子
func (r *PostRepository) UnlockPost(postID uint, userID uint) error {
	return r.db.Model(&models.Post{}).
		Where("id = ? AND locked_by = ?", postID, userID).
		Updates(map[string]interface{}{
			"locked_by": nil,
			"locked_at": nil,
		}).Error
}

// IsPostLocked 检查帖子是否被锁定（且锁定未过期，且不是当前用户锁定的）
func (r *PostRepository) IsPostLocked(postID uint, userID uint) (bool, error) {
	var count int64
	now := time.Now()
	lockExpiry := now.Add(-LockTimeout * time.Minute)

	err := r.db.Model(&models.Post{}).
		Where("id = ?", postID).
		Where("locked_by IS NOT NULL AND locked_by != ? AND locked_at > ?", userID, lockExpiry).
		Count(&count).Error

	return count > 0, err
}

// CountForSecondReview 统计待二级审核的帖子数量
func (r *PostRepository) CountForSecondReview() (int64, error) {
	var count int64
	err := r.db.Model(&models.Post{}).
		Where("status = ?", models.StatusSecondReview).
		Count(&count).Error
	return count, err
}
