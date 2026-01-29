<template>
  <div class="review-page">
    <header class="header">
      <div class="header-content">
        <router-link to="/" class="logo">
          <span class="logo-icon">🚀</span>
          <span class="logo-text">Linux.do</span>
          <span class="logo-badge">邀请码申请</span>
        </router-link>
        
        <nav class="nav">
          <router-link to="/posts" class="nav-link">申请列表</router-link>
          <router-link to="/review" class="nav-link active">二级审核</router-link>
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
              <span class="user-badge certified">认证用户</span>
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
      <div class="page-header slide-up">
        <div class="page-header-content">
          <div class="page-header-left">
            <h1 class="page-title">
              <SafetyCertificateOutlined class="title-icon" />
              二级审核
            </h1>
            <p class="page-subtitle">审核通过社区投票的申请，为合格的申请者提交邀请码</p>
          </div>
          <div class="header-stats">
            <div class="stat-item">
              <span class="stat-value">{{ total }}</span>
              <span class="stat-label">待审核</span>
            </div>
          </div>
        </div>
      </div>

      <div class="review-container">
        <a-spin :spinning="loading" tip="加载中...">
          <div v-if="posts.length > 0" class="posts-list">
            <div
              v-for="(post, index) in posts"
              :key="post.id"
              class="review-card fade-in"
              :style="{ animationDelay: `${index * 0.08}s` }"
            >
              <div class="card-header">
                <div class="applicant-info">
                  <a-avatar :size="48" class="applicant-avatar">
                    {{ (post.user?.username || 'U').charAt(0).toUpperCase() }}
                  </a-avatar>
                  <div class="applicant-details">
                    <span class="applicant-name">{{ post.user?.username || '匿名用户' }}</span>
                    <span class="applicant-email">{{ post.user?.email }}</span>
                  </div>
                </div>
                <div class="vote-stats">
                  <div class="vote-item up">
                    <LikeOutlined />
                    <span>{{ post.up_votes }}</span>
                  </div>
                  <div class="vote-item down">
                    <DislikeOutlined />
                    <span>{{ post.down_votes }}</span>
                  </div>
                  <div class="approval-rate" :class="getApprovalClass(post)">
                    {{ calculateApprovalRate(post) }}%
                  </div>
                </div>
              </div>

              <div class="card-content" @click="showDetail(post)">
                <h3 class="post-title">{{ post.title }}</h3>
                <p class="post-excerpt">{{ truncateContent(post.content) }}</p>
              </div>

              <div class="card-meta">
                <span class="meta-item">
                  <ClockCircleOutlined />
                  {{ formatTime(post.created_at) }}
                </span>
                <span class="meta-item">
                  <UserOutlined />
                  共 {{ post.up_votes + post.down_votes }} 票
                </span>
              </div>

              <div class="card-actions">
                <a-button 
                  type="primary" 
                  size="large" 
                  class="approve-btn"
                  @click="openApproveModal(post)"
                >
                  <template #icon><CheckOutlined /></template>
                  通过并发放邀请码
                </a-button>
                <a-button 
                  danger 
                  size="large" 
                  class="reject-btn"
                  @click="openRejectModal(post)"
                >
                  <template #icon><CloseOutlined /></template>
                  拒绝
                </a-button>
              </div>
            </div>
          </div>

          <a-empty v-else-if="!loading" description="暂无待审核的申请">
            <p class="empty-hint">所有通过社区投票的申请都会出现在这里</p>
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

    <!-- 帖子详情弹窗 -->
    <a-modal
      v-model:open="detailModalVisible"
      :title="null"
      :footer="null"
      width="800px"
      class="detail-modal"
      centered
    >
      <div v-if="selectedPost" class="detail-content">
        <div class="detail-header">
          <div class="detail-applicant">
            <a-avatar :size="56" class="applicant-avatar">
              {{ (selectedPost.user?.username || 'U').charAt(0).toUpperCase() }}
            </a-avatar>
            <div class="applicant-details">
              <span class="applicant-name">{{ selectedPost.user?.username || '匿名用户' }}</span>
              <span class="applicant-email">{{ selectedPost.user?.email }}</span>
              <span class="applicant-time">申请于 {{ formatTime(selectedPost.created_at) }}</span>
            </div>
          </div>
          <div class="detail-stats">
            <div class="stat-box">
              <span class="stat-number up">{{ selectedPost.up_votes }}</span>
              <span class="stat-label">赞成</span>
            </div>
            <div class="stat-box">
              <span class="stat-number down">{{ selectedPost.down_votes }}</span>
              <span class="stat-label">反对</span>
            </div>
            <div class="stat-box">
              <span class="stat-number" :class="getApprovalClass(selectedPost)">
                {{ calculateApprovalRate(selectedPost) }}%
              </span>
              <span class="stat-label">赞成率</span>
            </div>
          </div>
        </div>

        <h2 class="detail-title">{{ selectedPost.title }}</h2>
        <div class="detail-body" v-html="formatContent(selectedPost.content)"></div>

        <div class="detail-actions">
          <a-button 
            type="primary" 
            size="large" 
            class="approve-btn"
            @click="openApproveModal(selectedPost); detailModalVisible = false"
          >
            <template #icon><CheckOutlined /></template>
            通过并发放邀请码
          </a-button>
          <a-button 
            danger 
            size="large" 
            class="reject-btn"
            @click="openRejectModal(selectedPost); detailModalVisible = false"
          >
            <template #icon><CloseOutlined /></template>
            拒绝申请
          </a-button>
        </div>
      </div>
    </a-modal>

    <!-- 通过审核弹窗 -->
    <a-modal
      v-model:open="approveModalVisible"
      title="通过审核并发放邀请码"
      :ok-text="'确认发放'"
      :cancel-text="'取消'"
      :confirm-loading="approving"
      @ok="handleApprove"
      centered
      class="approve-modal"
    >
      <div class="approve-form">
        <a-alert
          message="请确认操作"
          description="通过审核后，邀请码将通过邮件发送给申请者。请确保邀请码有效且仅使用一次。"
          type="info"
          show-icon
          class="approve-alert"
        />
        
        <div v-if="targetPost" class="target-info">
          <span class="label">申请者：</span>
          <span class="value">{{ targetPost.user?.username }} ({{ targetPost.user?.email }})</span>
        </div>

        <a-form-item label="邀请码" :required="true" class="invite-code-item">
          <a-input
            v-model:value="inviteCode"
            placeholder="请输入要发放的 Linux.do 邀请码"
            size="large"
            class="invite-code-input"
          >
            <template #prefix>
              <GiftOutlined />
            </template>
          </a-input>
          <p class="input-hint">邀请码将发送至申请者的注册邮箱</p>
        </a-form-item>
      </div>
    </a-modal>

    <!-- 拒绝审核弹窗 -->
    <a-modal
      v-model:open="rejectModalVisible"
      title="拒绝申请"
      :ok-text="'确认拒绝'"
      :cancel-text="'取消'"
      :ok-button-props="{ danger: true }"
      :confirm-loading="rejecting"
      @ok="handleReject"
      centered
      class="reject-modal"
    >
      <div class="reject-form">
        <a-alert
          message="请谨慎操作"
          description="拒绝后申请者将无法获得邀请码。请在下方说明拒绝原因。"
          type="warning"
          show-icon
          class="reject-alert"
        />
        
        <div v-if="targetPost" class="target-info">
          <span class="label">申请者：</span>
          <span class="value">{{ targetPost.user?.username }}</span>
        </div>

        <a-form-item label="拒绝原因（可选）" class="reject-reason-item">
          <a-textarea
            v-model:value="rejectReason"
            placeholder="请输入拒绝原因，将通知申请者"
            :auto-size="{ minRows: 3, maxRows: 6 }"
            class="reject-reason-input"
          />
        </a-form-item>
      </div>
    </a-modal>

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
  SafetyCertificateOutlined,
  LikeOutlined,
  DislikeOutlined,
  ClockCircleOutlined,
  UserOutlined,
  CheckOutlined,
  CloseOutlined,
  GiftOutlined,
} from '@ant-design/icons-vue'
import { useUserStore } from '@/stores/user'
import { useThemeStore } from '@/stores/theme'
import { getReviewPosts, approvePost, rejectPost } from '@/api/post'
import type { Post } from '@/types'

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

