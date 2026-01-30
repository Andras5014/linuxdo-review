<template>
  <div class="profile-page">
    <header class="header">
      <div class="header-content">
        <div class="logo" @click="$router.push('/')">
          <span class="logo-icon">🚀</span>
          <span class="logo-text">Linux.do</span>
          <span class="logo-badge">邀请码申请</span>
        </div>
        
        <nav class="nav">
          <router-link to="/posts" class="nav-link">申请列表</router-link>
          <router-link v-if="userStore.canReview" to="/review" class="nav-link">二级审核</router-link>
          <router-link v-if="userStore.isAdmin" to="/admin" class="nav-link">管理后台</router-link>
          
          <!-- 主题切换按钮 -->
          <button class="theme-toggle" @click="toggleTheme" :title="themeStore.theme === 'light' ? '切换到暗色模式' : '切换到亮色模式'">
            <Transition name="theme-icon" mode="out-in">
              <span v-if="themeStore.theme === 'dark'" key="sun" class="theme-icon">☀️</span>
              <span v-else key="moon" class="theme-icon">🌙</span>
            </Transition>
          </button>
          
          <a-dropdown>
            <div class="user-info">
              <a-avatar :size="32" class="user-avatar" :src="userStore.user?.avatar_url">
                {{ !userStore.user?.avatar_url ? userStore.username.charAt(0).toUpperCase() : '' }}
              </a-avatar>
              <span class="user-name">{{ userStore.username }}</span>
              <DownOutlined />
            </div>
            <template #overlay>
              <a-menu>
                <a-menu-item key="profile" @click="$router.push('/profile')">
                  <UserOutlined />
                  <span>个人资料</span>
                </a-menu-item>
                <a-menu-item key="my-posts" @click="$router.push('/my-posts')">
                  <FileTextOutlined />
                  <span>我的申请</span>
                </a-menu-item>
                <a-menu-divider />
                <a-menu-item key="logout" @click="handleLogout">
                  <LogoutOutlined />
                  <span>退出登录</span>
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </nav>
      </div>
    </header>

    <main class="main">
      <div class="profile-container">
        <div class="page-header slide-up">
          <h1 class="page-title">个人资料</h1>
          <p class="page-desc">管理您的账号信息和绑定设置</p>
        </div>

        <div class="profile-content">
          <!-- 用户基本信息卡片 -->
          <div class="profile-card slide-up" style="animation-delay: 0.1s">
            <div class="card-header">
              <div class="card-icon">
                <UserOutlined />
              </div>
              <div class="card-title-wrapper">
                <h3 class="card-title">基本信息</h3>
                <span class="card-subtitle">您的账户基本信息</span>
              </div>
            </div>
            
            <div class="card-body">
              <div class="info-item">
                <div class="info-avatar">
                  <a-avatar :size="80" class="avatar-large" :src="profile?.avatar_url">
                    {{ !profile?.avatar_url ? (profile?.username?.charAt(0)?.toUpperCase() || '?') : '' }}
                  </a-avatar>
                  <a-button type="link" size="small" class="avatar-edit-btn" @click="showChangeAvatarModal = true">
                    <EditOutlined /> 修改头像
                  </a-button>
                </div>
              </div>
              
              <div class="info-grid">
                <div class="info-item">
                  <span class="info-label">用户名</span>
                  <div class="info-value-row">
                    <span class="info-value">{{ profile?.username }}</span>
                    <a-button type="link" size="small" @click="showEditUsernameModal = true">
                      <EditOutlined /> 修改
                    </a-button>
                  </div>
                </div>
                
                <div class="info-item">
                  <span class="info-label">邮箱</span>
                  <div class="info-value-row">
                    <span class="info-value" :class="{ 'placeholder-email': isPlaceholderEmail }">
                      {{ isPlaceholderEmail ? '未绑定' : profile?.email }}
                    </span>
                    <a-button v-if="isPlaceholderEmail" type="link" size="small" @click="showBindEmailModal = true">
                      <LinkOutlined /> 绑定邮箱
                    </a-button>
                    <a-button v-else type="link" size="small" @click="showChangeEmailModal = true">
                      <EditOutlined /> 修改
                    </a-button>
                  </div>
                </div>
                
                <div class="info-item">
                  <span class="info-label">角色</span>
                  <div class="info-value">
                    <a-tag :color="getRoleColor(profile?.role)">{{ profile?.role_text }}</a-tag>
                  </div>
                </div>
                
                <div class="info-item">
                  <span class="info-label">注册时间</span>
                  <span class="info-value">{{ profile?.created_at }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- LinuxDO 绑定卡片 -->
          <div class="profile-card slide-up" style="animation-delay: 0.2s">
            <div class="card-header">
              <div class="card-icon linuxdo">
                <GlobalOutlined />
              </div>
              <div class="card-title-wrapper">
                <h3 class="card-title">Linux.do 账号绑定</h3>
                <span class="card-subtitle">绑定后可获得认证用户权限</span>
              </div>
            </div>
            
            <div class="card-body">
              <template v-if="profile?.linuxdo_id">
                <!-- 已绑定状态 -->
                <div class="bind-status bound">
                  <div class="status-icon">
                    <CheckCircleFilled />
                  </div>
                  <div class="status-content">
                    <h4 class="status-title">已绑定 Linux.do 账号</h4>
                    <div class="linuxdo-info">
                      <a-avatar :size="40" :src="profile?.avatar_url" class="linuxdo-avatar">
                        {{ !profile?.avatar_url ? profile?.linuxdo_username?.charAt(0)?.toUpperCase() : '' }}
                      </a-avatar>
                      <div class="linuxdo-details">
                        <span class="linuxdo-username">{{ profile?.linuxdo_username }}</span>
                        <span class="linuxdo-trust">
                          信任等级: 
                          <a-tag :color="getTrustLevelColor(profile?.trust_level)">
                            Level {{ profile?.trust_level }}
                          </a-tag>
                        </span>
                      </div>
                    </div>
                  </div>
                </div>
                
                <!-- 只有邮箱注册用户才能解绑 -->
                <div class="bind-actions" v-if="canUnbind">
                  <a-popconfirm
                    title="确定要解绑 Linux.do 账号吗？"
                    description="解绑后可能会失去认证用户权限"
                    @confirm="handleUnbind"
                    ok-text="确定"
                    cancel-text="取消"
                  >
                    <a-button danger :loading="unbindLoading">
                      <template #icon><DisconnectOutlined /></template>
                      解除绑定
                    </a-button>
                  </a-popconfirm>
                </div>
                <div class="bind-note" v-else>
                  <InfoCircleOutlined /> 您是通过 Linux.do OAuth 登录的用户，无法解绑
                </div>
              </template>
              
              <template v-else>
                <!-- 未绑定状态 -->
                <div class="bind-status unbound">
                  <div class="status-icon">
                    <ExclamationCircleFilled />
                  </div>
                  <div class="status-content">
                    <h4 class="status-title">尚未绑定 Linux.do 账号</h4>
                    <p class="status-desc">
                      绑定 Linux.do 账号后，如果您的信任等级 ≥ 2，将自动获得认证用户权限，可以参与二级审核。
                    </p>
                  </div>
                </div>
                
                <div class="bind-actions">
                  <a-button type="primary" size="large" @click="handleBind" :loading="bindLoading">
                    <template #icon><LinkOutlined /></template>
                    绑定 Linux.do 账号
                  </a-button>
                </div>
              </template>
            </div>
          </div>

          <!-- 权限说明卡片 -->
          <div class="profile-card slide-up" style="animation-delay: 0.3s">
            <div class="card-header">
              <div class="card-icon info">
                <SafetyCertificateOutlined />
              </div>
              <div class="card-title-wrapper">
                <h3 class="card-title">权限说明</h3>
                <span class="card-subtitle">不同角色的权限差异</span>
              </div>
            </div>
            
            <div class="card-body">
              <div class="permission-list">
                <div class="permission-item">
                  <div class="permission-role">
                    <a-tag color="default">普通用户</a-tag>
                  </div>
                  <div class="permission-desc">
                    <ul>
                      <li>可以提交邀请码申请</li>
                      <li>可以参与社区投票</li>
                      <li>查看申请列表</li>
                    </ul>
                  </div>
                </div>
                
                <div class="permission-item">
                  <div class="permission-role">
                    <a-tag color="blue">认证用户</a-tag>
                  </div>
                  <div class="permission-desc">
                    <ul>
                      <li>拥有普通用户所有权限</li>
                      <li>可以参与二级审核</li>
                      <li>可以提交邀请码给申请者</li>
                    </ul>
                  </div>
                </div>
                
                <div class="permission-item">
                  <div class="permission-role">
                    <a-tag color="purple">管理员</a-tag>
                  </div>
                  <div class="permission-desc">
                    <ul>
                      <li>拥有所有权限</li>
                      <li>管理用户和配置</li>
                      <li>查看系统统计</li>
                    </ul>
                  </div>
                </div>
              </div>
              
              <div class="certification-note">
                <InfoCircleOutlined />
                <span>获得认证用户权限：绑定 Linux.do 账号且信任等级 ≥ 2</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- 修改用户名弹窗 -->
    <a-modal
      v-model:open="showEditUsernameModal"
      title="修改用户名"
      @ok="handleUpdateUsername"
      :confirmLoading="updateLoading"
    >
      <a-form layout="vertical">
        <a-form-item label="新用户名">
          <a-input v-model:value="newUsername" placeholder="请输入新用户名" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 绑定邮箱弹窗 -->
    <a-modal
      v-model:open="showBindEmailModal"
      title="绑定邮箱"
      @ok="handleBindEmail"
      :confirmLoading="bindEmailLoading"
    >
      <p class="bind-email-tip">绑定邮箱后，您可以使用邮箱密码登录，并可以在将来解绑 Linux.do 账号。</p>
      <a-form layout="vertical">
        <a-form-item label="邮箱">
          <a-input v-model:value="bindEmailForm.email" placeholder="请输入邮箱" type="email" />
        </a-form-item>
        <a-form-item label="设置密码">
          <a-input-password v-model:value="bindEmailForm.password" placeholder="请输入密码（至少6位）" />
        </a-form-item>
        <a-form-item label="确认密码">
          <a-input-password v-model:value="bindEmailForm.confirmPassword" placeholder="请再次输入密码" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 修改邮箱弹窗 -->
    <a-modal
      v-model:open="showChangeEmailModal"
      title="修改邮箱"
      @ok="handleChangeEmail"
      :confirmLoading="changeEmailLoading"
    >
      <p class="bind-email-tip">修改邮箱需要验证新邮箱，验证码将发送到新邮箱。</p>
      <a-form layout="vertical">
        <a-form-item label="当前邮箱">
          <a-input :value="profile?.email" disabled />
        </a-form-item>
        <a-form-item label="新邮箱">
          <a-input v-model:value="changeEmailForm.newEmail" placeholder="请输入新邮箱" type="email" />
        </a-form-item>
        <a-form-item label="验证码">
          <a-input-group compact>
            <a-input 
              v-model:value="changeEmailForm.code" 
              placeholder="请输入6位验证码" 
              style="width: calc(100% - 120px)" 
              maxlength="6"
            />
            <a-button 
              type="primary" 
              :disabled="codeCountdown > 0 || sendCodeLoading"
              :loading="sendCodeLoading"
              @click="handleSendCode"
              style="width: 120px"
            >
              {{ codeCountdown > 0 ? `${codeCountdown}s后重发` : '发送验证码' }}
            </a-button>
          </a-input-group>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 修改头像弹窗 -->
    <a-modal
      v-model:open="showChangeAvatarModal"
      title="修改头像"
      @ok="handleChangeAvatar"
      :confirmLoading="changeAvatarLoading"
    >
      <p class="bind-email-tip">请输入头像图片的URL地址。推荐使用可靠的图床服务。</p>
      <a-form layout="vertical">
        <a-form-item label="当前头像">
          <div class="current-avatar-preview">
            <a-avatar :size="64" :src="profile?.avatar_url">
              {{ !profile?.avatar_url ? (profile?.username?.charAt(0)?.toUpperCase() || '?') : '' }}
            </a-avatar>
          </div>
        </a-form-item>
        <a-form-item label="新头像URL">
          <a-input v-model:value="newAvatarUrl" placeholder="请输入图片URL，如 https://example.com/avatar.jpg" />
        </a-form-item>
        <a-form-item v-if="newAvatarUrl" label="预览">
          <div class="avatar-preview">
            <a-avatar :size="64" :src="newAvatarUrl">
              ?
            </a-avatar>
          </div>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { message } from 'ant-design-vue'
import {
  DownOutlined,
  UserOutlined,
  FileTextOutlined,
  LogoutOutlined,
  GlobalOutlined,
  EditOutlined,
  CheckCircleFilled,
  ExclamationCircleFilled,
  LinkOutlined,
  DisconnectOutlined,
  InfoCircleOutlined,
  SafetyCertificateOutlined,
} from '@ant-design/icons-vue'
import { useUserStore } from '@/stores/user'
import { useThemeStore } from '@/stores/theme'
import { getProfile, updateProfile, getBindLinuxDoUrl, unbindLinuxDo, bindEmail, sendEmailCode, changeEmail, updateAvatar } from '@/api/auth'
import type { User } from '@/types'
import { UserRole } from '@/types'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const themeStore = useThemeStore()

const profile = ref<User | null>(null)
const loading = ref(false)
const bindLoading = ref(false)
const unbindLoading = ref(false)
const updateLoading = ref(false)
const showEditUsernameModal = ref(false)
const newUsername = ref('')
const showBindEmailModal = ref(false)
const bindEmailLoading = ref(false)
const bindEmailForm = ref({
  email: '',
  password: '',
  confirmPassword: ''
})

// 修改邮箱相关
const showChangeEmailModal = ref(false)
const changeEmailLoading = ref(false)
const sendCodeLoading = ref(false)
const codeCountdown = ref(0)
const changeEmailForm = ref({
  newEmail: '',
  code: ''
})

// 修改头像相关
const showChangeAvatarModal = ref(false)
const changeAvatarLoading = ref(false)
const newAvatarUrl = ref('')

// 判断是否是占位邮箱（LinuxDO登录用户未绑定真实邮箱）
const isPlaceholderEmail = computed(() => {
  return profile.value?.email?.endsWith('@linuxdo.user') ?? false
})

// 判断是否可以解绑（只有邮箱注册用户才能解绑）
const canUnbind = computed(() => {
  // 如果用户是通过OAuth登录的（没有设置密码的邮箱），则不能解绑
  // 这里我们通过检查email是否包含@linuxdo.user来判断
  return profile.value?.email && !profile.value.email.endsWith('@linuxdo.user')
})

const toggleTheme = () => {
  themeStore.toggleTheme()
}

const handleLogout = () => {
  userStore.logout()
  message.success('已退出登录')
  router.push('/')
}

const fetchProfile = async () => {
  loading.value = true
  try {
    const response = await getProfile()
    profile.value = response.data.data
    newUsername.value = profile.value?.username || ''
  } catch {
    message.error('获取用户信息失败')
  } finally {
    loading.value = false
  }
}

const getRoleColor = (role?: UserRole) => {
  switch (role) {
    case UserRole.Admin:
      return 'purple'
    case UserRole.Certified:
      return 'blue'
    default:
      return 'default'
  }
}

const getTrustLevelColor = (level?: number) => {
  if (!level) return 'default'
  if (level >= 4) return 'gold'
  if (level >= 3) return 'green'
  if (level >= 2) return 'blue'
  return 'default'
}

const handleBind = async () => {
  bindLoading.value = true
  try {
    // 获取绑定URL（后端会自动处理回调地址和用户信息）
    const response = await getBindLinuxDoUrl()
    const { url } = response.data.data
    
    // 跳转到OAuth授权页面
    window.location.href = url
  } catch {
    message.error('获取授权链接失败')
  } finally {
    bindLoading.value = false
  }
}

const handleUnbind = async () => {
  unbindLoading.value = true
  try {
    const response = await unbindLinuxDo()
    profile.value = response.data.data
    userStore.setAuth(localStorage.getItem('token') || '', response.data.data)
    message.success('解绑成功')
  } catch {
    message.error('解绑失败')
  } finally {
    unbindLoading.value = false
  }
}

const handleUpdateUsername = async () => {
  if (!newUsername.value || newUsername.value.length < 2) {
    message.error('用户名至少2个字符')
    return
  }
  
  updateLoading.value = true
  try {
    const response = await updateProfile({ username: newUsername.value })
    profile.value = response.data.data
    userStore.setAuth(localStorage.getItem('token') || '', response.data.data)
    showEditUsernameModal.value = false
    message.success('用户名修改成功')
  } catch {
    message.error('修改失败')
  } finally {
    updateLoading.value = false
  }
}

const handleBindEmail = async () => {
  // 验证表单
  if (!bindEmailForm.value.email) {
    message.error('请输入邮箱')
    return
  }
  
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(bindEmailForm.value.email)) {
    message.error('请输入有效的邮箱地址')
    return
  }
  
  if (!bindEmailForm.value.password || bindEmailForm.value.password.length < 6) {
    message.error('密码至少6个字符')
    return
  }
  
  if (bindEmailForm.value.password !== bindEmailForm.value.confirmPassword) {
    message.error('两次输入的密码不一致')
    return
  }
  
  bindEmailLoading.value = true
  try {
    const response = await bindEmail({
      email: bindEmailForm.value.email,
      password: bindEmailForm.value.password
    })
    profile.value = response.data.data
    userStore.setAuth(localStorage.getItem('token') || '', response.data.data)
    showBindEmailModal.value = false
    // 重置表单
    bindEmailForm.value = { email: '', password: '', confirmPassword: '' }
    message.success('邮箱绑定成功')
  } catch {
    message.error('绑定邮箱失败')
  } finally {
    bindEmailLoading.value = false
  }
}

