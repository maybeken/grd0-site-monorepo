# Tasks

## 1. Public article read

- [x] 1.1 Replace the published-only lookup in the public article URI branch with the publish-state-agnostic lookup; verify `go build ./...` passes and requests for a draft, scheduled, and published URI all return `200`
- [x] 1.2 Make the public scene read resolve any existing article; verify a scene belonging to a draft article returns the stored JSON and an unknown article URI returns `404`
- [x] 1.3 Remove `findPublishedArticle` once it has no remaining callers; verify `go build ./...` and `go vet ./...` pass

## 2. Listing invariants

- [x] 2.1 Confirm the public list query is unchanged; verify draft and future-dated articles are absent from `GET /article` while published past articles are present
- [x] 2.2 Confirm the sitemap worker needs no change; verify it still derives article URLs from `GET /article` and therefore excludes unlisted articles

## 3. Integration verification

- [x] 3.1 Verify an unlisted article URL renders in the web client without appearing on `/article`, covering both a draft and a future-dated article
- [x] 3.2 Confirm rollback expectations: redeploying the previous API build restores `404` for unlisted URLs with no data cleanup

## 4. Unlisted article UI

- [x] 4.1 Show the public-page link in the editor sidebar regardless of publication state; verify the link is present for a draft and a scheduled article
- [x] 4.2 Show a not-published notice on the public article page when the publish timestamp is absent or in the future, with a localized string; verify drafts and scheduled articles show the notice, published articles do not, and `pnpm type-check` passes

## Workflow follow-up

- Archive the change after review.
