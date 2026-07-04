import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import type { UserInfo } from '@/types/api'
import * as authApi from '@/api/auth'

const TOKEN_KEY = 'life-tracker.token'
const REFRESH_KEY = 'life-tracker.refresh-token'
const USER_KEY = 'life-tracker.user'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem(TOKEN_KEY) || '')
  const refreshToken = ref(localStorage.getItem(REFRESH_KEY) || '')
  const user = ref<UserInfo | null>(readUser())
  const isLoggedIn = computed(() => Boolean(token.value && refreshToken.value))

  function readUser() {
    const raw = localStorage.getItem(USER_KEY)
    if (!raw) return null
    try {
      return JSON.parse(raw) as UserInfo
    } catch {
      return null
    }
  }

  function persist(nextToken: string, nextRefreshToken: string, nextUser?: UserInfo) {
    token.value = nextToken
    refreshToken.value = nextRefreshToken
    localStorage.setItem(TOKEN_KEY, nextToken)
    localStorage.setItem(REFRESH_KEY, nextRefreshToken)
    if (nextUser) {
      user.value = nextUser
      localStorage.setItem(USER_KEY, JSON.stringify(nextUser))
    }
  }

  async function login(target: string, password: string) {
    const resp = await authApi.login({ target, channel: 1, password })
    persist(resp.token, resp.refresh_token, resp.user_info)
  }

  async function refresh() {
    const resp = await authApi.refreshToken(refreshToken.value)
    persist(resp.token, resp.refresh_token)
  }

  function logout() {
    token.value = ''
    refreshToken.value = ''
    user.value = null
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(REFRESH_KEY)
    localStorage.removeItem(USER_KEY)
  }

  function patchUserInfo(payload: Partial<UserInfo>) {
    if (!user.value) return
    user.value = { ...user.value, ...payload }
    localStorage.setItem(USER_KEY, JSON.stringify(user.value))
  }

  return { token, refreshToken, user, isLoggedIn, login, refresh, logout, patchUserInfo }
})
