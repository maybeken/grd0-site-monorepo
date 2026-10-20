# Proposal

## Why

Blog posts are a single markdown string, so a post cannot mix prose with gallery photography, music, or a drawing without pasting raw links or compressed scene text. The existing Excalidraw tool already proves the drawing workflow, and the gallery and music pages already own their media, but none of it is composable in an article. A block model makes each of those a first-class part of a post while reusing the renderers and services that already exist.

## What Changes

- Introduce a block-based `Article` entity with its own `/article` API, independent of `/blog`.
- Make the article index the public navigation entry and sitemap source. The legacy `/blog` routes, editor, and API remain reachable for existing links but are no longer linked or sitemapped.
- Support four block types: markdown (with image uploads), gallery image, music player, and Excalidraw drawing.
- Add a new article editor that creates an empty, unpublished article first and then edits its blocks. The editor keeps an always-visible preview beside the editing surface, refreshes it 3 seconds after changes, saves with Ctrl/Cmd+S, and guards unsaved changes on navigation, reload, and article switches.
- Add a public article view that renders blocks, presents drawings in read-only mode from first paint, and does not let embedded drawings capture page scrolling until activated.
- Store Excalidraw scenes as JSON files behind a storage abstraction with local and S3 backends. The article block references the stored path.
- Reject duplicate article URIs on create and rename. Renaming a URI is destructive to the old link: no aliases or redirects.
- Extract the music player out of `MusicView` so the article block and the music page share one implementation, and fix the YouTube command target so the player becomes ready.
- Extend the Excalidraw iframe app with a view mode and a ready handshake, force view mode after scene loads, and forward Ctrl+S to the editor.
- Replace blocking alerts with non-blocking toasts across the app.
- Add an explicit local-development auth bypass (`APP_ENV=LOCAL`) that skips OIDC and injects a synthetic author.
- Keep the legacy blog rows, routes, editor, and endpoints intact; no data migration is included.

## Capabilities

### New Capabilities

- `article`: the Article entity, block schema, URI identity, create/read/update/unpublish behaviour, uniqueness, and publishing rules.
- `article-editor`: the create-then-blocks authoring flow, block palette, per-block editing, reordering, uploads, and scene authoring.
- `article-view`: public rendering of each block type, including read-only drawings and music playback.
- `article-asset-storage`: the blob store abstraction, scene file lifecycle keyed by article UUID, and presigned attachment uploads for markdown.

### Modified Capabilities

None. Existing specs are unaffected: `/blog`, the music page's visible behaviour, and the UI string extraction process stay as specified.

## Impact

- `api/`: new `Article` schema and handler, migration entry, and a `BlobStore` abstraction with local and S3 implementations. New route group under `/article`, plus an `APP_ENV=LOCAL` auth bypass and synthetic author guard.
- `vue-web/`: new interfaces, service, views, article summary and block components, block editors, pickers, and a shared toast system. Router entries, locale strings, article navigation, unsaved-change guards, and fixes to `MarkdownDisplay` and the YouTube player.
- `vue-web/src/components/music/`: `MusicPlayer` extracted from `MusicView`; `YouTubeEmbed` message handling scoped to its own iframe with the correct cross-origin command target.
- `excalidraw-web/src/App.tsx`: view mode forced after scene loads, ready message, autosave suppression in view mode, Ctrl+S forwarding, and no Excalidraw shortcut handling.
- `workers/sitemap-generator/index.js`: article listing; the blog entry is removed.
- `api/data/`: a gitignored directory for scene files.
- Dependencies: no new runtime dependencies; `lz-string` is not used for scene files.
