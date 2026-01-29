import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getSystemStatus } from '@/api/auth'

// 系统初始化状态缓存
let systemInitialized: boolean | null = null

const routes: RouteRecordRaw[] = [
  {
    path: '/setup',
    name: 'Setup',
    component: () => import('@/views/Setup.vue'),
    meta: { title: '系统初始化', setup: true },
  },
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue'),
    meta: { title: '首页' },
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { title: '登录', guest: true },
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { title: '注册', guest: true },
  },
  {
    path: '/oauth/callback',
    name: 'OAuthCallback',
    component: () => import('@/views/OAuthCallback.vue'),
    meta: { title: 'OAuth回调' },
  },
  {
    path: '/posts',
    name: 'PostList',
    component: () => import('@/views/PostList.vue'),
    meta: { title: '申请列表' },
  },
  {
    path: '/posts/create',
    name: 'CreatePost',
    component: () => import('@/views/CreatePost.vue'),
    meta: { title: '发布申请', requiresAuth: true },
  },
  {
    path: '/my-posts',
    name: 'MyPosts',
    component: () => import('@/views/MyPosts.vue'),
    meta: { title: '我的申请', requiresAuth: true },
  },
  {
    path: '/review',
    name: 'Review',
    component: () => import('@/views/Review.vue'),
    meta: { title: '二级审核', requiresAuth: true, requiresCertified: true },
  },
  {
    path: '/admin',
    name: 'Admin',
    component: () => import('@/views/Admin.vue'),
    meta: { title: '管理后台', requiresAuth: true, requiresAdmin: true },
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue'),
    meta: { title: '页面未找到' },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// 检查系统是否已初始化
async function checkSystemInitialized(): Promise<boolean> {
  if (systemInitialized !== null) {
    return systemInitialized
  }
  
  try {
    const response = await getSystemStatus()
    systemInitialized = response.data.data.initialized
    return systemInitialized
  } catch {
    // 如果请求失败，默认已初始化（避免阻塞）
    return true
  }
}

// 重置系统初始化状态缓存（供初始化完成后使用）
export function resetSystemInitializedCache() {
  systemInitialized = null
}

// 路由守卫
router.beforeEach(async (to, _from, next) => {
  // 设置页面标题
  document.title = `${to.meta.title || 'Linux.do'} - 邀请码申请系统`

  // 检查系统是否已初始化
  const initialized = await checkSystemInitialized()
  
  // 如果系统未初始化，强制跳转到初始化页面
  if (!initialized && to.name !== 'Setup') {
    next({ name: 'Setup' })
    return
  }
  
  // 如果系统已初始化，不允许访问初始化页面
  if (initialized && to.meta.setup) {
    next({ name: 'Home' })
    return
  }

  const userStore = useUserStore()
  
  // 需要认证的页面
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
    return
  }

  // 游客页面（已登录用户不能访问）
  if (to.meta.guest && userStore.isLoggedIn) {
    next({ name: 'Home' })
    return
  }

  // 需要认证用户权限的页面
  if (to.meta.requiresCertified && !userStore.isCertified) {
    next({ name: 'Home' })
    return
  }

  // 需要管理员权限的页面
  if (to.meta.requiresAdmin && !userStore.isAdmin) {
    next({ name: 'Home' })
    return
  }

  next()
})

export default router
