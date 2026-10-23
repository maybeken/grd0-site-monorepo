# Spec Delta

## MODIFIED Requirements

### Requirement: Sidebar section navigation

The sidebar SHALL list the available editor sections as Articles and Gallery under Content; Beans, Equipment, and Tastings under Coffee; Playlists under Music; and Map Locations under Travel. Each entry SHALL navigate to its section's list page, and the entry for the active section SHALL be visually distinguished.

#### Scenario: Section entries
- **WHEN** the sidebar renders
- **THEN** it SHALL offer Articles, Gallery, Beans, Equipment, Tastings, Playlists, and Map Locations entries
- **THEN** the entries SHALL be grouped under Content, Coffee, Music, and Travel headings

#### Scenario: Active section
- **WHEN** the author is viewing a section's list or edit page
- **THEN** that section's sidebar entry SHALL be marked as active

#### Scenario: Navigate between sections
- **WHEN** the author activates a sidebar entry
- **THEN** the system SHALL navigate to that section's list page

### Requirement: List and edit routes

Each section SHALL expose a list route and an edit route whose path identifies one item. The Articles section SHALL identify items by URI. Gallery items SHALL be identified by their collection directory name, map locations SHALL be identified by slug, and the remaining sections SHALL identify items by identifier. An edit route SHALL load its item without visiting the list page first.

#### Scenario: Direct edit link
- **WHEN** the author opens an edit route directly for an existing item
- **THEN** the editor SHALL load that item and its form
- **THEN** the list page SHALL NOT need to be visited first

#### Scenario: Browser navigation
- **WHEN** the author moves between an item's edit route and other routes
- **THEN** the browser back and forward controls SHALL restore the corresponding pages

#### Scenario: Unknown item
- **WHEN** an edit route names an item that does not exist
- **THEN** the editor SHALL show a not-found message and a link back to the section's list page
