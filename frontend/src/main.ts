import { createApp } from 'vue'
import { createPinia } from 'pinia'
import Antd from 'ant-design-vue'
import App from './App.vue'
import router from './router'
import { useUserStore } from './stores/user'

import 'ant-design-vue/dist/reset.css'
import './styles/main.css'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.use(Antd)

// 初始化用户状态（在挂载前）
const userStore = useUserStore()
userStore.initUser().then(() => {
  app.mount('#app')
})
