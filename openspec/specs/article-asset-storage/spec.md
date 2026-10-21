# article-asset-storage Specification

## Purpose

Defines storage for article-owned assets: Excalidraw scene files, uploaded markdown images, and the backend abstraction that allows local files today and object storage later.

## Requirements

### Requirement: Pluggable blob storage

Article assets SHALL be read and written through a storage abstraction that supports a local filesystem backend and an object storage backend selected by configuration. Article behaviour MUST NOT depend on which backend is active.

#### Scenario: Local backend
- **WHEN** the system starts with the local backend configured
- **THEN** assets SHALL be stored under the configured local directory
- **THEN** article endpoints SHALL read and write them normally

#### Scenario: Object storage backend
- **WHEN** the system starts with an object storage backend configured
- **THEN** article endpoints SHALL behave identically
- **THEN** the local filesystem SHALL NOT be required for asset access

#### Scenario: Backend failure
- **WHEN** the configured backend cannot store or retrieve an asset
- **THEN** the request SHALL fail with a server error and SHALL NOT record a reference to a missing asset

### Requirement: Scene file paths

An Excalidraw scene SHALL be stored as a JSON file whose path uses the article's UUID and the scene's UUID, for example `scenes/<article-uuid>/<scene-uuid>.json`. The article block SHALL reference that stored path.

#### Scenario: Scene stored under UUIDs
- **WHEN** a scene is saved for an article
- **THEN** its file path SHALL contain the article UUID and a scene UUID
- **THEN** the block SHALL reference that path

#### Scenario: URI rename leaves scenes in place
- **WHEN** an article's URI is renamed
- **THEN** existing scene files SHALL NOT move or be rewritten
- **THEN** scene references in blocks SHALL remain valid

### Requirement: Scene access endpoints

Scenes SHALL be addressed through the article URI and the scene identifier. Writing and deleting MUST require authentication; reading SHALL be public for published articles.

#### Scenario: Write a scene
- **WHEN** an authenticated author writes a scene for their article
- **THEN** the system SHALL store the JSON at the article's scene path
- **THEN** the system SHALL respond with success

#### Scenario: Read a scene
- **WHEN** a client requests a scene referenced by a published article
- **THEN** the system SHALL return the stored JSON

#### Scenario: Unauthenticated write
- **WHEN** an unauthenticated client attempts to write or delete a scene
- **THEN** the system SHALL reject the request and SHALL NOT modify storage

### Requirement: Scene content

Scene files SHALL contain the raw Excalidraw scene JSON and SHALL NOT be compressed or transformed by the storage layer.

#### Scenario: Round trip
- **WHEN** a scene is written and then read
- **THEN** the returned content SHALL be byte-equivalent JSON to what was written

### Requirement: Scene size limit

The system SHALL enforce a maximum size on scene writes and MUST reject bodies that exceed it without storing a partial file.

#### Scenario: Oversized scene
- **WHEN** a scene write exceeds the configured maximum
- **THEN** the system SHALL reject the request with an error indicating the limit
- **THEN** no partial file SHALL remain in storage

### Requirement: Attachment uploads

Markdown image uploads SHALL use a presigned upload flow scoped to the article asset prefix, separate from the legacy blog attachment prefix.

#### Scenario: Request an upload target
- **WHEN** an authenticated author requests an upload target for a filename
- **THEN** the system SHALL return a presigned URL and an object key under the article prefix

#### Scenario: Uploaded image referenced in markdown
- **WHEN** the author uploads a file and inserts the returned relative key into markdown
- **THEN** the article view SHALL resolve that relative reference through the CDN resizing helper

#### Scenario: Unauthenticated upload request
- **WHEN** an unauthenticated client requests an upload target
- **THEN** the system SHALL reject the request

### Requirement: Key validation

Asset identifiers accepted by the API SHALL be validated so that clients cannot read or write outside the article asset namespace.

#### Scenario: Valid identifiers only
- **WHEN** a request supplies an article URI or scene identifier
- **THEN** the system SHALL accept only values matching the documented formats
- **THEN** path traversal sequences and arbitrary paths SHALL be rejected

### Requirement: Asset lifecycle

Unpublishing an article SHALL retain its assets. Removing an excalidraw block SHALL delete that block's scene on a best-effort basis. Purging an article SHALL delete every asset under its UUID prefix.

#### Scenario: Unpublish retains assets
- **WHEN** an article is unpublished
- **THEN** its scene and uploaded image assets SHALL remain in storage

#### Scenario: Block removal deletes its scene
- **WHEN** an excalidraw block is removed and the article saved
- **THEN** the system SHALL attempt to delete that block's scene file
- **THEN** a deletion failure SHALL NOT prevent the article from saving

#### Scenario: Article purge
- **WHEN** an article is purged
- **THEN** every scene file under its article UUID prefix SHALL be deleted

#### Scenario: Abandoned assets
- **WHEN** an asset is no longer referenced by any article
- **THEN** the system MAY report or collect it without affecting live articles
