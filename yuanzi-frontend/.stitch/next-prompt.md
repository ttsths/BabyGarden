---
page: product-records-mobile
issue: 84
deviceType: MOBILE
frame: 390x844
---

Create the critical 390px mobile-browser records history page for BabyGarden.
Preserve the Product Web hierarchy of the desktop records list while making
date filtering and record scanning comfortable on a narrow screen. This is a
responsive Web page, not a native app screen.

**DESIGN SYSTEM (REQUIRED):**

- Platform: responsive Web; 390px mobile-browser breakpoint of the Product Web
  shell. Do not render a native status bar, bottom tab bar or App navigation.
- Personality: warm, trustworthy, calm, practical; never childish or
  decorative for its own sake.
- Density: comfortable. Filters, primary actions and list actions are at least
  44px high; use 16px side margins and clear scan paths.
- Canvas: Warm Porcelain `#FFFBF7`.
- Primary action: Deep Warm Coral `#C84B42` with white text.
- Brand accent only: Soft Coral `#FF998A`; never use it behind white body text.
- Supporting accent: Soft Mint `#A2D5C6`.
- Primary text: Deep Cocoa `#23110F`; secondary text: Slate Graphite `#344054`.
- Surface: white with quiet `#E7D8D2` borders and restrained shadow.
- Focus: 2px adjacent-surface separator plus 2px outer Focus Blue `#2E90FA`.
- Radius: 10-12px controls, 16px cards, 20px large panels.
- Typography: Chinese-first system sans; 24/32 mobile page title, 18/26
  section title, 16/26 product body, 14/22 supporting text.
- Motion: 120-240ms strong ease-out; no page-wide entrance animation, scale
  from zero, bounce, or `transition: all`.

**MVP DATA CONTRACT:**

- Keep `小满 · 8个月`, family privacy cue and a clear route to `记录一笔`.
- Use explicit Chinese sample records: 喂养·配方奶 10:30 120ml, 小睡
  08:45–10:00 1小时15分, 尿布·嘘嘘 08:30.
- Provide a compact date-range control and horizontally scrollable or wrapped
  type filters: `全部`, `喂养`, `睡眠`, `尿布`.
- Group records by date and show type, timestamp, amount/duration and detail
  affordance in stacked cards. Avoid dense table rows.

**MOBILE LAYOUT:**

1. Compact header with BabyGarden wordmark, `小满 · 8个月` switcher and
   account menu; no marketing or admin navigation.
2. Heading `最近记录`, last-sync/privacy cue and a full-width or prominent
   `记录一笔` action near the top.
3. Put date range and type filters in a scroll-friendly region without hiding
   active selections.
4. Render date-grouped record cards with 44px detail targets and enough room
   for Chinese labels, units and notes.
5. Add a small state reference strip below the list for loading, empty,
   filtered-empty, partial/stale data and retryable error; do not duplicate the
   whole page four times.

**REQUIRED STATES:**

- Loading skeleton, ready data, empty new-baby state, filtered-empty,
  partial/stale data, retryable core error, signed-out and session-expired.
- Do not leak private records in signed-out or forbidden states; explain the
  next safe action.

**ANTI-SLOP CONSTRAINTS:**

- No admin table density, compact 36px controls, purple gradients,
  glassmorphism, giant hero metrics, fake testimonials, decorative charts,
  emoji icons, lorem ipsum or stock photography.
- Keep labels Chinese-first and do not introduce English navigation or system
  preview copy unless it is an unavoidable technical token.
- Generated HTML is visual reference only, not production React code.
