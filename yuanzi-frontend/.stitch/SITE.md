# BabyGarden Web and Admin Design Roadmap

## 1. Vision

Create a coherent, reviewable Stitch source for BabyGarden's MVP Product Web
and Admin. The output guides React implementation; it does not replace the
application architecture.

## 2. Stitch Project

- Project ID: pending Stitch MCP connection
- Expected project title: `BabyGarden Web + Admin MVP`
- Visibility: private
- Current blocker: this Codex session has no Stitch MCP tools, so prompts are
  prepared locally but cloud screens cannot yet be generated or synchronized.

When MCP becomes available, create or select the project, write
`.stitch/metadata.json`, generate the current baton, then refresh metadata after
every screen.

## 3. Existing Evidence

Thirteen legacy mobile-oriented exports exist below
`.stitch/designs/stitch_yuanzi_baby_app/`. They preserve useful brand DNA but do
not count as the desktop Product or Admin masters required by issue #84.

## 4. Sitemap

### Product P0

- [ ] `/login` — username/password login, desktop and 390px
- [ ] `/baby/setup` — first family and baby setup, desktop and 390px
- [ ] `/baby-profile` — baby profile
- [ ] `/` — resilient home dashboard, desktop and 390px
- [ ] `/record` — quick-record modal and mobile drawer
- [ ] `/records` — record list, desktop and 390px
- [ ] `/record-detail/:id` — record detail
- [ ] `/stats` — core statistics
- [ ] `/family` — members and invitation
- [ ] `/settings` — account and product settings

### Admin P0

- [ ] `/admin/login`
- [ ] `/admin/dashboard`
- [ ] `/admin/users`
- [ ] `/admin/babies`
- [ ] `/admin/families`
- [ ] `/admin/records`

### Shared State Boards

- [ ] Data states
- [ ] Form states
- [ ] External-service states
- [ ] Identity and permission states
- [ ] Media upload states

## 5. Roadmap

1. Product login desktop master and 390px variant.
2. Product shell, home desktop master and 390px variant.
3. Setup and core-record vertical slice.
4. Remaining Product P0 masters.
5. Admin shell, dashboard and account provisioning.
6. Remaining Admin P0 masters and 1024px compact anchors.
7. Five shared state boards.
8. Human visual, responsive and accessibility review.

## 6. P1 Design Coverage

After P0 review, produce non-blocking variants for AI, photos, realtime family
state, dark mode and elder mode. Do not implement those features as part of
issue #84.

## 7. Review Gates

- Each prompt is self-contained and includes `.stitch/DESIGN.md` semantics.
- Every master has real Chinese labels and realistic BabyGarden data.
- Product and Admin share semantics but preserve their separate densities.
- Required state packs are visible and not hidden in prose.
- The human reviewer accepts visual consistency before a sitemap item is
  checked.
