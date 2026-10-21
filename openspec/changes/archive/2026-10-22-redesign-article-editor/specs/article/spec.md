# Spec Delta

## ADDED Requirements

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

## MODIFIED Requirements

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
