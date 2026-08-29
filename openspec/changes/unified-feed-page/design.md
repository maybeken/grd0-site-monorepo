## Context

The current site has separate blog and music pages with distinct data models. Blog posts are markdown content with metadata, music entries are YouTube video references. Both are underutilized — users want to create richer content that combines text and music in a single narrative.

The feed system introduces a block-based composition model where a FeedEntry contains an ordered list of FeedBlocks. Each block references either a Blog or Music row. The existing Blog and Music tables remain unchanged; the feed layer is additive.

Constraints:
- SQLite database via GORM (existing pattern)
- Go API with Echo framework (existing pattern)
- Vue 3 + Pinia + TypeScript frontend (existing pattern)
- No breaking changes to existing blog/music functionality
- Must support infinite scroll for feed listing
- Multiple inline music players must coordinate (only one plays at a time)

## Goals / Non-Goals

**Goals:**
- Enable block-based content composition mixing blog posts and music tracks
- Provide a Notion-style editor for creating feed entries with multiple blocks
- Support inline expansion of feed entries in the listing view
- Coordinate multiple music players so only one plays at a time
- Maintain backward compatibility with existing blog and music pages/APIs
- Support infinite scroll for feed listing

**Non-Goals:**
- Drag-and-drop block reordering (up/down buttons are sufficient)
- Rich text editing beyond markdown (no WYSIWYG)
- Block types beyond blog and music (no images, videos, embeds, etc.)
- Real-time collaboration or multi-user editing
- Version history or draft management beyond published_at flag
- Migration of existing blog/music content into feed entries

## Decisions

### 1. Database Schema: Separate FeedEntry and FeedBlock tables

**Decision**: Two tables with FK constraints, not a single table with JSON blocks.

```sql
feed_entries:
  id (uuid, PK)
  slug (text, unique, indexed)
  title (text)
  subtitle (text)
  author_id (uuid, FK → authors.id, not null)
  published_at (datetime, nullable)
  created_at (datetime)
  updated_at (datetime)

feed_blocks:
  id (uuid, PK)
  feed_entry_id (uuid, FK → feed_entries.id, indexed)
  type (text: "blog" | "music")
  order (integer)
  blog_id (uuid, FK → blogs.id, nullable, indexed)
  music_id (uuid, FK → music.id, nullable, indexed)
```

**Rationale**: 
- Relational integrity: FK constraints ensure blocks can't reference non-existent content
- Query efficiency: Can JOIN to get full block content without parsing JSON
- Flexibility: Can add new block types by adding columns (e.g., `image_id`) without schema migration
- GORM compatibility: Works naturally with GORM's relationship preloading

**Alternative considered**: Single table with JSON `blocks` array. Rejected because it loses relational integrity, makes queries harder, and complicates the upsert logic.

### 2. Upsert Semantics: Replace-all strategy for blocks

**Decision**: When upserting a feed entry, delete all existing blocks and recreate them.

**Rationale**:
- Simpler API contract: Client sends the complete desired state
- Avoids complex diffing logic (which blocks to add/update/delete)
- Matches the Notion mental model: you're composing a document, not editing individual blocks
- Transaction-safe: wrap in a single DB transaction

**Trade-off**: Loses block IDs across saves, but blocks are ephemeral (just references to Blog/Music content).

### 3. Preview vs Full Content: Server-side summary extraction

**Decision**: `GET /v2/feed` returns a `summary` field (first block content up to 1024 chars) and `cover_image` field (first markdown image URL found across all blog blocks). Blocks are not returned in the list response. `GET /feed/:slug` returns full content with all blocks.

**Rationale**:
- Reduces payload size for listing by not transferring block data
- Summary provides enough context for preview cards
- Cover image enables rich card display
- Matches existing blog API pattern (truncates to 1024 chars)
- Simple to implement in Go handler (extract before returning)

**Alternative considered**: Client-side truncation. Rejected because it transfers more data and duplicates logic across clients.

**Alternative considered**: Returning truncated blocks. Rejected because summary + cover_image is more efficient and provides better UX for preview cards.

### 4. Pagination: Cursor-based for infinite scroll

**Decision**: Use cursor-based pagination with `created_at` as the cursor. The paginated list endpoint is under the `/v2` prefix.

```
GET /v2/feed?limit=20&before=2026-07-19T00:00:00Z
```

**Rationale**:
- Stable results: New feed entries don't shift the cursor
- Efficient: Can use indexed `created_at` column
- Standard pattern for infinite scroll
- Matches the existing blog API pattern (though it uses offset, cursor is better for infinite scroll)
- `/v2` prefix follows the existing gallery API pattern for paginated endpoints

### 5. Player Coordination: Pinia store with activePlayerId

**Decision**: Centralized Pinia store tracks which player is active. Each player watches the store and pauses when another becomes active.

```typescript
// stores/playerCoordinator.ts
export const usePlayerCoordinator = defineStore('playerCoordinator', () => {
  const activePlayerId = ref<string | null>(null)
  
  const play = (id: string) => {
    activePlayerId.value = id
  }
  
  const pause = () => {
    activePlayerId.value = null
  }
  
  return { activePlayerId, play, pause }
})
```

Each music block component:
- Has a unique ID (the feed block ID)
- On play: calls `coordinator.play(myId)`
- Watches `activePlayerId`: if it changes to a different ID, pauses itself

