# Spec Delta

## MODIFIED Requirements

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
