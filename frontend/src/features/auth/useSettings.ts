import { ref } from 'vue'
import { useMutation } from '@tanstack/vue-query'
import { useMessage } from 'naive-ui'
import { updateUserInfo, changePassword, sendChangePasswordCode } from '@/entities/user/api/user.api'
import { useAuthStore } from '@/entities/user/stores/auth.store'

export function useSettings() {
  const authStore = useAuthStore()
  const msg = useMessage()
  const username = ref(authStore.user?.username || '')
  const code = ref('')
  const newPassword = ref('')

  const updateInfoMut = useMutation({
    mutationFn: () => updateUserInfo({ username: username.value }),
    onSuccess: () => msg.success('已更新'),
    onError: (e: any) => msg.error(e.response?.data?.msg || '更新失败'),
  })

  const changePwdMut = useMutation({
    mutationFn: () => changePassword({ code: code.value, new_password: newPassword.value }),
    onSuccess: () => { msg.success('密码已修改，请重新登录'); authStore.logout() },
    onError: (e: any) => msg.error(e.response?.data?.msg || '修改失败'),
  })

  async function handleSendCode() {
    try { await sendChangePasswordCode(authStore.user!.email); msg.success('验证码已发送') }
    catch (e: any) { msg.error(e.response?.data?.msg || '发送失败') }
  }

  return { username, code, newPassword, updateInfoMut, changePwdMut, handleSendCode }
}
