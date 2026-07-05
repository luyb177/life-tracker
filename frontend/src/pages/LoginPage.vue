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
        <h1>登录</h1>
        <p>进入你的每日记录工作台。</p>
        <n-input v-model:value="target" placeholder="邮箱" />
        <n-input
          v-model:value="password"
          type="password"
          show-password-on="click"
          placeholder="密码"
        />
        <n-button type="primary" block :loading="loading" @click="submit">登录</n-button>
        <p class="auth-switch">
          还没有账号？
          <RouterLink to="/register">去注册</RouterLink>
        </p>
      </n-form>
    </section>
  </main>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const message = useMessage()
const auth = useAuthStore()
const target = ref('')
const password = ref('')
const loading = ref(false)

async function submit() {
  if (!target.value || !password.value) return
  loading.value = true
  try {
    await auth.login(target.value, password.value)
    await router.push('/today')
  } catch (error) {
    message.error(error instanceof Error ? error.message : '登录失败')
  } finally {
    loading.value = false
  }
}
</script>
