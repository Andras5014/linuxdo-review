<template>
  <div class="home-page">
    <header class="header">
      <div class="header-content">
        <div class="logo">
          <span class="logo-icon">🚀</span>
          <span class="logo-text">Linux.do</span>
          <span class="logo-badge">邀请码申请</span>
        </div>
        
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
          
          <template v-if="userStore.isLoggedIn">
            <a-dropdown>
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
          </template>
          <template v-else>
            <router-link to="/login" class="nav-link">登录</router-link>
            <router-link to="/register">
              <a-button type="primary" class="nav-btn">注册</a-button>
            </router-link>
          </template>
        </nav>
      </div>
    </header>

    <main class="main">
      <section class="hero">
        <div class="hero-background">
          <div class="hero-glow"></div>
          <div class="floating-shapes">
            <div class="shape shape-1"></div>
            <div class="shape shape-2"></div>
            <div class="shape shape-3"></div>
          </div>
        </div>
        
        <div class="hero-content slide-up">
          <h1 class="hero-title">
            <span class="gradient-text">Linux.do</span> 邀请码申请系统
          </h1>
          <p class="hero-description">
            通过社区投票和认证用户审核的两级机制，公平透明地分发邀请码
          </p>
          
          <!-- 装饰性卡片 -->
          <div class="hero-card float">
            <div class="card-header">
              <span class="card-icon">📊</span>
              <span class="card-title">申请统计</span>
            </div>
            <div class="card-stats">
              <div class="card-stat">
                <span class="card-stat-value">{{ stats.totalApplications }}</span>
                <span class="card-stat-label">总申请</span>
              </div>
              <div class="card-stat">
                <span class="card-stat-value success">{{ stats.approved }}</span>
                <span class="card-stat-label">已通过</span>
              </div>
              <div class="card-stat">
                <span class="card-stat-value pending">{{ stats.pending }}</span>
                <span class="card-stat-label">审核中</span>
              </div>
            </div>
            <div class="card-notification">
              <span class="notification-icon">✨</span>
              <span class="notification-text">+ 新申请</span>
            </div>
          </div>

          <div class="hero-actions">
            <router-link v-if="!userStore.isLoggedIn" to="/register">
              <a-button type="primary" size="large" class="action-btn primary">
                <template #icon><RocketOutlined /></template>
                立即开始
              </a-button>
            </router-link>
            <router-link v-else to="/posts/create">
              <a-button type="primary" size="large" class="action-btn primary">
                <template #icon><FormOutlined /></template>
                发布申请
              </a-button>
            </router-link>
            <router-link to="/posts">
              <a-button size="large" class="action-btn secondary">
                了解更多
              </a-button>
            </router-link>
          </div>

          <div class="hero-features">
            <div class="hero-feature">
              <ThunderboltOutlined />
              <span>极速审核</span>
            </div>
            <div class="hero-feature">
              <GlobalOutlined />
              <span>全球覆盖</span>
            </div>
            <div class="hero-feature">
              <SafetyOutlined />
              <span>安全加密</span>
            </div>
          </div>
        </div>
      </section>

      <section class="features">
        <div class="features-grid">
          <div class="feature-card fade-in" style="animation-delay: 0.1s">
            <div class="feature-icon">
              <TeamOutlined />
            </div>
            <h3 class="feature-title">社区投票</h3>
            <p class="feature-desc">所有注册用户都可以参与投票，共同决定申请是否通过初审</p>
          </div>
          
          <div class="feature-card fade-in" style="animation-delay: 0.2s">
            <div class="feature-icon">
              <SafetyCertificateOutlined />
            </div>
            <h3 class="feature-title">认证审核</h3>
            <p class="feature-desc">通过初审的申请由 Linux.do 认证用户进行二级审核</p>
          </div>
          
          <div class="feature-card fade-in" style="animation-delay: 0.3s">
            <div class="feature-icon">
              <MailOutlined />
            </div>
            <h3 class="feature-title">邮件通知</h3>
            <p class="feature-desc">审核通过后自动发送邀请码到您的注册邮箱</p>
          </div>
        </div>
      </section>

      <section class="how-it-works">
        <h2 class="section-title">申请流程</h2>
        <div class="steps">
          <div class="step">
            <div class="step-number">1</div>
            <div class="step-content">
              <h4>提交申请</h4>
              <p>撰写小作文说明申请理由</p>
            </div>
          </div>
          <div class="step-arrow">→</div>
          <div class="step">
            <div class="step-number">2</div>
            <div class="step-content">
              <h4>社区投票</h4>
              <p>获得足够赞成票进入下一轮</p>
            </div>
          </div>
          <div class="step-arrow">→</div>
          <div class="step">
            <div class="step-number">3</div>
            <div class="step-content">
              <h4>认证审核</h4>
              <p>认证用户提交邀请码</p>
            </div>
          </div>
          <div class="step-arrow">→</div>
          <div class="step">
            <div class="step-number">4</div>
            <div class="step-content">
              <h4>获取邀请码</h4>
              <p>邮件收到邀请码完成注册</p>
            </div>
          </div>
        </div>
      </section>
    </main>

    <footer class="footer">
      <p>© 2024 Linux.do 邀请码申请系统 · 社区驱动的公平分发平台</p>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import {
  DownOutlined,
  FileTextOutlined,
  LogoutOutlined,
  RocketOutlined,
  FormOutlined,
  TeamOutlined,
  SafetyCertificateOutlined,
  MailOutlined,
  ThunderboltOutlined,
  GlobalOutlined,
  SafetyOutlined,
} from '@ant-design/icons-vue'
import { useUserStore } from '@/stores/user'
import { useThemeStore } from '@/stores/theme'

