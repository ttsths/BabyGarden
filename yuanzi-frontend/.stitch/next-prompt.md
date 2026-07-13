---
page: product-records-desktop
issue: 84
deviceType: DESKTOP
frame: 1440x1100
---

Create the canonical 1440px Product Web records list for BabyGarden. It is the
history view reached from the home dashboard and must make core feeding, sleep
and diaper records scannable without depending on AI, photos, realtime state,
object storage or any external service. This is a consumer Product Web page,
not an Admin table.

**DESIGN SYSTEM (REQUIRED):**

- Platform: responsive Web, desktop-first at 1440px with a later 390px mobile
  list variant.
- Personality: warm, trustworthy, calm, practical; never childish or
  decorative for its own sake.
- Density: medium. Product controls and filters are at least 44px high; use
  generous grouping and clear scan paths instead of operational table density.
- Canvas: Warm Porcelain `#FFFBF7`.
- Primary action: Deep Warm Coral `#C84B42` with white text.
- Brand accent only: Soft Coral `#FF998A`; never use it behind white body text.
- Supporting accent: Soft Mint `#A2D5C6`.
- Primary text: Deep Cocoa `#23110F`; secondary text: Slate Graphite `#344054`.
- Surface: white with quiet `#E7D8D2` borders and restrained shadow.
- Focus: 2px adjacent-surface separator plus 2px outer Focus Blue `#2E90FA`.
- Radius: 10-12px controls, 16px cards, 20px large panels.
- Typography: Chinese-first system sans; 32/40 desktop page title, 20/28
  section title, 16/26 product body, 14/22 supporting text.
- Motion: 120-240ms strong ease-out; no page-wide entrance animation, scale
  from zero, bounce, or `transition: all`.

**MVP DATA CONTRACT:**

- Keep `小满 · 8个月` and family privacy context visible in the product header.
- Use realistic records with explicit dates, times and units, for example:
  喂养·配方奶 120ml at 10:30, 小睡 08:45–10:00, 尿布·嘘嘘 at 08:30.
- Provide a date-range control and lightweight type filters for `全部`, `喂养`,
  `睡眠`, `尿布`; do not add complex reporting or fake charts.
- Each row/card exposes type, timestamp, amount/duration and a clear route to
  view detail. Include a primary `记录一笔` action that opens quick record.
- Empty/new-baby state and partial/degraded core data state must remain useful;
  optional service failure cannot replace the list with an error page.

**PAGE STRUCTURE:**

1. Calm Product Web header with BabyGarden wordmark, current baby switcher and
   account menu. No admin sidebar, dense table or operational wording.
2. Heading `最近记录` with a concise last-sync/privacy cue and primary
   `记录一笔` action.
3. Filter row with date range and type chips. Keep active filter obvious and
   keyboard/focus friendly.
4. Group the vertical record list by date. Each white card has a visible type
   icon, Chinese label, time, amount/duration and an action to open detail.
5. Add a small state reference area for loading, empty, partial/degraded,
   retryable core error and session expired without duplicating the full page.

**REQUIRED STATES:**

- Loading skeleton, ready data, empty (new baby), filtered-empty, partial data
  with stale/sync note, retryable core error and session expired.
- Do not leak private records in signed-out or forbidden states; explain the
  next safe action.

**ANTI-SLOP CONSTRAINTS:**

- This is Product Web, not Admin: no compact 36px controls, spreadsheet-like
  dense rows, purple gradients, glassmorphism, giant hero metrics, fake
  testimonials, decorative charts, emoji icons, lorem ipsum or stock photos.
- Keep labels Chinese-first and do not introduce English navigation or system
  preview copy unless it is an unavoidable technical token.
- Generated HTML is visual reference only, not production React code.
