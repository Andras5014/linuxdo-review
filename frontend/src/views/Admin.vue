<template>
  <div class="admin-page">
    <header class="header">
      <div class="header-content">
        <router-link to="/" class="logo">
          <span class="logo-icon">🚀</span>
          <span class="logo-text">Linux.do</span>
          <span class="logo-badge">管理后台</span>
        </router-link>
        
        <nav class="nav">
          <router-link to="/posts" class="nav-link">申请列表</router-link>
          <router-link to="/review" class="nav-link">二级审核</router-link>
          <router-link to="/admin" class="nav-link active">管理后台</router-link>
          
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
              <span class="user-badge admin">管理员</span>
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
      <div class="admin-container slide-up">
        <!-- Stats Section -->
        <div class="stats-section">
          <div class="section-header">
            <h2 class="section-title">
              <DashboardOutlined />
              系统概览
            </h2>
          </div>
          <div class="stats-grid">
            <div class="stat-card">
              <div class="stat-icon users">
                <TeamOutlined />
              </div>
              <div class="stat-info">
                <span class="stat-value">{{ stats.total_users }}</span>
                <span class="stat-label">总用户</span>
              </div>
            </div>
            <div class="stat-card">
              <div class="stat-icon certified">
                <SafetyCertificateOutlined />
              </div>
              <div class="stat-info">
                <span class="stat-value">{{ stats.certified_users }}</span>
                <span class="stat-label">认证用户</span>
              </div>
            </div>
            <div class="stat-card">
              <div class="stat-icon posts">
                <FileTextOutlined />
              </div>
              <div class="stat-info">
                <span class="stat-value">{{ stats.total_posts }}</span>
                <span class="stat-label">总申请</span>
              </div>
            </div>
            <div class="stat-card">
              <div class="stat-icon pending">
                <ClockCircleOutlined />
              </div>
              <div class="stat-info">
                <span class="stat-value">{{ stats.pending_posts }}</span>
                <span class="stat-label">待审核</span>
              </div>
            </div>
            <div class="stat-card">
              <div class="stat-icon approved">
                <CheckCircleOutlined />
              </div>
              <div class="stat-info">
                <span class="stat-value">{{ stats.approved_posts }}</span>
                <span class="stat-label">已通过</span>
              </div>
            </div>
            <div class="stat-card">
              <div class="stat-icon rejected">
                <CloseCircleOutlined />
              </div>
              <div class="stat-info">
                <span class="stat-value">{{ stats.rejected_posts }}</span>
                <span class="stat-label">已拒绝</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Tabs Section -->
        <div class="tabs-section">
          <a-tabs v-model:activeKey="activeTab" class="admin-tabs">
            <!-- Users Tab -->
            <a-tab-pane key="users">
              <template #tab>
                <span class="tab-label">
                  <UserOutlined />
                  用户管理
                </span>
              </template>
              
              <div class="tab-content">
                <div class="tab-header">
                  <div class="filter-group">
                    <a-radio-group v-model:value="userRoleFilter" button-style="solid" @change="handleUserFilterChange">
                      <a-radio-button :value="-1">全部</a-radio-button>
                      <a-radio-button :value="0">普通用户</a-radio-button>
                      <a-radio-button :value="1">认证用户</a-radio-button>
                      <a-radio-button :value="2">管理员</a-radio-button>
                    </a-radio-group>
                  </div>
                </div>

                <a-spin :spinning="loadingUsers">
                  <a-table
                    :columns="userColumns"
                    :data-source="users"
                    :pagination="false"
                    row-key="id"
                    class="admin-table"
                  >
                    <template #bodyCell="{ column, record }">
                      <template v-if="column.key === 'user'">
                        <div class="user-cell">
                          <a-avatar :size="36" class="user-cell-avatar">
                            {{ (record.username || 'U').charAt(0).toUpperCase() }}
                          </a-avatar>
                          <div class="user-cell-info">
                            <span class="user-cell-name">{{ record.username }}</span>
                            <span class="user-cell-email">{{ record.email }}</span>
                          </div>
                        </div>
                      </template>
                      <template v-else-if="column.key === 'role'">
                        <a-tag :color="getRoleColor(record.role)" class="role-tag">
                          {{ getRoleText(record.role) }}
                        </a-tag>
                      </template>
                      <template v-else-if="column.key === 'linuxdo'">
                        <span v-if="record.linuxdo_username" class="linuxdo-info">
                          <img src="https://linux.do/uploads/default/optimized/1X/3a18b4b0da3e8cf96f7eea15241c3d251f28a39b_2_32x32.png" alt="Linux.do" class="linuxdo-icon" />
                          {{ record.linuxdo_username }}
                        </span>
                        <span v-else class="no-linuxdo">-</span>
                      </template>
                      <template v-else-if="column.key === 'created_at'">
                        {{ formatDate(record.created_at) }}
                      </template>
                      <template v-else-if="column.key === 'action'">
                        <a-dropdown>
                          <a-button type="text" size="small">
                            <MoreOutlined />
                          </a-button>
                          <template #overlay>
                            <a-menu @click="(e: any) => handleUserAction(e, record)">
                              <a-menu-item key="normal" :disabled="record.role === 0">
                                <UserOutlined />
                                设为普通用户
                              </a-menu-item>
                              <a-menu-item key="certified" :disabled="record.role === 1">
                                <SafetyCertificateOutlined />
                                设为认证用户
                              </a-menu-item>
                              <a-menu-item key="admin" :disabled="record.role === 2">
                                <CrownOutlined />
                                设为管理员
                              </a-menu-item>
                            </a-menu>
                          </template>
                        </a-dropdown>
                      </template>
                    </template>
                  </a-table>
                </a-spin>

                <div v-if="usersTotal > usersPageSize" class="pagination-container">
                  <a-pagination
                    v-model:current="usersPage"
                    :total="usersTotal"
                    :page-size="usersPageSize"
                    show-quick-jumper
                    @change="handleUsersPageChange"
                  />
                </div>
              </div>
            </a-tab-pane>

            <!-- Configs Tab -->
            <a-tab-pane key="configs">
              <template #tab>
                <span class="tab-label">
                  <SettingOutlined />
                  系统配置
                </span>
              </template>
              
              <div class="tab-content">
                <a-spin :spinning="loadingConfigs">
                  <div class="config-list">
                    <div v-for="config in configs" :key="config.key" class="config-item">
                      <div class="config-info">
                        <span class="config-key">{{ config.key }}</span>
                        <span class="config-desc">{{ config.description }}</span>
                      </div>
                      <div class="config-value">
                        <a-input
                          v-model:value="configValues[config.key]"
                          :placeholder="config.value"
                          class="config-input"
                        />
                      </div>
                    </div>
                  </div>

                  <div v-if="configs.length > 0" class="config-actions">
                    <a-button 
                      type="primary" 
                      size="large"
                      :loading="savingConfigs"
                      @click="handleSaveConfigs"
                      class="save-config-btn"
                    >
                      <template #icon><SaveOutlined /></template>
                      保存配置
                    </a-button>
                  </div>

                  <a-empty v-else-if="!loadingConfigs" description="暂无配置项" />
                </a-spin>
              </div>
            </a-tab-pane>
          </a-tabs>
        </div>
      </div>
    </main>

    <footer class="footer">
      <p>© 2024 Linux.do 邀请码申请系统 · 管理后台</p>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import type { TableColumnsType } from 'ant-design-vue'
