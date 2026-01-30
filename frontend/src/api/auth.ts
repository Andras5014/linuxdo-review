import request from './request'
import type { ApiResponse, LoginRequest, RegisterRequest, LoginResponse, User, SystemStatusResponse, SetupAdminRequest, UpdateProfileRequest, OAuthURLResponse, BindEmailRequest, SendEmailCodeRequest, ChangeEmailRequest, UpdateAvatarRequest } from '@/types'

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

// Linux.do OAuth登录 - 使用重定向方式
export const getLinuxdoOAuthUrl = () => {
  return `${window.location.origin}/api/auth/oauth/linuxdo/redirect`
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

// 获取用户资料
export const getProfile = () => {
  return request.get<ApiResponse<User>>('/user/profile')
}

// 更新用户资料
export const updateProfile = (data: UpdateProfileRequest) => {
  return request.put<ApiResponse<User>>('/user/profile', data)
}

// 获取绑定LinuxDO的OAuth URL
export const getBindLinuxDoUrl = () => {
  return request.get<ApiResponse<OAuthURLResponse>>('/user/bindlinuxdo')
}

// 解绑LinuxDO
export const unbindLinuxDo = () => {
  return request.post<ApiResponse<User>>('/user/unbindlinuxdo')
}

// 绑定邮箱（LinuxDO用户专用）
export const bindEmail = (data: BindEmailRequest) => {
  return request.post<ApiResponse<User>>('/user/bindmail', data)
}

// 发送邮箱验证码
export const sendEmailCode = (data: SendEmailCodeRequest) => {
  return request.post<ApiResponse<null>>('/user/email/code', data)
}

// 修改邮箱
export const changeEmail = (data: ChangeEmailRequest) => {
  return request.post<ApiResponse<User>>('/user/email/change', data)
}

// 更新头像
export const updateAvatar = (data: UpdateAvatarRequest) => {
  return request.put<ApiResponse<User>>('/user/avatar', data)
}
