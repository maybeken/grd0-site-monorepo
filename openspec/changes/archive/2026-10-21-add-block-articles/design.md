# Design

## Context

See `proposal.md` for motivation. The relevant current state:

- `Blog` stores one markdown string (`api/schema/blog.go`), edited by `SideBySideEditor` with OverType and rendered by `MarkdownDisplay`. Attachments go to R2 through a presigned PUT under `blog/YYYY/MM/`.
- `api/handler/handler.go` registers public routes on `v1` and authenticated writes on `restrictV1`, guarded by `echojwt`. Handlers use GORM soft deletes (`BaseColumns.DeletedAt`).
- `vue-web` is Vue 3 + Vite + Tailwind v4 + alova, with gettext-style YAML produced by a custom Vite plugin. There are no tests and the lint script is broken; `vue-tsc` and `vite build` are the verification gates.
- `excalidraw-web` is a React shell around `@excalidraw/excalidraw`. It autosaves to `localStorage['grd0.draw.scene']` and accepts `save`/`load` postMessage commands, but has no view mode and no readiness signal.
- `MusicView` contains the only music player implementation; `YouTubeEmbed` listens to every window message and polls YouTube forever.
- `api/storage/s3.go` wraps MinIO only for presigned URLs; there is no general read/write blob layer.

## Goals / Non-Goals

**Goals:**

- A new `Article` surface that is additive and can replace `/blog` later without touching existing blog rows.
- Make the article index the public navigation and sitemap entry while keeping `/blog` reachable as a legacy route.
- One block schema shared by the API, editor, and public renderer.
- Reuse of the existing markdown pipeline, CDN helpers, gallery services, music playback, and Excalidraw app.
- An asset layer that works on disk now and moves to object storage without changing article code.

**Non-Goals:**

- Migrating existing blog posts or removing the legacy `/blog` routes. The navigation and sitemap now point at articles, but existing `/blog` links still render.
- Comments, revisions, collaboration, or concurrent editing.
- A generic media library or arbitrary-file blocks.
- Draft sharing via signed links.
- Deleting the legacy blog attachment prefix or endpoint.

## Decisions

### New entity, parallel to Blog

`Article` is a new table and route group. Evolving `Blog` in place would force the 1024-character list truncation, the single-string editor, and the existing `BlogSummary` rendering to change at the same time, risking live posts. A parallel entity lets the new surface stabilise first. The public navigation and sitemap now point at `/article`; the legacy `/blog` routes remain reachable, and a later content import (`[{type:'markdown', md: content}]`) can retire them.

Alternative rejected: a `format` column on `Blog`. It still requires touching every blog read path and list projection.

### URI for the API, UUID for storage

The URI is the public handle and appears in every endpoint. The UUID (`BaseColumns.ID`) is internal and appears only in asset paths, so renaming a URI never moves a file. Writes go through the current URI; a rename is an update whose body URI differs from the path URI.

Duplicate handling is explicit: `POST /article` returns `409` when the URI exists, including unpublished rows. Since `DELETE` unpublishes rather than soft-deletes, the row keeps its URI and republishing reuses it. A unique index on `uri` remains as a backstop against check-then-insert races.

### Blocks stored as a versioned JSON envelope

`Article.Blocks` is a `TEXT` column containing:

```json
{ "version": 1, "blocks": [ { "id": "uuid", "type": "markdown", "md": "..." } ] }
```

Reasons: block order is inherent, the whole article is written atomically, and no block needs to be queried independently. The envelope leaves room for a future migration. Validation happens in the handler before persistence: array shape, known `type`, required fields per type, and unique block IDs.

Alternatives rejected: a relational `article_blocks` table (ordering and reordering cost, no query need) and a schemaless pass-through (no validation, silent renderer failures).

### Block payloads

```ts
type ArticleBlock =
  | { id: string; type: 'markdown'; md: string }
  | { id: string; type: 'image'; collection: string; filename: string; caption?: string }
  | { id: string; type: 'music'; v: string; title: string; artist?: string }
  | { id: string; type: 'excalidraw'; scene_path: string }
```

- `image` stores exactly what `CDNImage` and `GalleryCard` already consume: `collection` is `/gallery/<id>` and `filename` is the asset name.
- `music` snapshots `title` and `artist` so an article does not change when the catalogue does.
- `excalidraw` stores the file path, not the scene. The client derives the scene identifier from the basename for API calls.

### Scene files behind a blob abstraction

```go
type BlobStore interface {
    Put(ctx context.Context, key string, r io.Reader, contentType string) error
    Get(ctx context.Context, key string) (io.ReadCloser, error)
    Exists(ctx context.Context, key string) (bool, error)
    Delete(ctx context.Context, key string) error
    DeletePrefix(ctx context.Context, prefix string) error
}
```

