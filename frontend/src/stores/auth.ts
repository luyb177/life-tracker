import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import type { UserInfo } from '@/types/api'
import * as authApi from '@/api/auth'

const TOKEN_KEY = 'life-tracker.token'
const REFRESH_KEY = 'life-tracker.refresh-token'
const USER_KEY = 'life-tracker.user'
const REFRESH_SKEW_MS = 5 * 60 * 1000
const FALLBACK_REFRESH_MS = 90 * 60 * 1000
const MIN_REFRESH_DELAY_MS = 30 * 1000

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem(TOKEN_KEY) || '')
  const refreshToken = ref(localStorage.getItem(REFRESH_KEY) || '')
  const user = ref<UserInfo | null>(readUser())
  const isLoggedIn = computed(() => Boolean(token.value && refreshToken.value))
  let refreshTimer: number | null = null
  let refreshPromise: Promise<void> | null = null

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
    startRefreshTimer()
  }

  async function login(target: string, password: string) {
    const resp = await authApi.login({ target, channel: 1, password })
    persist(resp.token, resp.refresh_token, resp.user_info)
  }

  async function refresh() {
    if (refreshPromise) return refreshPromise
    if (!refreshToken.value) throw new Error('缺少 refresh token')

    refreshPromise = authApi
      .refreshToken(refreshToken.value)
      .then((resp) => {
        persist(resp.token, resp.refresh_token)
      })
      .finally(() => {
        refreshPromise = null
      })
    return refreshPromise
  }

  function logout() {
    stopRefreshTimer()
    token.value = ''
    refreshToken.value = ''
    user.value = null
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(REFRESH_KEY)
    localStorage.removeItem(USER_KEY)
  }

  function stopRefreshTimer() {
    if (!refreshTimer) return
    window.clearTimeout(refreshTimer)
    refreshTimer = null
  }

  function startRefreshTimer() {
    stopRefreshTimer()
    if (!token.value || !refreshToken.value) return

    refreshTimer = window.setTimeout(async () => {
      try {
        await refresh()
      } catch {
        logout()
        redirectToLogin()
      }
    }, nextRefreshDelay(token.value))
  }

  function nextRefreshDelay(accessToken: string) {
    const expMs = readJwtExpMs(accessToken)
    if (!expMs) return FALLBACK_REFRESH_MS
    return Math.max(MIN_REFRESH_DELAY_MS, expMs - Date.now() - REFRESH_SKEW_MS)
  }

  function readJwtExpMs(jwt: string) {
    const payload = jwt.split('.')[1]
    if (!payload) return 0
    try {
      const normalized = payload.replace(/-/g, '+').replace(/_/g, '/')
      const padded = normalized.padEnd(normalized.length + ((4 - (normalized.length % 4)) % 4), '=')
      const decoded = JSON.parse(window.atob(padded)) as { exp?: number }
      return decoded.exp ? decoded.exp * 1000 : 0
    } catch {
      return 0
    }
  }

  function redirectToLogin() {
    if (window.location.pathname !== '/login') {
      window.location.assign('/login')
    }
  }

  function patchUserInfo(payload: Partial<UserInfo>) {
    if (!user.value) return
    user.value = { ...user.value, ...payload }
    localStorage.setItem(USER_KEY, JSON.stringify(user.value))
  }

  return {
    token,
    refreshToken,
    user,
    isLoggedIn,
    login,
    refresh,
    logout,
    patchUserInfo,
    startRefreshTimer,
    stopRefreshTimer
  }
})
