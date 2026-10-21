# Tasks

## 1. Article entity and API

- [x] 1.1 Add `schema.Article` and register it in `database/migration.go`; verify `cd api && go build ./...` and that a fresh `DB_PATH` run creates the `articles` table.
- [x] 1.2 Validate the versioned block envelope on write (array shape, known types, required payload fields, unique block IDs); verify an unknown type and a missing `scene_path` are each rejected with `400` and the stored article is unchanged.
- [x] 1.3 Implement `POST /article` with URI format validation, `409` on any existing URI, and author resolution from JWT claims; verify the first call returns `201` and the repeat returns `409`.
- [x] 1.4 Implement `GET /article`, `GET /article/:uri`, `GET /article/:uri/draft`, and `GET /article/all`; verify drafts return `404` publicly, the draft route returns blocks with a token, and list responses carry no block JSON.
- [x] 1.5 Implement `PUT /article/:uri` including rename via `body.uri`; verify a rename makes the old URI return `404` with no redirect, and renaming onto a taken URI returns `409` without changing the record.
- [x] 1.6 Implement `DELETE /article/:uri` as unpublish; verify the public read returns `404`, the block list and stored assets remain, and republishing restores access.
- [x] 1.7 Derive and store summary and cover on save from the first markdown and first image block; verify list entries include both and that an article with neither returns an empty summary and no cover.
- [x] 1.8 Register the `/article` routes and document the endpoints and new environment variables in `api/README.md`; verify the documented local run starts and the route table matches the document.

## 2. Article asset storage

- [x] 2.1 Define `storage.BlobStore` and implement `LocalStore` with `Put`, `Get`, `Exists`, `Delete`, and `DeletePrefix`; verify `cd api && go build ./...` and that keys containing `..` or absolute paths are rejected.
- [x] 2.2 Add `BLOB_BACKEND`, `BLOB_LOCAL_PATH`, and `ARTICLE_SCENE_MAX_BYTES` configuration with local defaults; verify a locally run API writes a scene beneath `data/scenes`.
- [x] 2.3 Implement `PUT`, `GET`, and `DELETE` scene endpoints plus the authenticated draft scene read, resolving the article URI to its UUID and enforcing publication; verify a public read of a draft scene is `404` and the draft route succeeds with a token.
- [x] 2.4 Enforce the scene size cap while streaming without leaving a partial file; verify an oversized body is rejected and the target directory is unchanged.
- [x] 2.5 Implement `PUT /article/attachment/:key` as a presigned upload under `article/YYYY/MM/`; verify it returns a URL and key with authentication and is rejected without it.
- [x] 2.6 Ignore the local scene directory in `api/data/.gitignore`; verify `git status` reports a test scene as ignored.
- [x] 2.7 Delete a block's scene on a best-effort basis when the block is removed during an article update, and document the S3 backend contract in `api/README.md`; verify the file disappears after saving without the block and that a simulated delete failure still persists the article.

## 3. Shared frontend types, services, and strings

- [x] 3.1 Add `interfaces/Article.d.ts` with the block union and article types; verify `cd vue-web && pnpm type-check`.
- [x] 3.2 Add `services/article.ts` covering create, reads, draft read, list, update, unpublish, scene read/write/delete, and attachment presign; verify `pnpm type-check`.
- [x] 3.3 Add all new interface strings and route title keys to `src/locales/en.yaml` in gettext form; verify `pnpm build` transforms them and no raw `msgid` appears in the UI.

## 4. Shared media groundwork

- [x] 4.1 Add `?view=1` mode to `excalidraw-web/src/App.tsx`: `viewModeEnabled`, no `localStorage` read or autosave, a `{ready:true}` message, and `UIOptions` if editing affordances remain; verify in a browser that view mode hides editing and never writes `localStorage`.
- [x] 4.2 Build `excalidraw-web` and verify `vue-web/public/excalidraw` resolves to the new `dist`.
- [x] 4.3 Extract `MusicPlayer.vue` from `MusicView` and refactor `MusicView` to use it; verify `/music` still plays, seeks, and advances at the end of a song.
- [x] 4.4 Scope `YouTubeEmbed` message handling to its own iframe, parse messages defensively, send commands to the cross-origin player with the correct target origin, and clear the polling timer on unmount; verify two players on one page do not interfere or leak timers.
- [x] 4.5 Fix `MarkdownDisplay`'s inverted watch arguments; verify typing updates the preview immediately without a save.

## 5. Public article experience

