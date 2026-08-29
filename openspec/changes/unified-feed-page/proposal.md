## Why

Blog and music content currently live on separate pages with separate APIs, data models, and UIs. This creates friction for content creation — a user who wants to share a blog post alongside a music track has no way to combine them. A unified feed page with block-based composition lets users create rich, multi-media entries that mix writing and music in a single cohesive experience, replacing the underutilized standalone blog and music pages.

## What Changes

- New `FeedEntry` and `FeedBlock` database tables. FeedEntry has `subtitle` and `author_id` (FK to authors). FeedBlock references existing Blog and Music rows via FK.
- New feed API endpoints: `GET /v2/feed` (listing with summary and cover image, paginated for infinite scroll), `GET /feed/:slug` (full content), `PUT /feed/:slug` (upsert), `DELETE /feed/:slug` (soft-delete).
- Upsert API auto-creates Blog/Music rows for inline blocks, or references existing entries via searchable lookup.
- New `/feed` route — listing page with preview cards (cover image, summary, author), infinite scroll, inline expansion on "Continue Reading."
- New `/feed/:slug` route — full feed entry detail with subtitle, author, and all blocks rendered.
- New `/feed/editor` route — Notion-style block editor. Supports adding blog blocks (markdown editor), music blocks (YouTube ID form), and reference blocks (searchable dropdown for existing Blog/Music). Up/down buttons for block reordering. Save draft / Publish actions.
- New Pinia store `usePlayerCoordinator` to manage multiple inline music players on a page — when one plays, all others pause.
- Existing `/blog`, `/blog/:slug`, `/music` routes and their APIs remain functional (no breaking changes).

## Capabilities

### New Capabilities
- `feed-data-model`: FeedEntry and FeedBlock database tables, Go schema definitions, GORM migrations.
- `feed-api`: REST endpoints for feed CRUD (list, get, upsert, delete) with preview vs full content modes and pagination.
- `feed-listing-page`: `/feed` route with preview cards, infinite scroll, inline expansion, and player coordination.
- `feed-detail-page`: `/feed/:slug` route rendering full feed entry with all block types.
- `feed-editor`: `/feed/editor` route with Notion-style block composer, supporting blog/music/reference blocks with reordering.
- `player-coordinator`: Pinia store for managing multiple inline music players with mutual exclusion.

### Modified Capabilities
<!-- No existing specs are modified. Blog and music APIs/pages remain unchanged. -->

## Impact

- **Database**: Two new tables (`feed_entries`, `feed_blocks`) with FK constraints to existing `blogs` and `music` tables.
- **API (Go)**: New handler file `handler/feed.go`, new schema in `schema/feed.go`, route registration in `handler/handler.go`.
- **Vue frontend**: New views (`FeedView.vue`, `FeedDetailView.vue`, `FeedEditorView.vue`), new components for block rendering, new router entries, new Pinia store, new service file for feed API calls.
- **i18n**: New translation keys for feed-related UI strings.
- **No breaking changes**: Existing blog and music functionality is preserved.
