---
page: product-family-desktop
issue: 84
deviceType: DESKTOP
frame: 1440x1152
---

Create the critical 1440px desktop Product Web family-members page for
BabyGarden. This follows the completed Product Web core-recording and
statistics flow. It is a responsive Web page, not a native App and not an
Admin console.

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

- Active context is `小满的家 · 3位成员`, with baby `小满 · 8个月`.
- Show three realistic members: `林悦` as `家庭管理员 · 妈妈`, `陈昊` as
  `可编辑 · 爸爸`, and `王阿姨` as `仅查看 · 照护者`.
- MVP invitations must not depend on SMS, email, contacts or social login.
  Provide an internal `生成邀请口令` action and a ready-state invitation card
  with copyable code `XM-8M-2741`, expiry `24小时后失效`, role selection and
  `复制口令`. Explain that the inviter shares it manually.
- Allow the administrator to edit roles and remove a member. Removing a member
  requires a visible irreversible-action confirmation. A user cannot remove
  themselves while they are the only administrator.

**DESKTOP LAYOUT:**

1. Use the established BabyGarden Product Web shell: compact brand header,
   active baby context, account control and breadcrumb `首页 / 家庭成员`.
2. Start with title `家庭成员` and a concise privacy explanation. Use a
   two-column workspace: member list as the main column and invitation/role
   guidance as the supporting column.
3. Member rows must show name, relationship, role, last activity, explicit
   `调整权限` and overflow/removal actions. Clearly distinguish `家庭管理员`,
   `可编辑` and `仅查看` without using color alone.
4. Show the generated invitation code card as a separate ready state. Keep
   invitation generation/copying internal and manual; do not add delivery
   channels or external services.
5. Add a compact state-reference section below: loading, no other members,
   read-only current user, invitation expired, permission denied, session
   expired, role update failure and remove-member confirmation.

**REQUIRED STATES:**

- Ready, loading, only-member empty state, read-only, invitation ready,
  invitation expired, role update failure, forbidden, session expired and
  destructive removal confirmation.
- Never expose family membership or invitation codes to signed-out or
  forbidden users; explain the safe next action.

**ANTI-SLOP CONSTRAINTS:**

- No bottom navigation, native status bar, Admin table density, purple
  gradients, glassmorphism, giant hero metrics, fake testimonials, emoji icons,
  lorem ipsum, stock photography or realtime presence.
- No SMS/email/contact integration, public registration, AI recommendations or
  external dependency.
- Keep every visible label Chinese-first; no English navigation, state labels
  or footer text.
- Generated HTML is visual reference only, not production React code.
