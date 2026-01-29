package dto

import "linuxdo-review/models"

// UserResponse 用户响应
type UserResponse struct {
	ID              uint            `json:"id"`
	Email           string          `json:"email"`
	Username        string          `json:"username"`
	Role            models.UserRole `json:"role"`
	RoleText        string          `json:"role_text"`
	LinuxDoID       string          `json:"linuxdo_id,omitempty"`
	LinuxDoUsername string          `json:"linuxdo_username,omitempty"`
	AvatarURL       string          `json:"avatar_url,omitempty"`
	TrustLevel      int             `json:"trust_level,omitempty"`
	IsCertified     bool            `json:"is_certified"`
	IsAdmin         bool            `json:"is_admin"`
	CreatedAt       string          `json:"created_at"`
}

// GetRoleText 获取角色文本
func GetRoleText(role models.UserRole) string {
	switch role {
	case models.RoleNormal:
		return "普通用户"
	case models.RoleCertified:
		return "认证用户"
	case models.RoleAdmin:
		return "管理员"
	default:
		return "未知"
	}
}

// ToUserResponse 转换为用户响应
func ToUserResponse(user *models.User) *UserResponse {
	return &UserResponse{
		ID:              user.ID,
		Email:           user.Email,
		Username:        user.Username,
		Role:            user.Role,
		RoleText:        GetRoleText(user.Role),
		LinuxDoID:       user.LinuxDoID,
		LinuxDoUsername: user.LinuxDoUsername,
		AvatarURL:       user.AvatarURL,
		TrustLevel:      user.TrustLevel,
		IsCertified:     user.IsCertified(),
		IsAdmin:         user.IsAdmin(),
		CreatedAt:       user.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ToUserResponseList 批量转换为用户响应列表
func ToUserResponseList(users []*models.User) []*UserResponse {
	list := make([]*UserResponse, len(users))
	for i, user := range users {
		list[i] = ToUserResponse(user)
	}
	return list
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string        `json:"token"`
	User  *UserResponse `json:"user"`
}

// PostResponse 帖子响应
type PostResponse struct {
	ID           uint              `json:"id"`
	UserID       uint              `json:"user_id"`
	User         *UserResponse     `json:"user,omitempty"`
	Title        string            `json:"title"`
	Content      string            `json:"content"`
	Status       models.PostStatus `json:"status"`
	StatusText   string            `json:"status_text"`
	UpVotes      int               `json:"up_votes"`
	DownVotes    int               `json:"down_votes"`
	TotalVotes   int               `json:"total_votes"`
	ApprovalRate float64           `json:"approval_rate"`
	ReviewerID   *uint             `json:"reviewer_id,omitempty"`
	Reviewer     *UserResponse     `json:"reviewer,omitempty"`
	RejectReason string            `json:"reject_reason,omitempty"`
	ReviewedAt   string            `json:"reviewed_at,omitempty"`
	CreatedAt    string            `json:"created_at"`
	UpdatedAt    string            `json:"updated_at"`
	MyVote       int               `json:"my_vote,omitempty"` // 当前用户的投票: 1赞, -1踩, 0未投票
	CanVote      bool              `json:"can_vote"`          // 是否可以投票
	CanApprove   bool              `json:"can_approve"`       // 是否可以通过
}

// GetStatusText 获取状态文本
func GetStatusText(status models.PostStatus) string {
	switch status {
	case models.StatusPending:
		return "待审核"
	case models.StatusFirstReview:
		return "社区投票中"
	case models.StatusSecondReview:
		return "等待邀请码"
	case models.StatusApproved:
		return "已通过"
	case models.StatusRejected:
		return "已拒绝"
	default:
		return "未知"
	}
}

// ToPostResponse 转换为帖子响应
func ToPostResponse(post *models.Post) *PostResponse {
	resp := &PostResponse{
		ID:           post.ID,
		UserID:       post.UserID,
		Title:        post.Title,
		Content:      post.Content,
		Status:       post.Status,
		StatusText:   GetStatusText(post.Status),
		UpVotes:      post.UpVotes,
		DownVotes:    post.DownVotes,
		TotalVotes:   post.TotalVotes(),
		ApprovalRate: post.ApprovalRate(),
		ReviewerID:   post.ReviewerID,
		RejectReason: post.RejectReason,
		CreatedAt:    post.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    post.UpdatedAt.Format("2006-01-02 15:04:05"),
		CanVote:      post.CanVote(),
		CanApprove:   post.CanApprove(),
	}

	if post.ReviewedAt != nil {
		resp.ReviewedAt = post.ReviewedAt.Format("2006-01-02 15:04:05")
	}

	if post.User != nil {
		resp.User = ToUserResponse(post.User)
	}

	if post.Reviewer != nil {
		resp.Reviewer = ToUserResponse(post.Reviewer)
	}

	return resp
}

// ToPostResponseList 批量转换为帖子响应列表
func ToPostResponseList(posts []*models.Post) []*PostResponse {
	list := make([]*PostResponse, len(posts))
	for i, post := range posts {
		list[i] = ToPostResponse(post)
	}
	return list
}

// VoteResponse 投票响应
type VoteResponse struct {
	PostID    uint             `json:"post_id"`
	VoteType  models.VoteType  `json:"vote_type"`
	UpVotes   int              `json:"up_votes"`
	DownVotes int              `json:"down_votes"`
	Message   string           `json:"message"`
}

// ConfigResponse 配置响应
type ConfigResponse struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

// ToConfigResponse 转换为配置响应
func ToConfigResponse(config *models.SystemConfig) *ConfigResponse {
	return &ConfigResponse{
		Key:         config.Key,
		Value:       config.Value,
		Description: config.Description,
	}
}

// ToConfigResponseList 批量转换为配置响应列表
func ToConfigResponseList(configs []*models.SystemConfig) []*ConfigResponse {
	list := make([]*ConfigResponse, len(configs))
	for i, config := range configs {
		list[i] = ToConfigResponse(config)
	}
	return list
}

// PaginationResponse 分页响应
type PaginationResponse struct {
	List       interface{} `json:"list"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// NewPaginationResponse 创建分页响应
func NewPaginationResponse(list interface{}, total int64, page, pageSize int) *PaginationResponse {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return &PaginationResponse{
		List:       list,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

// StatsResponse 统计响应
type StatsResponse struct {
	TotalUsers       int64 `json:"total_users"`
	CertifiedUsers   int64 `json:"certified_users"`
	TotalPosts       int64 `json:"total_posts"`
	PendingPosts     int64 `json:"pending_posts"`
	FirstReviewPosts int64 `json:"first_review_posts"`
	SecondReviewPosts int64 `json:"second_review_posts"`
	ApprovedPosts    int64 `json:"approved_posts"`
	RejectedPosts    int64 `json:"rejected_posts"`
	TotalVotes       int64 `json:"total_votes"`
	TodayNewUsers    int64 `json:"today_new_users"`
	TodayNewPosts    int64 `json:"today_new_posts"`
	TodayApproved    int64 `json:"today_approved"`
}

// OAuthURLResponse OAuth跳转URL响应
type OAuthURLResponse struct {
	URL   string `json:"url"`
	State string `json:"state,omitempty"`
}

// SystemStatusResponse 系统状态响应
type SystemStatusResponse struct {
	Initialized bool `json:"initialized"` // 是否已初始化（是否有管理员）
}