package dto

// RegisterRequest 注册请求
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Username string `json:"username" binding:"required,min=2,max=50"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// CreatePostRequest 创建帖子请求
type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=5,max=200"`
	Content string `json:"content" binding:"required,min=50"`
}

// UpdatePostRequest 更新帖子请求
type UpdatePostRequest struct {
	Title   string `json:"title" binding:"omitempty,min=5,max=200"`
	Content string `json:"content" binding:"omitempty,min=50"`
}

// VoteRequest 投票请求
type VoteRequest struct {
	VoteType int `json:"vote_type" binding:"required,oneof=1 -1"`
}

// ApproveRequest 审核通过请求(提交邀请码)
type ApproveRequest struct {
	InviteCode string `json:"invite_code" binding:"required,min=5"`
}

// RejectRequest 拒绝申请请求
type RejectRequest struct {
	Reason string `json:"reason" binding:"omitempty,max=500"`
}

// UpdateUserRoleRequest 更新用户角色请求
type UpdateUserRoleRequest struct {
	Role int `json:"role" binding:"oneof=0 1 2"`
}

// UpdateConfigRequest 更新配置请求
type UpdateConfigRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

// BatchUpdateConfigRequest 批量更新配置请求
type BatchUpdateConfigRequest struct {
	Configs []UpdateConfigRequest `json:"configs" binding:"required,dive"`
}

// PaginationRequest 分页请求
type PaginationRequest struct {
	Page     int `form:"page" binding:"min=1"`
	PageSize int `form:"page_size" binding:"min=1,max=100"`
}

// GetPage 获取页码(默认为1)
func (p *PaginationRequest) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

// GetPageSize 获取每页数量(默认为20)
func (p *PaginationRequest) GetPageSize() int {
	if p.PageSize <= 0 {
		return 20
	}
	return p.PageSize
}

// GetOffset 获取偏移量
func (p *PaginationRequest) GetOffset() int {
	return (p.GetPage() - 1) * p.GetPageSize()
}

// PostListRequest 帖子列表请求
// Status: -1=投票中+待审核, 0=待审核, 1=投票中, 2=二级审核, 3=已通过, 4=已拒绝
type PostListRequest struct {
	PaginationRequest
	Status *int   `form:"status" binding:"omitempty,oneof=-1 0 1 2 3 4"`
	Search string `form:"search" binding:"omitempty,max=100"`
}

// UserListRequest 用户列表请求
type UserListRequest struct {
	PaginationRequest
	Role   *int   `form:"role" binding:"omitempty,oneof=0 1 2"`
	Search string `form:"search" binding:"omitempty,max=100"`
}

// SetupAdminRequest 初始化管理员请求
type SetupAdminRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Username string `json:"username" binding:"required,min=2,max=50"`
}
