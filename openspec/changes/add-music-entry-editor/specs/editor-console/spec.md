# Spec Delta

## MODIFIED Requirements

### Requirement: Sidebar section navigation

The sidebar SHALL list the available editor sections as Articles and Gallery under Content; Beans, Equipment, and Tastings under Coffee; Songs and Playlists under Music; and Map Locations under Travel. Each entry SHALL navigate to its section's list page, and the entry for the active section SHALL be visually distinguished.

#### Scenario: Section entries
- **WHEN** the sidebar renders
- **THEN** it SHALL offer Articles, Gallery, Beans, Equipment, Tastings, Songs, Playlists, and Map Locations entries
- **THEN** the entries SHALL be grouped under Content, Coffee, Music, and Travel headings

#### Scenario: Music group entries
- **WHEN** the sidebar renders the Music group
- **THEN** it SHALL offer Songs and Playlists entries
- **THEN** the Songs entry SHALL navigate to `/editor/music/songs`

#### Scenario: Active section
- **WHEN** the author is viewing a section's list or edit page
- **THEN** that section's sidebar entry SHALL be marked as active
- **WHEN** the author is viewing a songs route
- **THEN** the Songs entry SHALL be marked as active and the Playlists entry SHALL NOT be

#### Scenario: Navigate between sections
- **WHEN** the author activates a sidebar entry
- **THEN** the system SHALL navigate to that section's list page

### Requirement: List and edit routes

Each section SHALL expose a list route. A section SHALL either expose an edit route whose path identifies one item or allow editing that item inline on its list page. The music catalogue SHALL edit entries inline. Sections that expose an edit route SHALL identify items by URI, gallery collection directory name, slug, or identifier, and SHALL load the item without visiting the list page first.

#### Scenario: Direct edit link
- **WHEN** the author opens an edit route directly for an existing item in a section that has one
- **THEN** the editor SHALL load that item and its form
- **THEN** the list page SHALL NOT need to be visited first

#### Scenario: Inline editing
- **WHEN** a section edits entries inline
- **THEN** its list route SHALL serve as the editing surface
- **THEN** no edit route SHALL be required for that section

#### Scenario: Browser navigation
- **WHEN** the author moves between an item's edit route and other routes
- **THEN** the browser back and forward controls SHALL restore the corresponding pages

#### Scenario: Unknown item
- **WHEN** an edit route names an item that does not exist
- **THEN** the editor SHALL show a not-found message and a link back to the section's list page
