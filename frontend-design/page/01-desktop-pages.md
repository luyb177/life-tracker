# Desktop Pages

## Global Layout

Desktop viewport target: `>= 1024px`.

Layout:

- Fixed left sidebar, about `232px`.
- Main content width uses available space with `24px` page padding.
- Top row contains page title, date picker, and quick actions.
- Cards should use `8px` radius or less.
- Avoid nested cards. Use section bands, tables, and compact panels.

Visual tone:

- Neutral background.
- White content surfaces.
- Teal for primary actions.
- Coral for expense actions.
- Amber or green for secondary status.

## `/today`

This is the default page after login.

Desktop structure:

- Header: "今天", date switcher, generate AI summary button.
- Quick record row:
  - Life log input: content textarea, tag selector, occurred time, save button.
  - Expense input: amount, category, note, occurred time, save button.
- Main grid:
  - Left: today's timeline mixing life logs and expenses.
  - Right: today's spending total, category distribution, AI daily summary preview.

States:

- Empty life logs: show one quiet prompt and keep the quick record input focused.
- Empty expenses: show `¥0.00` and category shortcuts.
- AI summary missing: show generate action.
- AI summary generating: show progress state, disable duplicate generate.

## `/life-logs`

Purpose: manage historical records.

Desktop structure:

- Filter bar: date range, tag filter, keyword input.
- Cursor-paginated list grouped by date.
- Each item shows time, content preview, tags, edit/delete actions.
- Edit opens a modal or side drawer.

## `/expenses`

Purpose: manage spending and review expense trends.

Desktop structure:

- Header summary: this month total, today total, average daily spend.
- Tabs: records, categories, stats.
- Records tab: filter bar, table/list, refund/edit/delete actions.
- Categories tab: system and custom category management.
- Stats tab: trend chart, category share chart.

## `/summaries`

Purpose: review AI and manual summaries.

Desktop structure:

- Period type tabs: day, week, month, year.
- Date range filter.
- Summary list with source label: AI or user.
- Detail drawer shows summary content, suggestion content, tags, and edit/delete actions.
- Generate button appears when the selected period has no AI summary.

## `/settings`

Purpose: profile and account security.

Desktop structure:

- Profile section: username, avatar, email display.
- Security section: send verification code, change password.
- Session behavior: refresh token rotation is handled silently by the API client.

