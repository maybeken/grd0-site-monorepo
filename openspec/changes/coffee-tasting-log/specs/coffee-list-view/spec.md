## ADDED Requirements

### Requirement: Coffee route
The system SHALL provide a top-level route `/coffee` that displays the public tasting log.

#### Scenario: Navigate to coffee page
- **WHEN** user navigates to `/coffee`
- **THEN** system displays the tasting log view with page title "Coffee"

### Requirement: Tasting card layout
Each tasting SHALL be displayed as a 1:1 square card with the following layout:
- Row 1: Equipment image (full width)
- Row 2 Left: Date (with pin indicator if pinned), rating, bean info, brew recipe
- Row 2 Right: Radar chart showing taste profile (7 dimensions)
- Row 3: Overall notes (full width)

#### Scenario: Card displays square aspect ratio
- **WHEN** a tasting card is rendered
- **THEN** the card maintains a 1:1 aspect ratio regardless of content

#### Scenario: Card shows equipment image
- **WHEN** a tasting has an associated equipment with type "Pour Over"
- **THEN** the card displays the local static image mapped to "Pour Over"

#### Scenario: Card shows pin indicator
- **WHEN** a tasting has `pinned = true`
- **THEN** the card displays a pin icon before the date

#### Scenario: Card shows IDK for null fields
- **WHEN** a tasting has `grind_size = null`
- **THEN** the card displays "IDK" for that field

### Requirement: Radar chart visualization
The tasting card SHALL display a radar chart showing the 7 taste dimensions: fruity, sour, sweetness, nutty, spice, floral, green. Each axis SHALL range from 0 to 10.

#### Scenario: Radar chart displays taste data
- **WHEN** a tasting has taste values (e.g., fruity=6, sweet=7, floral=8, others=lower)
- **THEN** the radar chart renders a polygon connecting the data points on each axis

#### Scenario: Radar chart handles all null values
- **WHEN** all 7 taste dimensions are null
- **THEN** the card displays a placeholder message instead of an empty chart

### Requirement: Equipment image mapping
The system SHALL maintain a constant mapping array that associates equipment types with local static image paths. The mapping SHALL be defined in code as a TypeScript constant.

#### Scenario: Equipment type has mapped image
- **WHEN** equipment type is "Pour Over" and the mapping includes an entry for "Pour Over"
- **THEN** the card displays the corresponding local image

#### Scenario: Equipment type has no mapped image
- **WHEN** equipment type has no entry in the mapping
- **THEN** the card displays a default fallback image

### Requirement: Tasting log sorting
The tasting log SHALL be sorted with pinned tastings first, then by `tasted_at` descending (newest first).

#### Scenario: Pinned tastings appear first
- **WHEN** the log contains pinned and unpinned tastings
- **THEN** pinned tastings are displayed before unpinned tastings

#### Scenario: Within same pin status, newest first
- **WHEN** two unpinned tastings have dates Jul 10 and Jul 5
- **THEN** the Jul 10 tasting appears before the Jul 5 tasting

### Requirement: Loading state
The tasting log SHALL display skeleton loading indicators while data is being fetched.

#### Scenario: Loading tastings
- **WHEN** the tasting data is being fetched
- **THEN** the view displays skeleton cards matching the card layout

### Requirement: Empty state
The tasting log SHALL display an appropriate message when no tastings exist.

#### Scenario: No tastings
- **WHEN** the tasting list is empty
- **THEN** the view displays a message indicating no tastings have been recorded