const router = useRouter()
const userStore = useUserStore()
const themeStore = useThemeStore()

// 统计数据（后续从API获取）
const stats = reactive({
  totalApplications: 128,
  approved: 45,
  pending: 23,
})

const toggleTheme = () => {
  themeStore.toggleTheme()
}

const handleLogout = () => {
  userStore.logout()
  message.success('已退出登录')
  router.push('/')
}
</script>

<style scoped>
.home-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
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

.nav-btn {
  border-radius: 10px !important;
  font-weight: 500 !important;
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
  flex: 1;
  padding-top: 72px;
}

/* Hero */
.hero {
  position: relative;
  min-height: calc(100vh - 72px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 60px 24px;
  overflow: hidden;
}

.hero-background {
  position: absolute;
  inset: 0;
  overflow: hidden;
}

.hero-glow {
  position: absolute;
  top: -200px;
  left: 50%;
  transform: translateX(-50%);
  width: 1000px;
  height: 600px;
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
  opacity: 0.5;
  filter: blur(80px);
}

.shape-1 {
  width: 400px;
  height: 400px;
  background: rgba(99, 102, 241, 0.2);
  top: 10%;
  right: 10%;
  animation: float 8s ease-in-out infinite;
}

.shape-2 {
  width: 300px;
  height: 300px;
  background: rgba(168, 85, 247, 0.15);
  bottom: 20%;
  left: 5%;
  animation: float 10s ease-in-out infinite reverse;
}

.shape-3 {
  width: 200px;
  height: 200px;
  background: rgba(59, 130, 246, 0.15);
  top: 40%;
  left: 30%;
  animation: float 12s ease-in-out infinite;
}

.hero-content {
  position: relative;
  text-align: center;
  max-width: 800px;
}

.hero-title {
  font-size: 52px;
  font-weight: 700;
  line-height: 1.2;
  margin-bottom: 20px;
  color: var(--text-primary);
}

.hero-description {
  font-size: 18px;
  color: var(--text-secondary);
  margin-bottom: 40px;
  line-height: 1.7;
}

/* Hero Card - 类似图片中的装饰卡片 */
.hero-card {
  position: absolute;
  right: -280px;
  top: 50%;
  transform: translateY(-50%);
  width: 260px;
  background: var(--bg-card);
  backdrop-filter: blur(20px);
  border: 1px solid var(--border-color-light);
  border-radius: 20px;
  padding: 20px;
  box-shadow: var(--shadow-card);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--border-color-light);
}

.card-icon {
  font-size: 20px;
}

.card-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.card-stats {
  display: flex;
  justify-content: space-between;
  margin-bottom: 16px;
}

