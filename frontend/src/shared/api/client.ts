import axios from 'axios'
import { useAuthStore } from '@/entities/user/stores/auth.store'

const client = axios.create({
  baseURL: '/api/v1',
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' },
})

client.interceptors.request.use((config) => {
  const authStore = useAuthStore()
  if (authStore.token) {
    config.headers.Authorization = `Bearer ${authStore.token}`
  }
  return config
})

client.interceptors.response.use(
  (res) => res,
  async (error) => {
    if (error.response?.status === 401) {
      const authStore = useAuthStore()
      // try refresh
      const success = await authStore.tryRefresh()
      if (success && error.config) {
        error.config.headers.Authorization = `Bearer ${authStore.token}`
        return client.request(error.config)
      }
      authStore.logout()
      window.location.href = '/login'
    }
    return Promise.reject(error)
  },
)

export default client
