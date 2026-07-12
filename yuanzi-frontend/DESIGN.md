# Design System: BabyGarden Web + Admin

- **Stitch source:** `stitch_yuanzi_baby_app` local export
- **Project ID:** Not available in the local export
- **Status:** Approved design direction; canonical for new design work, implementation pending
- **Scope:** Desktop-first responsive Web and desktop Admin; native App is deferred
- **Design direction:** Shared semantic foundation, Warm Trust product surface, Neutral Efficient admin surface

This document is the source of truth for prompting Stitch and reviewing BabyGarden UI. Existing files such as `DESIGN_TOKENS.md`, the legacy mobile pages, and hard-coded values in `open-design.css` are implementation evidence, not competing design authorities.

## 1. Visual Theme & Atmosphere

BabyGarden should feel warm, calm, trustworthy, and quietly competent. It serves families recording repeated childcare events, so it must be emotionally considerate without becoming childish, decorative, or slow.

The product Web surface uses a warm porcelain canvas, crisp white surfaces, restrained coral actions, soft mint status accents, and deep navy text. Cards feel gently curved and lightly elevated. Information is neither sparse nor dashboard-dense: the primary task and current baby context are obvious, while recent records and trends remain scannable.

The Admin surface shares the same semantic color and status language but uses a cooler neutral canvas, tighter spacing, smaller radii, flatter elevation, and denser tables. It should feel operational rather than playful.

The design is intentionally not:

- a native mobile app interface stretched onto desktop;
- a clinical hospital dashboard;
- a pastel illustration gallery;
- a collection of rounded cards without hierarchy;
- a default AI aesthetic built from purple gradients, glass panels, or decorative statistics.

## 2. Color Palette & Roles

### Light foundation

| Descriptive name     | Value     | Semantic role                                                                                 |
| -------------------- | --------- | --------------------------------------------------------------------------------------------- |
| Deep Warm Coral      | `#C84B42` | Primary action, selected control, important text link; supports white text at 4.62:1          |
| Pressed Brick Coral  | `#B93F36` | Hover and pressed state for primary actions                                                   |
| Soft Coral Bloom     | `#FF998A` | Brand illustration, chart accent, avatar ring, non-text decoration; never a white-text button |
| Petal Wash           | `#FDE7E2` | Soft selected background, product highlight surface                                           |
| Porcelain Warm White | `#FFFCF9` | Product Web page canvas                                                                       |
| Warm Linen           | `#F8F3EF` | Product secondary canvas and quiet grouped region                                             |
| Clean White          | `#FFFFFF` | Cards, dialogs, inputs, tables                                                                |
| Admin Mist           | `#F6F7F9` | Admin page canvas                                                                             |
| Deep Ink Navy        | `#172033` | Primary headings and body text                                                                |
| Slate Graphite       | `#344054` | Strong secondary text and labels                                                              |
| Quiet Slate          | `#667085` | Secondary text; supports white background at 4.97:1                                           |
| Warm Hairline        | `#E8E1DB` | Product borders and separators                                                                |
| Cool Hairline        | `#DFE3E8` | Admin borders and table separators                                                            |
| Trustworthy Mint     | `#2F7D68` | Success, healthy, connected, completed; supports white text at 4.94:1                         |
| Mint Veil            | `#DFF2EB` | Success/status background                                                                     |
| Warning Amber        | `#A15C00` | Warning text, pending attention                                                               |
| Danger Red           | `#B42318` | Destructive action and error text                                                             |
| Information Blue     | `#175CD3` | Informational status and utility link                                                         |
| Focus Blue           | `#2E90FA` | Outer keyboard focus ring used with a 2px surface-colored separator                           |

`#FF998A`, the dominant color in the original Stitch export, has only 2.06:1 contrast against white. Keep it as a recognizable brand accent, but use Deep Warm Coral (`#C84B42`) for accessible filled controls and links.

### Dark mode

Dark mode is a P1 cross-cutting variant. It should remain warm-neutral, not brown-on-brown.

| Descriptive name | Value     | Semantic role                                                        |
| ---------------- | --------- | -------------------------------------------------------------------- |
| Midnight Canvas  | `#17181C` | Dark page canvas                                                     |
| Charcoal Surface | `#22242A` | Cards, navigation, dialogs                                           |
| Raised Charcoal  | `#2B2E35` | Elevated or selected surface                                         |
| Dark Hairline    | `#373A42` | Borders and separators                                               |
| Frosted Text     | `#F7F7F8` | Primary text                                                         |
| Silver Text      | `#AEB4BF` | Secondary text                                                       |
| Luminous Coral   | `#FF998A` | Dark-mode primary action with Deep Ink Navy text; contrast is 7.88:1 |
| Luminous Mint    | `#A2D5C6` | Dark-mode success and healthy state with Deep Ink Navy text          |

