# BabyGarden Stitch Generation Contract

This file is the compact, prompt-ready companion to `../DESIGN.md`. Copy the
required block below into every Stitch request. Do not submit a local path as
design context: Stitch prompts must be self-contained.

## Product Web

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
- Typography: Chinese-first system sans; 32/40 desktop page title, 20/28
  section title, 16/26 product body, 14/22 supporting text.
- Motion: 120-240ms strong ease-out; no page-wide entrance animation, scale
  from zero, bounce, or `transition: all`.

## Admin

- Platform: desktop Web from 1024px upward; canonical frame 1440px and compact
  anchor 1024px.
- Personality: neutral, operational, scan-focused; shares semantic colors with
  Product but uses less decoration.
- Density: compact. Controls 36-40px high, 8px spacing rhythm, 6-8px control
  radius and 12px panel radius.
- Tables: sticky identifier/action context, contained horizontal scrolling at
  1024px, visible filters, pagination and explicit empty/error states.
- Use Deep Warm Coral only for semantic primary actions; do not turn the admin
  shell into a coral-themed dashboard.

## Shared State Packs

- Data: loading, empty, ready, degraded, error, retry.
- Form: default, focus, validation, submitting, success, failure, disabled.
- External: not configured, unavailable, quota exhausted, offline, retry.
- Identity: signed out, forbidden, read-only, role-limited, session expired.
- Media: idle, selected, uploading, success, failure, invalid size/type.

## Content and MVP Boundaries

- Use real, concise Chinese product copy. Name the affected baby or family
  when context could be ambiguous.
- Login contains username and password only. No SMS login, phone login, public
  registration or social login.
- AI, photos, realtime, dark mode and elder mode may appear only as explicit P1
  variants or unavailable states. They must not dominate a P0 screen.
- Do not use emoji as functional icons. Use a consistent outline icon family
  and visible labels where meaning could be ambiguous.
- Do not use purple AI gradients, glassmorphism, decorative analytics, fake
  charts, placeholder lorem ipsum, remote stock-photo URLs or native-App chrome.
- Generated HTML is reference output, not production React code.
