<template>
  <div class="my-posts-page">
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
                <a-menu-item key="my-posts" class="active-menu">
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
      <div class="page-header slide-up">
        <div class="page-header-content">
          <div class="page-header-left">
            <h1 class="page-title">
              <FileTextOutlined class="title-icon" />
              我的申请
            </h1>
            <p class="page-subtitle">查看你提交的所有邀请码申请及状态</p>
          </div>
          <router-link to="/posts/create">
            <a-button type="primary" size="large" class="create-btn">
              <template #icon><PlusOutlined /></template>
              发布新申请
            </a-button>
          </router-link>
        </div>
      </div>

      <div class="posts-container">
        <a-spin :spinning="loading" tip="加载中...">
          <div v-if="posts.length > 0" class="posts-list">
            <div
              v-for="(post, index) in posts"
              :key="post.id"
              class="post-card fade-in"
              :style="{ animationDelay: `${index * 0.05}s` }"
            >
              <div class="post-main">
                <div class="post-status-indicator" :class="getStatusClass(post.status)"></div>
                <div class="post-content">
                  <div class="post-header">
                    <h3 class="post-title">{{ post.title }}</h3>
                    <a-tag :color="getStatusColor(post.status)" class="status-tag">
                      {{ getStatusText(post.status) }}
                    </a-tag>
                  </div>
                  <p class="post-excerpt">{{ truncateContent(post.content) }}</p>
                  <div class="post-meta">
                    <span class="post-time">
                      <ClockCircleOutlined />
                      {{ formatTime(post.created_at) }}
                    </span>
                    <span class="post-votes">
                      <LikeOutlined />
                      {{ post.up_votes }} 赞成
                    </span>
                    <span class="post-votes">
                      <DislikeOutlined />
                      {{ post.down_votes }} 反对
                    </span>
                    <span class="approval-rate" :class="getApprovalClass(post)">
                      {{ calculateApprovalRate(post) }}% 赞成率
                    </span>
                  </div>
                </div>
              </div>

              <div v-if="post.status === PostStatus.Approved" class="post-result success">
                <CheckCircleOutlined class="result-icon" />
                <div class="result-info">
                  <h4>申请已通过</h4>
                  <p>邀请码已发送至您的注册邮箱，请查收</p>
                </div>
              </div>

              <div v-else-if="post.status === PostStatus.Rejected" class="post-result rejected">
                <CloseCircleOutlined class="result-icon" />
                <div class="result-info">
                  <h4>申请未通过</h4>
                  <p>很遗憾，您的申请未能通过审核</p>
                </div>
              </div>

              <div v-else-if="post.status === PostStatus.SecondReview" class="post-result pending">
                <SafetyCertificateOutlined class="result-icon" />
                <div class="result-info">
                  <h4>等待二级审核</h4>
                  <p>您的申请已通过社区投票，正在等待认证用户审核</p>
                </div>
              </div>

              <div v-else-if="post.status === PostStatus.FirstReview" class="post-result voting">
                <TeamOutlined class="result-icon" />
                <div class="result-info">
                  <h4>社区投票中</h4>
                  <p>等待更多用户投票，达到阈值后进入下一轮</p>
                </div>
              </div>
            </div>
          </div>

          <a-empty v-else-if="!loading" description="你还没有提交过申请">
            <router-link to="/posts/create">
              <a-button type="primary">提交第一个申请</a-button>
            </router-link>
          </a-empty>
        </a-spin>

        <div v-if="total > pageSize" class="pagination-container">
          <a-pagination
            v-model:current="currentPage"
            :total="total"
            :page-size="pageSize"
            show-quick-jumper
            @change="handlePageChange"
          />
        </div>
      </div>
    </main>

    <footer class="footer">
      <p>© 2024 Linux.do 邀请码申请系统 · 社区驱动的公平分发平台</p>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import {
  DownOutlined,
  FileTextOutlined,
  LogoutOutlined,
  PlusOutlined,
  ClockCircleOutlined,
  LikeOutlined,
  DislikeOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  SafetyCertificateOutlined,
  TeamOutlined,
} from '@ant-design/icons-vue'
import { useUserStore } from '@/stores/user'
import { useThemeStore } from '@/stores/theme'
import { getMyPosts } from '@/api/post'
import type { Post } from '@/types'
import { PostStatus } from '@/types'

const router = useRouter()
const userStore = useUserStore()
const themeStore = useThemeStore()

const toggleTheme = () => {
  themeStore.toggleTheme()
}