import {
  DownOutlined,
  FileTextOutlined,
  LogoutOutlined,
  DashboardOutlined,
  TeamOutlined,
  SafetyCertificateOutlined,
  ClockCircleOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  UserOutlined,
  SettingOutlined,
  MoreOutlined,
  CrownOutlined,
  SaveOutlined,
} from '@ant-design/icons-vue'
import { useUserStore } from '@/stores/user'
import { useThemeStore } from '@/stores/theme'
import { getUsers, getConfigs, batchUpdateConfigs, getStats, updateUserRole } from '@/api/admin'
import type { SystemConfigItem, SystemStats } from '@/api/admin'
import type { User } from '@/types'
import { UserRole } from '@/types'

const router = useRouter()
const userStore = useUserStore()
const themeStore = useThemeStore()

const activeTab = ref('users')

const toggleTheme = () => {
  themeStore.toggleTheme()
}

// Stats
const stats = reactive<SystemStats>({
  total_users: 0,
  certified_users: 0,
  total_posts: 0,
  pending_posts: 0,
  approved_posts: 0,
  rejected_posts: 0,
})

// Users
const loadingUsers = ref(false)
const users = ref<User[]>([])
const usersPage = ref(1)
const usersPageSize = ref(10)
const usersTotal = ref(0)
const userRoleFilter = ref(-1)

const userColumns: TableColumnsType = [
  {
    title: '用户',
    key: 'user',
    width: 280,
  },
  {
    title: '角色',
    key: 'role',
    width: 120,
  },
  {
    title: 'Linux.do',
    key: 'linuxdo',
    width: 160,
  },
  {
    title: '注册时间',
    key: 'created_at',
    width: 180,
  },
  {
    title: '操作',
    key: 'action',
    width: 80,
    align: 'center',
  },
]

// Configs
const loadingConfigs = ref(false)
const configs = ref<SystemConfigItem[]>([])
const configValues = reactive<Record<string, string>>({})
const savingConfigs = ref(false)

