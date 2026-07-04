<template>
  <div class="app-shell">
    <AppSidebar />
    <main class="app-main">
      <router-view />
    </main>
    <MobileNav />
  </div>
</template>

<script setup lang="ts">
import AppSidebar from '@/components/navigation/AppSidebar.vue'
import MobileNav from '@/components/navigation/MobileNav.vue'
</script>

