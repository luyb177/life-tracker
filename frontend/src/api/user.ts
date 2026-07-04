import { http, unwrap } from '@/api/client'

export interface UpdateUserInfoPayload {
  username?: string
  avatar?: string
}

export interface ChangePasswordPayload {
  code: string
  new_password: string
}

export function updateUserInfo(payload: UpdateUserInfoPayload) {
  return unwrap<Record<string, never>>(http.post('/user/update_info', payload))
}

export function changePassword(payload: ChangePasswordPayload) {
  return unwrap<Record<string, never>>(http.post('/user/change_password', payload))
}

