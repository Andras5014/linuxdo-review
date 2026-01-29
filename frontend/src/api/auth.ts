import request from './request'
import type { ApiResponse, LoginRequest, RegisterRequest, LoginResponse, User, SystemStatusResponse, SetupAdminRequest } from '@/types'

// 用户登录
export const login = (data: LoginRequest) => {
  return request.post<ApiResponse<LoginResponse>>('/auth/login', data)
}

// 用户注册
export const register = (data: RegisterRequest) => {
  return request.post<ApiResponse<User>>('/auth/register', data)
}

// 获取当前用户信息
export const getCurrentUser = () => {
  return request.get<ApiResponse<User>>('/auth/me')
}

// Linux.do OAuth登录
export const getLinuxdoOAuthUrl = () => {
  return `${window.location.origin}/api/auth/oauth/linuxdo`
}

// 退出登录
export const logout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
}

// 获取系统状态
export const getSystemStatus = () => {
  return request.get<ApiResponse<SystemStatusResponse>>('/system/status')
}

// 初始化管理员
export const setupAdmin = (data: SetupAdminRequest) => {
  return request.post<ApiResponse<User>>('/system/setup', data)
}
