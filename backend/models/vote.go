package models

import (
	"time"
)

// VoteType 投票类型
type VoteType int

const (
	VoteUp   VoteType = 1  // 赞
	VoteDown VoteType = -1 // 踩
)

// Vote 投票模型
type Vote struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PostID    uint      `gorm:"index;uniqueIndex:idx_post_user" json:"post_id"`
	Post      *Post     `gorm:"foreignKey:PostID" json:"post,omitempty"`
	UserID    uint      `gorm:"index;uniqueIndex:idx_post_user" json:"user_id"`
	User      *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	VoteType  VoteType  `json:"vote_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Vote) TableName() string {
	return "votes"
}

// IsUpVote 是否为赞
func (v *Vote) IsUpVote() bool {
	return v.VoteType == VoteUp
}

// IsDownVote 是否为踩
func (v *Vote) IsDownVote() bool {
	return v.VoteType == VoteDown
}

// GetVoteTypeText 获取投票类型文本
func (v *Vote) GetVoteTypeText() string {
	switch v.VoteType {
	case VoteUp:
		return "赞"
	case VoteDown:
		return "踩"
	default:
		return "未知"
	}
}