### Color discipline

- Color communicates action or state before decoration.
- Status never relies on color alone; pair it with a label or icon.
- Product charts use coral, mint, blue, and amber in that order. Do not invent a rainbow palette per chart.
- Admin tables use quiet neutrals; reserve semantic colors for statuses, exceptions, and actions.
- Gradients are not a default background treatment. A subtle tonal wash is allowed only in a rare hero or explanatory empty state.

## 3. Typography Rules

### Families

- Display and Latin numerals: `Plus Jakarta Sans`, then the body stack.
- Chinese body: `Noto Sans SC`, `PingFang SC`, `Microsoft YaHei`, sans-serif.
- Do not introduce Inter, Roboto, a serif display face, or a handwritten face without a separate brand decision.

### Type scale

| Role             | Size / line height | Weight  | Use                                                           |
| ---------------- | ------------------ | ------- | ------------------------------------------------------------- |
| Display          | `40 / 48px`        | 700     | Rare product hero or major onboarding message                 |
| Page title       | `32 / 40px`        | 700     | Desktop product page title                                    |
| Admin page title | `24 / 32px`        | 700     | Admin route heading                                           |
| Section title    | `20 / 28px`        | 650–700 | Product section and major card heading                        |
| Component title  | `16 / 24px`        | 600–650 | Card, dialog, table panel                                     |
| Product body     | `16 / 26px`        | 400–500 | Primary explanatory copy                                      |
| Standard body    | `14 / 22px`        | 400–500 | Forms, records, admin content                                 |
| Caption          | `12 / 18px`        | 500     | Metadata and secondary time labels; not for essential actions |

Use sentence case for Chinese labels and English helper text. Avoid all caps except short technical eyebrow labels. Use tabular numerals for time, measurement, and dashboard metrics.

## 4. Component Stylings

### Buttons

- Product primary: Deep Warm Coral fill, white label, 44px minimum height, 10–12px radius.
- Admin primary: same semantic primary color, 36–40px height, 6–8px radius.
- Secondary: white surface, visible border, Deep Ink Navy label.
- Quiet/ghost: transparent background; reveal a subtle tonal background on hover.
- Destructive: Danger Red is reserved for irreversible actions; require confirmation when data loss is possible.
- Press feedback: scale to `0.97–0.98` for 120ms. Never animate from `scale(0)`.
- Every variant supports default, hover, focus-visible, active, disabled, loading, and error-adjacent states.

### Cards and containers

- Product cards use 12–16px gently curved corners. Important grouped surfaces may use 20px; do not use 28px everywhere.
- Admin cards and table panels use 6–10px corners.
- Product elevation is whisper-soft and diffused; borders remain visible so cards do not float without structure.
- Admin is primarily flat with hairline borders. Shadow is reserved for overlays and sticky regions.
- Do not place every section inside a card. Use spacing, headings, and background regions to establish hierarchy.

### Inputs and forms

- Product controls are at least 44px high; Admin controls are 36–40px high.
- Inputs use a clean white surface, visible border, and a dual focus treatment: 2px surface-colored separation plus a 2px Focus Blue outer ring.
- Labels remain visible above the control; placeholders are examples, not replacements for labels.
- Validation appears near the field and in an optional form summary for multi-field failures.
- Required, optional, unit, and format expectations are explicit.

### Navigation

- Desktop product uses a persistent, restrained side navigation when five or more primary destinations are present.
- At narrow widths, use a compact top bar and bottom navigation only for the highest-frequency destinations.
- Admin uses collapsible side navigation and a stable header. Route changes do not animate.
- Active navigation is indicated through shape, label weight, and color, not color alone.

### Data display

- Product metrics combine a label, value, unit, period, and optional trend. No decorative metric is allowed.
- Record timelines prioritize event type, time, amount/duration, caregiver, and abnormal status.
- Admin tables keep filters close to the affected table, preserve visible sort state, and offer bulk actions only after selection.
- Charts always include a textual summary and meaningful empty/no-data state.

### Overlays and feedback

