# Mobile Pages

## Global Layout

Mobile viewport target: `320px` to `767px`.

Layout:

- Single-column pages.
- Bottom navigation with five tabs: Today, Logs, Expenses, Summary, Me.
- Page padding around `12px` to `16px`.
- Inputs must be large enough for one-handed use.
- Primary actions should sit near the thumb zone.

Responsive behavior:

- Desktop sidebar becomes bottom navigation.
- Multi-column dashboard sections stack vertically.
- Tables become grouped lists.
- Drawers become full-screen or bottom-sheet panels.

## `/today`

Mobile structure:

- Compact header with current date and user action.
- Segmented control: "记录" and "统计".
- Quick actions: life log, expense.
- Quick input:
  - Life log: textarea, tags, time, save.
  - Expense: amount keypad-friendly field, category chips, note, save.
- Timeline list sorted by occurred time.
- AI summary preview under today's records.

Important mobile details:

- Keep the record action visible after scrolling a short distance.
- Use bottom sheet for category selection.
- Use chips for common categories and tags.
- Avoid long horizontal controls.

## `/life-logs`

Mobile structure:

- Date picker row.
- Search input.
- Grouped timeline list.
- Item actions are in a contextual menu or swipe-style action area.
- Edit uses full-screen sheet.

## `/expenses`

Mobile structure:

- Top total for today or selected period.
- Category shortcuts.
- List grouped by date.
- Stats are under a separate tab to avoid crowding.
- Amount entry should use numeric keyboard.

## `/summaries`

Mobile structure:

- Period tabs: 日, 周, 月, 年.
- Summary cards grouped by period.
- Detail view opens as full-screen page.
- AI generation button is visible only for the selected period.

## `/settings`

Mobile structure:

- Profile block.
- Account security list.
- Change password as a step-by-step form.

