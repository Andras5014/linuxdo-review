<template>
  <div class="auth-page">
    <div class="auth-background">
      <div class="auth-glow"></div>
      <div class="floating-shapes">
        <div class="shape shape-1"></div>
        <div class="shape shape-2"></div>
      </div>
    </div>
    
    <!-- 主题切换按钮 -->
    <button class="theme-toggle" @click="toggleTheme" :title="themeStore.theme === 'light' ? '切换到暗色模式' : '切换到亮色模式'">
      <Transition name="theme-icon" mode="out-in">
        <span v-if="themeStore.theme === 'dark'" key="sun" class="theme-icon">☀️</span>
        <span v-else key="moon" class="theme-icon">🌙</span>
      </Transition>
    </button>
    
    <div class="auth-container slide-up">
      <div class="auth-header">
        <router-link to="/" class="logo">
          <span class="logo-icon">🚀</span>
          <span class="logo-text">Linux.do</span>
        </router-link>
        <h1 class="auth-title">创建账号</h1>
        <p class="auth-subtitle">开始您的邀请码申请之旅</p>
      </div>

      <a-form
        :model="formState"
        :rules="rules"
        @finish="handleSubmit"
        layout="vertical"
        class="auth-form"
      >
        <a-form-item name="username" label="用户名">
          <a-input
            v-model:value="formState.username"
            placeholder="请输入用户名"
            size="large"
            :prefix="h(UserOutlined)"
          />
        </a-form-item>

        <a-form-item name="email" label="邮箱">
          <a-input
            v-model:value="formState.email"
            placeholder="请输入邮箱"
            size="large"
            :prefix="h(MailOutlined)"
          />
        </a-form-item>

        <a-form-item name="password" label="密码">
          <a-input-password
            v-model:value="formState.password"
            placeholder="请输入密码（至少6位）"
            size="large"
            :prefix="h(LockOutlined)"
          />
        </a-form-item>

        <a-form-item name="confirmPassword" label="确认密码">
          <a-input-password
            v-model:value="formState.confirmPassword"
            placeholder="请再次输入密码"
            size="large"
            :prefix="h(LockOutlined)"
          />
        </a-form-item>

        <a-form-item>
          <a-button
            type="primary"
            html-type="submit"
            size="large"
            :loading="loading"
            block
            class="submit-btn"
          >
            注册
          </a-button>
        </a-form-item>
      </a-form>

      <div class="divider">
        <span>或</span>
      </div>

      <a-button
        size="large"
        block
        class="oauth-btn"
        @click="handleOAuthLogin"
      >
        <template #icon>
          <GlobalOutlined />
        </template>
        使用 Linux.do 账号注册
      </a-button>

      <div class="auth-footer">
        <span class="footer-text">已有账号？</span>
        <router-link to="/login" class="footer-link">立即登录</router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, h } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { UserOutlined, MailOutlined, LockOutlined, GlobalOutlined } from '@ant-design/icons-vue'
import { register, getLinuxdoOAuthUrl } from '@/api/auth'
import { useThemeStore } from '@/stores/theme'
import type { Rule } from 'ant-design-vue/es/form'

const router = useRouter()
const themeStore = useThemeStore()

const formState = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
})

const loading = ref(false)

const validatePassword = async (_rule: Rule, value: string) => {
  if (value && value !== formState.password) {
    return Promise.reject('两次输入的密码不一致')
  }
  return Promise.resolve()
}

const rules: Record<string, Rule[]> = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 2, max: 20, message: '用户名长度为2-20个字符', trigger: 'blur' },
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少6个字符', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { validator: validatePassword, trigger: 'blur' },
  ],
}

const toggleTheme = () => {
  themeStore.toggleTheme()
}

const handleSubmit = async () => {
  loading.value = true
  try {
    await register({
      username: formState.username,
      email: formState.email,
      password: formState.password,
    })
    
    message.success('注册成功，请登录')
    router.push('/login')
  } catch {
    // 错误已在拦截器中处理
  } finally {
    loading.value = false
  }
}

const handleOAuthLogin = () => {
  window.location.href = getLinuxdoOAuthUrl()
}
</script>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  position: relative;
  overflow: hidden;
}

.auth-background {
  position: absolute;
  inset: 0;
  z-index: 0;
}

.auth-glow {
  position: absolute;
  top: -150px;
  left: 50%;
  transform: translateX(-50%);
  width: 800px;
  height: 500px;
  background: var(--bg-hero-gradient);
  filter: blur(60px);
}

