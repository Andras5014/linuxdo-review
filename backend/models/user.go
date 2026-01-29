package models

import (
	"time"
)

// UserRole 用户角色
type UserRole int

const (
	RoleNormal    UserRole = 0 // 普通用户
	RoleCertified UserRole = 1 // 认证用户
	RoleAdmin     UserRole = 2 // 管理员
)

// User 用户模型
type User struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Email           string    `gorm:"uniqueIndex;size:255" json:"email"`
	Password        string    `gorm:"size:255" json:"-"`
	Username        string    `gorm:"size:100" json:"username"`
	Role            UserRole  `gorm:"default:0" json:"role"`
	LinuxDoID       string    `gorm:"size:100;index" json:"linuxdo_id,omitempty"`
	LinuxDoUsername string    `gorm:"size:100" json:"linuxdo_username,omitempty"`
	AvatarURL       string    `gorm:"size:500" json:"avatar_url,omitempty"`
	TrustLevel      int       `gorm:"default:0" json:"trust_level"` // LinuxDo信任等级
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// IsCertified 是否为认证用户
func (u *User) IsCertified() bool {
	return u.Role >= RoleCertified
}

// IsAdmin 是否为管理员
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsLinuxDoUser 是否为LinuxDo OAuth用户
func (u *User) IsLinuxDoUser() bool {
	return u.LinuxDoID != ""
}

// CanVote 是否可以投票(所有登录用户都可以投票)
func (u *User) CanVote() bool {
	return u.ID > 0
}

// CanApprove 是否可以通过审核(只有认证用户可以)
func (u *User) CanApprove() bool {
	return u.IsCertified()
}

// CanManage 是否可以管理(只有管理员可以)
func (u *User) CanManage() bool {
	return u.IsAdmin()
}
