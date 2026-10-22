# Proposal

## Why

Authors need to share a work-in-progress article through an unlisted URL without publishing it to the public index. Today an unauthenticated read of a draft returns `404`, which forces premature publishing just to get a shareable link.

## What Changes

- The public article read (`GET /article/:uri`) resolves any article by URI regardless of its publish timestamp, so drafts and future-dated articles are readable by anyone who has the URL.
- The public article list (`GET /article`) keeps returning only published, non-future articles, so drafts and future-dated articles remain unlisted.
- The public scene read (`GET /article/:uri/scene/:sid`) resolves scenes for any article regardless of publish timestamp, so drawing blocks render on a shared unlisted article.
- The sitemap keeps excluding drafts and future-dated articles, because it is generated from the public list.
- The editor sidebar always offers its "Open public page" link, so authors can open a draft or scheduled article without publishing it.
- The public article page shows a notice when the article is a draft or scheduled, stating that it is not published yet.
- **BREAKING**: the previous contract that an unauthenticated read of an unpublished or future-dated article returns `404` is removed. Those URLs become world-readable to anyone who knows them.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `article`: The Draft visibility and Publishing and unpublishing requirements change so every URI resolves through the public article endpoint regardless of publish state, while only published, non-future articles appear in the public list.
- `article-view`: The Article page requirement changes so draft and future-dated URIs render the article instead of the not-found state, and a new notice requirement states that an unlisted article is not published yet; the index and sitemap exclusions stay as they are.
- `article-editor`: The Metadata sidebar requirement changes so the public-page link is offered regardless of publication state.
- `article-asset-storage`: The Scene access endpoints requirement changes so scenes are publicly readable for any article, not only published ones.

## Impact

- `api/handler/article.go`: `GetArticle` and `GetArticleScene` stop using the published-only lookup; `findPublishedArticle` becomes unused if no other caller remains.
- `vue-web`: `ArticleView` renders an additional not-published notice for drafts and scheduled articles, and `ArticleSidebar` always shows the public-page link. `ArticleIndexView` still consumes the filtered list.
- `workers/sitemap-generator`: no change; the sitemap derives article URLs from the public list, which continues to exclude drafts and future-dated articles.
- No database migration and no new API surface.
- A draft or future-dated URL is an unlisted, not access-controlled, link: anyone with it can read the article.
