# editor-console Specification

## Purpose

Provides the authenticated full-screen shell that hosts every content editor, with sidebar navigation between editor sections and addressable routes for individual items.

## Requirements

### Requirement: Editor shell layout

The system SHALL provide a full-screen editor shell at `/editor` that renders a persistent left sidebar beside the active editor page. While the shell is active, the public navbar, site footer, and particle background SHALL NOT be rendered.

#### Scenario: Opening the shell
- **WHEN** an authenticated author opens any route under `/editor`
- **THEN** the shell SHALL render the sidebar and the active editor page
- **THEN** the public navbar and footer SHALL be absent

#### Scenario: Leaving the shell
- **WHEN** the author navigates from an `/editor` route to a public route
- **THEN** the public navbar, footer, and background SHALL render again

#### Scenario: Return to the public site
- **WHEN** the author activates the sidebar's site link
- **THEN** the system SHALL navigate to the public home page

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

### Requirement: Editor entry and section defaults

The shell root SHALL redirect to the Articles listing, and the coffee root SHALL redirect to the Beans listing, so every section has one canonical entry page.

#### Scenario: Shell root
- **WHEN** the author opens `/editor`
- **THEN** the system SHALL redirect to the Articles listing

#### Scenario: Coffee root
- **WHEN** the author opens `/editor/coffee`
- **THEN** the system SHALL redirect to the Beans listing

### Requirement: List and edit routes

Each section SHALL expose a list route and an edit route whose path identifies one item. The Articles section SHALL identify items by URI. Gallery items SHALL be identified by their collection directory name; map locations SHALL be identified by slug, and the remaining sections SHALL identify items by identifier. An edit route SHALL load its item without visiting the list page first.

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

### Requirement: Creation from the list page

Each section's list page SHALL own item creation through an inline create form. A create action MUST NOT navigate to an edit route. After a successful create, the list SHALL refresh and the form SHALL close; article creation SHALL instead create a draft and open that article's edit route.

#### Scenario: Create a section item
- **WHEN** the author saves a valid create form on a list page
- **THEN** the item SHALL be created and the listing SHALL refresh
- **THEN** the author SHALL remain on the list page

#### Scenario: Create an article
- **WHEN** the author confirms a new article from the Articles listing
- **THEN** the system SHALL create a draft article
- **THEN** the system SHALL open that article's edit route

#### Scenario: Edit never creates
- **WHEN** an edit route is opened for an unknown identifier
- **THEN** the editor SHALL NOT create a record

### Requirement: Centralized editor authentication

The system SHALL require an authenticated session for every route under `/editor` and SHALL redirect unauthenticated visitors to the login flow. The local development environment SHALL bypass this requirement.

#### Scenario: Unauthenticated access
- **WHEN** a visitor without an application token opens any `/editor` route
- **THEN** the system SHALL redirect the visitor to the login flow

#### Scenario: Authenticated access
- **WHEN** a visitor with an application token opens an `/editor` route
- **THEN** the shell SHALL render without a login redirect

#### Scenario: Local development
- **WHEN** the editor is opened in the local development environment
- **THEN** the shell SHALL render without requiring a token

### Requirement: Unsaved-change protection

The system SHALL warn the author before discarding unsaved edits on any editor page. In-application navigation, including sidebar entries, SHALL require confirmation. Reloading or closing the tab SHALL trigger the browser's confirmation prompt. Cancelling SHALL keep the author on the current page with their edits intact.

#### Scenario: Navigate away with edits
- **WHEN** the author has unsaved edits and activates any in-application navigation
- **THEN** the system SHALL ask for confirmation
- **THEN** cancelling SHALL keep the current page and the edits

#### Scenario: Reload or close
- **WHEN** the author has unsaved edits and reloads or closes the page
- **THEN** the browser SHALL prompt for confirmation

#### Scenario: No edits
- **WHEN** the author has no unsaved edits
- **THEN** in-application navigation SHALL proceed without a prompt

### Requirement: Responsive sidebar

Below the large breakpoint the sidebar SHALL collapse so the editor remains usable, and SHALL be reopenable from a control in the editor page.

#### Scenario: Narrow viewport
- **WHEN** the shell renders on a narrow viewport
- **THEN** the sidebar SHALL be collapsed by default
- **THEN** a control SHALL be available to open it

#### Scenario: Wide viewport
- **WHEN** the shell renders on a wide viewport
- **THEN** the sidebar SHALL be visible beside the editor page

### Requirement: Legacy editor redirects

The previous editor paths SHALL redirect to their shell equivalents so existing links continue to work.

#### Scenario: Legacy article editor paths
- **WHEN** a visitor opens `/article/editor` or `/article/editor/:uri`
- **THEN** the system SHALL redirect to the corresponding Articles list or edit route

#### Scenario: Legacy coffee editor path
- **WHEN** a visitor opens `/coffee/editor`
- **THEN** the system SHALL redirect to the Beans listing

#### Scenario: Legacy music editor path
- **WHEN** a visitor opens `/music/editor`
- **THEN** the system SHALL redirect to the Playlists listing
