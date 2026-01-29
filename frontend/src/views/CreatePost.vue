<template>
  <div class="create-post-page">
    <header class="header">
      <div class="header-content">
        <router-link to="/" class="logo">
          <span class="logo-icon">🚀</span>
          <span class="logo-text">Linux.do</span>
          <span class="logo-badge">邀请码申请</span>
        </router-link>
        
        <nav class="nav">
          <router-link to="/posts" class="nav-link">申请列表</router-link>
          <router-link v-if="userStore.isCertified" to="/review" class="nav-link">二级审核</router-link>
          <router-link v-if="userStore.isAdmin" to="/admin" class="nav-link">管理后台</router-link>
          
          <!-- 主题切换按钮 -->
          <button class="theme-toggle" @click="toggleTheme" :title="themeStore.theme === 'light' ? '切换到暗色模式' : '切换到亮色模式'">
            <Transition name="theme-icon" mode="out-in">
              <span v-if="themeStore.theme === 'dark'" key="sun" class="theme-icon">☀️</span>
              <span v-else key="moon" class="theme-icon">🌙</span>
            </Transition>
          </button>
          
          <a-dropdown v-if="userStore.isLoggedIn">
            <div class="user-info">
              <a-avatar :size="32" class="user-avatar">
                {{ userStore.username.charAt(0).toUpperCase() }}
              </a-avatar>
              <span class="user-name">{{ userStore.username }}</span>
              <DownOutlined />
            </div>
            <template #overlay>
              <a-menu>
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
      <div class="create-container slide-up">
        <div class="create-header">
          <router-link to="/posts" class="back-link">
            <ArrowLeftOutlined />
            返回列表
          </router-link>
          <h1 class="create-title">
            <EditOutlined class="title-icon" />
            发布邀请码申请
          </h1>
          <p class="create-subtitle">
            写一篇真诚的小作文，说明你为什么想加入 Linux.do 社区
          </p>
        </div>

        <div class="form-card">
          <a-form
            :model="formState"
            :rules="rules"
            layout="vertical"
            @finish="handleSubmit"
            class="create-form"
          >
            <a-form-item name="title" label="申请标题">
              <a-input
                v-model:value="formState.title"
                placeholder="给你的申请起一个吸引人的标题"
                size="large"
                :maxlength="100"
                show-count
                class="title-input"
              />
            </a-form-item>

            <a-form-item name="content" label="申请内容">
              <div class="content-editor">
                <a-textarea
                  v-model:value="formState.content"
                  placeholder="在这里写下你的申请理由...

建议包含以下内容：
• 自我介绍
• 技术背景和专长
• 为什么想加入 Linux.do
• 你能为社区贡献什么"
                  :auto-size="{ minRows: 12, maxRows: 24 }"
                  :maxlength="5000"
                  show-count
                  class="content-textarea"
                />
              </div>
            </a-form-item>

            <div class="form-tips">
              <h4>
                <BulbOutlined />
                小贴士
              </h4>
              <ul>
                <li>真诚是最重要的，不要使用模板或套话</li>
                <li>展示你的技术热情和学习态度</li>
                <li>分享你对开源社区的理解和贡献意愿</li>
                <li>保持内容简洁，重点突出</li>
              </ul>
            </div>

            <div class="form-actions">
              <a-button size="large" @click="handleCancel" class="cancel-btn">
                取消
              </a-button>
              <a-button
                type="primary"
                html-type="submit"
                size="large"
                :loading="submitting"
                class="submit-btn"
              >
                <template #icon><SendOutlined /></template>
                提交申请
              </a-button>
            </div>
          </a-form>
        </div>

        <div class="process-info">
          <h3>审核流程说明</h3>
          <div class="process-steps">
            <div class="process-step">
              <div class="step-icon">
                <FormOutlined />
              </div>
              <div class="step-info">
                <h4>1. 提交申请</h4>
                <p>填写申请表单并提交</p>
              </div>
            </div>
            <div class="process-arrow">→</div>
            <div class="process-step">
              <div class="step-icon">
                <TeamOutlined />
              </div>
              <div class="step-info">
                <h4>2. 社区投票</h4>
                <p>其他用户对你的申请进行投票</p>
              </div>
            </div>
            <div class="process-arrow">→</div>
            <div class="process-step">
              <div class="step-icon">
                <SafetyCertificateOutlined />
              </div>
              <div class="step-info">
                <h4>3. 认证审核</h4>
                <p>认证用户进行最终审核</p>
              </div>
            </div>
            <div class="process-arrow">→</div>
            <div class="process-step">
              <div class="step-icon">
                <MailOutlined />
              </div>
              <div class="step-info">
                <h4>4. 获取邀请码</h4>
                <p>审核通过后邮件通知</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>

    <footer class="footer">
      <p>© 2024 Linux.do 邀请码申请系统 · 社区驱动的公平分发平台</p>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import {
  DownOutlined,
  FileTextOutlined,
  LogoutOutlined,
  ArrowLeftOutlined,
  EditOutlined,
  BulbOutlined,
  SendOutlined,
  FormOutlined,
  TeamOutlined,
  SafetyCertificateOutlined,
  MailOutlined,
} from '@ant-design/icons-vue'
import { useUserStore } from '@/stores/user'
import { useThemeStore } from '@/stores/theme'
import { createPost } from '@/api/post'
import type { Rule } from 'ant-design-vue/es/form'

const router = useRouter()
const userStore = useUserStore()
const themeStore = useThemeStore()

const toggleTheme = () => {
  themeStore.toggleTheme()
}

const formState = reactive({
  title: '',
  content: '',
})

const submitting = ref(false)

