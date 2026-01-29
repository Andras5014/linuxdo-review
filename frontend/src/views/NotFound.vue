<template>
  <div class="not-found-page">
    <div class="content">
      <h1 class="error-code">404</h1>
      <h2 class="error-title">页面未找到</h2>
      <p class="error-desc">抱歉，您访问的页面不存在或已被移除</p>
      <router-link to="/">
        <a-button type="primary" size="large">
          <template #icon><HomeOutlined /></template>
          返回首页
        </a-button>
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { HomeOutlined } from '@ant-design/icons-vue'
</script>

<style scoped>
.not-found-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-primary);
}

.content {
  text-align: center;
}

.error-code {
  font-size: 120px;
  font-weight: 700;
  font-family: var(--font-mono);
  background: linear-gradient(135deg, var(--color-primary) 0%, #3b82f6 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  line-height: 1;
  margin-bottom: 16px;
}

.error-title {
  font-size: 28px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 12px;
}

.error-desc {
  font-size: 16px;
  color: var(--text-secondary);
  margin-bottom: 32px;
}
</style>