**Rationale**:
- Decoupled: Players don't need to know about each other
- Scalable: Works for any number of players on the page
- Testable: Store logic is isolated from components
- Matches Vue 3 patterns (reactive state, watchers)

**Alternative considered**: Event bus with `emit('pause-others')`. Rejected because it creates tight coupling and doesn't scale well.

### 6. Feed Listing: Inline expansion vs navigation

**Decision**: Feed listing shows preview cards. "Continue Reading" expands the entry inline (no navigation). `/feed/:slug` route exists for direct linking / full-page view.

**Rationale**:
- Better UX: Users can browse multiple entries without losing their place
- Matches social media patterns (Twitter, Facebook)
- Reduces navigation overhead
- Still supports direct linking via `/feed/:slug`

**Implementation**: 
- Feed listing component maintains an `expandedSlugs` Set
- When "Continue Reading" is clicked, fetch full entry via `GET /feed/:slug` and render inline
- Expanded entry shows all blocks with full content

### 7. Editor: Block state management

**Decision**: Editor maintains a local array of block objects. Each block has a `mode: 'create' | 'reference'` and the corresponding data.

```typescript
interface EditorBlock {
  id: string // temporary client-side ID
  type: 'blog' | 'music'
  mode: 'create' | 'reference'
  order: number
  // For create mode:
  createData?: {
    title?: string
    content?: string
    v?: string // YouTube ID
    artist?: string
  }
  // For reference mode:
  refId?: string // blog URI or music video ID
}
```

On save:
- Transform blocks into the API payload format
- Send `PUT /feed/:slug` with the complete desired state
- Server handles creating Blog/Music rows for `create` blocks, linking to existing rows for `ref` blocks

**Rationale**:
- Simple mental model: Editor is a form that collects block data
- Matches the upsert API contract
- Easy to serialize/deserialize for editing existing entries

### 8. Component Strategy: Custom feed block components

**Decision**: 
- Created dedicated feed block components (`BlogBlockPreview`, `MusicBlockPreview`, `BlogBlockFull`, `MusicBlockFull`) rather than reusing `BlogSummary`/`BlogDetail`
- Feed listing uses `CDNImage` for cover images and `MarkdownDisplay` for summary rendering
- Feed detail shows subtitle and author at the entry level, not per-block
- Blog blocks in detail view only render title (if present) and markdown content
- Music blocks use `YouTubeEmbed` with player coordinator integration

**Rationale**:
- Feed entries have their own subtitle and author at the entry level, not per-block
- Preview cards need cover image + summary, not full block structure
- Cleaner separation: feed components handle feed-specific layout concerns
- Still reuses core primitives: `MarkdownDisplay`, `CDNImage`, `YouTubeEmbed`, `ProfileIcon`

**Trade-off**: More components to maintain, but each has a clear single responsibility.

## Risks / Trade-offs

**[Risk] Upsert replaces all blocks, losing block IDs** → Blocks are ephemeral references, not first-class entities. Losing IDs is acceptable. If we later need stable block IDs (e.g., for comments), we'd need to redesign.

**[Risk] Inline expansion in listing could be slow for large entries** → Mitigation: Only fetch full content on demand. Lazy-load music players (don't initialize YouTube iframe until user interacts).

**[Risk] Multiple music players could cause performance issues** → Mitigation: Only initialize the active player. Pause unmounts the iframe (or at least stops playback). Monitor memory usage.

**[Risk] Cursor-based pagination breaks if `created_at` is not unique** → Mitigation: Use `created_at` + `id` as composite cursor. Or accept that entries with the same timestamp may appear in arbitrary order.

**[Risk] Editor state could get out of sync with server** → Mitigation: On save, reload the entry from the server to get the canonical state. Show a "saved" indicator.

**[Trade-off] No drag-and-drop reordering** → Simpler implementation, but less polished UX. Up/down buttons are sufficient for now. Can add drag-and-drop later if needed.

**[Trade-off] Server-side summary extraction for preview** → Less flexible than client-side (can't change summary length per client), but reduces payload and matches existing pattern. Summary is first block content up to 1024 chars; cover_image is first markdown image URL.

**[Trade-off] Inline expansion vs full-page view** → Inline is better for browsing, but full-page is better for reading. We support both via `/feed` and `/feed/:slug`, but inline is the default.

## Migration Plan

**Phase 1: Database schema**
- Add `feed_entries` and `feed_blocks` tables via GORM auto-migration
- No data migration needed (existing blog/music data stays as-is)

**Phase 2: API**
- Implement feed handlers (list, get, upsert, delete)
- Register routes in `handler/handler.go`
- Test with manual API calls

**Phase 3: Frontend - Feed listing and detail**
- Create `FeedView.vue` with infinite scroll
- Create `FeedDetailView.vue` for full-page view
- Implement player coordinator store
- Add router entries

**Phase 4: Frontend - Editor**
- Create `FeedEditorView.vue` with block composer
- Implement block add/remove/reorder
- Implement save/publish actions
- Test end-to-end

**Phase 5: i18n**
- Add translation keys for feed UI strings

**Rollback**: If issues arise, disable the `/feed` routes. Existing blog/music functionality is unaffected.

## Open Questions

None at this time. All major decisions have been resolved.
