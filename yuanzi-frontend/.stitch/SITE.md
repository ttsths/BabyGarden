# BabyGarden Web and Admin Design Roadmap

## 1. Vision

Create a coherent, reviewable Stitch source for BabyGarden's MVP Product Web
and Admin. The output guides React implementation; it does not replace the
application architecture.

## 2. Stitch Project

- Project ID: `6507375090366546067`
- Project title: `Yuanzi Baby App` (existing private project selected for this
  design cycle)
- Visibility: private
- Design system asset: `assets/f56839ef90a64ca1898687653ea64867`
- Current status: Product Web login, resilient home, first-family setup,
  quick-record and records-list masters are generated. The first records-list
  mobile draft was superseded because it introduced App navigation; a refined
  Web-only variant is now generated. The next baton is record detail desktop.

The project and screen identifiers are persisted in `.stitch/metadata.json`.
Refresh that file after every accepted Stitch generation.

## 3. Existing Evidence

Thirteen legacy mobile-oriented exports exist below
`.stitch/designs/stitch_yuanzi_baby_app/`. They preserve useful brand DNA but do
not count as the desktop Product or Admin masters required by issue #84.

## 4. Sitemap

### Product P0

- [x] `/login` — username/password login, 1440px desktop master
- [ ] `/login` — 390px mobile-browser variant generated; final review pending
- [ ] `/baby-profile` — baby profile
- [ ] `/` — 1440px resilient home dashboard generated; final review pending
- [ ] `/` — 390px mobile-browser home variant generated; final review pending
- [ ] `/baby/setup` — first family and baby setup, 1440px desktop master
- [ ] `/baby/setup` — first family and baby setup, 390px mobile-browser variant
- [ ] `/record` — quick-record modal, 1440px desktop master
- [ ] `/record` — quick-record drawer, 390px mobile-browser variant
- [ ] `/records` — record list, 1440px desktop master
- [ ] `/records` — record list, 390px mobile-browser variant (Web-only refined)
- [ ] `/record-detail/:id` — record detail
- [ ] `/record-detail/:id` — 1440px desktop master
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

1. Product login desktop and mobile masters — generated; mobile review pending.
2. Product shell and resilient home desktop/390px masters — generated; review
   pending.
3. First-family setup desktop/390px masters — generated; review pending.
4. Quick-record desktop/mobile vertical slice for feeding, sleep and diaper —
   generated; review pending.
5. Product records list desktop/mobile, then detail, statistics, family and
   settings masters.
6. Admin shell, dashboard and account provisioning.
7. Remaining Admin P0 masters and 1024px compact anchors.
8. Five shared state boards.
9. Human visual, responsive and accessibility review.

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
