import client from '@/shared/api/client'

export function updateUserInfo(data: { username?: string; avatar?: string }) {
  return client.post('/user/update_info', data).then(r => r.data.data)
}

export function changePassword(data: { code: string; new_password: string }) {
  return client.post('/user/change_password', data).then(r => r.data.data)
}

export function sendChangePasswordCode(email: string) {
  return client.post('/auth/send_code', { target: email, channel: 1, purpose: 2 }).then(r => r.data.data)
}
