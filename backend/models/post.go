package models

import (
	"time"
)

// PostStatus 帖子状态
type PostStatus int

const (
	StatusPending      PostStatus = 0 // 待审核
	StatusFirstReview  PostStatus = 1 // 一级审核中(社区投票)
	StatusSecondReview PostStatus = 2 // 二级审核(等待认证用户提交邀请码)
	StatusApproved     PostStatus = 3 // 已通过
	StatusRejected     PostStatus = 4 // 已拒绝
)

// Post 帖子/申请模型
type Post struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	UserID       uint       `gorm:"index" json:"user_id"`
	User         *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Title        string     `gorm:"size:255" json:"title"`
	Content      string     `gorm:"type:text" json:"content"`
	Status       PostStatus `gorm:"default:1;index" json:"status"`
	UpVotes      int        `gorm:"default:0" json:"up_votes"`
	DownVotes    int        `gorm:"default:0" json:"down_votes"`
	ReviewerID   *uint      `gorm:"index" json:"reviewer_id,omitempty"`
	Reviewer     *User      `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
	InviteCode   string     `gorm:"size:255" json:"-"` // 邀请码(加密存储,不返回给前端)
	RejectReason string     `gorm:"size:500" json:"reject_reason,omitempty"` // 拒绝原因
	ReviewedAt   *time.Time `json:"reviewed_at,omitempty"` // 审核时间
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// TableName 指定表名
func (Post) TableName() string {
	return "posts"
}

// TotalVotes 获取总票数
func (p *Post) TotalVotes() int {
	return p.UpVotes + p.DownVotes
}

// ApprovalRate 计算赞率(百分比)
func (p *Post) ApprovalRate() float64 {
	total := p.TotalVotes()
	if total == 0 {
		return 0
	}
	return float64(p.UpVotes) / float64(total) * 100
}

// CanVote 是否可以投票(只有一级审核中的帖子可以投票)
func (p *Post) CanVote() bool {
	return p.Status == StatusFirstReview
}

// CanApprove 是否可以通过审核(只有二级审核的帖子可以通过)
func (p *Post) CanApprove() bool {
	return p.Status == StatusSecondReview
}

// CanReject 是否可以拒绝(一级或二级审核中的帖子可以被管理员拒绝)
func (p *Post) CanReject() bool {
	return p.Status == StatusFirstReview || p.Status == StatusSecondReview
}

// ShouldPromoteToSecondReview 检查是否应该进入二级审核
func (p *Post) ShouldPromoteToSecondReview(minVotes int, approvalRate float64) bool {
	return p.TotalVotes() >= minVotes && p.ApprovalRate() >= approvalRate
}

// ShouldReject 检查是否应该因投票不足被拒绝
func (p *Post) ShouldReject(minVotes int, approvalRate float64) bool {
	return p.TotalVotes() >= minVotes && p.ApprovalRate() < approvalRate
}
