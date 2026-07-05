import { createApp } from 'vue'
import { createPinia } from 'pinia'
import naive from 'naive-ui'
import App from './App.vue'
import router from './router'
import { useAuthStore } from './stores/auth'
import './styles/main.css'

const pinia = createPinia()
const app = createApp(App)

app.use(pinia).use(router).use(naive).mount('#app')

useAuthStore().startRefreshTimer()
