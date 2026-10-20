# Spec Delta

## Purpose

Defines a block-based article: its URI identity, draft-by-default lifecycle, block composition, and the public/admin read and write contract.

## ADDED Requirements

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

The system SHALL create articles only through an explicit create request. A newly created article MUST have no blocks and MUST be unpublished.

#### Scenario: Create an empty draft
- **WHEN** an authenticated author posts a new URI
- **THEN** the system SHALL create an article with an empty block list and a null publish timestamp
- **THEN** the system SHALL respond with `201 Created` and the created article

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

An article with no publish timestamp SHALL be treated as a draft. Drafts MUST NOT be readable through the public article endpoint and MUST NOT appear in the public article list.

#### Scenario: Public read of a draft
- **WHEN** an unauthenticated client requests a draft article by URI
- **THEN** the system SHALL respond with `404 Not Found`

#### Scenario: Authenticated read of a draft
- **WHEN** an authenticated client requests the draft variant of the article endpoint
- **THEN** the system SHALL return the article with its blocks

#### Scenario: Public list excludes drafts
- **WHEN** an unauthenticated client requests the article list
- **THEN** the response SHALL contain only articles whose publish timestamp is set and not in the future

### Requirement: Publishing and unpublishing

The system SHALL publish an article by setting its publish timestamp and unpublish it by clearing that timestamp. Unpublishing MUST NOT delete the article, its blocks, or its stored assets.

#### Scenario: Publish an article
- **WHEN** an author sets a publish timestamp that is now or in the future
- **THEN** the article SHALL become publicly readable no earlier than that timestamp

#### Scenario: Unpublish an article
- **WHEN** an author unpublishes an article
- **THEN** the article SHALL stop appearing in the public list and the public read SHALL return `404`
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

### Requirement: List projection

The public article list SHALL return article metadata, a summary, and a cover reference without returning block bodies or scene contents.

#### Scenario: List payload excludes block content
- **WHEN** a client requests the public article list
- **THEN** each entry SHALL include the URI, title, subtitle, author, timestamps, summary, and cover
- **THEN** no entry SHALL include the stored block JSON or the scene file contents

### Requirement: Summary and cover derivation

The system SHALL derive the list summary from the article's first markdown block and the cover from its first image block, and SHALL persist both at save time.

#### Scenario: Summary from first markdown block
- **WHEN** an article with markdown blocks is saved
- **THEN** the summary SHALL be derived from the first markdown block and truncated for list display

#### Scenario: No markdown block
- **WHEN** an article has no markdown block
- **THEN** the summary SHALL be empty

#### Scenario: Cover from first image block
- **WHEN** an article contains an image block
- **THEN** the cover SHALL reference that block's gallery asset
- **WHEN** an article contains no image block
- **THEN** the cover SHALL be absent

### Requirement: Author attribution

The system SHALL attribute an article to the authenticated author rather than trusting an author supplied in the request body.

#### Scenario: Author taken from session
- **WHEN** an authenticated author creates or updates an article
- **THEN** the stored author SHALL be resolved from the session's email and display name
- **THEN** any author fields in the request body SHALL be ignored
