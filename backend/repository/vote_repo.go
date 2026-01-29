package repository

import (
	"linuxdo-review/models"

	"gorm.io/gorm"
)

// VoteRepository 投票仓库
type VoteRepository struct {
	db *gorm.DB
}

// NewVoteRepository 创建投票仓库
func NewVoteRepository(db *gorm.DB) *VoteRepository {
	return &VoteRepository{db: db}
}

// Create 创建投票
func (r *VoteRepository) Create(vote *models.Vote) error {
	return r.db.Create(vote).Error
}

// FindByPostAndUser 根据帖子ID和用户ID查找投票
func (r *VoteRepository) FindByPostAndUser(postID, userID uint) (*models.Vote, error) {
	var vote models.Vote
	err := r.db.Where("post_id = ? AND user_id = ?", postID, userID).First(&vote).Error
	if err != nil {
		return nil, err
	}
	return &vote, nil
}

// Update 更新投票
func (r *VoteRepository) Update(vote *models.Vote) error {
	return r.db.Save(vote).Error
}

// Delete 删除投票
func (r *VoteRepository) Delete(vote *models.Vote) error {
	return r.db.Delete(vote).Error
}

// CountByPost 统计帖子的投票数
func (r *VoteRepository) CountByPost(postID uint) (upVotes int64, downVotes int64, err error) {
	err = r.db.Model(&models.Vote{}).
		Where("post_id = ? AND vote_type = ?", postID, models.VoteUp).
		Count(&upVotes).Error
	if err != nil {
		return 0, 0, err
	}

	err = r.db.Model(&models.Vote{}).
		Where("post_id = ? AND vote_type = ?", postID, models.VoteDown).
		Count(&downVotes).Error
	if err != nil {
		return 0, 0, err
	}

	return upVotes, downVotes, nil
}

// GetUserVotesForPosts 获取用户对多个帖子的投票情况
func (r *VoteRepository) GetUserVotesForPosts(userID uint, postIDs []uint) (map[uint]models.VoteType, error) {
	var votes []models.Vote
	result := make(map[uint]models.VoteType)

	if len(postIDs) == 0 {
		return result, nil
	}

	err := r.db.Where("user_id = ? AND post_id IN ?", userID, postIDs).Find(&votes).Error
	if err != nil {
		return nil, err
	}

	for _, v := range votes {
		result[v.PostID] = v.VoteType
	}

	return result, nil
}

// CountAll 统计所有投票数量
func (r *VoteRepository) CountAll() (int64, error) {
	var count int64
	err := r.db.Model(&models.Vote{}).Count(&count).Error
	return count, err
}
