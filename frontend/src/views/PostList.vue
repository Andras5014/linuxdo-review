<template>
  <div class="post-list-page">
    <header class="header">
      <div class="header-content">
        <router-link to="/" class="logo">
          <span class="logo-icon">🚀</span>
          <span class="logo-text">Linux.do</span>
          <span class="logo-badge">邀请码申请</span>
        </router-link>
        
        <nav class="nav">
          <router-link to="/posts" class="nav-link active">申请列表</router-link>
          <router-link v-if="userStore.canReview" to="/review" class="nav-link">二级审核</router-link>
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
      <div class="page-header slide-up">
        <div class="page-header-content">
          <div class="page-header-left">
            <h1 class="page-title">
              <FireOutlined class="title-icon" />
              社区投票区
            </h1>
            <p class="page-subtitle">查看所有邀请码申请，参与投票决定谁能获得邀请码</p>
          </div>
          <router-link v-if="userStore.isLoggedIn" to="/posts/create">
            <a-button type="primary" size="large" class="create-btn">
              <template #icon><PlusOutlined /></template>
              发布申请
            </a-button>
          </router-link>
        </div>

        <div class="filter-bar">
          <a-radio-group v-model:value="statusFilter" button-style="solid" @change="handleFilterChange">
            <a-radio-button value="all">全部</a-radio-button>
            <a-radio-button :value="1">投票中</a-radio-button>
            <a-radio-button :value="2">待审核</a-radio-button>
            <a-radio-button :value="3">已通过</a-radio-button>
            <a-radio-button :value="4">已拒绝</a-radio-button>
          </a-radio-group>
        </div>
      </div>

      <div class="posts-container">
        <a-spin :spinning="loading" tip="加载中...">
          <div v-if="posts.length > 0" class="posts-grid">
            <div
              v-for="(post, index) in posts"
              :key="post.id"
              class="post-card fade-in"
              :style="{ animationDelay: `${index * 0.05}s` }"
            >
              <div class="post-header">
                <div class="post-author">
                  <a-avatar :size="40" class="author-avatar">
                    {{ (post.user?.username || 'U').charAt(0).toUpperCase() }}
                  </a-avatar>
                  <div class="author-info">
                    <span class="author-name">{{ post.user?.username || '匿名用户' }}</span>
                    <span class="post-time">{{ formatTime(post.created_at) }}</span>
                  </div>
                </div>
                <a-tag :color="getStatusColor(post.status)" class="status-tag">
                  {{ getStatusText(post.status) }}
                </a-tag>
              </div>

              <div class="post-content" @click="showPostDetail(post)">
                <h3 class="post-title">{{ post.title }}</h3>
                <p class="post-excerpt">{{ truncateContent(post.content) }}</p>
              </div>

              <div class="post-footer">
                <div class="vote-section">
                  <a-tooltip :title="getVoteTooltip(1)">
                    <a-button
                      :type="post.userVote === 1 ? 'primary' : 'default'"
                      :class="{ 'vote-up-active': post.userVote === 1 }"
                      shape="round"
                      size="small"
                      :disabled="!canVote(post)"
                      @click="handleVote(post.id, 1)"
                      class="vote-btn up"
                    >
                      <LikeOutlined />
                      <span class="vote-count">{{ post.up_votes }}</span>
                    </a-button>
                  </a-tooltip>
                  <a-tooltip :title="getVoteTooltip(-1)">
                    <a-button
                      :type="post.userVote === -1 ? 'primary' : 'default'"
                      :class="{ 'vote-down-active': post.userVote === -1 }"
                      shape="round"
                      size="small"
                      :disabled="!canVote(post)"
                      @click="handleVote(post.id, -1)"
                      class="vote-btn down"
                    >
                      <DislikeOutlined />
                      <span class="vote-count">{{ post.down_votes }}</span>
                    </a-button>
                  </a-tooltip>
                </div>
                
                <div class="post-stats">
                  <span class="approval-rate" :class="getApprovalClass(post)">
                    {{ calculateApprovalRate(post) }}% 赞成率
                  </span>
                </div>
              </div>
            </div>
          </div>

          <a-empty v-else-if="!loading" description="暂无申请">
            <router-link v-if="userStore.isLoggedIn" to="/posts/create">
              <a-button type="primary">成为第一个申请者</a-button>
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

    <!-- 帖子详情弹窗 -->
    <a-modal
      v-model:open="detailModalVisible"
      :title="null"
      :footer="null"
      width="700px"
      class="post-detail-modal"
      centered
    >
      <div v-if="selectedPost" class="detail-content">
        <div class="detail-header">
          <div class="detail-author">
            <a-avatar :size="48" class="author-avatar">
              {{ (selectedPost.user?.username || 'U').charAt(0).toUpperCase() }}
            </a-avatar>
            <div class="author-info">
              <span class="author-name">{{ selectedPost.user?.username || '匿名用户' }}</span>
              <span class="post-time">{{ formatTime(selectedPost.created_at) }}</span>
            </div>
          </div>
          <a-tag :color="getStatusColor(selectedPost.status)" class="status-tag">
            {{ getStatusText(selectedPost.status) }}
          </a-tag>
        </div>

        <h2 class="detail-title">{{ selectedPost.title }}</h2>
        <div class="detail-body" v-html="formatContent(selectedPost.content)"></div>

        <div class="detail-footer">
          <div class="vote-section">
            <a-button
              :type="selectedPost.userVote === 1 ? 'primary' : 'default'"
              :class="{ 'vote-up-active': selectedPost.userVote === 1 }"
              size="large"
              :disabled="!canVote(selectedPost)"
              @click="handleVote(selectedPost.id, 1)"
              class="vote-btn up"
            >
              <LikeOutlined />
              <span>赞成 {{ selectedPost.up_votes }}</span>
            </a-button>
            <a-button
              :type="selectedPost.userVote === -1 ? 'primary' : 'default'"
              :class="{ 'vote-down-active': selectedPost.userVote === -1 }"
              size="large"
              :disabled="!canVote(selectedPost)"
              @click="handleVote(selectedPost.id, -1)"
              class="vote-btn down"
            >
              <DislikeOutlined />
              <span>反对 {{ selectedPost.down_votes }}</span>
            </a-button>
          </div>
          <div class="approval-info">
            <span class="approval-rate" :class="getApprovalClass(selectedPost)">
              {{ calculateApprovalRate(selectedPost) }}% 赞成率
            </span>
            <span class="vote-total">共 {{ selectedPost.up_votes + selectedPost.down_votes }} 票</span>
          </div>
        </div>
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
  FireOutlined,
  PlusOutlined,
  LikeOutlined,
  DislikeOutlined,
  UserOutlined,
} from '@ant-design/icons-vue'
import { useUserStore } from '@/stores/user'
import { useThemeStore } from '@/stores/theme'
import { getPosts, votePost } from '@/api/post'
import type { Post } from '@/types'
import { PostStatus, VoteType } from '@/types'

