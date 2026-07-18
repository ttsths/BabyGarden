# BabyGarden Web + Admin Design System Specification

- **Date:** 2026-07-12
- **Status:** Approved design, ready for implementation planning
- **Canonical visual reference:** [`yuanzi-frontend/DESIGN.md`](../../../yuanzi-frontend/DESIGN.md)
- **Scope:** Web and Admin; native App deferred

## 1. Objective

Create one maintainable design system and one repeatable design workflow for BabyGarden Web and Admin. The system must preserve the warm coral-and-mint identity found in the existing Stitch work while correcting accessibility, desktop layout, component duplication, and design-to-code drift.

The design covers all current product and Admin routes. Development remains phased so the MVP can launch without waiting for AI, object storage, realtime updates, SMS, or other external services.

## 2. Confirmed Decisions

The following decisions were explicitly reviewed and approved:

1. Use OpenDesign as a method reference, not as an installed third-party runtime.
2. Keep Stitch as the high-fidelity visual design tool.
3. Use Emil Kowalski's skills as an interaction and motion quality gate, not as the base design system.
4. Prioritize desktop Web while keeping the product usable in a 390px mobile browser.
5. Do not design or implement native `/app/*` routes in this cycle.
6. Use one semantic foundation with two density expressions:
   - Warm Trust product Web;
   - Neutral Efficient Admin.
7. Design all current Web/Admin routes, but implement the MVP in tiers.
8. Use username/password login; SMS is out of scope.
9. AI, photos, realtime, dark mode, and elder mode receive complete design coverage but may ship later or in a degraded state.
10. MVP accounts are provisioned through Admin; public self-registration is outside this cycle. Admin provisioning must create a unique username and initial password without requiring SMS or a phone number.

## 3. Current-State Evidence

### 3.1 Stitch assets are mobile-only

The repository contains 13 local Stitch screen exports under:

```text
yuanzi-frontend/.stitch/designs/stitch_yuanzi_baby_app/
```

Every exported screenshot is 1600px tall and between 555px and 706px wide. The HTML frames generally cap content near 375–430px. These assets establish useful brand DNA, components, and mobile interaction patterns, but they do not define a desktop Web or Admin system.

Frequent source colors include:

- coral `#FF998A`;
- warm canvas `#FFFBF7`;
- dark brown canvas `#23110F`;
- mint `#A2D5C6`.

### 3.2 Active Web differs from legacy Stitch implementation

The live root route renders `OpenDesignApp`, not the older individual Stitch pages:

- `yuanzi-frontend/src/App.tsx` routes product traffic to `OpenDesignApp`;
- `OpenDesignApp.tsx` is a 984-line multipage component;
- `open-design.css` is a 1068-line parallel styling system with 138 hard-coded color/function occurrences.

The older token-driven components and reports therefore do not fully describe the active product UI.

### 3.3 Admin has a third theme path

Admin uses Ant Design with theme values embedded in `AdminLayout.tsx`, while `src/styles/antd-theme.ts` contains another theme object. Both still use the low-contrast `#FF9A8B` primary. This produces a separate source of truth from both `variables.css` and `open-design.css`.

### 3.4 Motion and accessibility drift

Current components contain `transition-all`, long 500ms transitions, and global keyframes without a consistent reduced-motion policy. The existing coral primary provides only 2.06:1 contrast with white text.

## 4. Design Architecture

The system has four canonical artifacts:

```text
DESIGN.md
  → semantic visual and content rules for people and Stitch

Semantic tokens
  → implementation variables and Ant Design mappings

Component contracts
  → supported variants, states, sizing, behavior, and accessibility

Page state matrix
  → required loading, empty, degraded, error, permission, and external states
```

The canonical system serves two expressions:

```text
Shared semantic foundation
├── Product Web: warm, trustworthy, medium density, responsive
└── Admin: neutral, efficient, compact density, desktop only
```

Product and Admin may change density aliases, but they must not independently redefine semantic meaning for action, success, warning, danger, text, border, or focus.

## 5. Visual Direction

### Product Web