const fetchStats = async () => {
  try {
    const response = await getStats()
    Object.assign(stats, response.data.data)
  } catch {
    // 错误已在拦截器中处理
  }
}

const fetchUsers = async () => {
  loadingUsers.value = true
  try {
    const params: { page: number; page_size: number; role?: number } = {
      page: usersPage.value,
      page_size: usersPageSize.value,
    }
    if (userRoleFilter.value !== -1) {
      params.role = userRoleFilter.value
    }
    const response = await getUsers(params)
    const data = response.data.data
    users.value = data.list || []
    usersTotal.value = data.total
  } catch {
    // 错误已在拦截器中处理
  } finally {
    loadingUsers.value = false
  }
}

const fetchConfigs = async () => {
  loadingConfigs.value = true
  try {
    const response = await getConfigs()
    configs.value = response.data.data || []
    // 初始化配置值
    configs.value.forEach(config => {
      configValues[config.key] = config.value
    })
  } catch {
    // 错误已在拦截器中处理
  } finally {
    loadingConfigs.value = false
  }
}

const handleUserFilterChange = () => {
  usersPage.value = 1
  fetchUsers()
}

const handleUsersPageChange = (page: number) => {
  usersPage.value = page
  fetchUsers()
}

const handleUserAction = async (e: { key: string }, user: User) => {
  const roleMap: Record<string, number> = {
    normal: 0,
    certified: 1,
    admin: 2,
  }
  const role = roleMap[e.key]
  if (role === undefined) return

  const roleNames: Record<number, string> = {
    0: '普通用户',
    1: '认证用户',
    2: '管理员',
  }

  Modal.confirm({
    title: '确认修改角色',
    content: `确定将用户 ${user.username} 的角色修改为 ${roleNames[role]} 吗？`,
    okText: '确认',
    cancelText: '取消',
    onOk: async () => {
      try {
        await updateUserRole(user.id, { role })
        message.success('角色修改成功')
        fetchUsers()
        fetchStats()
      } catch {
        // 错误已在拦截器中处理
      }
    },
  })
}

const handleSaveConfigs = async () => {
  // 收集有变化的配置
  const changedConfigs = configs.value
    .filter(config => configValues[config.key] !== config.value)
    .map(config => ({
      key: config.key,
      value: configValues[config.key],
    }))

  if (changedConfigs.length === 0) {
    message.info('没有需要保存的更改')
    return
  }

  savingConfigs.value = true
  try {
    await batchUpdateConfigs(changedConfigs)
    message.success('配置保存成功')
    fetchConfigs()
  } catch {
    // 错误已在拦截器中处理
  } finally {
    savingConfigs.value = false
  }
}

const handleLogout = () => {
  userStore.logout()
  message.success('已退出登录')
  router.push('/')
}

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

const getRoleText = (role: UserRole) => {
  const roleMap: Record<UserRole, string> = {
    [UserRole.Normal]: '普通用户',
    [UserRole.Certified]: '认证用户',
    [UserRole.Admin]: '管理员',
  }
  return roleMap[role] || '未知'
}

const getRoleColor = (role: UserRole) => {
  const colorMap: Record<UserRole, string> = {
    [UserRole.Normal]: 'default',
    [UserRole.Certified]: 'warning',
    [UserRole.Admin]: 'error',
  }
  return colorMap[role] || 'default'
}

onMounted(() => {
  // 检查是否有权限访问
  if (!userStore.isAdmin) {
    message.warning('需要管理员权限才能访问')
    router.push('/posts')
    return
  }
  fetchStats()
  fetchUsers()
  fetchConfigs()
})
</script>

<style scoped>
.admin-page {
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
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
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

.user-badge.admin {
  font-size: 11px;
  padding: 2px 6px;
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
  color: white;
  border-radius: 4px;
  font-weight: 500;
}

/* Main */
.main {
  flex: 1;
  padding: 72px 24px 40px;
}

.admin-container {
  max-width: 1400px;
  margin: 0 auto;
  padding-top: 24px;
}

/* Stats Section */
.stats-section {
  margin-bottom: 32px;
}

.section-header {
  margin-bottom: 20px;
}

.section-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  display: flex;
  align-items: center;
  gap: 10px;
}

.section-title :deep(svg) {
  color: var(--color-primary);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 24px;
  background: var(--bg-card);
  border: 1px solid var(--border-color-light);
  border-radius: 20px;
  transition: transform 0.2s, box-shadow 0.2s;
  backdrop-filter: blur(20px);
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
}

.stat-icon {
  width: 56px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 16px;
  font-size: 24px;
  flex-shrink: 0;
}

