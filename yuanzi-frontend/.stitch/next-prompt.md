---
page: product-record-mobile
issue: 84
deviceType: MOBILE
frame: 390x844
---

Create the critical 390px mobile-browser quick-record drawer for BabyGarden.
It is opened from the Product Web home dashboard and must preserve the current
baby context while allowing a caregiver to record one feeding, sleep session
or diaper change with one hand. This is responsive Web, not a native app
screen.

**DESIGN SYSTEM (REQUIRED):**

- Platform: responsive Web; 390px is the mobile-browser breakpoint of the same
  Product Web shell. Use a bottom sheet or full-height drawer, never native app
  chrome or a bottom tab bar.
- Personality: warm, trustworthy, calm, practical; never childish or
  decorative for its own sake.
- Density: comfortable. Inputs, tabs and primary actions are at least 44px
  high; use 16px side margins and preserve readable Chinese line lengths.
- Canvas: Warm Porcelain `#FFFBF7`.
- Primary action: Deep Warm Coral `#C84B42` with white text.
- Brand accent only: Soft Coral `#FF998A`; never use it behind white body text.
- Supporting accent: Soft Mint `#A2D5C6`.
- Primary text: Deep Cocoa `#23110F`; secondary text: Slate Graphite `#344054`.
- Surface: white with quiet `#E7D8D2` borders and restrained shadow.
- Focus: 2px adjacent-surface separator plus 2px outer Focus Blue `#2E90FA`.
- Radius: 10-12px controls, 16px cards, 20px drawer.
- Typography: Chinese-first system sans; 24/32 drawer title, 18/26 section
  title, 16/26 product body, 14/22 supporting text.
- Motion: 120-240ms strong ease-out; no page-wide entrance animation, scale
  from zero, bounce, or `transition: all`.

**MVP RECORD CONTRACT:**

- Keep `小满 · 8个月` and family privacy context visible in the drawer.
- Provide three Chinese modes: `喂养`, `睡眠`, `尿布`; only the relevant fields
  change when the mode switches.
- 喂养: 记录时间、喂养类型（母乳/配方奶/辅食）、喂奶量（ml）和可选备注。
- 睡眠: 开始时间、结束时间或时长和可选备注。
- 尿布: 记录时间、类型（尿布/嘘嘘/便便）和可选备注。
- Primary action `保存记录`; secondary action `取消`. A local-save/sync
  message must be honest and must not require an external service.

**MOBILE LAYOUT:**

1. Keep the underlying page dimmed only enough to preserve context; the drawer
   header should state `记录小满的喂养` and include a 44px close target.
2. Place the three mode tabs near the top, with clear active state and icons.
3. Stack visibly labelled fields, explicit units and accessible input affordances.
   Keep `保存记录` and `取消` reachable without hiding the last field.
4. Show a privacy cue `此记录对家庭成员可见` and local sync status near the
   actions.
5. Include a compact reference strip below the drawer for validation, saving,
   saved confirmation and retryable network error. Do not duplicate the whole
   page for each state.

**REQUIRED STATES:**

- Empty/default form, field validation error, saving, saved confirmation,
  retryable network error and session expired. Preserve entered values on error.
- Optional service failure must not block core local record creation.

**ANTI-SLOP CONSTRAINTS:**

- No SMS/OTP, phone verification, purple gradients, glassmorphism, giant hero
  metrics, fake testimonials, decorative charts, emoji icons, lorem ipsum,
  stock photography or invented AI conclusions.
- Keep labels Chinese-first and do not introduce English navigation or system
  preview copy unless it is an unavoidable technical token.
- Generated HTML is visual reference only, not production React code.