// 发送邮箱验证码
const handleSendCode = async () => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!changeEmailForm.value.newEmail || !emailRegex.test(changeEmailForm.value.newEmail)) {
    message.error('请输入有效的邮箱地址')
    return
  }

  sendCodeLoading.value = true
  try {
    await sendEmailCode({ email: changeEmailForm.value.newEmail })
    message.success('验证码已发送到新邮箱')
    // 开始倒计时
    codeCountdown.value = 60
    const timer = setInterval(() => {
      codeCountdown.value--
      if (codeCountdown.value <= 0) {
        clearInterval(timer)
      }
    }, 1000)
  } catch {
    message.error('发送验证码失败')
  } finally {
    sendCodeLoading.value = false
  }
}

// 修改邮箱
const handleChangeEmail = async () => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!changeEmailForm.value.newEmail || !emailRegex.test(changeEmailForm.value.newEmail)) {
    message.error('请输入有效的邮箱地址')
    return
  }
  
  if (!changeEmailForm.value.code || changeEmailForm.value.code.length !== 6) {
    message.error('请输入6位验证码')
    return
  }

  changeEmailLoading.value = true
  try {
    const response = await changeEmail({
      new_email: changeEmailForm.value.newEmail,
      code: changeEmailForm.value.code
    })
    profile.value = response.data.data
    userStore.setAuth(localStorage.getItem('token') || '', response.data.data)
    showChangeEmailModal.value = false
    changeEmailForm.value = { newEmail: '', code: '' }
    message.success('邮箱修改成功')
  } catch {
    message.error('修改邮箱失败')
  } finally {
    changeEmailLoading.value = false
  }
}