- warm porcelain canvas;
- crisp white content surfaces;
- accessible deep coral actions;
- mint success and healthy-state accents;
- deep navy typography;
- 12–16px card radii;
- soft borders plus restrained elevation;
- calm, direct Chinese copy;
- medium information density.

### Admin

- neutral gray canvas;
- white table and filter surfaces;
- shared action and status colors;
- 6–10px radii;
- compact 8/12/16px spacing;
- border-led hierarchy;
- no emotional hero treatment or decorative metric cards.

### Accessibility correction

`#FF998A` remains a brand accent. Filled primary controls use `#C84B42`, which supports white text at 4.62:1. Dark mode may use `#FF998A` with deep navy text at 7.88:1. Focus-visible controls use a 2px adjacent-surface separator plus a 2px `#2E90FA` outer ring; this prevents the focus color from sitting directly against a low-contrast coral or dark fill.

## 6. Responsive Strategy

### Product

- 1440px: canonical desktop review frame;
- 1024px: compact desktop/tablet landscape;
- 768px: navigation and grid reflow boundary;
- 390px: canonical mobile-browser review frame.

Desktop uses a restrained side navigation where the route count requires it. Mobile browser reflows into a compact header and limited bottom navigation. Quick record is a desktop modal or anchored panel and a mobile bottom drawer.

### Admin

Admin supports 1024px and above. A 1440px frame is canonical. Tables may scroll within their container at 1024px while preserving the entity identifier and action context. No phone layout is required.

## 7. Page and State Matrix

### State packs

| Pack                   | Required states                                                    |
| ---------------------- | ------------------------------------------------------------------ |
| S1 Data                | Loading, empty, ready, degraded, error, retry                      |
| S2 Form                | Default, focus, validation, submitting, success, failure, disabled |
| S3 External            | Not configured, unavailable, quota exhausted, offline, retry       |
| S4 Identity/permission | Signed out, forbidden, read-only, role-limited, session expired    |
| S5 Media               | Idle, selected, uploading, success, failure, invalid size/type     |

### Product Web

| Path or mode                           | Surface                                  | Tier | Frames              | Packs      |
| -------------------------------------- | ---------------------------------------- | ---- | ------------------- | ---------- |
| `/login`                               | Username/password login                  | P0   | 1440, 390           | S2, S3, S4 |
| `/baby/setup`                          | Family and baby setup                    | P0   | 1440, 390           | S2, S4     |
| `/baby-profile`                        | Baby profile                             | P0   | 1440, 390           | S1, S2     |
| `/`                                    | Home dashboard                           | P0   | 1440, 390           | S1, S3     |
| `/record`                              | Quick record                             | P0   | Modal, drawer       | S2         |
| `/records`                             | Record list                              | P0   | 1440, 390           | S1, S4     |
| `/record-detail/:id`                   | Record detail                            | P0   | 1440, 390           | S1, S4     |
| `/stats`                               | Statistics                               | P0   | 1440, 390           | S1         |
| `/family`                              | Family members and invitation            | P0   | 1440, 390           | S1, S2, S4 |
| `/family` realtime variant             | Realtime family panel                    | P1   | Family-page variant | S1, S3, S4 |
| `/settings`                            | Settings                                 | P0   | 1440, 390           | S2, S4     |
| `/ai`                                  | AI assistant                             | P1   | 1440, 390           | S1, S2, S3 |
| `/photos`                              | Photos                                   | P1   | 1440, 390           | S1, S3, S5 |
| `/elder`; dark mode on anchor surfaces | Elder mode and dark-mode anchor variants | P1   | Anchor variants     | S1, S2     |

### Admin

| Path               | Surface                        | Tier | Packs          |
| ------------------ | ------------------------------ | ---- | -------------- |
| `/admin/login`     | Admin login                    | P0   | S2, S3, S4     |
| `/admin/dashboard` | Dashboard                      | P0   | S1             |
| `/admin/users`     | Users and account provisioning | P0   | S1, S2, S4     |
| `/admin/babies`    | Babies                         | P0   | S1, S2, S4     |
| `/admin/families`  | Families                       | P0   | S1, S2, S4     |
| `/admin/records`   | Records                        | P0   | S1, S2, S4     |
| `/admin/photos`    | Photo moderation               | P1   | S1, S3, S4, S5 |
| `/admin/ai-usage`  | AI usage                       | P1   | S1, S3, S4     |

