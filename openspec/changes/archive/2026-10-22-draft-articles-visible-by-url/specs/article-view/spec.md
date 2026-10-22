# Spec Delta

## MODIFIED Requirements

### Requirement: Article page

The public article page SHALL resolve an article by URI regardless of publish state and render its header and blocks in stored order.

#### Scenario: Published article resolves
- **WHEN** a reader opens a URI whose article is published
- **THEN** the page SHALL show the title, subtitle, author, and timestamps
- **THEN** the blocks SHALL be rendered in their stored order

#### Scenario: Draft article resolves
- **WHEN** a reader opens the URI of an article with no publish timestamp
- **THEN** the page SHALL render the article and its blocks instead of the not-found state

#### Scenario: Scheduled article resolves
- **WHEN** a reader opens the URI of an article whose publish timestamp is in the future
- **THEN** the page SHALL render the article and its blocks instead of the not-found state

#### Scenario: Unknown or unpublished URI
- **WHEN** a reader opens a URI that no article uses
- **THEN** the page SHALL show the not-found state
- **THEN** a URI whose article is unpublished SHALL NOT show the not-found state

## ADDED Requirements

### Requirement: Unpublished article notice

The public article page SHALL display a notice when the article is not published yet, meaning its publish timestamp is absent or in the future. The notice MUST state that the article is unpublished.

#### Scenario: Draft article notice
- **WHEN** a reader opens the URI of an article with no publish timestamp
- **THEN** the page SHALL show a notice that the article is not published yet

#### Scenario: Scheduled article notice
- **WHEN** a reader opens the URI of an article whose publish timestamp is in the future
- **THEN** the page SHALL show a notice that the article is not published yet

#### Scenario: Published article has no notice
- **WHEN** a reader opens the URI of an article whose publish timestamp is in the past
- **THEN** the page SHALL NOT show the not-published notice
