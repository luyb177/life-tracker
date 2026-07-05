<template>
  <main class="auth-page">
    <section class="auth-panel">
      <div class="brand large">
        <div class="brand-mark">L</div>
        <div>
          <strong>Life Tracker</strong>
          <span>记录生活，复盘花销</span>
        </div>
      </div>

      <n-form class="auth-form" :show-label="false" @submit.prevent="submit">
        <h1>注册</h1>
        <p>创建账号，开始记录你的日常。</p>
        <n-input v-model:value="email" placeholder="邮箱" />
        <div class="verify-code-row">
          <n-input v-model:value="code" placeholder="验证码" />
          <n-button :disabled="!canSendCode" :loading="sendingCode" @click="sendCode">
            {{ codeCooldown > 0 ? `${codeCooldown}s` : '发送验证码' }}
          </n-button>
        </div>
        <n-input
          v-model:value="password"
          type="password"
          show-password-on="click"
          placeholder="密码"
        />
        <n-input
          v-model:value="confirmPassword"
          type="password"
          show-password-on="click"
          placeholder="确认密码"
        />
        <n-button type="primary" block :loading="loading" :disabled="!canSubmit" @click="submit">
          注册
        </n-button>
        <p class="auth-switch">
          已有账号？
          <RouterLink to="/login">去登录</RouterLink>
        </p>
      </n-form>
    </section>
  </main>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import { register, sendVerificationCode } from '@/api/auth'

const router = useRouter()
const message = useMessage()
const email = ref('')
const code = ref('')
const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const sendingCode = ref(false)
const codeCooldown = ref(0)
let cooldownTimer: number | null = null

const canSendCode = computed(() => Boolean(email.value.trim()) && codeCooldown.value === 0 && !sendingCode.value)
const canSubmit = computed(
  () =>
    Boolean(email.value.trim()) &&
    Boolean(code.value.trim()) &&
    Boolean(password.value) &&
    password.value === confirmPassword.value
)

async function sendCode() {
  if (!canSendCode.value) return
  sendingCode.value = true
  try {
    await sendVerificationCode(email.value.trim(), 1)
    startCooldown()
    message.success('验证码已发送')
  } catch (error) {
    message.error(error instanceof Error ? error.message : '发送验证码失败')
  } finally {
    sendingCode.value = false
  }
}

function startCooldown() {
  stopCooldown()
  codeCooldown.value = 60
  cooldownTimer = window.setInterval(() => {
    codeCooldown.value -= 1
    if (codeCooldown.value <= 0) stopCooldown()
  }, 1000)
}

function stopCooldown() {
  if (cooldownTimer) {
    window.clearInterval(cooldownTimer)
    cooldownTimer = null
  }
  if (codeCooldown.value < 0) codeCooldown.value = 0
}

async function submit() {
  if (!canSubmit.value) {
    if (password.value && confirmPassword.value && password.value !== confirmPassword.value) {
      message.error('两次输入的密码不一致')
    }
    return
  }
  loading.value = true
  try {
    await register({
      target: email.value.trim(),
      channel: 1,
      code: code.value.trim(),
      password: password.value
    })
    message.success('注册成功，请登录')
    await router.push('/login')
  } catch (error) {
    message.error(error instanceof Error ? error.message : '注册失败')
  } finally {
    loading.value = false
  }
}

onBeforeUnmount(stopCooldown)
</script>