### Stitch asset budget

Target approximately 40 design assets:

- 10 P0 product desktop masters;
- 6 critical product mobile-browser masters;
- 4 P1 AI/photo responsive designs plus one realtime family-state variant;
- 8 Admin desktop masters;
- 2 Admin compact-layout anchors;
- 4 dark/elder anchors;
- 5 shared state boards.

State boards are preferred over duplicating a complete page for every component state.

## 8. Component Contracts

### Shared primitives

- Button, IconButton;
- Input, TextArea, Select, DatePicker, Switch;
- Tabs, Badge, StatusIndicator;
- Modal, Drawer, Popover, Tooltip;
- Alert, Toast, Skeleton, EmptyState, ErrorState;
- Pagination, Upload, Avatar.

Every interactive primitive must cover default, hover, focus-visible, active, disabled, loading, keyboard, and reduced-motion behavior.

### Product composites

- ProductShell and responsive navigation;
- BabySwitcher;
- QuickRecordLauncher;
- DailyMetricCard;
- RecordTimeline and RecordDetail;
- GrowthChart with textual summary;
- FamilyMemberCard;
- PhotoTile and upload flow;
- AIComposer and external-service fallback.

### Admin composites

- AdminShell;
- FilterBar;
- DataTable;
- BulkActionBar;
- EntityDetailDrawer;
- AuditPanel;
- AdminStatisticCard;
- PermissionBadge;
- ExternalServiceStatus.

## 9. Motion Contract

| Token       | Value                             | Use                             |
| ----------- | --------------------------------- | ------------------------------- |
| Press       | 120ms                             | Press feedback                  |
| Fast        | 160ms                             | Tooltip, hover, compact control |
| Base        | 200ms                             | Dropdown and state transition   |
| Overlay     | 240ms                             | Modal and drawer                |
| Ease out    | `cubic-bezier(0.23, 1, 0.32, 1)`  | Responsive entry/exit           |
| Ease in/out | `cubic-bezier(0.77, 0, 0.175, 1)` | On-screen movement              |

Rules:

- no `transition: all`;
- no `ease-in` entrance;
- no `scale(0)` entrance;
- no route, keyboard-action, or table-row entrance animation;
- transform and opacity are preferred;
- popovers originate from their trigger; modals remain centered;
- movement is removed or reduced under `prefers-reduced-motion`;
- hover motion is limited to precise pointing devices.

No new animation dependency is required for MVP. CSS transitions and limited Web Animations API usage are sufficient.

## 10. Content and Error Handling

### Voice

Product copy is warm, concise, and respectful. It uses caregiver language without infantilizing the adult. Admin copy is operational and explicit.

### Failure rules

- Never use an endless spinner for a failed dependency.
- State whether a record was saved before showing retry.
- Preserve entered form data after recoverable failure.
- AI and photos display “暂未启用” or “服务暂不可用” with an available alternative.
- The P0 home dashboard loads stats and records independently from P1 photos, AI, and realtime services. Optional-service failure degrades only its own card or removes that card; it never turns the complete home page into an error state.
- Permission errors explain the required role without exposing sensitive policy details.
- Empty means no data; error means data could not be loaded. The two states must not share copy.

## 11. Design Workflow

```text
Current React, Ant Design, API domain and Stitch exports
  → extract/update DESIGN.md and semantic tokens
  → explore 3–5 structurally different wireframes for uncertain flows
  → approve one structure
  → compile page-specific DESIGN.md rules into a self-contained Stitch prompt
  → generate Stitch high-fidelity master and required state board
  → hand off dimensions, tokens, component states, copy and assets
  → rebuild with existing React/AntD components
  → functional E2E and visual regression
  → Emil-style motion review
  → independent acceptance against the brief
```

OpenDesign contributes context-first intake, structural exploration, verifier separation, and explicit handoff. It is not installed or executed in the project. Stitch remains the visual authoring surface. Generated HTML is never production source.

## 12. Stitch Prompt Rules

Every prompt must include:

