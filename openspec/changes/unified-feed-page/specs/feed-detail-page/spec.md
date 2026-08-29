## ADDED Requirements

### Requirement: Feed detail page route
The system SHALL provide a `/feed/:slug` route that renders a single feed entry with full content.

#### Scenario: Navigate to feed detail
- **WHEN** a user navigates to `/feed/summer-update`
- **THEN** the system SHALL fetch the feed entry with slug "summer-update" and render the detail page

#### Scenario: Feed detail page title
- **WHEN** the feed detail page is rendered
- **THEN** the document title SHALL be set to "{DEFAULT_TITLE} | {entry title}" using the feed entry's title

#### Scenario: Non-existent feed entry
- **WHEN** a user navigates to `/feed/non-existent` and the entry does not exist or is not published
- **THEN** the system SHALL render a NotFound component

### Requirement: Feed entry header rendering
The system SHALL render the feed entry's title, subtitle, and author at the top of the detail page.

#### Scenario: Render feed entry header
- **WHEN** a feed entry is rendered on the detail page
- **THEN** the system SHALL display the entry title, subtitle (if present), author avatar and display name, and creation date

### Requirement: Full content rendering
The system SHALL render all blocks in the feed entry with full content (no truncation).

#### Scenario: Render blog blocks
- **WHEN** a feed entry contains blog blocks
- **THEN** the system SHALL render each blog block with its title (if present) and full markdown content using the MarkdownDisplay component. Blog blocks SHALL NOT display subtitle, author, or timestamps (those are shown at the feed entry level).

#### Scenario: Render music blocks
- **WHEN** a feed entry contains music blocks
- **THEN** the system SHALL render each music block with an interactive YouTube player, album art, title, artist, progress bar, and transport controls

#### Scenario: Block order preserved
- **WHEN** a feed entry has multiple blocks
- **THEN** the system SHALL render them in the order specified by the order field

### Requirement: Player coordination on feed detail
The system SHALL ensure that only one music player can play at a time on the feed detail page.

#### Scenario: Play music in one block
- **WHEN** a user clicks play on a music block
- **THEN** the system SHALL start playback and pause any other playing music blocks on the page

#### Scenario: Auto-advance to next song
- **WHEN** a music block finishes playing
- **THEN** the system SHALL automatically start playing the next music block in the feed entry, if one exists

### Requirement: Loading state
The system SHALL display a loading state while the feed entry is being fetched.

#### Scenario: Loading feed entry
- **WHEN** the feed detail page is fetching the entry
- **THEN** the system SHALL display skeleton placeholders for the entry content

### Requirement: Router pre-fetch for title
The system SHALL pre-fetch the feed entry in the router beforeEach guard to extract the title for the document title suffix.

#### Scenario: Pre-fetch feed entry title
- **WHEN** a user navigates to a `/feed/:slug` route
- **THEN** the router SHALL fetch the feed entry before rendering, extract the title, and set it as meta.titleSuffix so the document title includes the entry title