- Desktop quick record uses a centered modal or anchored panel; mobile browser uses a bottom drawer.
- Popovers grow from their trigger origin. Modals remain centered.
- Toasts confirm transient success; durable failures remain visible in-page.
- External-service failure is never represented as an endless spinner.

## 5. Layout Principles

### Product Web

- Desktop-first canvas at 1440px; primary content width is 1200–1320px.
- Use a two-region workbench only when the secondary region materially supports the main task.
- The current baby/family context is always visible before record creation.
- Keep the primary action in a predictable location; do not scatter multiple coral buttons with equal weight.
- Responsive checkpoints: 390px, 768px, 1024px, and 1440px.
- Mobile browser is a reflowed Web experience, not a native App replica.

### Admin

- Supported from 1024px upward; 1440px is the primary review frame.
- Favor compact filters, scan-friendly tables, stable columns, and detail drawers.
- Avoid product-style emotional heroes and oversized empty space.

### Spacing

Use a 4px base rhythm: 4, 8, 12, 16, 24, 32, and 48px. A screen should normally use no more than three adjacent spacing levels. Dense Admin layouts use 8/12/16px; Product pages use 12/16/24/32px.

## 6. Density Modes

The system has one semantic foundation and two density expressions.

| Property       | Product Web                    | Admin                                 |
| -------------- | ------------------------------ | ------------------------------------- |
| Mood           | Warm, reassuring, task-focused | Neutral, operational, scan-focused    |
| Body size      | 14–16px                        | 14px                                  |
| Control height | 44–48px                        | 36–40px                               |
| Card radius    | 12–16px                        | 6–10px                                |
| Typical gap    | 16–24px                        | 8–16px                                |
| Elevation      | Soft border plus subtle shadow | Mostly border-only                    |
| Motion         | Restrained feedback            | Crisp feedback; no route/table motion |

## 7. Motion & Interaction

Use these shared motion values:

| Name                | Value                             | Role                            |
| ------------------- | --------------------------------- | ------------------------------- |
| Press               | `120ms`                           | Button and pressable feedback   |
| Fast                | `160ms`                           | Hover, tooltip, compact control |
| Base                | `200ms`                           | Dropdown, state transition      |
| Overlay             | `240ms`                           | Modal and drawer                |
| Responsive ease-out | `cubic-bezier(0.23, 1, 0.32, 1)`  | Enter/exit and direct response  |
| Natural ease-in-out | `cubic-bezier(0.77, 0, 0.175, 1)` | On-screen movement              |

Rules:

- Do not use `transition: all` or `ease-in` for UI entry.
- Prefer CSS transitions for predictable UI and Web Animations API only when programmatic control is required.
- Animate transform and opacity wherever possible.
- Do not animate keyboard-driven actions, route changes, table rows, or frequently repeated navigation.
- Respect `prefers-reduced-motion`: remove movement while retaining useful opacity and color feedback.
- Gate hover motion behind `(hover: hover) and (pointer: fine)`.

## 8. Page State Packs

Every screen references one or more reusable packs instead of inventing one-off states.

| Pack                       | Required states                                                    |
| -------------------------- | ------------------------------------------------------------------ |
| S1 Data                    | Loading, empty, ready, degraded, error, retry                      |
| S2 Form                    | Default, focus, validation, submitting, success, failure, disabled |
| S3 External service        | Not configured, unavailable, quota exhausted, offline, retry       |
| S4 Identity and permission | Signed out, forbidden, read-only, role-limited, session expired    |
| S5 Media upload            | Idle, selected, uploading, success, failure, invalid size/type     |

- P0 product routes: login, baby/family setup, baby profile, home, quick record, records, record detail, stats, family, and settings.
- P1 designed routes: AI assistant, photos, a realtime panel within the Web family surface, dark mode, and elder mode.
- P0 Admin routes: login, dashboard, users, babies, families, and records.
- P1 designed Admin routes: photo moderation and AI usage.

Native `/app/*` routes are outside this design cycle. Realtime Web behavior belongs inside the responsive family surface rather than reviving the deferred native-App route.

## 9. Content Voice