- `LocalStore` writes under `BLOB_LOCAL_PATH` (default `data/scenes`) with `..` and absolute-path rejection.
- A local store and an S3 store are both implemented behind the abstraction. Reads proxy through the API so the client contract is unchanged, with presigned redirects as a later optimisation.
- Scenes are raw JSON: no lz-string, no transformation. Compression existed for the clipboard textarea, not for files.
- Scene path: `scenes/<article-uuid>/<scene-uuid>.json`. Prefix deletion therefore purges an article's drawings in one call.
- The scene write response returns the stored path, so the editor records exactly what the server wrote.

Alternatives rejected: inline scene JSON in the block (bloats every article read and save), a separate SQLite table (same bandwidth without object-storage portability), and presigned direct uploads for scenes (CORS and orphan handling for no benefit at this scale).

### Scene endpoints follow the article URI

```
PUT    /article/:uri/scene/:sid          [auth]   write, returns { scene_path }
GET    /article/:uri/scene/:sid          public   published articles only
GET    /article/:uri/draft/scene/:sid    [auth]   drafts, used by the editor
DELETE /article/:uri/scene/:sid          [auth]   best-effort block removal
DELETE /article/:uri/purge               [auth]   permanently deletes the article and its scene prefix
```

The server resolves the URI, checks publication, builds `scenes/<article.ID>/<sid>.json`, and validates `sid` as a UUID. Clients never supply arbitrary paths. Purge hard-deletes the article row and every scene under its UUID prefix; uploaded markdown images under `article/YYYY/MM/` are not keyed by article and are not removed.

### Excalidraw view mode and handshake

`App.tsx` gains:

- `viewMode` read from `?view=1`.
- `viewModeEnabled` on `<Excalidraw>`, plus `UIOptions` to hide editing affordances if the prop leaves menu items reachable.
- In view mode, skip reading `localStorage` and skip the autosave `onChange` path entirely, so a reader never overwrites their own drawing.
- `post({ ready: true })` once `excalidrawAPI` is set.

After a scene is loaded, the app state is reapplied with `viewModeEnabled` forced to the frame's mode, so a saved `appState` cannot switch a public drawing back to editing. Ctrl/Cmd+S is intercepted in the capture phase; default and propagation are stopped and the frame asks its parent to save, so Excalidraw never handles the shortcut. Embedded drawings use a click-to-interact shield: wheel and trackpad scrolling go to the page until the reader activates the canvas, and leaving the canvas deactivates it.

The parent waits for `ready`, verifies `event.source === iframe.contentWindow` and `event.origin`, then posts `{event:'command', func:'load', payload: sceneJson}`. Without source filtering, two drawings on one article would receive each other's messages.

The existing `/tool/draw` page keeps its current behaviour and textarea; only the article block uses view mode.

### Renderer and editor split

- Public: `BlockRenderer.vue` switches on `type`; `ArticleHeader.vue` renders metadata. The same renderer is used by the editor's draft preview. The index uses `ArticleSummary.vue`, which mirrors the blog listing's summary-card presentation.
- Editor: a local `useArticleDraft` composable owns the article snapshot, dirty state, scene uploads, and save. A composable rather than a Pinia store keeps the state scoped to the one view.
- The editor is a two-column layout: editing controls on the left and an always-visible preview on the right. The preview is a snapshot that refreshes 3 seconds after the last edit and immediately on open or save.
- Markdown blocks are full-width editing surfaces; the article preview replaces the former per-block preview.
- The block palette and save action sit below the last block. Reordering scrolls the moved block into view.
- Ctrl/Cmd+S saves. Dirty state guards route changes, reload/close, and switching to another article.
- A Publish now action beside the publish date sets the timestamp to the current time and saves the article immediately.
- The editor list separates Unpublish from Delete: unpublish keeps the article and its assets; delete calls the purge endpoint after confirmation and removes the article and its scene files.
- Save, conflict, upload, and scene feedback use non-blocking toasts rather than alerts.
- Both sides use the same `ArticleBlock` union from `interfaces/Article.d.ts`.

### Reuse fixes required by the design

- `MarkdownDisplay.vue` currently watches props with argument names reversed (`(oldVal, newVal)` but Vue passes new first), so previews lag one edit. Fix the callback; the article preview depends on it.
- `YouTubeEmbed.vue` must ignore messages from other iframes, parse defensively, clear its polling timer on unmount, and send commands to the cross-origin YouTube frame with `'*'` as the target origin. Targeting the host origin silently drops commands, leaving the player unready and the play button disabled.
- `MusicPlayer.vue` is extracted from `MusicView` (cover, title/artist, progress, play/pause, seek, time). `MusicView` keeps the playlist, prev/next, and ends-of-song advancement.
- The gallery picker must not use `GalleryCard`, which issues one `getGalleryDetail` request per tile. It uses `getGalleryCollectionRaw` plus `loadMoreAssets` and lightweight `CDNImage` tiles.

