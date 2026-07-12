---
page: product-login-desktop
issue: 84
deviceType: DESKTOP
frame: 1440x1024
---

Create the canonical 1440px desktop login screen for BabyGarden, a private
family baby-care record product. The experience should feel warm, trustworthy,
calm and practical. This is a real product entry screen, not a marketing page.

**DESIGN SYSTEM (REQUIRED):**

- Responsive Web, desktop-first at 1440px; Product density is medium.
- Warm Porcelain `#FFFBF7` page canvas.
- Deep Warm Coral `#C84B42` primary action with white text.
- Soft Coral `#FF998A` is a brand accent only and must not sit behind white body
  text.
- Soft Mint `#A2D5C6` is a quiet supporting accent.
- Deep Cocoa `#23110F` primary text; Slate Graphite `#344054` secondary text.
- White surfaces use quiet `#E7D8D2` borders and restrained shadows.
- Focus-visible uses a 2px adjacent-surface separator and a 2px outer Focus
  Blue `#2E90FA` ring.
- Controls are at least 44px high; use 10-12px control radius and a 16px card
  radius.
- Chinese-first system sans typography: 32/40 page title, 16/26 body, 14/22
  labels and supporting text.
- Interaction feedback is subtle 120-240ms strong ease-out. No page entrance
  animation, bounce, scale-from-zero or `transition: all`.

**Page Structure:**

1. A restrained top-left BabyGarden wordmark with a tiny coral-and-mint garden
   mark. No large navigation or marketing links.
2. A balanced two-column desktop composition. The left side contains one calm
   caregiver message and a simple abstract growth illustration made from flat
   shapes, not a remote photo. The right side contains the login card.
3. Login card title `欢迎回到小园子` and supporting copy
   `登录后继续记录宝宝的成长。`.
4. Labeled username input `用户名` and password input `密码`, with visible
   password toggle and realistic focus/validation treatment.
5. Full-width primary button `登录`. Include default, hover, focus-visible,
   submitting and disabled button examples in a compact annotation strip.
6. Account help copy `没有账号或忘记密码？请联系管理员。` Do not add public
   registration, SMS login, phone login, social login or forgot-password link.
7. A quiet privacy note at the card bottom: `仅家庭成员可查看园子里的记录。`.

**Required States:**

- Default ready state is the main screen.
- Show a compact adjacent state panel for invalid credentials, disabled account
  and expired session; each state explains the next safe action.
- Error copy must not reveal whether a username exists.

**Responsive Handoff Notes:**

- Preserve a clear route to a later 390px variant: stack illustration copy
  above the card, remove nonessential decoration and keep controls full width.
- Do not render native mobile status bars or App navigation.

**Anti-slop Constraints:**

- No purple gradients, glassmorphism, giant hero metrics, fake testimonials,
  decorative charts, emoji icons, lorem ipsum or stock photography.
- Do not redesign unrelated screens.
- Generated HTML is visual reference only.