const loading = ref(false)
const posts = ref<Post[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const fetchPosts = async () => {
  loading.value = true
  try {
    const response = await getMyPosts({
      page: currentPage.value,
      page_size: pageSize.value,
    })
    const data = response.data.data
    posts.value = data.list || []
    total.value = data.total
  } catch {
    // 错误已在拦截器中处理
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page: number) => {
  currentPage.value = page
  fetchPosts()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const handleLogout = () => {
  userStore.logout()
  message.success('已退出登录')
  router.push('/')
}

const formatTime = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

const truncateContent = (content: string, maxLength = 150) => {
  if (content.length <= maxLength) return content
  return content.slice(0, maxLength) + '...'
}

const getStatusText = (status: PostStatus) => {
  const statusMap: Record<PostStatus, string> = {
    [PostStatus.Pending]: '待审核',
    [PostStatus.FirstReview]: '投票中',
    [PostStatus.SecondReview]: '二级审核',
    [PostStatus.Approved]: '已通过',
    [PostStatus.Rejected]: '已拒绝',
  }
  return statusMap[status] || '未知'
}

const getStatusColor = (status: PostStatus) => {
  const colorMap: Record<PostStatus, string> = {
    [PostStatus.Pending]: 'default',
    [PostStatus.FirstReview]: 'processing',
    [PostStatus.SecondReview]: 'warning',
    [PostStatus.Approved]: 'success',
    [PostStatus.Rejected]: 'error',
  }
  return colorMap[status] || 'default'
}

const getStatusClass = (status: PostStatus) => {
  const classMap: Record<PostStatus, string> = {
    [PostStatus.Pending]: 'status-pending',
    [PostStatus.FirstReview]: 'status-voting',
    [PostStatus.SecondReview]: 'status-review',
    [PostStatus.Approved]: 'status-approved',
    [PostStatus.Rejected]: 'status-rejected',
  }
  return classMap[status] || ''
}

const calculateApprovalRate = (post: Post) => {
  const total = post.up_votes + post.down_votes
  if (total === 0) return 0
  return Math.round((post.up_votes / total) * 100)
}

const getApprovalClass = (post: Post) => {
  const rate = calculateApprovalRate(post)
  if (rate >= 70) return 'rate-high'
  if (rate >= 50) return 'rate-medium'
  return 'rate-low'
}

onMounted(() => {
  fetchPosts()
})
</script>

<style scoped>
.my-posts-page {
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
  padding-top: 72px;
  padding-bottom: 40px;
}

/* Page Header */
.page-header {
  max-width: 1000px;
  margin: 0 auto;
  padding: 40px 24px 24px;
}

.page-header-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.page-title {
  font-size: 32px;
  font-weight: 700;
  color: var(--text-primary);
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.title-icon {
  color: var(--color-primary);
}

.page-subtitle {
  color: var(--text-secondary);
  font-size: 16px;
}

.create-btn {
  height: 44px !important;
  padding: 0 24px !important;
  border-radius: 12px !important;
  font-weight: 600 !important;
}

/* Posts List */
.posts-container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 0 24px;
}

.posts-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.post-card {
  background: var(--bg-card);
  border: 1px solid var(--border-color-light);
  border-radius: 20px;
  overflow: hidden;
  transition: transform 0.2s, box-shadow 0.2s;
  backdrop-filter: blur(20px);
}

.post-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
}

.post-main {
  display: flex;
  padding: 24px;
}

.post-status-indicator {
  width: 4px;
  border-radius: 4px;
  flex-shrink: 0;
  margin-right: 20px;
}

.status-pending {
  background: var(--text-muted);
}

.status-voting {
  background: var(--color-info);
}

.status-review {
  background: var(--color-warning);
}

.status-approved {
  background: var(--color-success);
}

.status-rejected {
  background: var(--color-error);
}

.post-content {
  flex: 1;
}

.post-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}

.post-title {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  line-height: 1.4;
  flex: 1;
  margin-right: 16px;
}

.status-tag {
  border-radius: 12px;
  padding: 2px 10px;
  font-size: 12px;
  font-weight: 500;
  flex-shrink: 0;
}

.post-excerpt {
  color: var(--text-secondary);
  font-size: 14px;
  line-height: 1.6;
  margin-bottom: 16px;
}

.post-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  font-size: 13px;
  color: var(--text-muted);
}

.post-meta span {
  display: flex;
  align-items: center;
  gap: 6px;
}

.approval-rate {
  font-weight: 600;
}

.rate-high {
  color: var(--color-success);
}

.rate-medium {
  color: var(--color-warning);
}

.rate-low {
  color: var(--color-error);
}

/* Post Result */
.post-result {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px 24px;
  border-top: 1px solid var(--border-color);
}

.post-result.success {
  background: rgba(16, 185, 129, 0.1);
}

.post-result.rejected {
  background: rgba(239, 68, 68, 0.1);
}

.post-result.pending {
  background: rgba(245, 158, 11, 0.1);
}

.post-result.voting {
  background: rgba(59, 130, 246, 0.1);
}

.result-icon {
  font-size: 32px;
  flex-shrink: 0;
}

.post-result.success .result-icon {
  color: var(--color-success);
}

.post-result.rejected .result-icon {
  color: var(--color-error);
}

.post-result.pending .result-icon {
  color: var(--color-warning);
}

.post-result.voting .result-icon {
  color: var(--color-info);
}

.result-info h4 {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.result-info p {
  font-size: 13px;
  color: var(--text-secondary);
  margin: 0;
}

/* Pagination */
.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 40px;
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
  .page-header-content {
    flex-direction: column;
    gap: 16px;
  }

  .page-title {
    font-size: 24px;
  }

  .post-main {
    padding: 20px;
  }

  .post-meta {
    flex-direction: column;
    gap: 8px;
  }

  .post-result {
    flex-direction: column;
    text-align: center;
  }
}
</style>
