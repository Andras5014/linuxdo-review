<template>
  <div class="callback-page">
    <div class="callback-container">
      <a-spin v-if="loading" size="large" />
      <div v-else-if="error" class="error-state">
        <CloseCircleOutlined class="error-icon" />
        <h2>登录失败</h2>
        <p>{{ error }}</p>
        <a-button type="primary" @click="goToLogin">返回登录</a-button>
      </div>
      <div v-else class="success-state">
        <CheckCircleOutlined class="success-icon" />
        <h2>登录成功</h2>
        <p>正在跳转...</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { CheckCircleOutlined, CloseCircleOutlined } from '@ant-design/icons-vue'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const loading = ref(true)
const error = ref('')

onMounted(async () => {
  try {
    // 从URL获取token和user信息
    const token = route.query.token as string
    const userStr = route.query.user as string

    if (!token || !userStr) {
      throw new Error('缺少登录信息')
    }

    // 解析用户信息
    const user = JSON.parse(decodeURIComponent(userStr))
    
    // 设置登录状态
    userStore.setAuth(token, user)
    
    loading.value = false
    
    // 延迟跳转，让用户看到成功状态
    setTimeout(() => {
      const redirect = route.query.redirect as string
      router.push(redirect || '/')
    }, 1000)
  } catch (e) {
    loading.value = false
    error.value = e instanceof Error ? e.message : 'OAuth登录失败，请重试'
  }
})

const goToLogin = () => {
  router.push('/login')
}
</script>

<style scoped>
.callback-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-primary);
}

.callback-container {
  text-align: center;
  padding: 48px;
}

.error-state,
.success-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.error-icon {
  font-size: 64px;
  color: var(--color-error);
}

.success-icon {
  font-size: 64px;
  color: var(--color-success);
}

h2 {
  font-size: 24px;
  color: var(--text-primary);
  margin: 0;
}

p {
  color: var(--text-secondary);
  margin: 0;
}
</style>
