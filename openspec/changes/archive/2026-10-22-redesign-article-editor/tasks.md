# Tasks

## 1. API and client contract

- [x] 1.1 Paginate `GetArticleByAdmin` with `parsePagination`, a `Count` on the same predicate, and `Offset`/`Limit`, returning `schema.PaginatedResponse[ArticleListEntry]` ordered by `updated_at DESC`; verify with `go build ./...` and a request to `/article/all?page=2&page_size=10` returning `data`, `total`, `page`, `page_size`, and `total_pages`.
- [x] 1.2 Add an `ArticleListPage` type in `interfaces/Article.d.ts` based on the existing `PaginatedResponse`, and make `listArticlesAdminRaw(page, page_size)` request the paginated endpoint; verify `pnpm type-check` passes.
- [x] 1.3 Replace `useArticleDraft.uriExists` and its list scan with conflict detection from `errorStatus(error) === 409`; verify a duplicate create and a duplicate rename report a conflict, and the network panel shows no extra admin list request.

## 2. Routes and listing page

- [x] 2.1 Register `/article/editor` as the listing page and `/article/editor/:uri` as the editing workspace, keeping `/article/:slug` intact; verify each route loads its own view on a direct visit.
- [x] 2.2 Build the listing view with server pagination, the page in the query string, a pager with totals, and loading/empty/error states; verify changing pages updates the URL and the rows, and reloading a paged URL restores that page.
- [x] 2.3 Move the create-article flow to the listing and navigate to the new editing route on success; verify duplicate, invalid, and failed create states still render inline without navigating.
- [x] 2.4 Implement listing row actions for open, unpublish, and permanent delete, refreshing the current page afterwards; verify deleting the last row on a page above page 1 moves back one page.
- [x] 2.5 Keep the unauthenticated access check on both editor routes; verify an unauthenticated visit does not render articles or accept block operations.
- [x] 2.6 Update `vue-web/README.md` with the two editor routes and new components, and verify every documented command still runs.

## 3. Editing workspace and sidebar

- [x] 3.1 Make `ArticleEditorView` route-driven: open the article from the `:uri` parameter, load on parameter change, and show an error with a return link for an unknown URI; verify a direct visit loads the article and a bad URI shows the error.
- [x] 3.2 Replace the rename flow's local state update with `router.replace` after a successful save; verify the address bar shows the new URI and the previous URI no longer resolves.
- [x] 3.3 Build the bounded three-pane layout with an independently scrolling editor, preview, and sidebar, plus an editor/preview toggle below the side-by-side breakpoint; verify the document does not scroll while each pane does.
- [x] 3.4 Build the metadata sidebar with URI editing, publish datetime, draft/scheduled/published state, Publish now, Unpublish, timestamps, a public-page link, and confirmed permanent delete; verify each control changes or displays the expected state.
- [x] 3.5 Remove the listing wiring from `ArticleEditor` and `ArticleList` (selected prop, open emit, stacked layout) and move creation to the listing; verify no unused state remains and `pnpm type-check` passes.
- [x] 3.6 Preserve dirty-state protection for returning to the listing, navigating to another article, and reload or close; verify each path prompts and cancelling keeps the edits.

## 4. Scroll synchronisation

- [x] 4.1 Add a wrapper carrying `data-block-id` around every rendered block and around the article header in `BlockRenderer`; verify the preview DOM exposes one anchor per block in article order.
- [x] 4.2 Implement `useScrollSync` with per-frame anchor measurement, intra-block interpolation, block-ID matching, clamping, and trailing-space handling; verify scrolling the editor aligns the preview and scrolling the preview alone leaves the editor unchanged.
- [x] 4.3 Enable OverType `autoResize` with a minimum height in `MarkdownBlockEditor`; verify a long markdown block grows instead of scrolling internally and the editor pane scrolls as one surface.
- [x] 4.4 Rebuild and realign anchors after the debounced preview refresh, add/remove/reorder, pane resize, and image or drawing load; verify the preview stays aligned through each change.
- [x] 4.5 Render markdown images eagerly in the editor preview and highlight the block currently aligned; verify images do not shift alignment while scrolling and the highlight follows the editor.
- [x] 4.6 Skip synchronisation in the narrow-screen toggle and align the preview to the active block when the preview tab opens; verify the mobile layout does not fight the user's scrolling.

## 5. Preview timing and polish

- [x] 5.1 Reduce the preview debounce to roughly 750 ms while keeping immediate refresh on open and save; verify the preview reflects an edit within one second and does not reparse on every keystroke.
- [x] 5.2 Add the new listing, pagination, sidebar, and workspace strings to `src/locales/en.yaml`; verify the YAML is parsed by the custom plugin and no rendered label falls back to a raw key.
- [ ] 5.3 Check keyboard and pointer behavior across the workspace: Ctrl/Cmd+S, block reordering, block removal, and the drawing iframe's scroll shield; verify each still works inside the new panes.

## 6. Integration verification

- [x] 6.1 Run `cd api && go build ./...` and confirm the API compiles.
- [x] 6.2 Run `cd vue-web && pnpm type-check` and confirm there are no type errors.
- [x] 6.3 Run `cd vue-web && pnpm build` and confirm the production build succeeds.
- [ ] 6.4 Manual end-to-end: create a draft from the listing, page through the listing, edit metadata and all four block types, scroll the editor and watch the preview follow, rename and confirm the old URI 404s, publish and open the public page, unpublish, delete, and confirm the listing page and totals stay consistent.

## 7. URI-less article creation

- [x] 7.1 Accept an empty URI in `CreateArticle` and assign a generated unique `untitled-<8 hex>` URI, retrying on collision; verify with `go build ./...` and a `POST /article` carrying `{}` returning `201` with a generated URI.
- [x] 7.2 Remove the URI input from the listing create form, send an empty create payload, navigate to the URI returned by the API, and add the confirmation string; verify `pnpm type-check`.
- [x] 7.3 Update the `article` and `article-editor` deltas and the design for URI-less creation; verify `openspec validate redesign-article-editor --strict`.
- [x] 7.4 Restart the local API and verify a blank create opens the workspace with the generated URI editable in the sidebar; run `go build ./...`, `pnpm type-check`, and `pnpm build`.

## 8. Block duplication

- [x] 8.1 Add a duplicate action to `BlockItem` and insert a new-id copy immediately after the source block in `ArticleEditor`; verify `pnpm type-check`.
- [x] 8.2 Copy a drawing block's scene to a new stored scene before inserting the duplicate, and report a toast if the copy fails; verify both drawings can be edited independently.
- [x] 8.3 Add the `article.editor.duplicateBlock` string and update the article-editor delta and design; verify `openspec validate redesign-article-editor --strict` and `pnpm build`.

## 9. Latest editor refinements

- [x] 9.1 Update the article-editor delta for the shared title/subtitle header, block-only preview, single-step creation, contextual sidebar actions, and save state; verify `openspec validate redesign-article-editor --strict`.
- [x] 9.2 Update the design for the shared header, pane chrome, block-only sync anchors, and save control state; verify `openspec validate redesign-article-editor --strict`.

## Workflow follow-up

- Work on the feature branch, stage and commit in increments, then push with upstream.
- Archive the change after verification.
