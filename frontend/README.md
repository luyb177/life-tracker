# Life Tracker Frontend

Vue 3 responsive web frontend for Life Tracker.

## Stack

- Vue 3 + TypeScript
- Vite
- Vue Router
- Pinia
- Axios
- Naive UI
- @lucide/vue

## Commands

```bash
npm install
npm run dev
npm run build
```

The dev server proxies `/api` to `http://localhost:8888`.

## Code Structure

```text
src/
  api/          API clients grouped by backend module
  components/   reusable UI building blocks
  layouts/      app shell, sidebar, mobile navigation
  pages/        route-level pages
  router/       route definitions and auth guards
  stores/       Pinia stores
  styles/       global responsive CSS
  types/        shared API types
  utils/        date, money, and formatting helpers
```

## Current Pages

- `/login`: email/password login.
- `/today`: daily workspace with quick life log, quick expense, timeline, metrics, and AI summary.
- `/life-logs`: date-based life log list.
- `/expenses`: date-based expense list and category overview.
- `/summaries`: period summary list and AI generation.
- `/settings`: profile shell and logout.

## Next Tasks

1. Add register and verification-code flows.
2. Add edit/delete/refund actions for records.
3. Add category create/delete UI.
4. Add ECharts trend and category charts.
5. Replace full Naive UI registration with on-demand imports or manual chunks.
6. Add loading skeletons and empty/error retry states per page.

