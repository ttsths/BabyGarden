---
page: product-home-desktop
issue: 84
deviceType: DESKTOP
frame: 1440x1200
---

Create the canonical 1440px Product Web home dashboard for BabyGarden, a
private family baby-care record product. This is the MVP home after a user has
logged in and completed family/baby setup. It must be useful with only core
record data; optional AI, photos, object storage and realtime services must not
block the page.

**DESIGN SYSTEM (REQUIRED):**

- Platform: responsive Web, desktop-first at 1440px with a critical 390px
  mobile-browser variant.
- Personality: warm, trustworthy, calm, practical; never childish or
  decorative for its own sake.
- Density: medium. Product controls are at least 44px high and touch targets
  are at least 44x44px.
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

**ADMIN CONTRAST RULES:**

- This is Product Web, not Admin. Do not use an admin sidebar, dense tables,
  compact 36px controls, or operational dashboard language.

**MVP DATA CONTRACT:**

- Use realistic Chinese sample data for one family and one baby, such as
  `小满`, 8 months old, with explicit timestamps and units.
- Core content is baby context, today's summary, recent records and quick
  record actions for feeding, sleep and diaper records.
- Do not invent AI insights, photo timelines, realtime presence or fake charts.
- Optional services may appear only as a clearly labelled unavailable card or
  be omitted; their failure must not turn the core page into an error screen.

**Page Structure:**

1. A calm product header with the BabyGarden wordmark, current family/baby
   context, a compact baby switcher and account menu. No marketing navigation.
2. A page heading such as `早上好，今天也一起记录小满的成长` with last-sync
   context and a clear privacy cue that records are visible to family members.
3. A prominent quick-record panel with three labelled primary actions: `喂养`,
   `睡眠`, and `尿布`. Each action must have a visible icon and a 44px target.
4. A `今日概览` section with three or four statistic cards: feeding count and
   volume, sleep duration, diaper count, and a compact growth signal only when
   real data exists. Show units and time range; do not use decorative metrics.
5. A `最近记录` section with a scannable list showing record type, timestamp,
   amount or duration, and a link to view all records. Include a clear empty
   state variant in the state panel.
6. A small secondary panel for `家庭成员` or recent collaboration status that
   remains useful without realtime; label stale or unavailable data explicitly.
7. Include a compact annotation strip showing the loading skeleton, empty,
   degraded optional-service and retry states without duplicating the complete
   page four times.

**Required States:**

- Data: loading, empty (new baby with no records), ready, partial/degraded and
  core error with local retry.
- Permission/session: signed out, forbidden and session expired must explain a
  safe next action without leaking private data.
- The degraded state must keep baby context, quick record and any successful
  core statistics usable when optional services are unavailable.

**Responsive Handoff Notes:**

- Preserve a clear route to a later 390px variant: stack the quick-record panel
  and statistic cards, keep the primary action visible above the fold, and use
  16px side margins.
- Do not render native mobile status bars, bottom tabs, or App navigation.
- Keep the information hierarchy and semantic colors identical across widths;
  only density and layout should change.

**Anti-slop Constraints:**

- No purple gradients, glassmorphism, giant hero metrics, fake testimonials,
  decorative charts, emoji icons, lorem ipsum, stock photography or invented
  AI conclusions.
- Do not make optional services visually more prominent than core records.
- Generated HTML is visual reference only, not production React code.
