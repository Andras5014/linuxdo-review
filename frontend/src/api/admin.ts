import request from './request'
import type { ApiResponse, PaginatedResponse, User } from '@/types'

// 系统配置类型
export interface SystemConfigItem {
  key: string
  value: string
  description: string
}

// 用户角色更新请求
export interface UpdateUserRoleRequest {
  role: number
}

// 系统统计数据
export interface SystemStats {
  total_users: number
  certified_users: number
  total_posts: number
  pending_posts: number
  approved_posts: number
  rejected_posts: number
}

// 获取用户列表
export const getUsers = (params?: { page?: number; page_size?: number; role?: number }) => {
  return request.get<ApiResponse<PaginatedResponse<User>>>('/admin/users', { params })
}

// 获取单个用户
export const getUser = (id: number) => {
  return request.get<ApiResponse<User>>(`/admin/users/${id}`)
}

// 更新用户角色
export const updateUserRole = (id: number, data: UpdateUserRoleRequest) => {
  return request.put<ApiResponse<null>>(`/admin/users/${id}`, data)
}

// 获取系统配置
export const getConfigs = () => {
  return request.get<ApiResponse<SystemConfigItem[]>>('/admin/configs')
}

// 更新单个配置
export const updateConfig = (key: string, value: string) => {
  return request.put<ApiResponse<null>>('/admin/configs', { key, value })
}

// 批量更新配置
export const batchUpdateConfigs = (configs: { key: string; value: string }[]) => {
  return request.put<ApiResponse<null>>('/admin/configs/batch', { configs })
}

// 获取系统统计
export const getStats = () => {
  return request.get<ApiResponse<SystemStats>>('/admin/stats')
}
