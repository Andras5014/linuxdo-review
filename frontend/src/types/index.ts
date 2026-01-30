// 用户相关类型
export interface User {
  id: number
  email: string
  username: string
  role: UserRole
  role_text?: string
  linuxdo_id?: string
  linuxdo_username?: string
  avatar_url?: string
  trust_level?: number
  is_certified?: boolean
  is_admin?: boolean
  created_at: string
}

export enum UserRole {
  Normal = 0,
  Certified = 1,
  Admin = 2,
}

// 帖子相关类型
export interface Post {
  id: number
  user_id: number
  title: string
  content: string
  status: PostStatus
  up_votes: number
  down_votes: number
  reviewer_id?: number
  created_at: string
  user?: User
  reviewer?: User
}

export enum PostStatus {
  Pending = 0,
  FirstReview = 1,
  SecondReview = 2,
  Approved = 3,
  Rejected = 4,
}

// 投票相关类型
export interface Vote {
  id: number
  post_id: number
  user_id: number
  vote_type: VoteType
  created_at: string
}

export enum VoteType {
  Up = 1,
  Down = -1,
}

// API响应类型
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

// 分页响应
export interface PaginatedResponse<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}

// 登录请求
export interface LoginRequest {
  email: string
  password: string
}

// 注册请求
export interface RegisterRequest {
  email: string
  password: string
  username: string
}

// 登录响应
export interface LoginResponse {
  token: string
  user: User
}

// 帖子创建请求
export interface CreatePostRequest {
  title: string
  content: string
}

// 投票请求
export interface VoteRequest {
  vote_type: VoteType
}

// 投票响应
export interface VoteResponse {
  post_id: number
  vote_type: number
  up_votes: number
  down_votes: number
  message: string
}

// 审核请求（提交邀请码）
export interface ApproveRequest {
  invite_code: string
}

// 系统配置
export interface SystemConfig {
  min_votes: number
  approval_rate: number
}

// 系统状态响应
export interface SystemStatusResponse {
  initialized: boolean
}

// 初始化管理员请求
export interface SetupAdminRequest {
  email: string
  password: string
  username: string
}

// 更新个人资料请求
export interface UpdateProfileRequest {
  username?: string
}

// 绑定LinuxDO回调请求
export interface BindLinuxDoCallbackRequest {
  code: string
  state: string
  redirect_uri: string
}

// 绑定邮箱请求（LinuxDO用户）
export interface BindEmailRequest {
  email: string
  password: string
}

// OAuth URL响应
export interface OAuthURLResponse {
  url: string
  state?: string
}

// 发送邮箱验证码请求
export interface SendEmailCodeRequest {
  email: string
}

// 修改邮箱请求
export interface ChangeEmailRequest {
  new_email: string
  code: string
}

// 更新头像请求
export interface UpdateAvatarRequest {
  avatar_url: string
}