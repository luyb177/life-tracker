<script setup lang="ts">
import { NForm, NFormItem, NInput, NButton } from 'naive-ui'
import PageHeader from '@/shared/ui/page-header/PageHeader.vue'
import { useAuthStore } from '@/entities/user/stores/auth.store'
import { useSettings } from '@/features/auth/useSettings'

const authStore = useAuthStore()
const { username, code, newPassword, updateInfoMut, changePwdMut, handleSendCode } = useSettings()
</script>
<template>
  <PageHeader title="设置" description="个人信息与安全" />
  <div style="max-width: 480px">
    <h4>个人信息</h4>
    <NForm>
      <NFormItem label="邮箱"><NInput :value="authStore.user?.email || ''" disabled /></NFormItem>
      <NFormItem label="用户名"><NInput v-model:value="username" placeholder="设置用户名" /></NFormItem>
      <NButton type="primary" :loading="updateInfoMut.isPending?.value" @click="updateInfoMut.mutate()">保存</NButton>
    </NForm>

    <h4 style="margin-top: 32px">修改密码</h4>
    <NForm>
      <NFormItem label="验证码">
        <div style="display: flex; gap: 8px">
          <NInput v-model:value="code" placeholder="验证码" style="flex: 1" />
          <NButton size="small" @click="handleSendCode">发送</NButton>
        </div>
      </NFormItem>
      <NFormItem label="新密码"><NInput v-model:value="newPassword" type="password" placeholder="新密码" /></NFormItem>
      <NButton type="primary" :loading="changePwdMut.isPending?.value" @click="changePwdMut.mutate()">修改密码</NButton>
    </NForm>
  </div>
</template>