const detailModalVisible = ref(false)
const selectedPost = ref<Post | null>(null)

const approveModalVisible = ref(false)
const rejectModalVisible = ref(false)
const targetPost = ref<Post | null>(null)
const inviteCode = ref('')
const rejectReason = ref('')
const approving = ref(false)
const rejecting = ref(false)

const fetchPosts = async () => {
  loading.value = true
  try {
    const response = await getReviewPosts({
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

const showDetail = (post: Post) => {
  selectedPost.value = post
  detailModalVisible.value = true
}

const openApproveModal = (post: Post) => {
  targetPost.value = post
  inviteCode.value = ''
  approveModalVisible.value = true
}

const openRejectModal = (post: Post) => {
  targetPost.value = post
  rejectReason.value = ''
  rejectModalVisible.value = true
}

const handleApprove = async () => {
  if (!inviteCode.value.trim()) {
    message.warning('请输入邀请码')
    return
  }

  if (!targetPost.value) return

  approving.value = true
  try {
    await approvePost(targetPost.value.id, { invite_code: inviteCode.value.trim() })
    message.success('审核通过，邀请码已发送至申请者邮箱')
    approveModalVisible.value = false
    // 从列表中移除
    posts.value = posts.value.filter(p => p.id !== targetPost.value?.id)
    total.value--
  } catch {
    // 错误已在拦截器中处理
  } finally {
    approving.value = false
  }
}

const handleReject = async () => {
  if (!targetPost.value) return

  rejecting.value = true
  try {
    await rejectPost(targetPost.value.id, rejectReason.value.trim() || undefined)
    message.success('已拒绝该申请')
    rejectModalVisible.value = false
    // 从列表中移除
    posts.value = posts.value.filter(p => p.id !== targetPost.value?.id)
    total.value--
  } catch {
    // 错误已在拦截器中处理
  } finally {
    rejecting.value = false
  }
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

const truncateContent = (content: string, maxLength = 200) => {
  if (content.length <= maxLength) return content
  return content.slice(0, maxLength) + '...'
}

const formatContent = (content: string) => {
  return content.replace(/\n/g, '<br>')
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
  // 检查是否有权限访问
  if (!userStore.isCertified) {
    message.warning('需要认证用户权限才能访问')
    router.push('/posts')
    return
  }
  fetchPosts()
})
</script>

<style scoped>
.review-page {
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
  max-width: 1400px;
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
  padding: 8px 12px;
  border-radius: 8px;
}

.nav-link:hover,
.nav-link.active {
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
  padding: 4px 12px;
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

.user-badge.certified {
  font-size: 11px;
  padding: 2px 6px;
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
  color: white;
  border-radius: 4px;
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
  max-width: 1200px;
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
  color: #f59e0b;
}

.page-subtitle {
  color: var(--text-secondary);
  font-size: 16px;
}

.header-stats {
  display: flex;
  gap: 20px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px 24px;
  background: var(--bg-card);
  border: 1px solid var(--border-color-light);
  border-radius: 16px;
  backdrop-filter: blur(20px);
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  color: #f59e0b;
}

.stat-label {
  font-size: 14px;
  color: var(--text-muted);
  margin-top: 4px;
}

/* Review Container */
.review-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
}

.posts-list {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.review-card {
  background: var(--bg-card);
  border: 1px solid var(--border-color-light);
  border-radius: 24px;
  padding: 28px;
  transition: transform 0.2s, box-shadow 0.2s, border-color 0.2s;
  backdrop-filter: blur(20px);
}

.review-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-lg);
  border-color: rgba(245, 158, 11, 0.4);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.applicant-info {
  display: flex;
  align-items: center;
  gap: 16px;
}

.applicant-avatar {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%) !important;
  color: white !important;
  font-weight: 600;
  font-size: 18px;
}

.applicant-details {
  display: flex;
  flex-direction: column;
}

.applicant-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.applicant-email {
  font-size: 13px;
  color: var(--text-muted);
  margin-top: 2px;
}

.vote-stats {
  display: flex;
  align-items: center;
  gap: 16px;
}

.vote-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 15px;
  font-weight: 600;
  padding: 8px 14px;
  border-radius: 10px;
  background: var(--bg-tertiary);
}

.vote-item.up {
  color: var(--color-success);
}

.vote-item.down {
  color: var(--color-error);
}

.approval-rate {
  font-size: 18px;
  font-weight: 700;
  padding: 8px 16px;
  border-radius: 10px;
  background: var(--bg-tertiary);
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

.card-content {
  cursor: pointer;
  margin-bottom: 16px;
}

.post-title {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 12px;
  line-height: 1.4;
}

.post-excerpt {
  color: var(--text-secondary);
  font-size: 15px;
  line-height: 1.7;
}

.card-meta {
  display: flex;
  gap: 24px;
  margin-bottom: 24px;
  padding-top: 16px;
  border-top: 1px solid var(--border-color);
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: var(--text-muted);
}

.card-actions {
  display: flex;
  gap: 16px;
}

.approve-btn {
  flex: 1;
  height: 48px !important;
  border-radius: 12px !important;
  font-weight: 600 !important;
  font-size: 15px !important;
  background: linear-gradient(135deg, var(--color-success) 0%, #059669 100%) !important;
  border: none !important;
}

.approve-btn:hover {
  background: linear-gradient(135deg, #059669 0%, #047857 100%) !important;
}

.reject-btn {
  flex: 1;
  height: 48px !important;
  border-radius: 12px !important;
  font-weight: 600 !important;
  font-size: 15px !important;
}

/* Detail Modal */
.detail-modal :deep(.ant-modal-content) {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 24px;
  overflow: hidden;
}

.detail-content {
  padding: 8px;
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 28px;
}

.detail-applicant {
  display: flex;
  align-items: center;
  gap: 16px;
}

.detail-applicant .applicant-details {
  gap: 4px;
}

.applicant-time {
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 4px;
}

.detail-stats {
  display: flex;
  gap: 20px;
}

.stat-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 12px 20px;
  background: var(--bg-tertiary);
  border-radius: 12px;
}

.stat-number {
  font-size: 24px;
  font-weight: 700;
}

.stat-number.up {
  color: var(--color-success);
}

.stat-number.down {
  color: var(--color-error);
}

.stat-box .stat-label {
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 4px;
}

.detail-title {
  font-size: 26px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 20px;
  line-height: 1.4;
}

.detail-body {
  color: var(--text-secondary);
  font-size: 15px;
  line-height: 1.9;
  margin-bottom: 28px;
  max-height: 400px;
  overflow-y: auto;
  padding-right: 8px;
}

.detail-actions {
  display: flex;
  gap: 16px;
  padding-top: 24px;
  border-top: 1px solid var(--border-color);
}

/* Approve Modal */
.approve-modal :deep(.ant-modal-content) {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 20px;
}

.approve-modal :deep(.ant-modal-header) {
  background: transparent;
  border-bottom: 1px solid var(--border-color);
}

.approve-modal :deep(.ant-modal-title) {
  color: var(--text-primary);
}

.approve-form {
  padding-top: 8px;
}

.approve-alert {
  margin-bottom: 20px;
  border-radius: 12px;
}

.target-info {
  padding: 12px 16px;
  background: var(--bg-tertiary);
  border-radius: 10px;
  margin-bottom: 20px;
}

.target-info .label {
  color: var(--text-muted);
  margin-right: 8px;
}

.target-info .value {
  color: var(--text-primary);
  font-weight: 500;
}

.invite-code-item {
  margin-bottom: 0;
}

.invite-code-item :deep(.ant-form-item-label > label) {
  color: var(--text-primary) !important;
  font-weight: 600;
}

.invite-code-input {
  height: 52px !important;
  border-radius: 12px !important;
  font-size: 16px !important;
}

.input-hint {
  font-size: 13px;
  color: var(--text-muted);
  margin-top: 8px;
}

/* Reject Modal */
.reject-modal :deep(.ant-modal-content) {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 20px;
}

.reject-modal :deep(.ant-modal-header) {
  background: transparent;
  border-bottom: 1px solid var(--border-color);
}

.reject-modal :deep(.ant-modal-title) {
  color: var(--text-primary);
}

.reject-form {
  padding-top: 8px;
}

.reject-alert {
  margin-bottom: 20px;
  border-radius: 12px;
}

.reject-reason-item :deep(.ant-form-item-label > label) {
  color: var(--text-primary) !important;
  font-weight: 600;
}

.reject-reason-input {
  border-radius: 12px !important;
}

/* Empty */
.empty-hint {
  color: var(--text-muted);
  font-size: 14px;
  margin-top: 8px;
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
    gap: 20px;
  }

  .page-title {
    font-size: 24px;
  }

  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .card-actions {
    flex-direction: column;
  }

  .detail-header {
    flex-direction: column;
    gap: 20px;
  }

  .detail-stats {
    width: 100%;
    justify-content: space-around;
  }

  .nav {
    gap: 12px;
  }

  .nav-link {
    display: none;
  }
}
</style>