- Warm and direct, never infantilizing the caregiver.
- Prefer action language: “记录一次喂养” over “添加数据”.
- Name the affected baby or family when ambiguity is possible.
- Errors explain what happened, whether data was saved, and the next safe action.
- External features that are disabled say “暂未启用” and explain the available alternative.
- MVP login contains username and password only. Do not show SMS login, phone login, or public registration. For a missing account or reset, direct the user to an administrator.
- Admin account provisioning collects a unique username, nickname, initial password, status, and role. Phone is not a required design field. Reveal an initial password only in the immediate creation confirmation and never in later lists or detail views.
- Do not use emoji as functional icons. Use the established icon system with visible text where meaning could be ambiguous.
- Measurements use consistent units: `ml`, `h`, `min`, `°C`, and localized date/time formatting.

## 10. Accessibility

- All essential text and controls meet WCAG AA contrast.
- Product touch targets are at least 44×44px; elder mode targets are at least 48×48px.
- Focus-visible state is obvious and never removed. Use `box-shadow: 0 0 0 2px var(--surface-primary), 0 0 0 4px #2E90FA`; map `--surface-primary` to the actual adjacent light or dark surface so the Focus Blue ring remains distinguishable beside coral, white, and dark controls.
- Forms, dialogs, tables, drawers, and menus support keyboard navigation.
- Status is communicated by icon or label in addition to color.
- Elder mode uses at least 16px body text, stronger contrast, simpler grouping, and no gesture-only action.
- Responsive layouts reflow content; they do not rely on horizontal page scrolling. Admin tables may use contained horizontal scrolling with sticky context columns.

## 11. Stitch Prompt Contract

Every Stitch request must include:

1. route, user role, and primary task;
2. P0 or P1 scope;
3. 1440px desktop, 1024px compact Admin, or 390px mobile-browser frame;
4. a self-contained compilation of the relevant rules from this `DESIGN.md`;
5. components to reuse;
6. required S1–S5 state packs;
7. realistic Chinese product copy and real domain fields;
8. accessibility and reduced-motion requirements;
9. a statement that native App UI is out of scope;
10. forbidden patterns and external-service fallback behavior.

### Product prompt base

```text
Design a desktop-first responsive childcare Web interface for BabyGarden.
Apply the compiled BabyGarden design rules included with this brief.

Personality: warm, calm, trustworthy, and quietly competent.
Use a warm porcelain canvas, accessible Deep Warm Coral actions,
soft mint status accents, and Deep Ink Navy text.
Use medium information density and gently curved 12–16px containers.

Prioritize recording speed, data clarity, and family confidence.
Provide the required loading, empty, degraded, error, and retry states.
When requested, provide a 390px mobile-browser reflow, not a native App screen.

Use subtle 120–240ms interaction feedback. Do not add page-wide entrance
animation, bounce, emoji icons, decorative gradients, or invented metrics.
```

### Admin prompt base

```text
Design a compact desktop administration interface for BabyGarden,
compatible with Ant Design component patterns and the compiled
BabyGarden design rules included with this brief.

Reuse the shared semantic color and status roles. Use a neutral canvas,
tighter spacing, restrained 6–10px radii, and scan-friendly tables.
Prioritize filtering, permissions, bulk operations, auditability,
and external-service health.

No mobile layout, route animation, row entrance animation, or bounce.
Provide loading, empty, error, forbidden, and service-unavailable states.
```

## 12. Design-to-Code Rules

- Stitch HTML is reference output, not production React code.
- Before submitting to Stitch, Codex reads this document and expands the page-specific palette, components, state packs, responsive rules, content rules, and prohibitions into a self-contained prompt. A local file path by itself is not a valid handoff.
- Implement with the existing React, Ant Design, and CSS architecture unless an approved plan changes it.
- Map all colors, spacing, radii, typography, elevation, and motion to semantic tokens before page migration.
- Do not copy remote image URLs from generated Stitch HTML into production.
- Verify every implemented screen at the intended breakpoint and against its required state packs.
- Run an interaction and animation review after functional and visual acceptance, not before the core flow works.

## 13. Evidence Consulted

- Local Stitch exports under `.stitch/designs/stitch_yuanzi_baby_app/`.
- Active product shell in `src/pages/OpenDesignApp.tsx` and `src/styles/open-design.css`.
- Existing token evidence in `src/styles/variables.css`, `tailwind.config.js`, and `DESIGN_TOKENS.md`.
- Admin theme and routes under `src/admin/`.
- OpenDesign context-first, design-system, wireframe, verifier, and handoff methods.
- Emil Kowalski design-engineering and animation review rules.

The strongest evidence is the current code and local Stitch exports. The live Stitch project ID and project-level metadata were not available in this environment, so this document must be reconciled with the live Stitch project before it is declared synchronized.
