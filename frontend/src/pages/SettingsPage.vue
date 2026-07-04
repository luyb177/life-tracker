<template>
  <div class="page">
    <PageHeader title="设置" description="管理账号信息和登录状态。" />

    <section class="panel settings-panel">
      <div class="profile-row">
        <n-avatar round :size="56" :src="avatarSrc" />
        <div>
          <strong>{{ auth.user?.username || 'Life Tracker 用户' }}</strong>
          <p>{{ auth.user?.email || '未读取到邮箱' }}</p>
        </div>
      </div>

      <div class="settings-actions">
        <n-button secondary @click="openProfileModal">修改资料</n-button>
        <n-button secondary @click="openPasswordModal">修改密码</n-button>
        <n-button type="error" ghost @click="logout">退出登录</n-button>
      </div>
    </section>

    <n-modal v-model:show="showProfileModal" preset="dialog" title="修改资料">
      <n-form :show-label="true" label-placement="top">
        <n-form-item label="用户名">
          <n-input v-model:value="profileUsername" placeholder="输入用户名" />
        </n-form-item>
        <n-form-item label="头像地址">
          <n-input v-model:value="profileAvatar" placeholder="https://..." />
        </n-form-item>
      </n-form>
      <template #action>
        <n-button @click="showProfileModal = false">取消</n-button>
        <n-button type="primary" :loading="savingProfile" @click="saveProfile">保存</n-button>
      </template>
    </n-modal>

    <n-modal v-model:show="showPasswordModal" preset="dialog" title="修改密码">
      <n-form :show-label="true" label-placement="top">
        <n-form-item label="邮箱">
          <n-input :value="auth.user?.email || ''" disabled />
        </n-form-item>
        <n-form-item label="验证码">
          <div class="verify-code-row">
            <n-input v-model:value="passwordCode" placeholder="6 位验证码" maxlength="6" />
            <n-button secondary :loading="sendingCode" :disabled="!canSendCode" @click="sendPasswordCode">
              {{ codeButtonText }}
            </n-button>
          </div>
        </n-form-item>
        <n-form-item label="新密码">
          <n-input
            v-model:value="newPassword"
            type="password"
            show-password-on="click"
            placeholder="6-128 位"
          />
        </n-form-item>
      </n-form>
      <template #action>
        <n-button @click="showPasswordModal = false">取消</n-button>
        <n-button type="primary" :loading="changingPassword" @click="savePassword">保存</n-button>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import PageHeader from '@/components/common/PageHeader.vue'
import { sendVerificationCode } from '@/api/auth'
import { changePassword, updateUserInfo } from '@/api/user'
import { useAuthStore } from '@/stores/auth'
import defaultAvatar from '@/assets/default-avatar.svg'

const auth = useAuthStore()
const router = useRouter()
const message = useMessage()
const showProfileModal = ref(false)
const showPasswordModal = ref(false)
const savingProfile = ref(false)
const sendingCode = ref(false)
const changingPassword = ref(false)
const profileUsername = ref('')
const profileAvatar = ref('')
const passwordCode = ref('')
const newPassword = ref('')
const codeCooldown = ref(0)
let codeTimer: number | undefined
const avatarSrc = computed(() => auth.user?.avatar || defaultAvatar)
const canSendCode = computed(() => Boolean(auth.user?.email) && codeCooldown.value === 0 && !sendingCode.value)
const codeButtonText = computed(() =>
  codeCooldown.value > 0 ? `${codeCooldown.value}s 后重发` : '发送验证码'
)

function openProfileModal() {
  profileUsername.value = auth.user?.username || ''
  profileAvatar.value = auth.user?.avatar || ''
  showProfileModal.value = true
}

function openPasswordModal() {
  passwordCode.value = ''
  newPassword.value = ''
  showPasswordModal.value = true
}

async function saveProfile() {
  savingProfile.value = true
  try {
    await updateUserInfo({
      username: profileUsername.value.trim() || undefined,
      avatar: profileAvatar.value.trim() || undefined
    })
    auth.patchUserInfo({
      username: profileUsername.value.trim() || auth.user?.username || '',
      avatar: profileAvatar.value.trim()
    })
    showProfileModal.value = false
    message.success('资料已更新')
  } catch (error) {
    message.error(error instanceof Error ? error.message : '更新资料失败')
  } finally {
    savingProfile.value = false
  }
}

async function sendPasswordCode() {
  if (!canSendCode.value || !auth.user?.email) return
  sendingCode.value = true
  try {
    await sendVerificationCode(auth.user.email, 2)
    startCodeCooldown()
    message.success('验证码已发送')
  } catch (error) {
    message.error(error instanceof Error ? error.message : '发送验证码失败')
  } finally {
    sendingCode.value = false
  }
}

function startCodeCooldown(seconds = 60) {
  codeCooldown.value = seconds
  if (codeTimer) window.clearInterval(codeTimer)
  codeTimer = window.setInterval(() => {
    codeCooldown.value -= 1
    if (codeCooldown.value <= 0 && codeTimer) {
      window.clearInterval(codeTimer)
      codeTimer = undefined
      codeCooldown.value = 0
    }
  }, 1000)
}

async function savePassword() {
  changingPassword.value = true
  try {
    await changePassword({
      code: passwordCode.value.trim(),
      new_password: newPassword.value
    })
    message.success('密码已修改，请重新登录')
    auth.logout()
    await router.push('/login')
  } catch (error) {
    message.error(error instanceof Error ? error.message : '修改密码失败')
  } finally {
    changingPassword.value = false
  }
}

async function logout() {
  auth.logout()
  await router.push('/login')
}

onBeforeUnmount(() => {
  if (codeTimer) window.clearInterval(codeTimer)
})
</script>
