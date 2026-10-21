# Design

## Context

See `proposal.md` for motivation. The relevant current state:

- `ArticleEditorView.vue` renders `ArticleList` and `ArticleEditor` on one page, using local `selected` state. `ArticleEditor` is a two-column grid with a sticky, internally scrolling preview. `useArticleDraft` owns the article snapshot, dirty state, and save.
- The admin list calls `GET /article/all`, which returns every row ordered by `updated_at DESC`. The public list is `GET /article`. The only pagination in the repo is the client-side `PostTable.vue`, used by the legacy blog editor, whose slicing and page-count logic are both off by one.
- `api/handler/gallery.go` already provides `parsePagination` (default 40, max 200) and `schema.PaginatedResponse[T]` with `data`, `total`, `page`, `page_size`, and `total_pages`. `interfaces/Gallery.d.ts` already types that envelope.
- Blocks have stable UUIDs in both `blocks` and the debounced preview clone, but `BlockRenderer` renders them without wrapper elements, so there is nothing to anchor a preview scroll position to.
- Markdown in the preview is rendered HTML, while the editor uses OverType. OverType 2.6.1 defaults to a fixed-height container with `overflow-y: auto`, so markdown blocks can scroll internally.
- The preview is a JSON clone refreshed three seconds after a change. Preview images are lazy-loaded, and drawing blocks are a fixed `h-[60vh]`.
- Duplicate URIs are detected client-side by scanning the admin list. The API returns `409 Conflict` on create and rename duplicates, and `errorStatus()` already reads status codes.
- Auth on the editor is a `sessionStorage` token check in `ArticleEditorView`, not a router guard.

## Goals / Non-Goals

**Goals:**

- Two routes with a route-owned article identity and a shareable editing URL.
- Server-side pagination reusing the existing gallery envelope and helpers.
- A preview that follows the editor by block identity with no feedback loop.
- A bounded three-pane workspace that never scrolls the document while editing.
- Correct duplicate-URI detection under pagination.
- Markdown blocks that grow so the editor pane is the only scroll surface.

**Non-Goals:**

- Paginating the public `/article` index or the legacy blog list.
- Search, status filters, or alternate sort orders in the admin listing.
- Virtualising either pane or adding a scroll-sync dependency.
- Changing the block schema or the public article view.
- Preserving the old combined `/article/editor` experience.

## Decisions

### Two routes, one listing, one workspace

```
/article                  public index        unchanged
/article/editor           admin listing       new page, paginated
/article/editor/:uri      editing workspace   new page, route identifies the article
/article/:slug            public article      unchanged
```

Vue Router ranks the static `editor` segment above `:slug`, and the editing route has three segments, so there is no collision. Opening an article becomes `router.push`; the listing's page lives in `?page=`, so browser back restores it. A successful rename uses `router.replace(/article/editor/<new-uri>)` to keep the address bar honest. An unknown URI renders an error with a link back to the listing.

Alternative rejected: keeping the list mounted beside the editor and using a child route or overlay. It preserves the cramped single-page feel and makes the listing's pagination state fight the editing route.

### URI-less creation gets a generated placeholder

An article's URI is unique, non-null, and used by the workspace route and draft reads. Rather than make it nullable and add ID-addressed draft endpoints, `CreateArticle` accepts an omitted URI and assigns `untitled-<8 hex>`, retrying on the unlikely collision. The generated value satisfies the URI rules, keeps the unique index and every existing route unchanged, and is editable in the sidebar immediately. The listing create form therefore has no URI field at all: one Create action produces a draft and opens it, and the author names it in the sidebar.

Alternatives rejected: a nullable URI with `/article/editor/draft/:id` (a second addressing scheme across the API and router), a client-generated URI (trusts the client and duplicates collision handling), and blocking creation until a URI is entered (the behaviour this change removes).

### Pagination reuses the gallery contract

`GetArticleByAdmin` gains `parsePagination(c)`, a `Count` on the same predicate, and `Offset((page-1)*pageSize).Limit(pageSize)`, returning `schema.PaginatedResponse[ArticleListEntry]`. The client sends `page_size=10`; the gallery default of 40 is too many table rows. Ordering stays `updated_at DESC`.

