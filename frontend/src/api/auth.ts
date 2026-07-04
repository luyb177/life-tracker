import { http, unwrap } from '@/api/client'
import type { LoginResp } from '@/types/api'

export interface LoginPayload {
  target: string
  channel: number
  password: string
}

export interface RegisterPayload {
  target: string
  channel: number
  code: string
  password: string
}

export function login(payload: LoginPayload) {
  return unwrap<LoginResp>(http.post('/auth/login', payload))
}

export function register(payload: RegisterPayload) {
  return unwrap<Record<string, never>>(http.post('/auth/register', payload))
}

export function sendVerificationCode(target: string, purpose: number) {
  return unwrap<Record<string, never>>(
    http.post('/auth/send_code', { target, channel: 1, purpose })
  )
}

export function refreshToken(refresh_token: string) {
  return unwrap<{ token: string; refresh_token: string }>(
    http.post('/auth/refresh', { refresh_token })
  )
}

