## 1. Database Schema

- [x] 1.1 Create `api/schema/feed.go` with FeedEntry and FeedBlock struct definitions, including GORM tags for relationships, indexes, and soft delete
- [x] 1.2 Add FeedEntry and FeedBlock to GORM auto-migration in API startup code
- [x] 1.3 Verify tables are created correctly on API startup

## 2. Feed API Handlers

- [x] 2.1 Create `api/handler/feed.go` with GetFeed handler (list with pagination, preview mode, truncation)
- [x] 2.2 Implement GetFeedBySlug handler (single entry with full content)
- [x] 2.3 Implement UpsertFeed handler (create/update with blocks, auto-create Blog/Music rows, transaction-safe)
- [x] 2.4 Implement DeleteFeed handler (soft-delete)
- [x] 2.5 Register feed routes in `api/handler/handler.go` (GET /feed, GET /feed/:slug, PUT /feed/:slug, DELETE /feed/:slug)
- [x] 2.6 Add authentication middleware to PUT and DELETE routes
- [x] 2.7 Test feed API endpoints manually with curl or Postman

## 3. Frontend Services

- [x] 3.1 Create `vue-web/src/services/feed.ts` with alova methods for feed API calls (listFeed, getFeedBySlug, upsertFeed, deleteFeed)
- [x] 3.2 Define TypeScript interfaces for FeedEntry, FeedBlock, and request/response types in `vue-web/src/interfaces/Feed.d.ts`
- [x] 3.3 Add feed service functions using useRequest hook for reactive data fetching

## 4. Player Coordinator Store

- [x] 4.1 Create `vue-web/src/stores/playerCoordinator.ts` with Pinia store (activePlayerId, play, pause actions)
- [x] 4.2 Verify store is registered in the app's Pinia instance

## 5. Feed Listing Page

- [x] 5.1 Create `vue-web/src/views/FeedView.vue` with infinite scroll pagination (fetch first page, scroll listener, cursor-based pagination)
- [x] 5.2 Implement feed entry preview card component showing title, date, and truncated block content
- [x] 5.3 Implement inline expansion logic (fetch full entry on "Continue Reading", render full blocks, collapse button)
- [x] 5.4 Create blog block preview component (reusing BlogSummary patterns)
- [x] 5.5 Create music block preview component (title and artist, no player)
- [x] 5.6 Create blog block full component for expanded view (reusing BlogDetail patterns)
- [x] 5.7 Create music block full component for expanded view (reusing YouTubeEmbed, with player coordinator integration)
- [x] 5.8 Implement loading states (skeleton placeholders for initial load, spinner for pagination)
- [x] 5.9 Implement empty state message when no feed entries exist
- [x] 5.10 Implement error state with retry button

## 6. Feed Detail Page

- [x] 6.1 Create `vue-web/src/views/FeedDetailView.vue` that fetches and renders a single feed entry by slug
- [x] 6.2 Render all blocks in order with full content (blog markdown, music players)
- [x] 6.3 Integrate player coordinator for music blocks (only one plays at a time)
- [x] 6.4 Implement loading state with skeleton placeholders
- [x] 6.5 Implement NotFound state for non-existent or unpublished entries
- [x] 6.6 Add router beforeEach guard to pre-fetch feed entry and extract title for document title suffix

## 7. Feed Editor

- [x] 7.1 Create `vue-web/src/views/FeedEditorView.vue` with route guard requiring authentication
- [x] 7.2 Implement feed entry title input
- [x] 7.3 Implement "Add Block" dropdown with options: Blog, Music, Reference Existing
- [x] 7.4 Implement blog block editor (reusing SideBySideEditor component for markdown)
- [x] 7.5 Implement music block editor (YouTube ID, title, artist inputs)
- [x] 7.6 Implement reference block editor with searchable dropdown (query blog/music APIs, display results, select)
- [x] 7.7 Implement remove block functionality
- [x] 7.8 Implement up/down buttons for block reordering
- [x] 7.9 Implement "Save Draft" button (PUT /feed/:slug with published_at = null)
- [x] 7.10 Implement "Publish" button (PUT /feed/:slug with published_at = current timestamp)
- [x] 7.11 Implement validation (title required, at least one block, blog content required, music YouTube ID required)
- [x] 7.12 Implement loading existing feed entry for editing (populate editor state from API response)
- [x] 7.13 Implement success/error messages for save operations

## 8. Router Configuration

- [x] 8.1 Add /feed route to `vue-web/src/router/index.ts` pointing to FeedView.vue
- [x] 8.2 Add /feed/:slug route pointing to FeedDetailView.vue with pre-fetch guard
- [x] 8.3 Add /feed/editor route pointing to FeedEditorView.vue with requiresAuth meta
- [x] 8.4 Add /feed/editor/:slug route for editing existing entries with requiresAuth meta
- [x] 8.5 Add i18n translation keys for feed-related route titles

## 9. i18n Translations

- [x] 9.1 Add translation keys to `vue-web/src/locales/en.yaml` for feed UI strings (page titles, button labels, validation messages, empty states, error messages)
- [x] 9.2 Add translation keys to `vue-web/src/locales/zh.yaml` (or other language files) for feed UI strings

## 10. Integration and Testing

- [x] 10.1 Test feed listing page with various scenarios (empty, loading, pagination, inline expansion)
- [x] 10.2 Test feed detail page with blog and music blocks
- [x] 10.3 Test player coordination (multiple music players, only one plays at a time)
- [x] 10.4 Test feed editor (add/remove/reorder blocks, save draft, publish)
- [x] 10.5 Test reference blocks (search and select existing blog/music)
- [x] 10.6 Verify existing blog and music pages still work (no breaking changes)
- [x] 10.7 Run type-check and build to verify no TypeScript errors
