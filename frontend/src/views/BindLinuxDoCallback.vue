<template>
  <div class="callback-page">
    <div class="callback-container">
      <a-spin v-if="loading" size="large">
        <template #tip>
          <span class="loading-text">正在处理...</span>
        </template>
      </a-spin>
      <div v-else-if="error" class="error-state">
        <CloseCircleOutlined class="error-icon" />
        <h2>绑定失败</h2>
        <p>{{ error }}</p>
        <a-button type="primary" @click="goToProfile">返回个人资料</a-button>
      </div>
      <div v-else class="success-state">
        <CheckCircleOutlined class="success-icon" />
        <h2>绑定成功</h2>
        <p>正在跳转到个人资料页面...</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { CheckCircleOutlined, CloseCircleOutlined } from '@ant-design/icons-vue'
import { useUserStore } from '@/stores/user'
import { getCurrentUser } from '@/api/auth'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const loading = ref(true)
const error = ref('')

onMounted(async () => {
  try {
    // 检查是否有错误参数
    const errorParam = route.query.error as string
    if (errorParam) {
      const errorDesc = route.query.error_description as string
      throw new Error(errorDesc || 'OAuth授权失败')
    }

    // 检查是否是绑定成功的回调
    const success = route.query.success as string
    const token = route.query.token as string

    if (success !== 'true') {
      throw new Error('无效的回调参数')
    }

    if (!token) {
      throw new Error('缺少登录凭证')
    }

    // 恢复 token
    localStorage.setItem('token', token)

    // 获取最新的用户信息（包含绑定后的 LinuxDO 信息）
    const response = await getCurrentUser()
    const user = response.data.data

    // 更新用户状态
    userStore.setAuth(token, user)

    loading.value = false

    // 延迟跳转，让用户看到成功状态
    setTimeout(() => {
      router.push('/profile')
    }, 1500)
  } catch (e: unknown) {
    loading.value = false
    if (e instanceof Error) {
      error.value = e.message
    } else {
      error.value = '绑定失败，请重试'
    }
  }
})

const goToProfile = () => {
  router.push('/profile')
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

.loading-text {
  display: block;
  margin-top: 16px;
  color: var(--text-secondary);
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
  max-width: 300px;
  line-height: 1.6;
}
</style>
