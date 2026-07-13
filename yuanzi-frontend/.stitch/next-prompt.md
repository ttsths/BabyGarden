---
page: product-setup-mobile
issue: 84
deviceType: MOBILE
frame: 390x844
---

Create the critical 390px mobile-browser variant of the Product Web
first-family setup page for BabyGarden. Preserve the desktop form semantics and
the Nurturing Foundation visual hierarchy while making the form comfortable to
complete with one hand. This is responsive Web, not a native app screen.

**DESIGN SYSTEM (REQUIRED):**

- Platform: responsive Web; 390px is the mobile-browser breakpoint of the same
  Product Web shell.
- Personality: warm, trustworthy, calm, practical; never childish or
  decorative for its own sake.
- Density: comfortable. Inputs, segmented choices and primary actions are at
  least 44px high; use 16px side margins and readable Chinese line lengths.
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

**MVP FORM CONTRACT:**

- Collect only family name, caregiver display name, baby name, birth date and
  optional gender choice. Do not introduce SMS, OTP, phone verification, email
  verification, photos or external-service dependencies.
- Use visible Simplified Chinese labels, required markers, helper text and
  explicit date format. Keep entered values intact when showing errors.
- Make `创建家庭并继续` the clear primary action. A secondary `稍后完善`
  route is allowed only when its destination remains recoverable.
- Explain that records are private to family members.

**MOBILE LAYOUT:**

1. Compact header with BabyGarden wordmark, `第 1 步，共 1 步` and account
   menu. Do not render native status bars, bottom tabs or App navigation.
2. Put the heading and privacy explanation above one scrollable white form card.
3. Stack `家庭信息` and `宝宝信息` groups with 44px+ inputs, date picker
   affordance and a full-width gender choice row.
4. Keep the primary CTA visible after the final field and repeat a compact
   privacy cue near it; avoid a sticky footer that hides content.
5. Include a small state reference strip below the form for submitting,
   validation error and retryable network error; it is not a duplicate page.

**REQUIRED STATES:**

- Empty form, field validation errors, submitting, recoverable network error
  with local retry, signed-out/session-expired and forbidden states.
- Error copy must be actionable, local to the field or operation, and must not
  discard entered data.

**ANTI-SLOP CONSTRAINTS:**

- No purple gradients, glassmorphism, giant hero metrics, fake testimonials,
  decorative charts, emoji icons, lorem ipsum, stock photography or invented
  AI conclusions.
- Keep labels Chinese-first and do not introduce English navigation or system
  preview copy unless it is an unavoidable technical token.
- Generated HTML is visual reference only, not production React code.
