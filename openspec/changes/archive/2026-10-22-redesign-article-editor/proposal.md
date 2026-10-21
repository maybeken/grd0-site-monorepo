# Proposal

## Why

The article editor currently stacks the admin list and the editing surface on one page: the list fetches every article at once, opening an article is local state rather than a route, and the preview is a detached snapshot that does not follow the blocks being edited. As the article count grows the list becomes unwieldy, and the preview provides little orientation while authoring a long article.

This change builds on the archived `add-block-articles` change. It redesigns the editor surface that change introduced and does not depend on any blog behaviour.

## What Changes

- Split the single editor page into two routes: a paginated admin listing at `/article/editor` and a dedicated editing workspace at `/article/editor/:uri`.
- Change `GET /article/all` to return the project's existing paginated envelope (`data`, `total`, `page`, `page_size`, `total_pages`) instead of a bare array. The article surface is not live yet, so the response shape can change freely without compatibility work.
- Paginate the listing server-side using the existing `parsePagination` helper and `schema.PaginatedResponse`, with the page reflected in the query string so returning from an article restores the list position.
- Move title and subtitle editing into a shared header above the editor and preview. The preview renders blocks only, the panes are separated by a single vertical divider, and the preview's scrollbar is hidden.
- Synchronise the preview to the editor's scroll position using block anchors, so the preview follows the block being edited. The editor is the only driver; the preview never scrolls the editor.
- Move URI and publish controls into the sidebar, alongside status, timestamps, an icon copy action, contextual publish/unpublish actions, and delete.
- Make markdown blocks auto-resize so each block grows with its content and the editor pane scrolls as one surface.
- Allow duplicating a block immediately after itself; drawing blocks receive an independent stored scene.
- Allow creating an article without a URI; the API generates a placeholder that the author can replace in the sidebar. The listing creates after a single confirmation.
- Integrate the saved/unsaved state into the save control and make publication actions contextual in the sidebar.
- Replace duplicate-URI detection based on scanning the admin list with the API's `409 Conflict` response, since a paginated list can no longer be scanned reliably.
- Keep the existing dirty-state guards, Ctrl/Cmd+S save, debounced preview, and public routes unchanged.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `article`: the admin article list becomes server-paginated and returns the shared paginated envelope, and creation may omit the URI; the public list is unchanged.
- `article-editor`: the authoring experience becomes a two-page workflow with a paginated listing, a route-addressed editing workspace, a shared title/subtitle header, a block-only scroll-synchronised preview, block duplication, and a metadata sidebar with contextual actions; duplicate-URI conflicts are reported from the API status.

## Impact

- `api/handler/article.go`: paginate `GetArticleByAdmin` with `Count`, `Offset`, and `Limit`, reusing `parsePagination` and `schema.PaginatedResponse[ArticleListEntry]`.
- `vue-web/src/interfaces/Article.d.ts`: an `ArticleListPage` type based on the existing `PaginatedResponse`.
- `vue-web/src/services/article.ts`: `listArticlesAdminRaw` takes page parameters; conflict checks use `errorStatus`.
- `vue-web/src/services/api.ts`: the admin instance rejects non-OK responses so HTTP status codes reach `errorStatus`.
- `vue-web/src/components/editor/article/useArticleDraft.ts`: remove `uriExists` and the extra list request it issues.
- `vue-web/src/router/index.ts`: add `/article/editor/:uri`; `/article/editor` becomes the listing.
- `vue-web/src/views/`: a new listing view; `ArticleEditorView.vue` becomes the editing workspace.
- `vue-web/src/components/editor/article/`: `ArticleList.vue` gains server pagination; `ArticleEditor.vue` loses the list wiring and gains the workspace layout; new sidebar and pager components; a new scroll-sync composable.
- `vue-web/src/components/article/BlockRenderer.vue`: per-block anchor wrappers.
- `vue-web/src/components/editor/article/blocks/MarkdownBlockEditor.vue`: enable OverType auto-resize.
- `vue-web/src/locales/en.yaml`: strings for the listing, sidebar, and pagination.
- No new dependencies. No database schema change.
