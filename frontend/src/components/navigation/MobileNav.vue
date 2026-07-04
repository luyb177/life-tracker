<template>
  <nav class="mobile-nav">
    <router-link v-for="item in navItems" :key="item.path" :to="item.path" class="mobile-nav-link">
      <component :is="item.icon" :size="20" />
      <span>{{ item.label }}</span>
    </router-link>
  </nav>
</template>

<script setup lang="ts">
import { navItems } from '@/components/navigation/navItems'
</script>

