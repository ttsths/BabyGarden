---
page: product-record-detail-mobile
issue: 84
deviceType: MOBILE
frame: 390x844
---

Create the critical 390px mobile-browser record-detail page for BabyGarden.
Preserve the desktop detail hierarchy while making the core record readable and
editable with one hand. This is responsive Web, not a native app screen.

**DESIGN SYSTEM (REQUIRED):**

- Platform: responsive Web; 390px mobile-browser breakpoint of the Product Web
  shell. Do not render native status bars, bottom tabs or App navigation.
- Personality: warm, trustworthy, calm, practical; never childish or
  decorative for its own sake.
- Density: comfortable. Buttons and destructive confirmations are at least
  44px high; use 16px side margins.
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

- Use one realistic record: `小满 · 8个月`, `今天 10:30`, `喂养 · 配方奶`,
  `120ml`, note `宝宝胃口很好`.
- Keep `此记录对家庭成员可见`, local save/stale sync status, and links to
  `最近记录` and `记录一笔` visible.
- Provide `编辑记录` and `删除记录`; deleting requires a visible confirmation
  with irreversible-action copy. Do not add AI, photos, SMS or external
  dependencies.

**MOBILE LAYOUT:**

1. Compact header with BabyGarden wordmark, `小满 · 8个月` context and account
   control. No marketing navigation or App chrome.
2. Back link `最近记录`, heading `喂养记录详情`, timestamp and a focused white
   detail card with type, amount, time, note, privacy and sync status.
3. Stack `编辑记录` and `删除记录` actions with clear hierarchy and 44px
   targets; keep the destructive confirmation state in a compact reference
   area below the ready state.
4. Add loading, not-found/forbidden, session-expired, stale/offline,
   validation-error and delete-confirmation snippets without duplicating the
   complete page.

**REQUIRED STATES:**

- Ready, loading skeleton, record not found, forbidden, session expired,
  local-only/stale sync, edit validation error and delete confirmation.
- Never leak private records in signed-out or forbidden states; explain the
  safe next action.

**ANTI-SLOP CONSTRAINTS:**

- No bottom navigation, native status bar, admin table density, purple
  gradients, glassmorphism, giant hero metrics, fake testimonials, decorative
  charts, emoji icons, lorem ipsum or stock photography.
- Keep every visible label Chinese-first; no English navigation or footer.
- Generated HTML is visual reference only, not production React code.