// 修改头像
const handleChangeAvatar = async () => {
  const urlRegex = /^https?:\/\/.+/
  if (!newAvatarUrl.value || !urlRegex.test(newAvatarUrl.value)) {
    message.error('请输入有效的图片URL')
    return
  }

  changeAvatarLoading.value = true
  try {
    const response = await updateAvatar({ avatar_url: newAvatarUrl.value })
    profile.value = response.data.data
    userStore.setAuth(localStorage.getItem('token') || '', response.data.data)
    showChangeAvatarModal.value = false
    newAvatarUrl.value = ''
    message.success('头像修改成功')
  } catch {
    message.error('修改头像失败')
  } finally {
    changeAvatarLoading.value = false
  }
}

onMounted(async () => {
  // 检查是否是绑定成功后的跳转
  if (route.query.bindSuccess === 'true') {
    message.success('Linux.do 账号绑定成功！')
    // 清除 URL 中的参数
    router.replace('/profile')
  }
  
  await fetchProfile()
  
  // 如果是绑定成功，更新 store 中的用户信息
  if (route.query.bindSuccess === 'true' && profile.value) {
    const token = localStorage.getItem('token')
    if (token) {
      userStore.setAuth(token, profile.value)
    }
  }
})
</script>

<style scoped>
.profile-page {
  min-height: 100vh;
  background: var(--bg-primary);
}