- [x] 5.1 Add the `/article` index view with title, subtitle, author, timestamps, summary, and cover; verify drafts and future-dated articles are absent.
- [x] 5.2 Add the `/article/:slug` view with `ArticleHeader` and `BlockRenderer`; verify blocks render in stored order and an unknown or unpublished slug shows the not-found state.
- [x] 5.3 Implement markdown and image block renderers; verify sanitisation, CDN-resized relative images, and captions.
- [x] 5.4 Implement the music block with `MusicPlayer` and the stored title and artist; verify playback and that two blocks control only their own song.
- [x] 5.5 Implement the drawing block: wait for `ready`, verify `event.source` and origin, load the scene, and allow view mode only; verify editing is unavailable, `localStorage` is untouched, and two drawings load independently.
- [x] 5.6 Render per-block failure states (missing scene, unavailable gallery asset) without breaking the rest of the article; verify by pointing a block at a missing scene and a missing asset.
- [x] 5.7 Add the article entry to `workers/sitemap-generator/index.js`; verify the generated sitemap contains published article URLs and no drafts.

## 6. Article editor

- [x] 6.1 Add `/article/editor` with the article list and the create-or-open flow; verify a duplicate URI shows conflict feedback and opens nothing.
- [x] 6.2 Implement metadata editing and rename through the current URI; verify the editor continues under the returned new URI and duplicate renames are rejected.
- [x] 6.3 Implement the block palette, ordering controls, block removal, and a single save of metadata plus blocks; verify a reload preserves order and removal.
- [x] 6.4 Implement the markdown block editor with a full-width OverType surface and drag or pick uploads; the always-visible article preview provides live rendering. Verify an uploaded image inserts a relative reference and renders in the preview.
- [x] 6.5 Implement the gallery picker over collections and `loadMoreAssets` without per-tile detail requests; verify selection stores a collection and filename and no N+1 detail calls appear in the network log.
- [x] 6.6 Implement the music picker over the music catalogue and snapshot title and artist; verify a later catalogue title change leaves the saved block unchanged.
- [x] 6.7 Implement the drawing block editor: generate a scene UUID, write on save with debounce, store the returned path, and reload the scene on open; verify reopening shows the same drawing.
- [x] 6.8 Implement draft preview through the authenticated draft endpoint using the public renderer; verify an unpublished article previews without publishing.
- [x] 6.9 Verify the editor fails closed without a token and exposes no draft content.

## 7. Integration and documentation

- [x] 7.1 Run `cd api && go build ./...` and `cd vue-web && pnpm type-check`; resolve every failure.
- [x] 7.2 Run `cd excalidraw-web && pnpm build` then `cd vue-web && pnpm build`; verify the build order succeeds and the symlinked drawing assets resolve.
- [x] 7.3 Update `vue-web/README.md` with the new routes and components, and the root `README.md` if the layout listing changes; verify every documented command still runs.
- [ ] 7.4 Manual end-to-end: create a draft, add all four block types, preview, publish, open publicly, rename and confirm the old URI 404s, unpublish and confirm public access stops while assets remain.
- [x] 7.5 Confirm `/blog`, its editor, and its API endpoints behave exactly as before.

## 8. Follow-up refinements

- [x] 8.1 Editor: always-visible, right-hand preview with a 3-second debounced snapshot; editor on the left; palette and save below the last block.
- [x] 8.2 Editor: full-width markdown surface; per-block preview removed in favour of the article preview.
- [x] 8.3 Editor: Ctrl/Cmd+S saves, with repeat/concurrency guards and drawing-frame forwarding.
- [x] 8.4 Editor: dirty-state guards for route changes, reload/close, and opening another article.
- [x] 8.5 Editor: reordering scrolls the moved block into view.
- [x] 8.6 Drawing: force view mode after scene load, stop Excalidraw from handling Ctrl+S, and add a click-to-interact scroll guard.
- [x] 8.7 Public listing: article index uses the blog listing summary-card design.
- [x] 8.8 Navigation and sitemap: `/article` replaces `/blog` in the nav and sitemap; legacy `/blog` routes remain.
- [x] 8.9 Feedback: blocking alerts replaced with dismissible toasts.
- [x] 8.10 Music: fix the cross-origin YouTube command target so the player becomes ready and play is enabled.
- [x] 8.11 Local development: `APP_ENV=LOCAL` API auth bypass and frontend local mode.
- [x] 8.12 Publication updates: omitted `published_at` retains state; explicit null unpublishes.
- [x] 8.13 API: `DELETE /article/:uri/purge` hard-deletes the article and its scene prefix; the URI becomes reusable.
- [x] 8.14 Editor list: separate Unpublish and permanent Delete actions, with confirmation and toast feedback.
- [x] 8.15 Editor: Publish now action beside the publish date sets the timestamp and saves immediately.
- [x] 8.16 Public index: remove the page subtitle; article cards carry their own metadata.

## Workflow follow-up

- Work on a feature branch created from up-to-date `origin/main`, stage and commit in increments, then push with upstream.
- Archive the change after verification.
