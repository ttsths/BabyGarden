# P0 Stitch Asset Plan

Target: approximately 40 review assets. Generate one screen per prompt and use
state boards instead of duplicating full pages for every component state.

Current progress: `product-login-desktop` was accepted to continue and
`product-login-mobile` is generated in Stitch, pending final visual review.
`product-home-desktop` is now generated with the MVP core/degraded-state
contract; it is pending visual review. The next baton is
`product-home-mobile`.

## Batch A — Product desktop masters (10)

Login, setup, baby profile, home, quick record, record list, record detail,
statistics, family and settings.

## Batch B — Critical 390px product masters (6)

Login, setup, home, quick-record drawer, record list and family.

## Batch C — Admin masters and compact anchors (10)

Admin login, dashboard, users, babies, families and records at 1440px, plus
1024px anchors for dashboard, users, families and records.

## Batch D — Shared state boards (5)

Data, forms, external services, identity/permission and media upload.

## Batch E — Non-blocking P1 design coverage (9)

AI and photos responsive anchors, realtime family state, dark anchors and elder
anchors. Generate only after all P0 review gates pass.

## Baton Order

After `product-login-desktop`, prepare `product-login-mobile`, then
`product-home-desktop`, followed by `product-home-mobile`. Update `SITE.md`,
metadata and this order after every generated Stitch screen; only mark a
sitemap item complete after human visual review.