/* Header */
.header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  background: var(--glass-bg);
  backdrop-filter: blur(20px);
  border-bottom: 1px solid var(--border-color-light);
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 16px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
}

.logo-icon {
  font-size: 28px;
}

.logo-text {
  font-size: 20px;
  font-weight: 700;
  color: var(--text-primary);
}

.logo-badge {
  font-size: 11px;
  padding: 3px 10px;
  background: var(--color-primary-gradient);
  color: white;
  border-radius: 20px;
  font-weight: 500;
}

.nav {
  display: flex;
  align-items: center;
  gap: 20px;
}

.nav-link {
  color: var(--text-secondary);
  font-weight: 500;
  transition: color 0.2s;
  font-size: 14px;
}

.nav-link:hover {
  color: var(--color-primary);
}

/* Theme Toggle */
.theme-toggle {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  border: 1px solid var(--border-color);
  background: var(--bg-secondary);
  color: var(--text-secondary);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  transition: all 0.2s ease;
}

.theme-toggle:hover {
  border-color: var(--color-primary);
  color: var(--color-primary);
  background: var(--color-primary-light);
}

.theme-icon {
  font-size: 18px;
  line-height: 1;
}

.theme-icon-enter-active,
.theme-icon-leave-active {
  transition: all 0.2s ease;
}

