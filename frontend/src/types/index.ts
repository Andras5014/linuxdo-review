// 用户相关类型
export interface User {
  id: number
  email: string
  username: string
  role: UserRole
  linuxdo_id?: string
  linuxdo_username?: string
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