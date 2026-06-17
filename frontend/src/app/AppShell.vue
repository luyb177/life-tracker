<script setup lang="ts">
import { ref, computed, h } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  NLayout, NLayoutHeader, NLayoutSider, NLayoutContent,
  NMenu, NButton, NIcon, NSpace, NAvatar, NDropdown, useMessage,
} from 'naive-ui'
import {
  HomeOutline, DocumentTextOutline, WalletOutline,
  BarChartOutline, SettingsOutline, LogOutOutline,
  MenuOutline,
} from '@vicons/ionicons5'
import { useAuthStore } from '@/entities/user/stores/auth.store'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const collapsed = ref(false)

const menuOptions = [
  { label: '仪表盘', key: '/dashboard', icon: HomeOutline },
  { label: '生活记录', key: '/records', icon: DocumentTextOutline },
  { label: '支出', key: '/expenses', icon: WalletOutline },
  { label: '总结', key: '/summaries', icon: BarChartOutline },
  { label: '分析', key: '/analytics', icon: BarChartOutline },
  { label: '设置', key: '/settings', icon: SettingsOutline },
]

const activeKey = computed(() => {
  const path = route.path
  if (path.startsWith('/dashboard')) return '/dashboard'
  if (path.startsWith('/records')) return '/records'
  if (path.startsWith('/expenses')) return '/expenses'
  if (path.startsWith('/summaries')) return '/summaries'
  if (path.startsWith('/analytics')) return '/analytics'
  return '/settings'
})

function handleMenuClick(key: string) {
  router.push(key)
}

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>

<template>
  <NLayout has-sider position="absolute" style="height: 100vh">
    <!-- Desktop sidebar -->
    <NLayoutSider
      bordered
      collapse-mode="width"
      :collapsed-width="64"
      :width="220"
      :collapsed="collapsed"
      show-trigger
      @collapse="collapsed = true"
      @expand="collapsed = false"
      class="hidden-mobile"
    >
      <div style="padding: 16px; text-align: center; font-weight: bold; font-size: 18px">
        {{ collapsed ? 'LT' : 'Life Tracker' }}
      </div>
      <NMenu
        :value="activeKey"
        :options="menuOptions.map(o => ({
          label: o.label,
          key: o.key,
          icon: () => h(NIcon, null, { default: () => h(o.icon) }),
        }))"
        @update:value="handleMenuClick"
      />
      <div style="position: absolute; bottom: 12px; left: 0; right: 0; padding: 0 8px">
        <NButton text @click="handleLogout" style="width: 100%">
          <NIcon><LogOutOutline /></NIcon>
          <span v-if="!collapsed" style="margin-left: 8px">退出</span>
        </NButton>
      </div>
    </NLayoutSider>

    <NLayout>
      <NLayoutHeader bordered style="height: 48px; padding: 0 16px; display: flex; align-items: center; justify-content: space-between">
        <NButton text @click="collapsed = !collapsed" class="hidden-mobile">
          <NIcon size="20"><MenuOutline /></NIcon>
        </NButton>
        <NSpace>
          <span>{{ authStore.user?.username || authStore.user?.email || '' }}</span>
          <NAvatar size="small" :src="authStore.user?.avatar" />
        </NSpace>
      </NLayoutHeader>
      <NLayoutContent :native-scrollbar="false" style="padding: 16px">
        <router-view />
      </NLayoutContent>
    </NLayout>

    <!-- Mobile bottom nav -->
    <div class="mobile-nav">
      <div
        v-for="o in menuOptions"
        :key="o.key"
        class="mobile-nav-item"
        :class="{ active: activeKey === o.key }"
        @click="handleMenuClick(o.key)"
      >
        <NIcon :size="20"><component :is="o.icon" /></NIcon>
        <span style="font-size: 10px">{{ o.label }}</span>
      </div>
    </div>
  </NLayout>
</template>

<style scoped>
.hidden-mobile {
  @media (max-width: 768px) { display: none; }
}
.mobile-nav {
  display: none;
  @media (max-width: 768px) {
    display: flex;
    position: fixed;
    bottom: 0; left: 0; right: 0;
    height: 56px;
    background: var(--n-color);
    border-top: 1px solid var(--n-border-color);
    z-index: 100;
  }
}
.mobile-nav-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--n-text-color-3);
  cursor: pointer;
}
.mobile-nav-item.active {
  color: var(--n-color-target);
}
</style>
