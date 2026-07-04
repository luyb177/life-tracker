# Information Architecture

## Product Goal

The frontend is a responsive Vue 3 web app for personal life logging, expense tracking, and AI-assisted periodic review.

The first usable screen should be the daily workspace. Users should be able to:

- Add a life log in a few seconds.
- Add an expense in a few seconds.
- Review today's timeline.
- See today's spending total and category split.
- Read or trigger an AI summary.

## Primary Routes

| Route | Page | Purpose | Backend APIs |
| --- | --- | --- | --- |
| `/login` | Login | Login with target, channel, password | `POST /api/v1/auth/login` |
| `/register` | Register | Send code and create account | `POST /api/v1/auth/send_code`, `POST /api/v1/auth/register` |
| `/today` | Today | Daily recording workspace | life log, expense, summary range |
| `/life-logs` | Life Logs | Search, edit, delete life records | `/api/v1/life_log/*` |
| `/expenses` | Expenses | Expense list, categories, stats | `/api/v1/expense/*` |
| `/summaries` | Summaries | Day, week, month, year summaries | `/api/v1/summary/*` |
| `/settings` | Settings | Profile, password, preferences | `/api/v1/user/*` |

## Navigation

Desktop navigation:

- Left sidebar: Today, Life Logs, Expenses, Summaries, Settings.
- Top bar: date picker, quick create button, user menu.

Mobile navigation:

- Bottom tabs: Today, Logs, Expenses, Summaries, Me.
- Floating or fixed quick action button on Today.

## Core Data Mapping

Life log:

- Create: `POST /api/v1/life_log/create`
- List: `GET /api/v1/life_log/list`
- By date: `GET /api/v1/life_log/by_date`
- Update: `POST /api/v1/life_log/update`
- Delete: `POST /api/v1/life_log/delete`

Expense:

- Create: `POST /api/v1/expense/create`
- List: `GET /api/v1/expense/list`
- By date: `GET /api/v1/expense/by_date`
- Categories: `GET /api/v1/expense/categories`
- Category stats: `GET /api/v1/expense/stats/category`
- Trend stats: `GET /api/v1/expense/stats/trend`
- Monthly stats: `GET /api/v1/expense/stats/monthly`

Summary:

- Create manual summary: `POST /api/v1/summary/create`
- Generate AI summary: `POST /api/v1/summary/generate`
- List: `GET /api/v1/summary/list`
- Range: `GET /api/v1/summary/range`
- Update: `POST /api/v1/summary/update`
- Delete: `POST /api/v1/summary/delete`

## MVP Page Priority

1. Auth: login, register, token refresh handling.
2. Today: quick life log, quick expense, daily timeline, daily expense total.
3. Expense: category list, create/update/delete/refund expense.
4. Life Logs: list, edit, delete, date filter.
5. Summaries: list and AI generation.
6. Settings: user info and password change.

