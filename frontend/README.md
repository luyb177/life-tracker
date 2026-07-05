# Life Tracker Frontend

Vue 3 + TypeScript 前端，提供桌面端和移动端 Web 体验。

## 技术栈

- Vue 3 Composition API
- TypeScript
- Vite
- Vue Router
- Pinia
- Axios
- Naive UI
- ECharts
- markdown-it
- @lucide/vue

## 快速开始

```bash
npm install
npm run dev
npm run build
npm run preview
```

开发服务器默认监听 `5173`，`/api/v1` 请求转发给后端。

## 目录结构

```text
frontend/
├── index.html
├── package.json
├── vite.config.ts
└── src/
    ├── api/                      # API 封装
    │   ├── auth.ts               # 注册、登录、验证码、刷新 token
    │   ├── client.ts             # Axios 实例、统一响应、401 刷新重试
    │   ├── expense.ts            # 支出、分类、统计
    │   ├── lifeLog.ts            # 生活记录
    │   ├── summary.ts            # 总结
    │   └── user.ts               # 用户资料、密码
    ├── assets/
    │   └── default-avatar.svg
    ├── components/
    │   ├── charts/
    │   │   └── ExpenseTrendChart.vue
    │   ├── common/
    │   │   ├── EmptyState.vue
    │   │   ├── MetricCard.vue
    │   │   └── PageHeader.vue
    │   ├── navigation/
    │   │   ├── AppSidebar.vue
    │   │   ├── MobileNav.vue
    │   │   └── navItems.ts
    │   ├── records/
    │   │   ├── ExpenseDetailModal.vue
    │   │   ├── LifeLogDetailModal.vue
    │   │   ├── QuickExpenseForm.vue
    │   │   ├── QuickLifeLogForm.vue
    │   │   └── TimelineList.vue
    │   └── summary/
    │       ├── SummaryDetailModal.vue
    │       └── SummaryPreview.vue
    ├── layouts/
    │   └── AppLayout.vue
    ├── pages/
    │   ├── ExpensesPage.vue
    │   ├── LifeLogsPage.vue
    │   ├── LoginPage.vue
    │   ├── RegisterPage.vue
    │   ├── SettingsPage.vue
    │   ├── SummariesPage.vue
    │   └── TodayPage.vue
    ├── router/
    │   └── index.ts
    ├── stores/
    │   └── auth.ts
    ├── styles/
    │   └── main.css
    ├── types/
    │   └── api.ts
    └── utils/
        ├── date.ts
        ├── markdown.ts
        ├── money.ts
        └── summary.ts
```

## 页面

| 路由 | 页面 | 功能 |
| --- | --- | --- |
| `/login` | 登录 | 邮箱密码登录 |
| `/register` | 注册 | 邮箱验证码注册，验证码 60 秒冷却 |
| `/today` | 今日 | 指标、快捷记录、快捷支出、时间线、AI 日报 |
| `/life-logs` | 生活记录 | 按日期查看、编辑、删除 |
| `/expenses` | 支出 | 按日期查看、编辑、删除、退款、分类、趋势图 |
| `/summaries` | 总结 | AI/用户总结列表、生成、手写、编辑、删除 |
| `/settings` | 设置 | 修改资料、头像、修改密码、退出登录 |

## 认证状态

认证逻辑在 `stores/auth.ts` 和 `api/client.ts`：

- 登录后保存 access token、refresh token 和用户信息。
- 应用启动时恢复本地 token，并启动 refresh 定时器。
- 定时器解析 access token 的 `exp`，提前 5 分钟刷新。
- 如果解析不到 `exp`，使用 90 分钟兜底刷新。
- 401 响应会先用 refresh token 刷新，再重试原请求。
- 多个请求同时 401 时共用一个 refresh promise，避免 refresh token 轮换冲突。
- refresh 失败时清空登录态并跳转 `/login`。

## 组件约定

- 页面负责加载数据、调用 API、处理提示和刷新。
- 可复用表单/弹层放在 `components/`。
- 生活记录、支出、总结详情分别复用：
  - `LifeLogDetailModal.vue`
  - `ExpenseDetailModal.vue`
  - `SummaryDetailModal.vue`
- 页面头部统一用 `PageHeader.vue`。
- 指标卡统一用 `MetricCard.vue`。
- 空状态统一用 `EmptyState.vue`。

## 数据约定

- 金额后端返回分，前端通过 `utils/money.ts` 转为元展示。
- 日期和时间格式通过 `utils/date.ts` 统一处理。
- AI/手写总结支持 Markdown，使用 `utils/markdown.ts` 渲染。
- 总结周期计算放在 `utils/summary.ts`。

## 响应式

- 桌面端使用侧边栏导航。
- 移动端使用底部导航。
- 列表和详情弹层适配窄屏。
- 支出趋势图使用 ECharts，并根据视口调整高度。

## 构建

```bash
npm run typecheck
npm run build
```

当前构建可能出现 Vite chunk 体积提示，主要来自 UI 库、ECharts 和 Markdown 渲染依赖；不影响构建结果。后续可通过按需导入或手动拆 chunk 优化。
