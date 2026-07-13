---
page: product-record-detail-desktop
issue: 84
deviceType: DESKTOP
frame: 1440x1000
---

Create the canonical 1440px Product Web record-detail page for BabyGarden. Use
one realistic feeding record for `小满 · 8个月` and show how a caregiver can
review, edit or delete a core record safely. This is an MVP consumer page; do
not add AI explanations, photo timelines, realtime presence, object storage,
SMS or other external dependencies.

**DESIGN SYSTEM (REQUIRED):**

- Platform: responsive Web, desktop-first at 1440px with a later 390px detail
  variant. Do not use native app chrome or an Admin sidebar.
- Personality: warm, trustworthy, calm, practical; never childish or
  decorative for its own sake.
- Density: medium. Product controls and destructive-action confirmations are at
  least 44px high.
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

- Record example: `今天 10:30`, `喂养 · 配方奶`, `120ml`, note `宝宝胃口很好`;
  include source/timezone or local-save context honestly.
- Show family privacy cue `此记录对家庭成员可见` and a route back to `最近
  记录` plus `记录一笔`.
- Allow editing core fields in the same semantic model as quick record; delete
  must use an explicit confirmation step and explain that the action is
  irreversible.
- Offline or partial sync must not hide the saved local record; show a stale or
  waiting-to-sync label rather than inventing success from an external service.

**PAGE STRUCTURE:**

1. Product Web header with BabyGarden wordmark, baby context switcher and
   account menu; no marketing navigation.
2. Breadcrumb or back link `最近记录` followed by a focused heading such as
   `喂养记录详情` and timestamp.
3. Main white detail card with type icon, amount, time, note, privacy and sync
   status. Keep the information scannable, not dashboard-like.
4. Secondary action group: `编辑记录`, `删除记录`; include a compact delete
   confirmation state in the reference area.
5. Include loading, not-found/forbidden, session-expired, stale/offline and
   retryable error treatments without duplicating the full page.

**REQUIRED STATES:**

- Ready, loading skeleton, record not found, forbidden, session expired,
  local-only/stale sync, edit validation error and delete confirmation.
- Never leak a private record in signed-out/forbidden states; explain the safe
  next action.

**ANTI-SLOP CONSTRAINTS:**

- No Admin table density, purple gradients, glassmorphism, giant hero metrics,
  fake testimonials, decorative charts, emoji icons, lorem ipsum or stock
  photography.
- Keep labels Chinese-first and do not introduce English navigation or system
  preview copy unless it is an unavoidable technical token.
- Generated HTML is visual reference only, not production React code.