.card-stat {
  text-align: center;
}

.card-stat-value {
  display: block;
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  font-family: var(--font-mono);
}

.card-stat-value.success {
  color: var(--color-success);
}

.card-stat-value.pending {
  color: var(--color-warning);
}

.card-stat-label {
  font-size: 11px;
  color: var(--text-muted);
}

.card-notification {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  background: var(--color-success-light);
  border-radius: 10px;
  font-size: 13px;
  font-weight: 500;
  color: var(--color-success);
}

.notification-icon {
  font-size: 14px;
}

.hero-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-bottom: 40px;
}

.action-btn {
  height: 50px !important;
  padding: 0 32px !important;
  border-radius: 14px !important;
  font-size: 15px !important;
  font-weight: 600 !important;
}

.action-btn.primary {
  box-shadow: 0 4px 14px rgba(99, 102, 241, 0.35) !important;
}

.action-btn.primary:hover {
  box-shadow: 0 6px 20px rgba(99, 102, 241, 0.45) !important;
  transform: translateY(-1px);
}

.action-btn.secondary {
  background: var(--bg-secondary) !important;
  border: 1px solid var(--border-color) !important;
  color: var(--text-primary) !important;
}

.action-btn.secondary:hover {
  border-color: var(--color-primary) !important;
  color: var(--color-primary) !important;
}

/* Hero Features */
.hero-features {
  display: flex;
  justify-content: center;
  gap: 40px;
}

.hero-feature {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text-muted);
  font-size: 14px;
}

.hero-feature :deep(.anticon) {
  color: var(--color-primary);
  font-size: 16px;
}

/* Features */
.features {
  padding: 80px 24px;
  background: var(--bg-secondary);
}

.features-grid {
  max-width: 1000px;
  margin: 0 auto;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 28px;
}

.feature-card {
  background: var(--bg-card);
  border: 1px solid var(--border-color-light);
  border-radius: 20px;
  padding: 32px;
  text-align: center;
  transition: all 0.3s ease;
  backdrop-filter: blur(20px);
}

.feature-card:hover {
  transform: translateY(-6px);
  box-shadow: var(--shadow-lg);
  border-color: var(--color-primary);
}

.feature-icon {
  width: 64px;
  height: 64px;
  margin: 0 auto 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-primary-gradient);
  border-radius: 18px;
  font-size: 26px;
  color: white;
  box-shadow: 0 8px 24px rgba(99, 102, 241, 0.25);
}

.feature-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 12px;
}

.feature-desc {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.7;
}

/* How it works */
.how-it-works {
  padding: 80px 24px;
}

.section-title {
  text-align: center;
  font-size: 32px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 50px;
}

.steps {
  max-width: 1000px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-wrap: wrap;
  gap: 16px;
}

.step {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px 24px;
  background: var(--bg-card);
  border: 1px solid var(--border-color-light);
  border-radius: 16px;
  backdrop-filter: blur(20px);
  transition: all 0.3s ease;
}

.step:hover {
  box-shadow: var(--shadow-md);
  border-color: var(--color-primary);
}

.step-number {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-primary-gradient);
  color: white;
  font-size: 18px;
  font-weight: 700;
  border-radius: 14px;
  flex-shrink: 0;
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.3);
}

.step-content h4 {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.step-content p {
  font-size: 13px;
  color: var(--text-secondary);
}

.step-arrow {
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
@media (max-width: 1200px) {
  .hero-card {
    display: none;
  }
}

@media (max-width: 768px) {
  .hero-title {
    font-size: 32px;
  }
  
  .hero-description {
    font-size: 15px;
  }
  
  .hero-actions {
    flex-direction: column;
    width: 100%;
    max-width: 300px;
    margin-left: auto;
    margin-right: auto;
  }
  
  .action-btn {
    width: 100% !important;
  }
  
  .hero-features {
    flex-direction: column;
    gap: 16px;
  }
  
  .features-grid {
    grid-template-columns: 1fr;
  }
  
  .steps {
    flex-direction: column;
  }
  
  .step-arrow {
    transform: rotate(90deg);
  }

  .nav {
    gap: 12px;
  }

  .nav-link {
    display: none;
  }
}
</style>
