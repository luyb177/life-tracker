import axios, { AxiosError, type AxiosRequestConfig } from 'axios'
import type { ApiEnvelope } from '@/types/api'
import { useAuthStore } from '@/stores/auth'

export const http = axios.create({
  baseURL: '/api/v1',
  timeout: 15000
})

http.interceptors.request.use((config) => {
  const auth = useAuthStore()
  if (auth.token) {
    config.headers.Authorization = `Bearer ${auth.token}`
  }
  return config
})

http.interceptors.response.use(
  (response) => {
    const envelope = response.data as ApiEnvelope<unknown>
    if (typeof envelope?.code === 'number' && envelope.code !== 0) {
      return Promise.reject(new Error(envelope.msg || '请求失败'))
    }
    return response
  },
  async (error: AxiosError) => {
    const auth = useAuthStore()
    const original = error.config as (AxiosRequestConfig & { _retried?: boolean }) | undefined
    if (error.response?.status === 401 && auth.refreshToken && original && !original._retried) {
      original._retried = true
      try {
        await auth.refresh()
        original.headers = {
          ...original.headers,
          Authorization: `Bearer ${auth.token}`
        }
        return http(original)
      } catch {
        auth.logout()
      }
    }
    return Promise.reject(error)
  }
)

export async function unwrap<T>(request: Promise<{ data: ApiEnvelope<T> }>): Promise<T> {
  const response = await request
  return response.data.data
}

