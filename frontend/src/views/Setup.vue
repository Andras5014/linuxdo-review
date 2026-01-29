<template>
  <div class="setup-page">
    <div class="setup-background">
      <div class="grid-pattern"></div>
      <div class="glow-orb glow-orb-1"></div>
      <div class="glow-orb glow-orb-2"></div>
      <div class="glow-orb glow-orb-3"></div>
    </div>
    
    <div class="setup-container slide-up">
      <div class="setup-header">
        <div class="logo">
          <span class="logo-icon">⚙️</span>
          <span class="logo-text">系统初始化</span>
        </div>
        <h1 class="setup-title">创建管理员账号</h1>
        <p class="setup-subtitle">首次使用系统，请创建一个管理员账号</p>
      </div>

      <div class="setup-notice">
        <div class="notice-icon">💡</div>
        <div class="notice-text">
          <strong>重要提示</strong>
          <p>这是系统的首次初始化，您正在创建的账号将拥有最高管理权限。请妥善保管您的账号信息。</p>
        </div>
      </div>

      <a-form
        :model="formState"
        :rules="rules"
        @finish="handleSubmit"
        layout="vertical"
        class="setup-form"
      >
        <a-form-item name="username" label="用户名">
          <a-input
            v-model:value="formState.username"
            placeholder="请输入管理员用户名"
            size="large"
            :prefix="h(UserOutlined)"
          />
        </a-form-item>

        <a-form-item name="email" label="邮箱">
          <a-input
            v-model:value="formState.email"
            placeholder="请输入管理员邮箱"
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
            <template #icon>
              <RocketOutlined />
            </template>
            完成初始化
          </a-button>
        </a-form-item>
      </a-form>

      <div class="setup-footer">
        <span class="footer-text">初始化完成后，您可以使用此账号登录管理后台</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, h } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { UserOutlined, MailOutlined, LockOutlined, RocketOutlined } from '@ant-design/icons-vue'
import { setupAdmin } from '@/api/auth'
import { resetSystemInitializedCache } from '@/router'
import type { Rule } from 'ant-design-vue/es/form'

const router = useRouter()

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

const handleSubmit = async () => {
  loading.value = true
  try {
    await setupAdmin({
      username: formState.username,
      email: formState.email,
      password: formState.password,
    })
    
    // 重置系统初始化状态缓存
    resetSystemInitializedCache()
    
    message.success('初始化成功！请使用管理员账号登录')
    router.push('/login')
  } catch {
    // 错误已在拦截器中处理
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.setup-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  position: relative;
  overflow: hidden;
}

.setup-background {
  position: absolute;
  inset: 0;
  z-index: 0;
}

.grid-pattern {
  position: absolute;
  inset: 0;
  background-image: 
    linear-gradient(rgba(245, 158, 11, 0.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(245, 158, 11, 0.03) 1px, transparent 1px);
  background-size: 50px 50px;
}

.glow-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(100px);
  opacity: 0.4;
}

.glow-orb-1 {
  width: 400px;
  height: 400px;
  background: #f59e0b;
  top: -100px;
  left: 50%;
  transform: translateX(-50%);
}

.glow-orb-2 {
  width: 300px;
  height: 300px;
  background: #8b5cf6;
  bottom: -50px;
  left: -50px;
}

.glow-orb-3 {
  width: 250px;
  height: 250px;
  background: var(--color-primary);
  bottom: 100px;
  right: -50px;
}

.setup-container {
  width: 100%;
  max-width: 480px;
  background: var(--bg-card);
  backdrop-filter: blur(20px);
  border: 1px solid var(--border-color);
  border-radius: 24px;
  padding: 48px 40px;
  position: relative;
  z-index: 1;
  box-shadow: var(--shadow-lg);
}

.setup-header {
  text-align: center;
  margin-bottom: 32px;
}

.logo {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 24px;
}

.logo-icon {
  font-size: 36px;
}

.logo-text {
  font-size: 26px;
  font-weight: 700;
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.setup-title {
  font-size: 28px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.setup-subtitle {
  font-size: 15px;
  color: var(--text-secondary);
}

.setup-notice {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px;
  background: rgba(245, 158, 11, 0.1);
  border: 1px solid rgba(245, 158, 11, 0.2);
  border-radius: 12px;
  margin-bottom: 28px;
}

.notice-icon {
  font-size: 24px;
  flex-shrink: 0;
}

.notice-text strong {
  display: block;
  color: #f59e0b;
  font-size: 14px;
  margin-bottom: 4px;
}

.notice-text p {
  color: var(--text-secondary);
  font-size: 13px;
  margin: 0;
  line-height: 1.5;
}

.setup-form {
  margin-bottom: 24px;
}

.setup-form :deep(.ant-form-item-label > label) {
  font-weight: 500;
}

.setup-form :deep(.ant-input-affix-wrapper) {
  background: var(--bg-tertiary) !important;
  border-color: var(--border-color) !important;
  border-radius: 12px;
  padding: 12px 16px;
}

.setup-form :deep(.ant-input-affix-wrapper input) {
  background: transparent !important;
}

.setup-form :deep(.ant-input-prefix) {
  color: var(--text-muted);
  margin-right: 12px;
}

.submit-btn {
  height: 52px !important;
  border-radius: 12px !important;
  font-size: 16px !important;
  font-weight: 600 !important;
  margin-top: 8px;
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%) !important;
  border: none !important;
}

.submit-btn:hover {
  background: linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%) !important;
}

.setup-footer {
  text-align: center;
  margin-top: 24px;
}

.footer-text {
  color: var(--text-muted);
  font-size: 13px;
}

/* 动画 */
.slide-up {
  animation: slideUp 0.5s ease-out;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 响应式 */
@media (max-width: 480px) {
  .setup-container {
    padding: 36px 24px;
    border-radius: 20px;
  }
  
  .setup-title {
    font-size: 24px;
  }
  
  .setup-notice {
    flex-direction: column;
    align-items: center;
    text-align: center;
  }
}
</style>
