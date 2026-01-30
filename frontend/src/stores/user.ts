import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types'
import { UserRole } from '@/types'
import { getCurrentUser, logout as apiLogout } from '@/api/auth'

export const useUserStore = defineStore('user', () => {
  // 状态
  const user = ref<User | null>(null)
  const token = ref<string | null>(localStorage.getItem('token'))
  const loading = ref(false)
  const initialized = ref(false)

  // 计算属性
  const isLoggedIn = computed(() => !!token.value && !!user.value)
  const isAdmin = computed(() => user.value?.role === UserRole.Admin)
  const isCertified = computed(() => 
    user.value?.role === UserRole.Certified || user.value?.role === UserRole.Admin
  )
  // 是否绑定了 Linux.do 账号
  const isLinuxDoBound = computed(() => !!user.value?.linuxdo_id)
  // 信任等级
  const trustLevel = computed(() => user.value?.trust_level || 0)
  // 是否可以进行二级审核 (trust_level >= 3 或管理员)
  const canReview = computed(() => 
    user.value?.role === UserRole.Admin || (user.value?.trust_level || 0) >= 3
  )
  const username = computed(() => user.value?.username || user.value?.email || '用户')

  // 初始化用户状态
  const initUser = async () => {
    if (initialized.value) return
    
    // 首先从 localStorage 恢复用户信息
    const storedUser = localStorage.getItem('user')
    if (storedUser) {
      try {
        user.value = JSON.parse(storedUser)
      } catch {
        user.value = null
      }
    }

    // 如果有token，尝试获取最新用户信息
    if (token.value) {
      await fetchUser()
    }
    
    initialized.value = true
  }

  // 获取用户信息
  const fetchUser = async () => {
    if (!token.value) return

    loading.value = true
    try {
      const response = await getCurrentUser()
      user.value = response.data.data
      localStorage.setItem('user', JSON.stringify(user.value))
    } catch {
      // 获取失败时清除登录状态
      clearAuth()
    } finally {
      loading.value = false
    }
  }

  // 设置登录信息
  const setAuth = (newToken: string, newUser: User) => {
    token.value = newToken
    user.value = newUser
    localStorage.setItem('token', newToken)
    localStorage.setItem('user', JSON.stringify(newUser))
  }

  // 清除登录信息
  const clearAuth = () => {
    token.value = null
    user.value = null
    apiLogout()
  }

  // 退出登录
  const logout = () => {
    clearAuth()
  }

  return {
    // 状态
    user,
    token,
    loading,
    
    // 计算属性
    isLoggedIn,
    isAdmin,
    isCertified,
    isLinuxDoBound,
    trustLevel,
    canReview,
    username,
    
    // 方法
    initUser,
    fetchUser,
    setAuth,
    clearAuth,
    logout,
  }
})
