# P0 Stitch Asset Plan

Target: approximately 40 review assets. Generate one screen per prompt and use
state boards instead of duplicating full pages for every component state.

Current progress: `product-login-desktop` was accepted to continue and
`product-login-mobile` is generated in Stitch, pending final visual review.
`product-home-desktop` and `product-home-mobile` are now generated with the
MVP core/degraded-state contract; both are pending visual review. The next
baton `product-setup-desktop` is now generated with the username/password-only
MVP onboarding contract and is pending visual review. The next baton is
`product-setup-mobile`, which is now generated and pending visual review. The
next baton `product-record-desktop` is now generated as the quick-record modal
vertical slice and is pending visual review. The next baton is
`product-record-mobile`, which is now generated and pending visual review. The
next baton `product-records-desktop` is now generated as the consumer records
history page and is pending visual review. The next baton is
`product-records-mobile`; its first draft was superseded for adding App
navigation, and a Web-only refinement is now generated pending review. The next
baton `product-record-detail-desktop` is now generated and pending visual
review. The next baton is `product-record-detail-mobile`.

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
`product-home-desktop`, `product-home-mobile`, `product-setup-desktop`,
`product-setup-mobile`, `product-record-desktop`,
`product-record-mobile`, `product-records-desktop`,
`product-records-mobile`, `product-record-detail-desktop`, followed by
`product-record-detail-mobile`. Update `SITE.md`, metadata and this order after
every generated Stitch screen; only mark a sitemap item complete after human
visual review.
