---
page: product-login-mobile
issue: 84
deviceType: MOBILE
frame: 390x844
---

Create the critical 390px mobile-browser login variant for BabyGarden, a private
family baby-care record product. This is a responsive Web screen, not a native
mobile App screen and not a marketing page. Preserve the accepted desktop login
master's hierarchy and content while making the first login action effortless
on a narrow viewport.

**DESIGN SYSTEM (REQUIRED):**

- Platform: responsive Web, desktop-first at 1440px with a critical 390px
  mobile-browser variant.
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
- Typography: Chinese-first system sans; 24/32 mobile page title, 16/26
  product body, 14/22 supporting text.
- Motion: 120-240ms strong ease-out; no page-wide entrance animation, scale
  from zero, bounce, or `transition: all`.

**ADMIN CONTRAST RULES:**

- This is Product Web, not Admin. Do not use an admin sidebar, dense table,
  compact 36px controls, or operational dashboard language.

**Page Structure:**

1. Keep a compact top header with the BabyGarden wordmark and no marketing
   navigation. The wordmark must not consume more than one row.
2. Stack the accepted desktop illustration message above the login card. Use a
   restrained flat coral-and-mint illustration; remove nonessential decoration
   rather than shrinking the desktop two-column composition.
3. Login card title `欢迎回到小园子` and supporting copy
   `登录后继续记录宝宝的成长。`.
4. Labeled username input `用户名` and password input `密码`, both full width
   and at least 44px high. Keep the visible password toggle inside the field.
5. Full-width primary button `登录` with a clear 44px touch target. Show the
   ready state as the primary composition and include a compact annotation strip
   for focus-visible, submitting and disabled states.
6. Account help copy `没有账号或忘记密码？请联系管理员。` Do not add public
   registration, SMS login, phone login, social login or forgot-password link.
7. Keep the privacy note `仅家庭成员可查看园子里的记录。` visible below the
   card without requiring a long scroll.

**Required States:**

- Default ready state is the main screen.
- Show a compact bottom state strip for invalid credentials, disabled account
  and expired session; it must remain readable at 390px without horizontal
  scrolling.
- Error copy must not reveal whether a username exists.

**Responsive Handoff Notes:**

- Preserve the desktop route: the mobile layout must be a stacked responsive
  variant, not a separate visual language.
- Use 16px side margins, full-width controls, and a minimum 12px vertical gap
  between labeled fields.
- Do not render native mobile status bars, bottom tabs, or App navigation.

**Anti-slop Constraints:**

- No purple gradients, glassmorphism, giant hero metrics, fake testimonials,
  decorative charts, emoji icons, lorem ipsum or stock photography.
- Do not redesign unrelated screens.
- Generated HTML is visual reference only.
