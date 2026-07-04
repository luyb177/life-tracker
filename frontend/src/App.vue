<template>
  <n-config-provider :theme-overrides="themeOverrides">
    <n-message-provider>
      <n-dialog-provider>
        <router-view />
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import type { GlobalThemeOverrides } from 'naive-ui'

const themeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: '#168a7a',
    primaryColorHover: '#1a9e8d',
    primaryColorPressed: '#0f6f63',
    borderRadius: '8px',
    fontFamily: 'Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif'
  },
  Button: {
    borderRadiusMedium: '8px'
  },
  Card: {
    borderRadius: '8px'
  }
}
</script>