.stat-icon.users {
  background: rgba(59, 130, 246, 0.15);
  color: #3b82f6;
}

.stat-icon.certified {
  background: rgba(245, 158, 11, 0.15);
  color: #f59e0b;
}

.stat-icon.posts {
  background: rgba(139, 92, 246, 0.15);
  color: #8b5cf6;
}

.stat-icon.pending {
  background: rgba(236, 72, 153, 0.15);
  color: #ec4899;
}

.stat-icon.approved {
  background: rgba(16, 185, 129, 0.15);
  color: #10b981;
}

.stat-icon.rejected {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
}

.stat-info {
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1.2;
}

.stat-label {
  font-size: 14px;
  color: var(--text-muted);
  margin-top: 4px;
}

/* Tabs Section */
.tabs-section {
  background: var(--bg-card);
  border: 1px solid var(--border-color-light);
  border-radius: 24px;
  padding: 24px;
  backdrop-filter: blur(20px);
}

.admin-tabs :deep(.ant-tabs-nav) {
  margin-bottom: 24px;
}

.admin-tabs :deep(.ant-tabs-tab) {
  padding: 12px 20px !important;
  font-size: 15px;
}

.admin-tabs :deep(.ant-tabs-tab-active .ant-tabs-tab-btn) {
  color: var(--color-primary) !important;
}

.admin-tabs :deep(.ant-tabs-ink-bar) {
  background: var(--color-primary) !important;
}

.tab-label {
  display: flex;
  align-items: center;
  gap: 8px;
}

.tab-content {
  min-height: 400px;
}

.tab-header {
  margin-bottom: 20px;
}

.filter-group :deep(.ant-radio-group) {
  display: flex;
  gap: 8px;
}

.filter-group :deep(.ant-radio-button-wrapper) {
  border-radius: 20px !important;
  border: 1px solid var(--border-color) !important;
  background: var(--bg-secondary) !important;
  color: var(--text-secondary) !important;
  padding: 0 16px;
  height: 36px;
  line-height: 34px;
  font-weight: 500;
}

.filter-group :deep(.ant-radio-button-wrapper::before) {
  display: none !important;
}

.filter-group :deep(.ant-radio-button-wrapper-checked) {
  background: var(--color-primary) !important;
  border-color: var(--color-primary) !important;
  color: white !important;
  box-shadow: 0 2px 8px rgba(99, 102, 241, 0.3) !important;
}

/* Admin Table */
.admin-table :deep(.ant-table) {
  background: transparent !important;
}

.admin-table :deep(.ant-table-thead > tr > th) {
  background: var(--bg-tertiary) !important;
  border-bottom: 1px solid var(--border-color) !important;
  color: var(--text-secondary) !important;
  font-weight: 600;
}

.admin-table :deep(.ant-table-tbody > tr > td) {
  border-bottom: 1px solid var(--border-color) !important;
  background: transparent !important;
}

.admin-table :deep(.ant-table-tbody > tr:hover > td) {
  background: var(--bg-tertiary) !important;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-cell-avatar {
  background: linear-gradient(135deg, var(--color-primary) 0%, #3b82f6 100%) !important;
  color: white !important;
  font-weight: 600;
}

.user-cell-info {
  display: flex;
  flex-direction: column;
}

.user-cell-name {
  font-weight: 600;
  color: var(--text-primary);
  font-size: 14px;
}

.user-cell-email {
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 2px;
}

.role-tag {
  border-radius: 8px;
  padding: 2px 10px;
  font-weight: 500;
}

.linuxdo-info {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text-primary);
}

.linuxdo-icon {
  width: 20px;
  height: 20px;
  border-radius: 4px;
}

.no-linuxdo {
  color: var(--text-muted);
}

/* Config List */
.config-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.config-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  background: var(--bg-tertiary);
  border-radius: 16px;
  gap: 40px;
}

.config-info {
  display: flex;
  flex-direction: column;
  flex: 1;
}

.config-key {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
  font-family: var(--font-mono);
}

.config-desc {
  font-size: 13px;
  color: var(--text-muted);
  margin-top: 4px;
}

.config-value {
  width: 300px;
  flex-shrink: 0;
}

.config-input {
  border-radius: 10px !important;
  height: 44px !important;
}

.config-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid var(--border-color);
}

.save-config-btn {
  height: 48px !important;
  padding: 0 32px !important;
  border-radius: 12px !important;
  font-weight: 600 !important;
  font-size: 15px !important;
}

/* Pagination */
.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 24px;
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
@media (max-width: 1024px) {
  .stats-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .config-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .config-value {
    width: 100%;
  }

  .nav-link {
    display: none;
  }
}

@media (max-width: 480px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>
