---
page: product-stats-desktop
issue: 84
deviceType: DESKTOP
frame: 1440x1152
---

Create the critical 1440px desktop Product Web statistics page for BabyGarden.
This follows the login → home → setup → quick record → records list → record
detail → baby profile flow. It is a responsive Web page, not a native App and
not an Admin dashboard.

**DESIGN SYSTEM (REQUIRED):**

- Platform: responsive Web, canonical desktop frame 1440x1152. Do not render
  native status bars, bottom tabs or App navigation.
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

**MVP DATA CONTRACT:**

- Active baby is `小满 · 8个月`. Use a fixed `最近7天` range by default and
  display realistic, clearly labeled data: 喂养 28次, 睡眠 63小时, 尿布 42次.
- Provide a range control for `最近7天` and `最近30天`, but do not add
  realtime, AI interpretation, medical advice or external analytics.
- Use one accessible trend visualization (simple bars/lines with visible
  labels and a text summary) for feeding, sleep and diaper activity. Do not
  create decorative charts or unsupported precision.
- Include a clear link to `记录一笔` and `最近记录`; explain that statistics
  are derived from saved records and may be incomplete when offline.

**DESKTOP LAYOUT:**

1. Use the established BabyGarden Product Web shell: compact brand header,
   active baby context, account control and breadcrumb `首页 / 统计`.
2. Start with a concise title `成长统计` and range selector. Show three summary
   cards for 喂养、睡眠、尿布 with period labels and a small comparison hint.
3. Place the main trend card below with tabs or a segmented control for the
   three record types. Keep legends and data labels readable; provide a text
   summary such as `近7天平均每天喂养4次` for screen-reader parity.
4. Add a compact degraded/offline note and the visible state-reference section
   below the ready view: loading skeleton, empty, no records, partial data,
   error with retry, and signed-out/session-expired. These are references only.

**REQUIRED STATES:**

- Ready, loading, empty/no records, partial/offline data, error with retry,
  forbidden/read-only and session expired.
- Never infer diagnosis or recommendations. Never leak private baby data in
  signed-out or forbidden states; explain the safe next action.

**ANTI-SLOP CONSTRAINTS:**

- No bottom navigation, native status bar, admin table density, purple
  gradients, glassmorphism, giant hero metrics, fake testimonials, emoji icons,
  lorem ipsum or stock photography.
- Keep every visible label Chinese-first; no English navigation, state labels
  or footer text.
- Generated HTML is visual reference only, not production React code.