.floating-shapes {
  position: absolute;
  inset: 0;
  overflow: hidden;
}

.shape {
  position: absolute;
  border-radius: 50%;
  opacity: 0.4;
  filter: blur(80px);
}

.shape-1 {
  width: 350px;
  height: 350px;
  background: rgba(139, 92, 246, 0.25);
  top: 10%;
  left: 10%;
  animation: float 8s ease-in-out infinite;
}

.shape-2 {
  width: 280px;
  height: 280px;
  background: rgba(99, 102, 241, 0.2);
  bottom: 10%;
  right: 5%;
  animation: float 10s ease-in-out infinite reverse;
}

/* Theme Toggle */
.theme-toggle {
  position: fixed;
  top: 24px;
  right: 24px;
  z-index: 10;
  width: 44px;
  height: 44px;
  border-radius: 14px;
  border: 1px solid var(--border-color);
  background: var(--glass-bg);
  backdrop-filter: blur(20px);
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

.auth-container {
  width: 100%;
  max-width: 420px;
  background: var(--bg-card);
  backdrop-filter: blur(24px);
  border: 1px solid var(--border-color-light);
  border-radius: 28px;
  padding: 48px 40px;
  position: relative;
  z-index: 1;
  box-shadow: var(--shadow-xl);
}

.auth-header {
  text-align: center;
  margin-bottom: 36px;
}

.logo {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 28px;
  text-decoration: none;
}

.logo-icon {
  font-size: 32px;
}

.logo-text {
  font-size: 24px;
  font-weight: 700;
  background: var(--color-primary-gradient);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.auth-title {
  font-size: 28px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.auth-subtitle {
  font-size: 15px;
  color: var(--text-secondary);
}

.auth-form {
  margin-bottom: 24px;
}

.auth-form :deep(.ant-form-item-label > label) {
  font-weight: 500;
  color: var(--text-secondary);
}

.auth-form :deep(.ant-input-affix-wrapper) {
  background: var(--bg-tertiary) !important;
  border-color: var(--border-color) !important;
  border-radius: 14px !important;
  padding: 14px 16px;
  transition: all 0.2s ease;
}

.auth-form :deep(.ant-input-affix-wrapper:hover) {
  border-color: var(--color-primary) !important;
}

.auth-form :deep(.ant-input-affix-wrapper-focused) {
  border-color: var(--color-primary) !important;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1) !important;
}

.auth-form :deep(.ant-input-affix-wrapper input) {
  background: transparent !important;
  color: var(--text-primary);
}

.auth-form :deep(.ant-input-prefix) {
  color: var(--text-muted);
  margin-right: 12px;
}

.submit-btn {
  height: 52px !important;
  border-radius: 14px !important;
  font-size: 16px !important;
  font-weight: 600 !important;
  margin-top: 8px;
  box-shadow: 0 4px 14px rgba(99, 102, 241, 0.35) !important;
}

.submit-btn:hover {
  box-shadow: 0 6px 20px rgba(99, 102, 241, 0.45) !important;
}

.divider {
  display: flex;
  align-items: center;
  margin: 24px 0;
  color: var(--text-muted);
  font-size: 14px;
}

.divider::before,
.divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: var(--border-color);
}

.divider span {
  padding: 0 16px;
}

.oauth-btn {
  height: 52px !important;
  border-radius: 14px !important;
  background: var(--bg-tertiary) !important;
  border-color: var(--border-color) !important;
  color: var(--text-primary) !important;
  font-size: 15px !important;
  font-weight: 500 !important;
  transition: all 0.2s ease !important;
}

.oauth-btn:hover {
  border-color: var(--color-primary) !important;
  color: var(--color-primary) !important;
  background: var(--color-primary-light) !important;
}

.oauth-btn :deep(.anticon) {
  font-size: 18px;
}

.auth-footer {
  text-align: center;
  margin-top: 32px;
  font-size: 15px;
}

.footer-text {
  color: var(--text-secondary);
}

.footer-link {
  color: var(--color-primary);
  font-weight: 600;
  margin-left: 4px;
}

.footer-link:hover {
  text-decoration: underline;
}

/* Animation */
@keyframes float {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-20px);
  }
}

/* 响应式 */
@media (max-width: 480px) {
  .auth-container {
    padding: 36px 24px;
    border-radius: 24px;
  }
  
  .auth-title {
    font-size: 24px;
  }

  .theme-toggle {
    top: 16px;
    right: 16px;
  }
}
</style>
