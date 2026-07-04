# Life Tracker Frontend Design

This folder contains frontend planning assets before writing the Vue 3 code.

## Structure

- `page/00-information-architecture.md`: product routes, navigation, and API mapping.
- `page/01-desktop-pages.md`: desktop web page descriptions.
- `page/02-mobile-pages.md`: mobile web page descriptions.
- `image/desktop-dashboard.png`: generated desktop UI reference.
- `image/mobile-today.png`: generated mobile UI reference.

## Design Direction

Life Tracker should feel like a practical personal dashboard, not a marketing site.
The first screen after login should let the user record life events and expenses immediately, then review today's timeline, spending, and AI summary.

Responsive rules:

- Desktop: sidebar navigation plus dense dashboard layout.
- Tablet: compact sidebar or top navigation, two-column content where space allows.
- Mobile: bottom navigation, single-column content, thumb-friendly quick actions.