.theme-icon-enter-from {
  opacity: 0;
  transform: rotate(-90deg) scale(0.5);
}

.theme-icon-leave-to {
  opacity: 0;
  transform: rotate(90deg) scale(0.5);
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 6px 12px;
  border-radius: 12px;
  transition: background 0.2s;
  color: var(--text-secondary);
}

.user-info:hover {
  background: var(--bg-tertiary);
}

.user-avatar {
  background: var(--color-primary-gradient) !important;
  color: white !important;
}

.user-name {
  color: var(--text-primary);
  font-weight: 500;
  font-size: 14px;
}

/* Main */
.main {
  padding-top: 90px;
  padding-bottom: 40px;
}

.profile-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 0 24px;
}

.page-header {
  margin-bottom: 32px;
}

.page-title {
  font-size: 32px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.page-desc {
  font-size: 15px;
  color: var(--text-secondary);
}

/* Profile Cards */
.profile-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.profile-card {
  background: var(--bg-card);
  border: 1px solid var(--border-color-light);
  border-radius: 20px;
  overflow: hidden;
  backdrop-filter: blur(20px);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 24px;
  border-bottom: 1px solid var(--border-color-light);
}

.card-icon {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-primary-gradient);
  border-radius: 14px;
  font-size: 22px;
  color: white;
  flex-shrink: 0;
}

