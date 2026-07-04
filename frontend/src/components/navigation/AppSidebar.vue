<template>
  <aside class="sidebar">
    <div class="brand">
      <div class="brand-mark">L</div>
      <div>
        <strong>Life Tracker</strong>
        <span>个人生活账本</span>
      </div>
    </div>

    <nav class="sidebar-nav">
      <router-link v-for="item in navItems" :key="item.path" :to="item.path" class="nav-link">
        <component :is="item.icon" :size="18" />
        <span>{{ item.label }}</span>
      </router-link>
    </nav>
  </aside>
</template>

<script setup lang="ts">
import { navItems } from '@/components/navigation/navItems'
</script>

