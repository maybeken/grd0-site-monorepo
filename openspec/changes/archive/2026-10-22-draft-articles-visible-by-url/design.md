# Design

## Context

See proposal.md for motivation and the spec deltas for the required behavior.

The public article read and the public scene read both use a publish-state-aware lookup (`findPublishedArticle`), while the public list already applies its own `published_at IS NOT NULL AND published_at <= now` filter. The web client needs no change: `ArticleView` renders any `200` response, and `ArticleHeader` already omits the published date when `published_at` is null. The sitemap worker builds article URLs from `GET /article`, so any change to the list filter would propagate there automatically.

## Goals / Non-Goals

**Goals:**
- Decouple URI resolution from publication state for article and scene reads.
- Keep listing and sitemap discovery governed solely by the publish timestamp.
- Preserve the existing response shape, admin endpoints, and editor flows.

**Non-Goals:**
- Access-controlled share links, tokens, or expiring URLs.
- A rich unpublished-page treatment beyond a single notice.
- Search-engine directives (`noindex`, `X-Robots-Tag`) for unlisted pages.
- Changing how the admin list or draft endpoints behave.

## Decisions

**Make the public reads publish-state agnostic.** `GetArticle`'s URI branch and `GetArticleScene` call `findArticle` instead of `findPublishedArticle`; `findPublishedArticle` is removed once it has no callers. Publication state then affects only list membership. Alternative considered: add a `published_only` flag threaded through both handlers. Rejected because the flag would always be false after this change, so it would add indirection without a second caller.

**Leave the list query and sitemap untouched.** The list filter is already the single place that decides discoverability, and the sitemap inherits it by consuming `GET /article`. Alternative considered: filter drafts in the sitemap worker instead. Rejected because it would duplicate the rule and allow the API and sitemap to drift.

**Keep the admin draft endpoints.** `GET /article/:uri/draft` and its scene counterpart still work and still require authentication, even though public reads now return the same content. Removing them would break the editor's request paths for no behavioral gain.

**Keep the article response shape unchanged.** A null `published_at` already means draft, and a future `published_at` means scheduled; the client uses the field for display only. No new field is needed to distinguish unlisted from listed articles because list membership is the distinction.

**Collapse the scene handler's published-only branch.** Both public and draft scene reads resolve the article by URI; the authenticated route remains as a compatibility alias. Alternative considered: drop the draft scene route. Rejected as unnecessary API surface removal in a change that is otherwise additive.

**Always show the editor's public-page link.** The sidebar renders the link whenever an article is open, using the saved URI, instead of gating it on `saved_status === 'published'`. Alternative considered: gate it on "published or shared", which would require a new sharing state the data model does not have.

**Derive the public notice from the publish timestamp.** `ArticleView` treats an article as not published when `published_at` is null or in the future, and renders one notice. Alternative considered: reuse the header's existing draft badge prop. Rejected because it only covers a null timestamp and reads as a small label rather than a notice.

## Risks / Trade-offs

- **Unlisted is not private** → Anyone with a draft or scheduled URL can read it. The spec states this; the mitigation is to keep URLs out of the list, sitemap, and navigation, not to add authentication.
- **URL enumeration** → Short, guessable URIs make unlisted articles findable. Authors control the URI; the system does not add unguessable tokens.
- **Search-engine indexing of shared links** → Sitemap exclusion limits discovery, but a crawler following a shared link may still index it. Deferred; if needed, add a `noindex` signal for articles whose publish timestamp is null or in the future.
- **Cache behaviour changes for drafts** → Any future CDN or response caching now sees draft content on the public path. Today no such cache exists; if one is added, draft pages must be treated like any other public page or explicitly excluded.
- **Publishing loses its "visibility gate" meaning** → Publish now controls listing and sitemap inclusion only. The spec and API documentation must describe it that way to avoid authors treating publish as the only way to make a URL work.

## Migration Plan

Deploy the API change; the web client and sitemap worker need no deployment and behave correctly against both old and new API versions for published content. There is no schema migration. Roll back by redeploying the previous API build; no data is mutated, so rollback restores the previous `404` behavior for unlisted articles without cleanup. Newly shared draft links stop resolving after rollback.

## Open Questions

- Should unlisted public pages emit a `noindex` signal, or is sitemap exclusion sufficient? This can be decided independently and does not change the specs or task breakdown.