`parsePagination` does not clamp an out-of-range page, so the endpoint returns an empty `data` with the correct totals. The listing treats an empty page above page 1 as a signal to navigate to the previous page.

Alternative rejected: a bespoke `{items, total}` shape. The gallery envelope is already typed on the frontend, consumed in two handlers, and validated by the project.

### Conflict detection moves to the status code

`useArticleDraft.uriExists()` scans `listArticlesAdminRaw()`. After pagination that scan sees one page, so a duplicate outside the page would be misreported as a generic failure. All four duplicate paths return `409`, so create and save distinguish conflict with `errorStatus(error) === 409`. This also removes an extra admin list request from every failed save.

alova's global `responded` currently resolves a non-2xx body as a value, so the admin instance must throw an error carrying `status` and `response` for non-OK responses. Every admin caller already uses try/catch, and `errorStatus` then reads the status; without this, a `409` would look like a malformed success.

### Bounded workspace with a shared header

```
+----------------------------------------------------------------------------------+
| back / status / save                                                             |
+----------------------------------------------+-----------------------------------+
| title / subtitle (shared, underline inputs)  |                                   |
+----------------------+-----------------------+ SIDEBAR   overflow-y             |
| EDITOR   overflow-y  | PREVIEW   overflow-y  | URI, publish actions              |
| blocks               | rendered blocks only  | timestamps, public link           |
| palette              | scrollbar hidden      | delete                            |
+----------------------+-----------------------+-----------------------------------+
```

The title and subtitle inputs sit in a shared header above the editor and preview, spanning their combined width but not the sidebar; they use an underline style rather than a bordered pill. Neither pane carries its own border: a single vertical divider on the preview's left separates them. The preview keeps `overflow-y-auto` but hides its scrollbar.

The workspace height is `calc(100vh - navbar)`. The shared header and toolbar do not scroll; each pane does. The document does not scroll while editing, which makes the editor's `scrollTop` a meaningful sync source. The editor and preview split the remaining width evenly, and the sidebar is a fixed-width column at `lg`; below `lg` the sidebar becomes a toggled panel and the editor and preview become an `Edit | Preview` toggle with sync disabled.

Alternative rejected: letting the page scroll and making only the preview sticky. The source scroll position would then be the window's, and mapping it to the preview container would require compensating for the page header and sidebar.

### Scroll sync by block anchor

Each block gets a wrapper carrying its block ID. The shared title/subtitle header sits outside both scroll panes, so it is not part of the anchor set:

```
editor pane                                  preview pane
+---------------------------+                +---------------------------+
| wrapper[data-block-id=A]  |  scroll --->   | wrapper[data-block-id=A]  |
| wrapper[data-block-id=B]  |                | wrapper[data-block-id=B]  |
+---------------------------+                +---------------------------+
```

On an editor scroll, throttled to animation frames:

1. Measure each wrapper's top and height relative to its pane, for both panes.
2. Find the editor anchor containing `scrollTop`.
3. Compute `t = (scrollTop - src.top) / src.height`.
4. Set `preview.scrollTop = dst.top + t * dst.height`, clamped to the preview's scrollable range.
5. If the editor is in trailing space past the last anchor, put the preview at its maximum.

Measurement happens each frame rather than being cached, so typing growth, image loads, add/remove, and preview refreshes cannot leave stale anchors behind. With a normal article's block count the cost is negligible.

Anchors are matched by block ID, not index. During the debounce window after adding a block, the preview may not have that ID yet; the sync falls back to the nearest preceding anchor and catches up when the snapshot lands. Scrolling the preview alone never writes to the editor, so there is no feedback loop and no echo lock.

After a debounced refresh, a structural change, or a preview resize, the preview realigns from the editor's current position. The editor preview renders markdown images eagerly so lazy loading does not change anchor heights during a scroll.

Alternatives rejected: proportional scroll sync, which drifts because raw markdown and rendered HTML have different heights; bidirectional sync (explicitly not wanted) which needs echo detection; and an IntersectionObserver that snaps block by block, which is chunkier than continuous alignment.