### Uploads for markdown

A new `PUT /article/attachment/:key` reuses the blog presign handler with an `article/YYYY/MM/` prefix. The editor's drop/pick flow requests the target, PUTs to R2, and inserts `![](/<key>)`. `MarkdownDisplay` already resizes relative image sources, so no renderer change is needed.

### Summary, cover, and author

- Summary and cover are computed in the editor at save and stored: summary from the first markdown block, cover from the first image block's `collection/filename`.
- The list endpoint selects only metadata, summary, and cover, so scene paths and block bodies never reach a list response.
- Author identity comes from JWT claims (`email`, `name`) and the existing `Author` table; request-body author fields are ignored.

### Frontend shape

```
interfaces/Article.d.ts
services/article.ts
views/ArticleIndexView.vue
views/ArticleView.vue
views/ArticleEditorView.vue
components/article/ArticleSummary.vue
components/article/BlockRenderer.vue + blocks/*
components/editor/article/ArticleEditor.vue, ArticleList.vue,
  BlockPalette.vue, BlockItem.vue, blocks/*Editor.vue,
  GalleryPicker.vue, MusicPicker.vue, useArticleDraft.ts
components/layout/ToastContainer.vue
components/music/MusicPlayer.vue
composables/useToast.ts
helpers/article.ts
helpers/environment.ts
```

Routes: `/article`, `/article/editor`, `/article/:slug`. Vue Router ranks the static editor segment above the parameterised slug; the existing `/blog/editor` versus `/blog/:slug` pair already relies on this. The main navigation points at `/article`; the legacy `/blog` routes remain registered. The sitemap lists articles and no longer lists blog URLs.

### Local development authentication

`APP_ENV=LOCAL` makes the API skip OIDC initialization and replace the JWT guard with `LocalAuthGuard`, which injects a synthetic author (`LOCAL_AUTHOR_EMAIL`, `LOCAL_AUTHOR_NAME`). Auth endpoints return 404 in this mode. The frontend treats `import.meta.env.DEV` or `VITE_APP_ENV=LOCAL` as local, skips the token-refresh redirect, and omits the bearer header. The bypass is explicit, logs a startup warning, and must never be set in a deployed environment. `local-test.sh` sets it for the local Docker workflow.

## Risks / Trade-offs

- [Photo-heavy scenes produce large JSON files] → enforce `ARTICLE_SCENE_MAX_BYTES` with `http.MaxBytesReader`, report the limit in the editor, and mount drawing iframes only when they scroll into view.
- [Scene files can be orphaned by an abandoned drawing] → delete on block removal, retain on unpublish, and leave a future sweep as an explicit open question rather than silent data loss.
- [Local blob storage is ephemeral on Lambda] → document that the S3 backend is required there; the abstraction is the migration path.
- [Multiple YouTube iframes cross-talk] → source-filter messages and clean up polling; verify with two music blocks on one page.
- [`viewModeEnabled` may not hide every menu entry] → `UIOptions` disables editing actions, and scene loads force `viewModeEnabled` so a saved app state cannot re-enable editing.
- [Wheel scrolling over an embedded canvas zooms the drawing instead of scrolling the page] → a click-to-interact shield passes scrolling to the page until the reader activates the canvas.
- [Unsaved edits are lost on navigation or reload] → track a dirty snapshot and guard route changes, `beforeunload`, and article switches.
- [A 3-second preview debounce can lag the editor] → show the previous snapshot rather than blocking editing; refresh immediately on open and save.
- [Block JSON has no schema enforcement at the database] → validate on write and treat unknown `version` values as an error.
- [No automated tests exist and lint is broken] → verification is `pnpm type-check`, `pnpm build`, `go build`, and manual browser checks. Do not invent test commands.
- [Rate limiting is 20 requests/second per IP] → debounce scene saves instead of writing on every drawing change.
- [Public scene reads reveal a draft if the UUID leaks] → draft scenes are served only through the authenticated draft route; published scenes are public by definition.

## Migration Plan

1. Deploy the API first: additive table, blob store, and `/article` route group. No existing route or table changes.
2. Add `scenes/` to `api/data/.gitignore` when using the default local path.
3. Build `excalidraw-web` before `vue-web` because of the symlinked `public/excalidraw`, then deploy `vue-web` with the new routes, the article navigation entry, and the hidden editor links; the legacy `/blog` routes remain available.
4. Deploy the sitemap worker with articles listed and blog removed.
5. Rollback: restore the navigation and sitemap entries, remove the frontend article routes and the `/article` route group; `/blog` remains available throughout.

## Open Questions

- When to migrate existing blog content and whether to redirect `/blog` to `/article`.
- Whether a music block should pause other players on the same article; the specs currently require independent control only.
- Whether to sweep orphaned uploaded markdown images, which are not keyed by article UUID.
