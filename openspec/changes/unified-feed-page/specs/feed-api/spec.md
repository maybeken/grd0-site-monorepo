## ADDED Requirements

### Requirement: List feed entries with pagination
The system SHALL provide a `GET /v2/feed` endpoint that returns a paginated list of published feed entries with summary and cover image. The endpoint SHALL accept `limit` (default 20, max 100) and `before` (ISO 8601 timestamp) query parameters for cursor-based pagination. Blocks SHALL NOT be returned in the list response.

#### Scenario: List first page of feed entries
- **WHEN** a client requests `GET /v2/feed?limit=20`
- **THEN** the system SHALL return up to 20 published feed entries ordered by created_at DESC, each with a `summary` field (first block content up to 1024 characters) and `cover_image` field (first markdown image URL found across blog blocks)

#### Scenario: List next page with cursor
- **WHEN** a client requests `GET /v2/feed?limit=20&before=2026-07-19T00:00:00Z`
- **THEN** the system SHALL return up to 20 published feed entries created before the specified timestamp, ordered by created_at DESC

#### Scenario: Summary extraction for blog blocks
- **WHEN** a feed entry's first block is a blog block with content longer than 1024 characters
- **THEN** the system SHALL set the summary to the first 1024 characters of the content followed by "..."

#### Scenario: Summary extraction for music blocks
- **WHEN** a feed entry's first block is a music block
- **THEN** the system SHALL set the summary to "title - artist" (or just "title" if artist is empty)

#### Scenario: Cover image extraction
- **WHEN** a feed entry contains blog blocks with markdown images
- **THEN** the system SHALL set the cover_image to the URL of the first markdown image found across all blog blocks

#### Scenario: Only published entries returned
- **WHEN** a feed entry has published_at set to null or a future timestamp
- **THEN** the system SHALL NOT include that entry in the listing response

### Requirement: Get single feed entry with full content
The system SHALL provide a `GET /feed/:slug` endpoint that returns a single feed entry with all blocks in full content mode (no truncation).

#### Scenario: Get existing published feed entry
- **WHEN** a client requests `GET /feed/summer-update` and the entry exists and is published
- **THEN** the system SHALL return the feed entry with all blocks, including full blog content (no truncation)

#### Scenario: Get non-existent feed entry
- **WHEN** a client requests `GET /feed/non-existent` and the entry does not exist
- **THEN** the system SHALL return a 404 Not Found response

#### Scenario: Get unpublished feed entry
- **WHEN** a client requests `GET /feed/draft-post` and the entry exists but published_at is null
- **THEN** the system SHALL return a 404 Not Found response

### Requirement: Upsert feed entry with blocks
The system SHALL provide a `PUT /feed/:slug` endpoint that creates or updates a feed entry with its blocks. The endpoint SHALL require authentication. The operation SHALL be atomic: if any part fails, the entire operation SHALL be rolled back. The request body SHALL include `title`, `subtitle`, `published_at`, and `blocks`. The feed entry's author SHALL be set to the authenticated user.

#### Scenario: Create new feed entry with inline blog block
- **WHEN** an authenticated user sends `PUT /feed/new-post` with title, subtitle, and a block containing type "blog" and create data (title, content)
- **THEN** the system SHALL create a new Blog row, a new FeedEntry row with the provided slug, title, subtitle, and author_id, and a FeedBlock row linking them, returning the created entry

#### Scenario: Create new feed entry with inline music block
- **WHEN** an authenticated user sends `PUT /feed/new-song` with a block containing type "music" and create data (v, title, artist)
- **THEN** the system SHALL create a new Music row, a new FeedEntry row, and a FeedBlock row linking them, returning the created entry

#### Scenario: Create new feed entry with reference to existing blog
- **WHEN** an authenticated user sends `PUT /feed/mixed-post` with a block containing type "blog" and ref "existing-blog-uri"
- **THEN** the system SHALL create a FeedEntry and a FeedBlock with blog_id pointing to the existing Blog row with that URI, returning 404 if the blog does not exist

#### Scenario: Update existing feed entry replaces all blocks
- **WHEN** an authenticated user sends `PUT /feed/existing-post` with new blocks
- **THEN** the system SHALL update the FeedEntry fields, delete all existing FeedBlock rows for that entry, create new FeedBlock rows (and Blog/Music rows if needed), all within a single transaction

#### Scenario: Upsert with draft status
- **WHEN** an authenticated user sends `PUT /feed/draft-post` with published_at set to null
- **THEN** the system SHALL create or update the feed entry with published_at set to null, making it a draft

#### Scenario: Upsert with published status
- **WHEN** an authenticated user sends `PUT /feed/published-post` with published_at set to a timestamp
- **THEN** the system SHALL create or update the feed entry with the provided published_at timestamp

#### Scenario: Unauthorized upsert
- **WHEN** an unauthenticated client sends `PUT /feed/post`
- **THEN** the system SHALL return a 401 Unauthorized response

### Requirement: Delete feed entry
The system SHALL provide a `DELETE /feed/:slug` endpoint that hard-deletes a feed entry and all associated feed blocks. The endpoint SHALL require authentication. The operation SHALL be atomic: if any part fails, the entire operation SHALL be rolled back.

#### Scenario: Delete existing feed entry
- **WHEN** an authenticated user sends `DELETE /feed/post-to-delete`
- **THEN** the system SHALL permanently remove the feed entry and all associated feed blocks from the database

#### Scenario: Delete non-existent feed entry
- **WHEN** an authenticated user sends `DELETE /feed/non-existent`
- **THEN** the system SHALL return a 404 Not Found response

#### Scenario: Unauthorized delete
- **WHEN** an unauthenticated client sends `DELETE /feed/post`
- **THEN** the system SHALL return a 401 Unauthorized response

### Requirement: Feed entry response format
The system SHALL return feed entries in a consistent JSON format. Each entry SHALL include id, slug, title, subtitle, author (email, display_name), published_at, created_at, updated_at. The `GET /feed/:slug` endpoint SHALL include a blocks array where each block contains id, type, order, and the referenced Blog or Music object. The `GET /v2/feed` endpoint SHALL include summary and cover_image fields instead of blocks.

#### Scenario: Feed entry response structure (detail)
- **WHEN** a feed entry is returned from `GET /feed/:slug`
- **THEN** the response SHALL include id, slug, title, subtitle, author, published_at, created_at, updated_at, and a blocks array where each block contains id, type, order, and the full Blog or Music object

#### Scenario: Feed entry response structure (list)
- **WHEN** a feed entry is returned from `GET /v2/feed`
- **THEN** the response SHALL include id, slug, title, subtitle, author, published_at, created_at, updated_at, summary, cover_image, and blocks set to null

#### Scenario: Feed entry includes author
- **WHEN** a feed entry is returned from any endpoint
- **THEN** the author object SHALL include email and display_name

#### Scenario: Block order is preserved
- **WHEN** a feed entry has multiple blocks with different order values
- **THEN** the blocks array SHALL be sorted by the order field in ascending order