interface PostWithVote extends Post {
  userVote?: number
}

const router = useRouter()
const userStore = useUserStore()
const themeStore = useThemeStore()

const loading = ref(false)
const posts = ref<PostWithVote[]>([])
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)
const statusFilter = ref<string | number>('all')

const detailModalVisible = ref(false)
const selectedPost = ref<PostWithVote | null>(null)

const toggleTheme = () => {
  themeStore.toggleTheme()
}

const fetchPosts = async () => {
  loading.value = true
  try {
    const params: { page: number; page_size: number; status?: number } = {
      page: currentPage.value,
      page_size: pageSize.value,
    }
    // 只有选择具体状态时才传 status 参数，"全部"时不传参数让后端返回所有帖子
    if (statusFilter.value !== 'all' && typeof statusFilter.value === 'number') {
      params.status = statusFilter.value
    }
    const response = await getPosts(params)
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

const handleFilterChange = () => {
  currentPage.value = 1
  fetchPosts()
}

// 检查是否可以投票
const canVote = (post: PostWithVote) => {
  if (!userStore.isLoggedIn) return false
  if (!userStore.isLinuxDoBound) return false
  if (post.status !== 1) return false
  return true
}

// 获取投票按钮的提示信息
const getVoteTooltip = (voteType: number) => {
  if (!userStore.isLoggedIn) {
    return '请先登录'
  }
  if (!userStore.isLinuxDoBound) {
    return '请先绑定 Linux.do 账号后再投票'
  }
  return voteType === 1 ? '赞成' : '反对'
}

const handleVote = async (postId: number, voteType: number) => {
  if (!userStore.isLoggedIn) {
    message.warning('请先登录')
    router.push('/login')
    return
  }

  if (!userStore.isLinuxDoBound) {
    message.warning('请先绑定 Linux.do 账号后再投票')
    router.push('/profile')
    return
  }

  try {
    const response = await votePost(postId, { vote_type: voteType as VoteType })
    const voteResponse = response.data.data
    message.success(voteResponse.message || '投票成功')
    
    // 使用后端返回的数据更新本地状态
    const post = posts.value.find(p => p.id === postId)
    if (post) {
      post.up_votes = voteResponse.up_votes
      post.down_votes = voteResponse.down_votes
      post.userVote = voteResponse.vote_type || 0
    }

    // 更新详情弹窗中的状态
    if (selectedPost.value?.id === postId) {
      selectedPost.value = { ...post! }
    }
  } catch {
    // 错误已在拦截器中处理
  }
}

const showPostDetail = (post: PostWithVote) => {
  selectedPost.value = post
  detailModalVisible.value = true
}

const handleLogout = () => {
  userStore.logout()
  message.success('已退出登录')
  router.push('/')
}

const formatTime = (dateStr: string) => {
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)

  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes} 分钟前`
  if (hours < 24) return `${hours} 小时前`
  if (days < 30) return `${days} 天前`
  return date.toLocaleDateString('zh-CN')
}

const truncateContent = (content: string, maxLength = 120) => {
  if (content.length <= maxLength) return content
  return content.slice(0, maxLength) + '...'
}

const formatContent = (content: string) => {
  return content.replace(/\n/g, '<br>')
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
.post-list-page {
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
  gap: 10px;
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
  padding: 8px 12px;
  border-radius: 8px;
  font-size: 14px;
}

.nav-link:hover,
.nav-link.active {
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
  padding-bottom: 40px;
}

/* Page Header */
.page-header {
  max-width: 1400px;
  margin: 0 auto;
  padding: 40px 24px 24px;
}

.page-header-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
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
  color: #f97316;
}

.page-subtitle {
  color: var(--text-secondary);
  font-size: 16px;
}

.create-btn {
  height: 48px !important;
  padding: 0 28px !important;
  border-radius: 14px !important;
  font-weight: 600 !important;
  box-shadow: 0 4px 14px rgba(99, 102, 241, 0.3) !important;
}

.create-btn:hover {
  box-shadow: 0 6px 20px rgba(99, 102, 241, 0.4) !important;
}

.filter-bar {
  margin-bottom: 16px;
}

.filter-bar :deep(.ant-radio-group) {
  display: flex;
  gap: 8px;
}

.filter-bar :deep(.ant-radio-button-wrapper) {
  border-radius: 20px !important;
  border: 1px solid var(--border-color) !important;
  background: var(--bg-secondary) !important;
  color: var(--text-secondary) !important;
  padding: 0 18px;
  height: 38px;
  line-height: 36px;
  font-weight: 500;
}

.filter-bar :deep(.ant-radio-button-wrapper::before) {
  display: none !important;
}

.filter-bar :deep(.ant-radio-button-wrapper-checked) {
  background: var(--color-primary) !important;
  border-color: var(--color-primary) !important;
  color: white !important;
  box-shadow: 0 2px 8px rgba(99, 102, 241, 0.3) !important;
}

/* Posts Grid */
.posts-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 24px;
}

.posts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(360px, 1fr));
  gap: 24px;
}

.post-card {
  background: var(--bg-card);
  border: 1px solid var(--border-color-light);
  border-radius: 20px;
  padding: 24px;
  transition: all 0.3s ease;
  cursor: default;
  backdrop-filter: blur(20px);
}

.post-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-lg);
  border-color: var(--color-primary);
}

.post-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
}

.post-author {
  display: flex;
  align-items: center;
  gap: 12px;
}

.author-avatar {
  background: var(--color-primary-gradient) !important;
  color: white !important;
  font-weight: 600;
}

.author-info {
  display: flex;
  flex-direction: column;
}

.author-name {
  font-weight: 600;
  color: var(--text-primary);
  font-size: 15px;
}

.post-time {
  font-size: 13px;
  color: var(--text-muted);
}

.status-tag {
  border-radius: 12px;
  padding: 2px 10px;
  font-size: 12px;
  font-weight: 500;
}

.post-content {
  margin-bottom: 20px;
  cursor: pointer;
}

.post-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 10px;
  line-height: 1.4;
}

.post-excerpt {
  color: var(--text-secondary);
  font-size: 14px;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.post-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 16px;
  border-top: 1px solid var(--border-color-light);
}

.vote-section {
  display: flex;
  gap: 8px;
}

.vote-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  border-color: var(--border-color) !important;
  background: var(--bg-secondary) !important;
  color: var(--text-secondary) !important;
  transition: all 0.2s ease !important;
}

.vote-btn.up:hover:not(:disabled) {
  border-color: var(--color-success) !important;
  color: var(--color-success) !important;
  background: rgba(16, 185, 129, 0.1) !important;
}

.vote-btn.down:hover:not(:disabled) {
  border-color: var(--color-error) !important;
  color: var(--color-error) !important;
  background: rgba(239, 68, 68, 0.1) !important;
}

/* 投票按钮激活状态 - 赞成 */
.vote-up-active,
.vote-btn.up.vote-up-active {
  background: var(--color-success) !important;
  border-color: var(--color-success) !important;
  color: white !important;
}

.vote-up-active:hover:not(:disabled),
.vote-btn.up.vote-up-active:hover:not(:disabled) {
  background: #059669 !important;
  border-color: #059669 !important;
  color: white !important;
}

/* 投票按钮激活状态 - 反对 */
.vote-down-active,
.vote-btn.down.vote-down-active {
  background: var(--color-error) !important;
  border-color: var(--color-error) !important;
  color: white !important;
}

.vote-down-active:hover:not(:disabled),
.vote-btn.down.vote-down-active:hover:not(:disabled) {
  background: #dc2626 !important;
  border-color: #dc2626 !important;
  color: white !important;
}

.vote-count {
  font-weight: 600;
}

.post-stats {
  text-align: right;
}

.approval-rate {
  font-size: 13px;
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

/* Pagination */
.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 40px;
}

/* Detail Modal */
.post-detail-modal :deep(.ant-modal-content) {
  background: var(--bg-card);
  backdrop-filter: blur(20px);
  border: 1px solid var(--border-color-light);
  border-radius: 24px;
  overflow: hidden;
}

.detail-content {
  padding: 8px;
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.detail-author {
  display: flex;
  align-items: center;
  gap: 12px;
}

.detail-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 20px;
  line-height: 1.4;
}

.detail-body {
  color: var(--text-secondary);
  font-size: 15px;
  line-height: 1.8;
  margin-bottom: 24px;
  max-height: 400px;
  overflow-y: auto;
  padding-right: 8px;
}

.detail-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 20px;
  border-top: 1px solid var(--border-color-light);
}

.detail-footer .vote-section {
  gap: 12px;
}

.detail-footer .vote-btn {
  height: 48px !important;
  padding: 0 24px !important;
  border-radius: 12px !important;
}

.approval-info {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4px;
}

.vote-total {
  font-size: 13px;
  color: var(--text-muted);
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

  .posts-grid {
    grid-template-columns: 1fr;
  }

  .filter-bar :deep(.ant-radio-group) {
    flex-wrap: wrap;
  }

  .nav {
    gap: 12px;
  }

  .nav-link {
    display: none;
  }
}
</style>
