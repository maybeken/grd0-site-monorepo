## ADDED Requirements

### Requirement: FeedEntry table with required fields
The system SHALL provide a `feed_entries` table with the following columns: `id` (UUID, primary key), `slug` (text, unique, indexed), `title` (text), `subtitle` (text), `author_id` (UUID, foreign key to authors.id, not null), `published_at` (datetime, nullable), `created_at` (datetime), `updated_at` (datetime).

#### Scenario: Create feed entry with all fields
- **WHEN** a new feed entry is created with title "Summer Update", subtitle "Thoughts on the season", and slug "summer-update"
- **THEN** the system SHALL insert a row with a UUID id, the provided title, subtitle, and slug, the authenticated user's author_id, current timestamps for created_at and updated_at, and null for published_at

#### Scenario: Slug uniqueness constraint
- **WHEN** attempting to create a feed entry with a slug that already exists
- **THEN** the system SHALL reject the operation with a unique constraint violation

### Requirement: FeedBlock table with foreign key references
The system SHALL provide a `feed_blocks` table with the following columns: `id` (UUID, primary key), `feed_entry_id` (UUID, foreign key to feed_entries.id, indexed), `type` (text, either "blog" or "music"), `order` (integer), `blog_id` (UUID, foreign key to blogs.id, nullable, indexed), `music_id` (UUID, foreign key to music.id, nullable, indexed).

#### Scenario: Create blog block
- **WHEN** a feed block is created with type "blog" and a valid blog_id
- **THEN** the system SHALL insert a row with the provided feed_entry_id, type, order, and blog_id, with music_id set to null

#### Scenario: Create music block
- **WHEN** a feed block is created with type "music" and a valid music_id
- **THEN** the system SHALL insert a row with the provided feed_entry_id, type, order, and music_id, with blog_id set to null

#### Scenario: Foreign key constraint on blog_id
- **WHEN** attempting to create a feed block with a blog_id that does not exist
- **THEN** the system SHALL reject the operation with a foreign key constraint violation

#### Scenario: Foreign key constraint on music_id
- **WHEN** attempting to create a feed block with a music_id that does not exist
- **THEN** the system SHALL reject the operation with a foreign key constraint violation

#### Scenario: Hard delete of feed entry and blocks
- **WHEN** a feed entry is deleted
- **THEN** the system SHALL first delete all associated feed blocks, then delete the feed entry, within a single atomic transaction

### Requirement: GORM schema definitions
The system SHALL provide Go struct definitions for FeedEntry and FeedBlock that map to the database tables, with proper GORM tags for relationships and indexes.

#### Scenario: FeedEntry struct with relationships
- **WHEN** the FeedEntry struct is used in GORM queries
- **THEN** it SHALL include a HasMany relationship to FeedBlock, a BelongsTo relationship to Author, and UUID primary key

#### Scenario: FeedBlock struct with relationships
- **WHEN** the FeedBlock struct is used in GORM queries
- **THEN** it SHALL include BelongsTo relationships to FeedEntry, Blog, and Music, with nullable foreign keys

### Requirement: Automatic table migration
The system SHALL automatically create the feed_entries and feed_blocks tables on application startup if they do not exist, using GORM's auto-migration.

#### Scenario: First startup creates tables
- **WHEN** the API starts for the first time with no existing feed tables
- **THEN** the system SHALL create feed_entries and feed_blocks tables with all required columns and indexes

#### Scenario: Subsequent startup preserves data
- **WHEN** the API starts with existing feed tables
- **THEN** the system SHALL NOT drop or modify existing tables, preserving all data
