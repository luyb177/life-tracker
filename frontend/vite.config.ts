import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

function getNaiveChunk(id: string): string | undefined {
  const normalized = id.replace(/\\/g, '/')

  if (
    normalized.includes('/node_modules/vueuc/') ||
    normalized.includes('/node_modules/vooks/') ||
    normalized.includes('/node_modules/vdirs/') ||
    normalized.includes('/node_modules/seemly/') ||
    normalized.includes('/node_modules/@css-render/')
  ) {
    return 'naive-ui-core'
  }

  const marker = '/node_modules/naive-ui/lib/'
  const index = normalized.indexOf(marker)
  if (index === -1) return undefined

  const segment = normalized.slice(index + marker.length).split('/')[0]
  const formComponents = new Set(['button', 'input', 'input-number', 'select', 'date-picker', 'form', 'modal'])
  const feedbackComponents = new Set(['tag', 'spin', 'alert', 'message', 'notification', 'dialog', 'popconfirm'])
  const layoutComponents = new Set(['space', 'grid'])

  if (formComponents.has(segment)) return 'naive-ui-form'
  if (feedbackComponents.has(segment)) return 'naive-ui-feedback'
  if (layoutComponents.has(segment)) return 'naive-ui-layout'
  return 'naive-ui-core'
}

function getEchartsChunk(id: string): string | undefined {
  const normalized = id.replace(/\\/g, '/')

  if (normalized.includes('/node_modules/vue-echarts/')) return 'echarts-core'
  if (normalized.includes('/node_modules/echarts/lib/chart/')) return 'echarts-charts'
  if (normalized.includes('/node_modules/echarts/lib/component/')) return 'echarts-components'
  if (normalized.includes('/node_modules/echarts/')) return 'echarts-core'
  return undefined
}

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: { '@': resolve(__dirname, 'src') },
  },
  server: {
    port: 3000,
    proxy: {
      '/api': { target: 'http://localhost:8888', changeOrigin: true },
    },
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks(id: string) {
          return getEchartsChunk(id) || getNaiveChunk(id)
        },
      },
    },
  },
})
