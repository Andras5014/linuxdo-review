package service

import (
	"errors"
	"strconv"

	"linuxdo-review/config"
	"linuxdo-review/dto"
	"linuxdo-review/models"
	"linuxdo-review/repository"

	"gorm.io/gorm"
)

// PostService 帖子服务
type PostService struct {
	postRepo   *repository.PostRepository
	voteRepo   *repository.VoteRepository
	configRepo *repository.ConfigRepository
	userRepo   *repository.UserRepository
	cfg        *config.Config
}

// NewPostService 创建帖子服务
func NewPostService(
	postRepo *repository.PostRepository,
	voteRepo *repository.VoteRepository,
	configRepo *repository.ConfigRepository,
	userRepo *repository.UserRepository,
	cfg *config.Config,
) *PostService {
	return &PostService{
		postRepo:   postRepo,
		voteRepo:   voteRepo,
		configRepo: configRepo,
		userRepo:   userRepo,
		cfg:        cfg,
	}
}

// Create 创建帖子
func (s *PostService) Create(userID uint, req *dto.CreatePostRequest) (*models.Post, error) {
	// 检查用户是否已有已通过的帖子
	hasApproved, err := s.postRepo.HasApprovedPost(userID)
	if err != nil {
		return nil, errors.New("检查用户状态失败")
	}
	if hasApproved {
		return nil, errors.New("您已有已通过的申请，不能再次发起申请")
	}

	// 检查用户是否有正在投票中或待审核的帖子
	hasPending, err := s.postRepo.HasPendingOrVotingPost(userID)
	if err != nil {
		return nil, errors.New("检查用户状态失败")
	}
	if hasPending {
		return nil, errors.New("您已有进行中的申请，请等待审核完成后再发起新申请")
	}

	post := &models.Post{
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
		Status:  models.StatusFirstReview, // 默认进入一级审核(社区投票)
	}

	if err := s.postRepo.Create(post); err != nil {
		return nil, errors.New("创建帖子失败")
	}

	return post, nil
}

// GetByID 根据ID获取帖子
func (s *PostService) GetByID(id uint) (*models.Post, error) {
	return s.postRepo.FindByIDWithReviewer(id)
}

// ListForFirstReview 获取一级审核列表(社区投票中)
func (s *PostService) ListForFirstReview(page, pageSize int) ([]*models.Post, int64, error) {
	offset := (page - 1) * pageSize
	return s.postRepo.ListForFirstReview(offset, pageSize)
}

// ListWithFilter 根据可选状态过滤获取帖子列表
func (s *PostService) ListWithFilter(status *int, page, pageSize int) ([]*models.Post, int64, error) {
	offset := (page - 1) * pageSize
	return s.postRepo.ListWithFilter(status, offset, pageSize)
}

// ListForSecondReview 获取二级审核列表
func (s *PostService) ListForSecondReview(page, pageSize int) ([]*models.Post, int64, error) {
	offset := (page - 1) * pageSize
	return s.postRepo.ListForSecondReview(offset, pageSize)
}

// ListByUserID 获取用户的帖子列表
func (s *PostService) ListByUserID(userID uint, page, pageSize int) ([]*models.Post, int64, error) {
	offset := (page - 1) * pageSize
	return s.postRepo.ListByUserID(userID, offset, pageSize)
}

// Vote 投票
func (s *PostService) Vote(postID, userID uint, voteType models.VoteType) error {
	// 检查用户
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}
	// 管理员可以直接投票，其他用户需要绑定 Linux.do
	if !user.IsAdmin() && user.LinuxDoID == "" {
		return errors.New("请先绑定 Linux.do 账号后再投票")
	}

	// 检查帖子是否存在且处于一级审核状态
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("帖子不存在")
		}
		return err
	}

	if post.Status != models.StatusFirstReview {
		return errors.New("帖子不在投票阶段")
	}

	// 不能给自己的帖子投票
	if post.UserID == userID {
		return errors.New("不能给自己的帖子投票")
	}

	// 检查是否已投票
	existingVote, err := s.voteRepo.FindByPostAndUser(postID, userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if existingVote != nil {
		// 已投票,更新投票
		if existingVote.VoteType == voteType {
			// 取消投票
			if err := s.voteRepo.Delete(existingVote); err != nil {
				return errors.New("取消投票失败")
			}
		} else {
			// 修改投票
			existingVote.VoteType = voteType
			if err := s.voteRepo.Update(existingVote); err != nil {
				return errors.New("修改投票失败")
			}
		}
	} else {
		// 未投票,创建投票
		vote := &models.Vote{
			PostID:   postID,
			UserID:   userID,
			VoteType: voteType,
		}
		if err := s.voteRepo.Create(vote); err != nil {
			return errors.New("投票失败")
		}
	}

	// 更新帖子票数
	upVotes, downVotes, err := s.voteRepo.CountByPost(postID)
	if err != nil {
		return err
	}

	if err := s.postRepo.UpdateVotes(postID, int(upVotes), int(downVotes)); err != nil {
		return err
	}

	// 检查是否需要更新帖子状态
	return s.checkAndUpdatePostStatus(postID, int(upVotes), int(downVotes))
}

