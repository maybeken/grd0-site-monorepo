# article Specification

## Purpose

Defines a block-based article: its URI identity, draft-by-default lifecycle, block composition, and the public/admin read and write contract.

## Requirements

### Requirement: URI identity

Every article SHALL have a URI that is unique across all articles and used as the handle for every article API call. The URI MUST match `[a-z0-9-]+` and MUST be a single URL path segment.

#### Scenario: Valid URI accepted
- **WHEN** an authenticated author creates an article with URI `brew-notes`
- **THEN** the system SHALL accept it and address the article as `/article/brew-notes`

#### Scenario: Invalid URI rejected
- **WHEN** a URI contains uppercase letters, spaces, slashes, or other characters outside `[a-z0-9-]`
- **THEN** the system SHALL reject the request with `400 Bad Request`
- **THEN** the system SHALL NOT create or modify an article

### Requirement: Article creation

The system SHALL create articles only through an explicit create request. A newly created article MUST have no blocks and MUST be unpublished. A create request that omits the URI SHALL receive a generated placeholder URI.

#### Scenario: Create an empty draft
- **WHEN** an authenticated author posts a new URI
- **THEN** the system SHALL create an article with an empty block list and a null publish timestamp
- **THEN** the system SHALL respond with `201 Created` and the created article

#### Scenario: Create without a URI
- **WHEN** an authenticated author posts a create request without a URI
- **THEN** the system SHALL create an unpublished article with an empty block list
- **THEN** the system SHALL assign a generated URI that matches `[a-z0-9-]+`
- **THEN** the response SHALL include the generated URI

#### Scenario: Blocks require an existing article
- **WHEN** a request attempts to write blocks for a URI that does not exist
- **THEN** the system SHALL respond with `404 Not Found`

### Requirement: Duplicate URI rejected

The system SHALL NOT create a second article with an existing URI, whether that article is published or unpublished. A create request for an existing URI MUST be rejected rather than returning the existing article.

#### Scenario: Create with an existing URI
- **WHEN** an authenticated author posts a URI already used by another article
- **THEN** the system SHALL respond with `409 Conflict`
- **THEN** the system SHALL NOT create a second row and SHALL NOT modify the existing article

#### Scenario: Rename onto an existing URI
- **WHEN** an author renames an article to a URI already used by another article
- **THEN** the system SHALL respond with `409 Conflict`
- **THEN** the original article SHALL keep its previous URI

### Requirement: Draft visibility

Every article SHALL be readable through the public article endpoint by URI regardless of its publish timestamp. An article with no publish timestamp is a draft; one whose timestamp is in the future is scheduled. Drafts and scheduled articles MUST NOT appear in the public article list or the sitemap.

#### Scenario: Public read of a draft
- **WHEN** an unauthenticated client requests a draft article by URI
- **THEN** the system SHALL respond with `200 OK` and return the article with its blocks

#### Scenario: Public read of a scheduled article
- **WHEN** an unauthenticated client requests an article whose publish timestamp is in the future
- **THEN** the system SHALL respond with `200 OK` and return the article with its blocks

#### Scenario: Public read of an unknown URI
- **WHEN** an unauthenticated client requests a URI that no article uses
- **THEN** the system SHALL respond with `404 Not Found`

#### Scenario: Authenticated read of a draft
- **WHEN** an authenticated client requests the draft variant of the article endpoint
- **THEN** the system SHALL return the article with its blocks

#### Scenario: Public list excludes drafts
- **WHEN** an unauthenticated client requests the article list
- **THEN** the response SHALL contain only articles whose publish timestamp is set and not in the future

#### Scenario: URL access is independent of listing
- **WHEN** an article is absent from the public list because it is a draft or scheduled
- **THEN** its URI SHALL still resolve through the public article endpoint

### Requirement: Publishing and unpublishing

The system SHALL publish an article by setting its publish timestamp and unpublish it by clearing that timestamp. Unpublishing MUST NOT delete the article, its blocks, or its stored assets.

#### Scenario: Publish an article
- **WHEN** an author sets a publish timestamp that is now or in the future
- **THEN** the article SHALL appear in the public list no earlier than that timestamp
- **THEN** the article's URI SHALL remain readable through the public article endpoint before that timestamp

#### Scenario: Unpublish an article
- **WHEN** an author unpublishes an article
- **THEN** the article SHALL stop appearing in the public list
- **THEN** the article's URI SHALL remain readable through the public article endpoint
- **THEN** the article's blocks and stored assets SHALL be retained

#### Scenario: Republish after unpublish
- **WHEN** an author republishes a previously unpublished article
- **THEN** the article SHALL reuse its existing URI

#### Scenario: Partial update retains publication
- **WHEN** an author updates an article without supplying a publish timestamp
- **THEN** the article SHALL retain its current publication state
- **WHEN** the author supplies an explicit null publish timestamp
- **THEN** the article SHALL be unpublished

### Requirement: URI rename

The system SHALL allow an authenticated author to rename an article by updating it with a different URI. After a rename, the previous URI MUST stop resolving and MUST NOT redirect anywhere.

#### Scenario: Rename succeeds
- **WHEN** an author updates an article through its current URI and supplies a new, unused URI
- **THEN** the article SHALL be addressable only by the new URI
- **THEN** the response SHALL contain the new URI

#### Scenario: Old URI after rename
- **WHEN** any client requests the previous URI after a rename
- **THEN** the system SHALL respond with `404 Not Found`
- **THEN** the system SHALL NOT issue a redirect

### Requirement: Article purge

The system SHALL provide an authenticated purge operation that permanently deletes an article and its stored scene files. The URI SHALL become available for reuse.

