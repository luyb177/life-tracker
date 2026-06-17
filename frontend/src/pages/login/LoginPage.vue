<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { NForm, NFormItem, NInput, NButton, NTabs, NTabPane, useMessage } from 'naive-ui'
import { useAuthStore } from '@/entities/user/stores/auth.store'

const router = useRouter()
const authStore = useAuthStore()
const message = useMessage()

const tab = ref<'login' | 'register'>('login')
const email = ref('')
const password = ref('')
const code = ref('')
const loading = ref(false)
const sending = ref(false)

async function handleLogin() {
  loading.value = true
  try {
    await authStore.login(email.value, password.value)
    router.push('/dashboard')
  } catch (e: any) {
    message.error(e.response?.data?.msg || '登录失败')
  } finally { loading.value = false }
}

async function handleRegister() {
  loading.value = true
  try {
    await authStore.register(email.value, password.value, code.value)
    message.success('注册成功，请登录')
    tab.value = 'login'
  } catch (e: any) {
    message.error(e.response?.data?.msg || '注册失败')
  } finally { loading.value = false }
}

async function handleSendCode() {
  sending.value = true
  try {
    await authStore.sendCode(email.value, 1)
    message.success('验证码已发送')
  } catch (e: any) {
    message.error(e.response?.data?.msg || '发送失败')
  } finally { sending.value = false }
}
</script>

<template>
  <div style="display: flex; justify-content: center; align-items: center; min-height: 100vh; background: #f5f5f5">
    <div style="width: 100%; max-width: 400px; padding: 32px; background: #fff; border-radius: 12px; box-shadow: 0 2px 12px rgba(0,0,0,0.08)">
      <h1 style="text-align: center; margin-bottom: 24px">Life Tracker</h1>
      <NTabs v-model:value="tab" type="segment" animated>
        <NTabPane name="login" tab="登录">
          <NForm>
            <NFormItem label="邮箱"><NInput v-model:value="email" placeholder="your@email.com" /></NFormItem>
            <NFormItem label="密码"><NInput v-model:value="password" type="password" placeholder="输入密码" @keyup.enter="handleLogin" /></NFormItem>
            <NButton type="primary" block :loading="loading" @click="handleLogin">登录</NButton>
          </NForm>
        </NTabPane>
        <NTabPane name="register" tab="注册">
          <NForm>
            <NFormItem label="邮箱"><NInput v-model:value="email" placeholder="your@email.com" /></NFormItem>
            <NFormItem label="验证码">
              <div style="display: flex; gap: 8px">
                <NInput v-model:value="code" placeholder="验证码" style="flex: 1" />
                <NButton :loading="sending" @click="handleSendCode">发送</NButton>
              </div>
            </NFormItem>
            <NFormItem label="密码"><NInput v-model:value="password" type="password" placeholder="设置密码" /></NFormItem>
            <NButton type="primary" block :loading="loading" @click="handleRegister">注册</NButton>
          </NForm>
        </NTabPane>
      </NTabs>
    </div>
  </div>
</template>