// checkAndUpdatePostStatus 检查并更新帖子状态
// 根据投票结果判断是否应该进入二级审核或被拒绝
func (s *PostService) checkAndUpdatePostStatus(postID uint, upVotes, downVotes int) error {
	// 获取审核配置
	minVotes, approvalRate := s.getReviewConfig()

	totalVotes := upVotes + downVotes
	if totalVotes < minVotes {
		return nil // 票数不足,不更新状态
	}

	// 计算赞率
	currentRate := float64(0)
	if totalVotes > 0 {
		currentRate = float64(upVotes) / float64(totalVotes) * 100
	}

	if currentRate >= float64(approvalRate) {
		// 赞率达标,进入二级审核
		return s.postRepo.PromoteToSecondReview(postID)
	} else {
		// 赞率不达标,拒绝
		return s.postRepo.Reject(postID, "社区投票未通过(赞率未达到阈值)")
	}
}

// getReviewConfig 获取审核配置
func (s *PostService) getReviewConfig() (minVotes int, approvalRate int) {
	// 优先从数据库获取配置
	if s.configRepo != nil {
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
	}

	// 使用配置文件中的默认值
	if minVotes == 0 && s.cfg != nil {
		minVotes = s.cfg.Review.MinVotes
	}
	if approvalRate == 0 && s.cfg != nil {
		approvalRate = s.cfg.Review.ApprovalRate
	}

	// 兜底默认值
	if minVotes == 0 {
		minVotes = 10
	}
	if approvalRate == 0 {
		approvalRate = 70
	}

	return minVotes, approvalRate
}

// GetUserVoteForPost 获取用户对帖子的投票
func (s *PostService) GetUserVoteForPost(postID, userID uint) (int, error) {
	vote, err := s.voteRepo.FindByPostAndUser(postID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return int(vote.VoteType), nil
}

// GetUserVotesForPosts 获取用户对多个帖子的投票
func (s *PostService) GetUserVotesForPosts(userID uint, postIDs []uint) (map[uint]int, error) {
	votes, err := s.voteRepo.GetUserVotesForPosts(userID, postIDs)
	if err != nil {
		return nil, err
	}

	result := make(map[uint]int)
	for postID, voteType := range votes {
		result[postID] = int(voteType)
	}
	return result, nil
}

// VoteWithResponse 投票并返回结果
func (s *PostService) VoteWithResponse(postID, userID uint, voteType models.VoteType) (*dto.VoteResponse, error) {
	if err := s.Vote(postID, userID, voteType); err != nil {
		return nil, err
	}

	// 重新获取帖子信息
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return nil, err
	}

	// 获取用户当前的投票状态
	currentVote, _ := s.voteRepo.FindByPostAndUser(postID, userID)

	resp := &dto.VoteResponse{
		PostID:    postID,
		UpVotes:   post.UpVotes,
		DownVotes: post.DownVotes,
	}

	if currentVote != nil {
		resp.VoteType = currentVote.VoteType
		resp.Message = "投票成功"
	} else {
		resp.VoteType = 0
		resp.Message = "已取消投票"
	}

	return resp, nil
}

// GetReviewConfig 获取审核配置(公开方法)
func (s *PostService) GetReviewConfig() (minVotes int, approvalRate int) {
	return s.getReviewConfig()
}

// HasUserApplied 检查用户是否已有进行中的申请
func (s *PostService) HasUserApplied(userID uint) (bool, error) {
	posts, _, err := s.postRepo.ListByUserID(userID, 0, 100)
	if err != nil {
		return false, err
	}

	for _, post := range posts {
		// 如果有待审核、一级审核中或二级审核的申请,则不能再发起新申请
		if post.Status == models.StatusPending ||
			post.Status == models.StatusFirstReview ||
			post.Status == models.StatusSecondReview {
			return true, nil
		}
	}

	return false, nil
}

// GetPostWithUserVote 获取帖子详情并包含用户投票信息
func (s *PostService) GetPostWithUserVote(postID, userID uint) (*models.Post, int, error) {
	post, err := s.postRepo.FindByIDWithReviewer(postID)
	if err != nil {
		return nil, 0, err
	}

	myVote := 0
	if userID > 0 {
		vote, err := s.voteRepo.FindByPostAndUser(postID, userID)
		if err == nil && vote != nil {
			myVote = int(vote.VoteType)
		}
	}

	return post, myVote, nil
}