#### Scenario: Purge an article
- **WHEN** an authenticated author purges an article
- **THEN** the article row and every scene file under its article UUID prefix SHALL be deleted
- **THEN** the URI SHALL no longer resolve
- **THEN** the URI SHALL be available for a new article

#### Scenario: Unauthenticated purge
- **WHEN** an unauthenticated client attempts to purge an article
- **THEN** the system SHALL reject the request and SHALL NOT delete anything

### Requirement: Block composition

An article SHALL store an ordered list of blocks. Each block MUST have a stable identifier and a type of `markdown`, `image`, `music`, or `excalidraw`, with a payload appropriate to that type.

#### Scenario: Ordered blocks returned
- **WHEN** an article with multiple blocks is read
- **THEN** the blocks SHALL be returned in their stored order

#### Scenario: Unknown block type rejected
- **WHEN** an update contains a block type outside the supported set
- **THEN** the system SHALL reject the update with `400 Bad Request`
- **THEN** the stored article SHALL remain unchanged

#### Scenario: Block-level references
- **WHEN** an image block is read
- **THEN** it SHALL reference a gallery collection and filename
- **WHEN** a music block is read
- **THEN** it SHALL reference a song identifier and SHALL carry the title and artist captured at authoring time
- **WHEN** an excalidraw block is read
- **THEN** it SHALL reference the stored scene path

#### Scenario: Music block authored from a YouTube link
- **WHEN** a music block is read that was authored from a YouTube link
- **THEN** it SHALL reference the link's video identifier and SHALL carry the title and artist resolved at authoring time
- **THEN** it SHALL be indistinguishable from a music block authored from the catalogue

### Requirement: List projection

The public article list SHALL return article metadata, a summary, and a cover reference without returning block bodies or scene contents. The paginated admin list SHALL use the same projection.

#### Scenario: List payload excludes block content
- **WHEN** a client requests the public article list
- **THEN** each entry SHALL include the URI, title, subtitle, author, timestamps, summary, and cover
- **THEN** no entry SHALL include the stored block JSON or the scene file contents

#### Scenario: Admin list payload excludes block content
- **WHEN** an authenticated client requests a page of the admin article list
- **THEN** each entry SHALL use the same metadata, summary, and cover projection as the public list
- **THEN** no entry SHALL include the stored block JSON or the scene file contents

### Requirement: Admin article list pagination

The authenticated admin article list SHALL be paginated and SHALL return the shared paginated envelope containing the entries, the total entry count, the current page, the page size, and the total page count.

#### Scenario: Default page
- **WHEN** an authenticated client requests the admin list without pagination parameters
- **THEN** the response SHALL contain the first page and a bounded default page size
- **THEN** the response SHALL include the total entry count and total page count

#### Scenario: Requested page
- **WHEN** an authenticated client requests page 2 with a page size of 10
- **THEN** the response SHALL contain at most 10 entries starting after the first 10
- **THEN** the response SHALL report the page, page size, total, and total page count that describe the returned slice

#### Scenario: Stable ordering
- **WHEN** the admin list returns any page
- **THEN** entries SHALL be ordered by most recently updated first

#### Scenario: Invalid paging parameters
- **WHEN** a page parameter is zero, negative, or not a number
- **THEN** the system SHALL use the first page
- **WHEN** a page size exceeds the maximum
- **THEN** the system SHALL clamp it to the maximum

#### Scenario: Page beyond the last page
- **WHEN** a client requests a page beyond the last page
- **THEN** the response SHALL contain an empty entry list
- **THEN** the response SHALL still report the total and total page count

### Requirement: Summary derivation

The system SHALL derive the list summary from the article's first markdown block and SHALL persist it at save time.

#### Scenario: Summary from first markdown block
- **WHEN** an article with markdown blocks is saved
- **THEN** the summary SHALL be derived from the first markdown block and truncated for list display

#### Scenario: No markdown block
- **WHEN** an article has no markdown block
- **THEN** the summary SHALL be empty

### Requirement: Cover derivation

The system SHALL derive the cover from the article's first image in block order, whether that image is an image block or an inline image in a markdown block. A markdown image cover MUST reference the underlying asset path: an absolute CDN image URL SHALL have its transform prefix removed.

#### Scenario: Cover from first image block
- **WHEN** the first image in block order is an image block
- **THEN** the cover SHALL reference that block's gallery asset

#### Scenario: Cover from markdown image
- **WHEN** the first image in block order is an inline image in a markdown block
- **THEN** the cover SHALL reference that image's asset path

#### Scenario: Cover from absolute CDN markdown image
- **WHEN** an inline markdown image uses an absolute CDN URL with an image transform prefix
- **THEN** the cover SHALL reference the underlying asset path without the transform prefix
- **THEN** the stored path SHALL NOT contain double-encoded characters

#### Scenario: Cover from external markdown image
- **WHEN** an inline markdown image uses an absolute URL outside the asset CDN
- **THEN** the cover SHALL reference that absolute URL unchanged

#### Scenario: No image
- **WHEN** an article contains no image block and no markdown image
- **THEN** the cover SHALL be absent

### Requirement: Cover persistence

The system SHALL persist the cover when an article is saved and SHALL derive it from stored blocks when a list entry has no persisted cover.

#### Scenario: Missing persisted cover
- **WHEN** a list entry has no persisted cover
- **THEN** the system SHALL derive the cover from the article's stored blocks before responding

### Requirement: Author attribution

The system SHALL attribute an article to the authenticated author rather than trusting an author supplied in the request body.

#### Scenario: Author taken from session
- **WHEN** an authenticated author creates or updates an article
- **THEN** the stored author SHALL be resolved from the session's email and display name
- **THEN** any author fields in the request body SHALL be ignored