### Block duplication copies data, and drawings copy their scene

Duplicating takes the source block, assigns a fresh block UUID, and inserts the copy at the source's current index plus one. Markdown, image, and music blocks are value copies, so the two blocks are independent immediately.

A drawing block would otherwise share its `scene_path`, making both blocks the same drawing. Instead, the editor fetches the source scene through the authenticated draft endpoint, writes it under a new scene UUID, and records that path on the copy before inserting it. If the scene copy fails, the duplicate is not inserted and a toast reports the failure. The source's last saved scene is used; unsaved canvas changes are not included until the article is saved.

Alternative rejected: sharing the scene path, which makes edits to either block change both.

### OverType auto-resize is a prerequisite

`MarkdownBlockEditor` passes `autoResize: true` with a minimum height and no maximum. Without it, each markdown block has an internal scrollbar, so a long block's text cannot be mapped to the preview and the editor pane would scroll by block only. The measured wrapper includes the toolbar and stats, so anchors remain correct.

### Debounce tightened for a live feel

The snapshot debounce drops from three seconds to roughly 750 ms. The preview still renders a snapshot rather than the live block array, so a large markdown document is not reparsed on every keystroke. It refreshes immediately on open and save, and realigns to the editor after every refresh.

### Save state lives on the save control

The toolbar shows publication status as plain text. The save control carries the save state itself: Save while there are unsaved changes, Saved (muted) when clean, and Saving while in flight. There is no separate unsaved marker.

### Listing behaviour

- Pager shows previous/next, a windowed set of page numbers, and `1-10 of N`, driven by `total_pages`.
- The page is written to the query string on change and read on mount.
- The New Article action creates after a single confirmation; there is no intermediate create form. The API returns the generated URI and the listing opens that workspace. Failure states remain inline.
- Unpublish and permanent delete stay on the listing, with confirmation and toast feedback. The sidebar mirrors delete but does not duplicate unpublish.
- Row columns: title, URI, publication status, updated, published.

### Sidebar contents

URI with the `/article/` prefix and an icon copy control on the URI label row; publish datetime and a draft/scheduled/published badge; contextual publication actions (Publish now hidden when published, Unpublish hidden when draft, both shown when scheduled); created and updated timestamps; an open-public-page link when published; and a permanent delete action. Title and subtitle live in the shared header above the panes, not in the sidebar.

## Risks / Trade-offs

- [One-way sync means a manually scrolled preview is pulled back on the next editor scroll] → This is the chosen behaviour and matches "the preview follows the editor"; the active-block highlight makes the relationship visible.
- [Measuring both panes every animation frame could cost on very long articles] → Keep measurement inside the rAF callback, skip when the editor has not scrolled, and fall back to cached anchors plus a `ResizeObserver` if a profile shows a problem.
- [A newly added block is absent from the preview until the debounce fires] → Match anchors by ID and align to the nearest preceding anchor; the next refresh catches up.
- [Images and drawings change preview heights] → Render markdown images eagerly in the editor preview, keep the drawing height fixed, and realign after resize observations.
- [An out-of-range page returns empty after deletions] → The listing moves to the previous page when the current page is empty and is not the first.
- [OverType auto-resize interacts with its toolbar and stats] → Verify growth and anchor measurement manually across empty, short, and long blocks.
- [Nested scrolling on touch devices] → Disable sync and use editor/preview tabs below the side-by-side breakpoint.
- [No automated tests and the lint script is broken] → Verification is `vue-tsc`, `vite build`, `go build`, and manual browser checks. Do not invent test commands.

## Migration Plan

None. The article surface has not been deployed, so there is no article data to migrate and the `/article/all` response shape can change without compatibility shims or deploy sequencing. Rollback is a plain revert of the change.

## Open Questions

- Whether the listing later gains status filters or title search. The paginated endpoint makes either cheap, but neither is required now.
- Whether the listing should show a derived status badge or the raw publish timestamp only. The sidebar needs the badge; the listing can start with the timestamp.
