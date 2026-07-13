---
page: product-home-mobile
issue: 84
deviceType: MOBILE
frame: 390x844
---

Create the critical 390px mobile-browser variant of the accepted Product Web
home dashboard for BabyGarden. Keep the desktop home hierarchy and semantics,
but adapt density and layout for a 390x844 logical viewport. This is a
responsive Web page, not a native mobile app screen.

**DESIGN SYSTEM (REQUIRED):**

- Platform: responsive Web; 390px is a mobile-browser breakpoint of the same
  Product Web shell.
- Personality: warm, trustworthy, calm, practical; never childish or
  decorative for its own sake.
- Density: comfortable. All controls and primary record actions are at least
  44px high; keep 16px side margins and readable Chinese line lengths.
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

- Use realistic Chinese sample data for one family and one baby: `小满`, 8
  个月, with explicit timestamps and units.
- Keep the baby context, today's summary, recent records and quick actions for
  `喂养`, `睡眠`, and `尿布` as the first-class content.
- Do not invent AI insights, photo timelines, realtime presence or fake
  charts. Optional services may be omitted or shown as a small clearly
  labelled unavailable notice.
- If optional services fail, the core record actions and any successful summary
  remain usable; never replace the whole page with an error state.

**MOBILE LAYOUT:**

1. Use a compact header with the BabyGarden wordmark, `小满 · 8个月` context,
   a baby switcher and account menu. Do not render a native status bar, bottom
   tab bar or App navigation.
2. Keep the greeting and privacy/last-sync cue concise above the fold.
3. Stack the quick-record panel near the top. Show three prominent 44px
   actions labelled exactly `喂养`, `睡眠`, and `尿布`, with visible icons and
   enough spacing for one-handed tapping.
4. Stack `今日概览` cards in a readable two-column-or-single-column rhythm;
   preserve units and the time range for feeding, sleep and diaper values.
5. Show `最近记录` as a scannable vertical list with type, timestamp and
   amount/duration, plus a clear `查看全部记录` route.
6. Keep a compact `家庭成员` card useful without realtime; explicitly mark
   stale or unavailable collaboration data.
7. Include a small state annotation section below the core content showing
   loading, empty, degraded optional-service and retry treatments. It is a
   reference strip, not four duplicated full pages.

**REQUIRED STATES:**

- Data: loading, empty (new baby with no records), ready, partial/degraded and
  core error with local retry.
- Permission/session: signed out, forbidden and session expired must explain a
  safe next action without leaking private data.

**ANTI-SLOP CONSTRAINTS:**

- No purple gradients, glassmorphism, giant hero metrics, fake testimonials,
  decorative charts, emoji icons, lorem ipsum, stock photography or invented
  AI conclusions.
- Keep labels Chinese-first and do not introduce English navigation or system
  preview copy unless it is an unavoidable technical token.
- Do not make optional services more prominent than core records.
- Generated HTML is visual reference only, not production React code.
