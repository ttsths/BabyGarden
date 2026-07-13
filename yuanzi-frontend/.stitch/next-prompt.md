---
page: product-setup-desktop
issue: 84
deviceType: DESKTOP
frame: 1440x1000
---

Create the canonical 1440px Product Web first-family setup page for BabyGarden.
It appears after username/password login when a new account has no family or
baby context yet. This is an MVP onboarding form for one family and one baby;
do not add SMS, invitations that require external services, AI, photos or
other optional dependencies.

**DESIGN SYSTEM (REQUIRED):**

- Platform: responsive Web, desktop-first at 1440px with a critical 390px
  mobile-browser variant later.
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

**MVP FORM CONTRACT:**

- Collect only the minimum data needed to unlock the product: family name,
  caregiver display name, baby name, birth date, and optional gender choice.
- Use realistic labels and helper text in Simplified Chinese. Use explicit
  date format and inline validation; never rely on placeholder-only labels.
- Make the primary route obvious: `创建家庭并继续`. Provide a safe secondary
  route such as `稍后完善` only if it leaves the user in a known recoverable
  state.
- Explain that the data is private to the family. Do not require SMS, email
  verification, object storage, realtime presence or external integrations.

**PAGE STRUCTURE:**

1. Calm product header with BabyGarden wordmark, a simple progress cue such as
   `第 1 步，共 1 步`, and account menu. No marketing navigation or admin
   sidebar.
2. Center a generous setup panel with a concise heading such as `先认识一下
   你的家庭` and supporting copy that sets expectations.
3. Group fields into `家庭信息` and `宝宝信息`; include visible labels,
   accessible required markers, and examples that do not masquerade as values.
4. Show a small privacy note and a compact preview of the next home experience
   without inventing charts or AI conclusions.
5. Include the ready, loading/submitting, validation-error and recoverable
   network-error treatments in a compact state reference strip.

**REQUIRED STATES:**

- Empty form, field validation errors, submitting, recoverable network error
  with local retry, signed-out/session-expired and forbidden states.
- Errors must be actionable, local to the field or operation, and must not
  discard already-entered values.

**ADMIN CONTRAST RULES:**

- This is Product Web, not Admin. Do not use dense tables, compact 36px
  controls, operational dashboard language or a fixed admin sidebar.

**ANTI-SLOP CONSTRAINTS:**

- No SMS/OTP flow, phone-number verification, purple gradients, glassmorphism,
  giant hero metrics, fake testimonials, decorative charts, emoji icons,
  lorem ipsum, stock photography or invented AI conclusions.
- Keep all labels Chinese-first; do not introduce English navigation or system
  preview copy unless it is an unavoidable technical token.
- Generated HTML is visual reference only, not production React code.
