<template>
  <a-config-provider :locale="zhCN" :theme="antTheme">
    <router-view />
  </a-config-provider>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { theme as antdTheme } from 'ant-design-vue'
import zhCN from 'ant-design-vue/es/locale/zh_CN'
import { useThemeStore } from '@/stores/theme'

const themeStore = useThemeStore()

// Ant Design Vue 主题配置
const antTheme = computed(() => ({
  algorithm: themeStore.theme === 'dark' 
    ? antdTheme.darkAlgorithm 
    : antdTheme.defaultAlgorithm,
  token: {
    colorPrimary: themeStore.theme === 'dark' ? '#818cf8' : '#6366f1',
    borderRadius: 10,
    fontFamily: "'Noto Sans SC', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif",
  },
}))

onMounted(() => {
  // 确保主题状态被初始化
  document.documentElement.setAttribute('data-theme', themeStore.theme)
})
</script>

<style>
#app {
  min-height: 100vh;
}
</style>