const rules: Record<string, Rule[]> = {
  title: [
    { required: true, message: '请输入申请标题', trigger: 'blur' },
    { min: 5, max: 100, message: '标题长度为5-100个字符', trigger: 'blur' },
  ],
  content: [
    { required: true, message: '请输入申请内容', trigger: 'blur' },
    { min: 50, max: 5000, message: '内容长度为50-5000个字符', trigger: 'blur' },
  ],
}

const handleSubmit = async () => {
  submitting.value = true
  try {
    await createPost({
      title: formState.title,
      content: formState.content,
    })
    message.success('申请提交成功！')
    router.push('/posts')
  } catch {
    // 错误已在拦截器中处理
  } finally {
    submitting.value = false
  }
}

const handleCancel = () => {
  if (formState.title || formState.content) {
    Modal.confirm({
      title: '确认离开？',
      content: '你有未保存的内容，确定要离开吗？',
      okText: '确定',
      cancelText: '取消',
      onOk: () => {
        router.push('/posts')
      },
    })
  } else {
    router.push('/posts')
  }
}

const handleLogout = () => {
  userStore.logout()
  message.success('已退出登录')
  router.push('/')
}
</script>

<style scoped>
.create-post-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
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
  gap: 8px;
  text-decoration: none;
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
  font-size: 12px;
  padding: 2px 8px;
  background: var(--color-primary);
  color: white;
  border-radius: 4px;
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
  padding: 4px 8px;
  border-radius: 8px;
  transition: background 0.2s;
}

.user-info:hover {
  background: var(--bg-tertiary);
}

.user-avatar {
  background: var(--color-primary) !important;
  color: white !important;
}

.user-name {
  color: var(--text-primary);
  font-weight: 500;
}

/* Main */
.main {
  flex: 1;
  padding: 72px 24px 40px;
}

.create-container {
  max-width: 800px;
  margin: 0 auto;
}

/* Header */
.create-header {
  padding: 40px 0 32px;
}

.back-link {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: var(--text-secondary);
  font-size: 14px;
  margin-bottom: 20px;
  transition: color 0.2s;
}

.back-link:hover {
  color: var(--color-primary);
}

.create-title {
  font-size: 36px;
  font-weight: 700;
  color: var(--text-primary);
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.title-icon {
  color: var(--color-primary);
}

.create-subtitle {
  color: var(--text-secondary);
  font-size: 16px;
}

/* Form Card */
.form-card {
  background: var(--bg-card);
  border: 1px solid var(--border-color-light);
  border-radius: 24px;
  padding: 40px;
  margin-bottom: 32px;
  backdrop-filter: blur(20px);
}

.create-form :deep(.ant-form-item-label > label) {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}

.title-input {
  height: 52px !important;
  border-radius: 12px !important;
  font-size: 16px !important;
}

.content-textarea {
  border-radius: 12px !important;
  font-size: 15px !important;
  line-height: 1.8 !important;
  padding: 16px !important;
}

.content-textarea:focus {
  box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.15) !important;
}

.form-tips {
  background: rgba(99, 102, 241, 0.08);
  border: 1px solid rgba(99, 102, 241, 0.15);
  border-radius: 16px;
  padding: 20px 24px;
  margin-bottom: 32px;
}

.form-tips h4 {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--color-primary);
  font-size: 15px;
  font-weight: 600;
  margin-bottom: 12px;
}

.form-tips ul {
  margin: 0;
  padding-left: 20px;
  color: var(--text-secondary);
  font-size: 14px;
  line-height: 1.8;
}

.form-tips li {
  margin-bottom: 4px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 16px;
  padding-top: 16px;
  border-top: 1px solid var(--border-color);
}

.cancel-btn {
  height: 48px !important;
  padding: 0 32px !important;
  border-radius: 12px !important;
  font-size: 15px !important;
  background: transparent !important;
  border-color: var(--border-color) !important;
  color: var(--text-secondary) !important;
}

.cancel-btn:hover {
  border-color: var(--text-muted) !important;
  color: var(--text-primary) !important;
}

.submit-btn {
  height: 48px !important;
  padding: 0 32px !important;
  border-radius: 12px !important;
  font-size: 15px !important;
  font-weight: 600 !important;
}

/* Process Info */
.process-info {
  background: var(--bg-card);
  border: 1px solid var(--border-color-light);
  border-radius: 20px;
  padding: 32px;
  backdrop-filter: blur(20px);
}

.process-info h3 {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 24px;
  text-align: center;
}

.process-steps {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-wrap: wrap;
  gap: 16px;
}

.process-step {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 20px;
  background: var(--bg-tertiary);
  border-radius: 12px;
  min-width: 180px;
}

.step-icon {
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-primary);
  border-radius: 10px;
  color: white;
  font-size: 20px;
  flex-shrink: 0;
}

.step-info h4 {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.step-info p {
  font-size: 12px;
  color: var(--text-muted);
  margin: 0;
}

.process-arrow {
  color: var(--text-muted);
  font-size: 20px;
}

/* Footer */
.footer {
  padding: 24px;
  text-align: center;
  border-top: 1px solid var(--border-color-light);
  color: var(--text-muted);
  font-size: 14px;
  background: var(--glass-bg);
  backdrop-filter: blur(10px);
}

/* Responsive */
@media (max-width: 768px) {
  .create-header {
    padding: 24px 0;
  }

  .create-title {
    font-size: 28px;
  }

  .form-card {
    padding: 24px;
    border-radius: 16px;
  }

  .form-actions {
    flex-direction: column;
  }

  .cancel-btn,
  .submit-btn {
    width: 100%;
  }

  .process-steps {
    flex-direction: column;
  }

  .process-arrow {
    transform: rotate(90deg);
  }

  .process-step {
    width: 100%;
    justify-content: flex-start;
  }
}
</style>
