import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import client from '@/shared/api/client'

interface UserInfo {
  user_id: number
  username: string
  email: string
  avatar: string
  created_at: string
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const refreshToken = ref(localStorage.getItem('refresh_token') || '')
  const user = ref<UserInfo | null>(null)
  const isLoggedIn = computed(() => !!token.value)

  function setTokens(t: string, rt: string) {
    token.value = t
    refreshToken.value = rt
    localStorage.setItem('token', t)
    localStorage.setItem('refresh_token', rt)
  }

  async function login(email: string, password: string) {
    const res = await client.post('/auth/login', { target: email, channel: 1, password })
    setTokens(res.data.data.token, res.data.data.refresh_token)
    user.value = res.data.data.user_info
  }

  async function register(email: string, password: string, code: string) {
    await client.post('/auth/register', { target: email, channel: 1, code, password })
  }

  async function sendCode(email: string, purpose: number) {
    await client.post('/auth/send_code', { target: email, channel: 1, purpose })
  }

  async function tryRefresh(): Promise<boolean> {
    if (!refreshToken.value) return false
    try {
      const res = await client.post('/auth/refresh', { refresh_token: refreshToken.value })
      setTokens(res.data.data.token, res.data.data.refresh_token)
      return true
    } catch { return false }
  }

  function logout() {
    token.value = ''
    refreshToken.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('refresh_token')
  }

  return { token, refreshToken, user, isLoggedIn, login, register, sendCode, tryRefresh, logout }
})
