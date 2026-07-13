---
page: product-record-desktop
issue: 84
deviceType: DESKTOP
frame: 1440x900
---

Create the canonical 1440px Product Web quick-record vertical slice for
BabyGarden. It is opened from the home dashboard primary actions and lets a
caregiver record one feeding, sleep session or diaper change without leaving
the current baby context. This is an MVP interaction design; do not add SMS,
AI, photos, object storage, realtime presence or other external dependencies.

**DESIGN SYSTEM (REQUIRED):**

- Platform: responsive Web, desktop-first at 1440px with a later 390px drawer
  variant. This is not a native mobile app.
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
- Radius: 10-12px controls, 16px cards, 20px dialogs.
- Typography: Chinese-first system sans; 24/32 page title, 20/28 section
  title, 16/26 product body, 14/22 supporting text.
- Motion: 120-240ms strong ease-out; no page-wide entrance animation, scale
  from zero, bounce, or `transition: all`.

**MVP RECORD CONTRACT:**

- Keep baby context explicit: `小满 · 8个月`; the record is owned by the
  current family and timestamped locally before any sync.
- Provide three modes with Chinese labels: `喂养`, `睡眠`, `尿布`. The active
  mode must change only the minimal relevant fields.
- 喂养 fields: record time, feeding type (母乳/配方奶/辅食), amount in ml when
  applicable, and optional note.
- 睡眠 fields: start time, end time or duration, and optional note.
- 尿布 fields: record time, type (尿布/嘘嘘/便便) and optional note.
- Make `保存记录` primary and `取消` secondary. Do not require external
  services, invite flows, media uploads or realtime confirmation.

**PAGE STRUCTURE:**

1. Preserve the Product Web header and current baby/family context behind a
   focused modal or right-side dialog; do not use an Admin sidebar.
2. Modal header states the action, e.g. `记录小满的喂养`, with a close button
   and a visible mode switcher for 喂养/睡眠/尿布.
3. Use grouped, visibly labelled inputs with explicit units and examples.
   Keep the primary action within the dialog's comfortable reach.
4. Show a privacy cue that the record is visible to family members and a small
   local-save/sync note that remains honest when the network is unavailable.
5. Include compact reference treatments for validation error, saving,
   successfully saved and retryable network error without duplicating the
   entire page.

**REQUIRED STATES:**

- Empty/default form, field validation error, saving, saved confirmation,
  retryable network error and session expired. Preserve entered values on error.
- A degraded optional service must not block local core record creation; do not
  present a fake realtime or AI result.

**ANTI-SLOP CONSTRAINTS:**

- No SMS/OTP, phone verification, purple gradients, glassmorphism, giant hero
  metrics, fake testimonials, decorative charts, emoji icons, lorem ipsum,
  stock photography or invented AI conclusions.
- Keep labels Chinese-first and do not introduce English navigation or system
  preview copy unless it is an unavoidable technical token.
- Generated HTML is visual reference only, not production React code.