- route, role, task, and tier;
- required breakpoint;
- a self-contained compilation of the relevant `yuanzi-frontend/DESIGN.md` rules;
- required existing components;
- required state packs;
- real Chinese labels and domain fields;
- native App exclusion;
- external-service degraded behavior;
- reduced-motion and keyboard requirements;
- anti-slop exclusions.

Product and Admin base prompts are maintained in `DESIGN.md`. Codex must compile the selected page's tokens, component contracts, state packs, responsive rules, copy rules, and prohibitions into the submitted prompt. Stitch must not be expected to resolve a local file path.

## 13. Implementation Sequence

This specification authorizes planning, not immediate business-code refactoring. The implementation plan should follow this order:

1. Establish the canonical token layer and map it into Ant Design.
2. Mark legacy token documents as superseded or reconcile them.
3. Implement the P0 account contract: Admin creates a unique username, nickname, securely hashed initial password, status, and role; phone is not required; public registration and SMS remain absent. The current phone-required Admin API and plaintext password handling must be replaced before acceptance.
4. Inventory active components and remove duplicate primitives only where the migrated pages require it.
5. Produce Stitch P0 masters and shared state boards using self-contained compiled prompts.
6. Implement the product shell, login, setup, and home as the first vertical slice. Separate core stats/records loading from photos, AI, and realtime requests so optional failures remain local.
7. Implement quick record, record list/detail, stats, family, and settings.
8. Implement Admin login, shell, dashboard, account provisioning, and core entity management.
9. Produce and optionally implement P1 AI, photos, realtime, dark, elder, photo moderation, and AI usage.
10. Add Playwright visual baselines and state-pack coverage.
11. Run motion and accessibility review before release acceptance.

Avoid a full-codebase visual rewrite in one merge request. Each slice should leave the shared system more canonical than before and remain deployable.

## 14. Acceptance Criteria

The design phase is complete when:

1. `DESIGN.md` is the documented source of truth.
2. Every in-scope route maps to a tier, frame, and state pack.
3. Product and Admin share semantic color/status roles.
4. Filled controls and essential text meet WCAG AA.
5. Stitch masters include desktop Web, critical mobile-browser frames, Admin, and state boards.
6. No native App work is included.
7. AI, photos, and realtime have explicit disabled/unavailable states.
8. Handoff identifies tokens, component contracts, copy, assets, and responsive behavior.
9. Stitch briefs are self-contained compiled prompts rather than unresolved local-file references.

The implementation phase is accepted when:

1. production pages contain no unapproved raw design colors;
2. Ant Design consumes the canonical semantic mappings;
3. P0 states have functional and visual tests;
4. keyboard focus and reduced motion are verified;
5. deployed Web and Admin pass screenshot comparison at canonical frames;
6. animation review finds no feel-breaking or accessibility-blocking regressions.
7. username/password accounts are admin-provisioned without SMS or a required phone number, and passwords are never stored or returned in plaintext.
8. the P0 home remains usable when AI, photos, realtime, or storage services are unavailable.

## 15. Risks and Controls

| Risk                                                   | Control                                                                      |
| ------------------------------------------------------ | ---------------------------------------------------------------------------- |
| Live Stitch metadata is unavailable                    | Reconcile `DESIGN.md` with the live project before declaring synchronization |
| Four visual vocabularies currently coexist             | Migrate by vertical slice and prohibit new raw colors                        |
| Full-route design expands scope                        | Keep P0/P1 development gates independent from design completion              |
| External services are unavailable                      | Design explicit S3 degraded states and avoid blocking P0                     |
| Generated Stitch code diverges from React architecture | Treat generated HTML as visual reference only                                |
| Third-party design skills change rapidly               | Reference reviewed methods and pin any future installation                   |
| Large all-at-once refactor destabilizes MVP            | Use small, deployable slices with visual regression                          |

## 16. Self-Review

- No unresolved placeholders or work markers are present.
- Native App exclusion is consistent across architecture, routes, and workflow.
- Full design coverage does not imply full MVP implementation.
- Color contrast corrections are explicit and measurable.
- Product and Admin differences are expressed through density, not duplicated semantics.
- External dependency failures have defined UI behavior.
- The implementation sequence preserves the existing stack and avoids an unnecessary new animation library.
