# Spec Delta

## Purpose

Presents the vehicles the author has driven as a full-height, horizontally paged showroom page at `/drive`, where each vehicle is introduced by its cutout, its country and city, and its specifications.

## ADDED Requirements

### Requirement: Grand Tour route

The system SHALL serve the Grand Tour at `/drive` as a public page. While the collection is being prepared, the page SHALL remain reachable by direct URL and MUST NOT be listed in the public navigation.

#### Scenario: Direct visit
- **WHEN** a visitor opens `/drive`
- **THEN** the Grand Tour renders with the page title taken from the route metadata

#### Scenario: Unlisted page
- **WHEN** the public navigation is rendered
- **THEN** it contains no entry linking to `/drive`

### Requirement: Vehicle card composition

Each slide SHALL present one vehicle as a self-contained card containing its cover as the dominant element, its make, model and year, its country flag and country label, its city, a compact set of technical details, and its notes when present. A card MUST NOT link to a separate vehicle page.

#### Scenario: Full card
- **WHEN** a vehicle has a cover, country, city, details and notes
- **THEN** all of them are shown on the card without further navigation

#### Scenario: Sparse card
- **WHEN** a vehicle has no notes and no technical details
- **THEN** the card renders the cover, identity, flag and city without empty placeholders

### Requirement: Cover rendering

The card SHALL render the vehicle's photographic cutout as its cover. Where a cutout is absent but a silhouette asset exists, the card SHALL render the silhouette instead.

#### Scenario: Cutout present
- **WHEN** a vehicle has a cutout key
- **THEN** the card shows the coloured photographic cutout

#### Scenario: Only a silhouette
- **WHEN** a vehicle has a silhouette key but no cutout key
- **THEN** the card shows the silhouette

### Requirement: Three-dimensional carousel

The page SHALL present the vehicles in a carousel that occupies the available viewport height beneath the site header. The active vehicle SHALL sit at the front of the deck, and the vehicles adjacent to it SHALL remain visible behind and beside it. Navigation SHALL advance by one vehicle and SHALL wrap around from the last vehicle to the first. Vertical page scrolling MUST NOT be intercepted.

#### Scenario: Neighbouring vehicles remain visible
- **WHEN** the carousel renders two or more vehicles
- **THEN** the active vehicle is in front and the adjacent vehicles are visible behind it

#### Scenario: One vehicle per step
- **WHEN** the visitor advances the carousel
- **THEN** the next vehicle becomes active and the deck moves by a single position

#### Scenario: Wrap around
- **WHEN** the visitor advances past the last vehicle or returns before the first
- **THEN** the carousel continues from the opposite end

#### Scenario: Keyboard navigation
- **WHEN** no form control has focus and the visitor presses the arrow keys
- **THEN** the carousel advances or returns by one vehicle

#### Scenario: Reduced motion
- **WHEN** the visitor prefers reduced motion
- **THEN** slide transitions are not animated

### Requirement: Carousel pointer input

Pointer input on the deck SHALL drive the carousel only. Clicking a neighbouring vehicle SHALL bring it forward. The cover MUST NOT be draggable, selectable, or receive pointer input.

#### Scenario: Drag across the deck
- **WHEN** the visitor drags horizontally across a slide
- **THEN** the carousel advances or returns by one vehicle, and the cover is not dragged or selected

#### Scenario: Click a neighbour
- **WHEN** the visitor clicks a vehicle that is not active
- **THEN** that vehicle becomes active

### Requirement: Country flag rendering

The card SHALL render the vehicle's country flag from its ISO country code using assets served by the application itself. The page MUST NOT request flag imagery from an external service.

#### Scenario: Known country
- **WHEN** a vehicle carries a valid ISO country code
- **THEN** the matching flag renders beside the country label

#### Scenario: Unknown or absent country
- **WHEN** the country code is missing or unrecognised
- **THEN** the country label renders alone without a broken image

### Requirement: Missing cover placeholder

Where a vehicle has neither a cutout nor a silhouette asset, the card SHALL render an "Unknown" placeholder in place of the cover so that no slide is left empty.

#### Scenario: No cover assets
- **WHEN** a vehicle has no cutout key and no silhouette key
- **THEN** the card shows the "Unknown" placeholder instead of a cover

### Requirement: Loading and empty states

The page SHALL show placeholder slides while the vehicle collection is loading, and SHALL show a descriptive message when no vehicles exist.

#### Scenario: While loading
- **WHEN** the collection request is in flight
- **THEN** the page renders placeholder slides of the same shape as real cards

#### Scenario: No vehicles
- **WHEN** the collection resolves empty
- **THEN** the page renders a message instead of an empty carousel and does not show an error