.card-icon.linuxdo {
  background: linear-gradient(135deg, #10b981, #059669);
}

.card-icon.info {
  background: linear-gradient(135deg, #8b5cf6, #6366f1);
}

.card-title-wrapper {
  flex: 1;
}

.card-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.card-subtitle {
  font-size: 13px;
  color: var(--text-muted);
}

.card-body {
  padding: 24px;
}

/* Info Items */
.info-avatar {
  display: flex;
  justify-content: center;
  margin-bottom: 24px;
}

.avatar-large {
  background: var(--color-primary-gradient) !important;
  color: white !important;
  font-size: 32px !important;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.info-label {
  font-size: 13px;
  color: var(--text-muted);
}

.info-value {
  font-size: 15px;
  color: var(--text-primary);
  font-weight: 500;
}

.info-value-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

/* Bind Status */
.bind-status {
  display: flex;
  gap: 16px;
  padding: 20px;
  border-radius: 14px;
  margin-bottom: 20px;
}

.bind-status.bound {
  background: rgba(16, 185, 129, 0.08);
  border: 1px solid rgba(16, 185, 129, 0.2);
}

.bind-status.unbound {
  background: rgba(245, 158, 11, 0.08);
  border: 1px solid rgba(245, 158, 11, 0.2);
}

.status-icon {
  font-size: 24px;
  flex-shrink: 0;
}

.bind-status.bound .status-icon {
  color: #10b981;
}

.bind-status.unbound .status-icon {
  color: #f59e0b;
}

.status-content {
  flex: 1;
}

.status-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.status-desc {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.6;
}

.linuxdo-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.linuxdo-avatar {
  background: linear-gradient(135deg, #10b981, #059669) !important;
  color: white !important;
}

.linuxdo-details {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.linuxdo-username {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}

.linuxdo-trust {
  font-size: 13px;
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  gap: 4px;
}

.bind-actions {
  display: flex;
  justify-content: flex-start;
}

.bind-note {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: var(--text-muted);
  padding: 12px 16px;
  background: var(--bg-tertiary);
  border-radius: 10px;
}

/* Permission List */
.permission-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.permission-item {
  display: flex;
  gap: 16px;
  padding: 16px;
  background: var(--bg-tertiary);
  border-radius: 12px;
}

.permission-role {
  flex-shrink: 0;
  width: 100px;
}

.permission-desc ul {
  margin: 0;
  padding-left: 20px;
  color: var(--text-secondary);
  font-size: 14px;
  line-height: 1.8;
}

.certification-note {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 20px;
  padding: 14px 16px;
  background: rgba(99, 102, 241, 0.08);
  border: 1px solid rgba(99, 102, 241, 0.2);
  border-radius: 10px;
  font-size: 14px;
  color: var(--color-primary);
}

/* 占位邮箱样式 */
.placeholder-email {
  color: var(--text-muted);
  font-style: italic;
}

/* 绑定邮箱弹窗提示 */
.bind-email-tip {
  color: var(--text-secondary);
  font-size: 14px;
  margin-bottom: 16px;
  padding: 12px;
  background: var(--bg-tertiary);
  border-radius: 8px;
  line-height: 1.6;
}

/* 头像编辑按钮 */
.info-avatar {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.avatar-edit-btn {
  margin-top: 4px;
}

/* 头像预览 */
.current-avatar-preview,
.avatar-preview {
  display: flex;
  justify-content: center;
  padding: 8px;
  background: var(--bg-tertiary);
  border-radius: 8px;
}

/* Responsive */
@media (max-width: 768px) {
  .info-grid {
    grid-template-columns: 1fr;
  }
  
  .permission-item {
    flex-direction: column;
    gap: 12px;
  }
  
  .permission-role {
    width: auto;
  }

  .nav-link {
    display: none;
  }

  .page-title {
    font-size: 24px;
  }
}
</style>
