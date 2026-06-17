import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/entities/user/stores/auth.store'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', name: 'login', component: () => import('@/pages/login/LoginPage.vue') },
    {
      path: '/',
      component: () => import('@/app/AppShell.vue'),
      meta: { requiresAuth: true },
      children: [
        { path: '', redirect: '/dashboard' },
        { path: 'dashboard', name: 'dashboard', component: () => import('@/pages/dashboard/DashboardPage.vue') },
        { path: 'records', name: 'records', component: () => import('@/pages/records/RecordsPage.vue') },
        { path: 'expenses', name: 'expenses', component: () => import('@/pages/expenses/ExpensesPage.vue') },
        { path: 'summaries', name: 'summaries', component: () => import('@/pages/summaries/SummariesPage.vue') },
        { path: 'analytics', name: 'analytics', component: () => import('@/pages/analytics/AnalyticsPage.vue') },
        { path: 'settings', name: 'settings', component: () => import('@/pages/settings/SettingsPage.vue') },
      ],
    },
  ],
})

router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore()
  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    next('/login')
  } else if (to.path === '/login' && authStore.isLoggedIn) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router
