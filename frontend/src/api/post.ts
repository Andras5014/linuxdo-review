import request from './request'
import type { ApiResponse, PaginatedResponse, Post, CreatePostRequest, VoteRequest, ApproveRequest } from '@/types'

// 获取帖子列表（一级审核区）
export const getPosts = (params?: { page?: number; page_size?: number; status?: number }) => {
  return request.get<ApiResponse<PaginatedResponse<Post>>>('/posts', { params })
}

// 获取二级审核列表（需认证用户）
export const getReviewPosts = (params?: { page?: number; page_size?: number }) => {
  return request.get<ApiResponse<PaginatedResponse<Post>>>('/posts/review', { params })
}

// 获取我的申请列表
export const getMyPosts = (params?: { page?: number; page_size?: number }) => {
  return request.get<ApiResponse<PaginatedResponse<Post>>>('/user/posts', { params })
}

// 获取帖子详情
export const getPostDetail = (id: number) => {
  return request.get<ApiResponse<Post>>(`/posts/${id}`)
}

// 发布申请
export const createPost = (data: CreatePostRequest) => {
  return request.post<ApiResponse<Post>>('/posts', data)
}

// 投票
export const votePost = (id: number, data: VoteRequest) => {
  return request.post<ApiResponse<null>>(`/posts/${id}/vote`, data)
}

// 获取下一个待审核的帖子
export const getNextReviewPost = (skipIds?: number[]) => {
  const params = skipIds?.length ? { skip_ids: skipIds.join(',') } : {}
  return request.get<ApiResponse<{ post: Post | null; total: number }>>('/review/next', { params })
}

// 跳过当前帖子
export const skipPost = (id: number) => {
  return request.post<ApiResponse<null>>(`/review/${id}/skip`)
}

// 通过审核并提交邀请码
export const approvePost = (id: number, data: ApproveRequest) => {
  return request.post<ApiResponse<null>>(`/review/${id}/approve`, data)
}

// 拒绝申请
export const rejectPost = (id: number, reason?: string) => {
  return request.post<ApiResponse<null>>(`/review/${id}/reject`, { reason })
}

// 获取系统统计数据
export const getStats = () => {
  return request.get<ApiResponse<{
    total_applications: number
    approved: number
    pending: number
    rejected: number
  }>>('/posts/stats')
}
