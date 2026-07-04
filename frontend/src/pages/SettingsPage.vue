<template>
  <div class="page">
    <PageHeader title="设置" description="管理账号信息和登录状态。" />

    <section class="panel settings-panel">
      <div class="profile-row">
        <n-avatar round :size="56" :src="auth.user?.avatar || undefined">
          {{ auth.user?.username?.slice(0, 1) || 'L' }}
        </n-avatar>
        <div>
          <strong>{{ auth.user?.username || 'Life Tracker 用户' }}</strong>
          <p>{{ auth.user?.email || '未读取到邮箱' }}</p>
        </div>
      </div>

      <div class="settings-actions">
        <n-button secondary>修改资料</n-button>
        <n-button secondary>修改密码</n-button>
        <n-button type="error" ghost @click="logout">退出登录</n-button>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import PageHeader from '@/components/common/PageHeader.vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()

async function logout() {
  auth.logout()
  await router.push('/login')
}
</script>

