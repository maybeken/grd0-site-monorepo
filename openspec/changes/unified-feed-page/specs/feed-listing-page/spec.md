## ADDED Requirements

### Requirement: Feed listing page route
The system SHALL provide a `/feed` route that renders the feed listing page with preview cards for each feed entry.

#### Scenario: Navigate to feed listing
- **WHEN** a user navigates to `/feed`
- **THEN** the system SHALL render the feed listing page showing published feed entries in reverse chronological order

#### Scenario: Feed listing page title
- **WHEN** the feed listing page is rendered
- **THEN** the document title SHALL be set to the translated "Feed" title via i18n

### Requirement: Infinite scroll pagination
The system SHALL implement infinite scroll on the feed listing page, automatically loading more entries as the user scrolls near the bottom.

#### Scenario: Initial load
- **WHEN** the feed listing page loads
- **THEN** the system SHALL fetch and display the first page of feed entries (up to 20)

#### Scenario: Scroll to load more
- **WHEN** the user scrolls within 200px of the bottom of the page
- **THEN** the system SHALL fetch the next page of feed entries using cursor-based pagination and append them to the list

#### Scenario: No more entries
- **WHEN** all feed entries have been loaded
- **THEN** the system SHALL stop making pagination requests and display a "No more entries" indicator

#### Scenario: Loading indicator
- **WHEN** a pagination request is in progress
- **THEN** the system SHALL display a loading spinner at the bottom of the list

### Requirement: Feed entry preview cards
The system SHALL render each feed entry as a preview card showing the entry title, subtitle, author, creation date, cover image, and summary.

#### Scenario: Preview card layout
- **WHEN** a feed entry is rendered in the listing
- **THEN** the card SHALL display the entry title, subtitle (if present), author avatar and display name, formatted creation date, cover image (if present, rendered via CDNImage), and summary (rendered via MarkdownDisplay)

#### Scenario: Cover image display
- **WHEN** a feed entry has a non-empty cover_image field
- **THEN** the preview card SHALL display the cover image using the CDNImage component with rounded styling

#### Scenario: Summary display
- **WHEN** a feed entry has a non-empty summary field
- **THEN** the preview card SHALL render the summary using the MarkdownDisplay component

### Requirement: Inline expansion on Continue Reading
The system SHALL allow users to expand a feed entry inline by clicking "Continue Reading," fetching and rendering the full content without navigating away.

#### Scenario: Expand feed entry
- **WHEN** a user clicks "Continue Reading" on a preview card
- **THEN** the system SHALL fetch the full feed entry via `GET /feed/:slug` and replace the preview with the full content inline

#### Scenario: Expanded entry shows all blocks
- **WHEN** a feed entry is expanded
- **THEN** the system SHALL render all blocks with full content (blog markdown fully rendered, music player interactive)

#### Scenario: Collapse expanded entry
- **WHEN** a user clicks a "Collapse" button on an expanded entry
- **THEN** the system SHALL revert to showing the preview card

### Requirement: Player coordination on feed listing
The system SHALL ensure that only one music player can play at a time across all expanded feed entries on the listing page.

#### Scenario: Play music in one entry
- **WHEN** a user clicks play on a music block in an expanded feed entry
- **THEN** the system SHALL start playback and pause any other playing music blocks on the page

#### Scenario: Play music in another entry
- **WHEN** a user clicks play on a music block in a different expanded feed entry while another is playing
- **THEN** the system SHALL pause the currently playing block and start the new one

### Requirement: Loading states
The system SHALL display appropriate loading states while feed entries are being fetched.

#### Scenario: Initial loading
- **WHEN** the feed listing page is loading the first page
- **THEN** the system SHALL display skeleton placeholders for feed entry cards

#### Scenario: Error loading
- **WHEN** the feed listing page fails to fetch entries
- **THEN** the system SHALL display an error message with a retry button

### Requirement: Empty state
The system SHALL display an appropriate message when no feed entries exist.

#### Scenario: No feed entries
- **WHEN** the feed listing page loads and there are no published feed entries
- **THEN** the system SHALL display a "No feed entries yet" message
