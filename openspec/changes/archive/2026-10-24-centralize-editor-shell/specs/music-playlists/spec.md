# Spec Delta

## MODIFIED Requirements

### Requirement: Playlist editor

The system SHALL provide an authenticated playlist editor inside the editor console at `/editor/music`. The listing SHALL show playlists with their title, description, and song count, and SHALL allow creating, editing, and deleting them. Creating a playlist SHALL use an inline form on the listing; editing a playlist SHALL open the playlist's edit route.

#### Scenario: Unauthenticated access
- **WHEN** an unauthenticated visitor opens the playlist editor
- **THEN** the visitor SHALL be redirected to the login flow

#### Scenario: List playlists
- **WHEN** the author opens `/editor/music`
- **THEN** every playlist SHALL be listed with its title, optional description, and song count

#### Scenario: Create a playlist
- **WHEN** the author chooses New Playlist
- **THEN** the editor SHALL show a playlist form on the listing
- **WHEN** the author saves a valid form
- **THEN** the playlist SHALL be created and the listing refreshed
- **THEN** the author SHALL remain on the listing

#### Scenario: Edit a playlist
- **WHEN** the author chooses Edit on a playlist
- **THEN** the editor SHALL navigate to `/editor/music/:id`
- **THEN** the form SHALL be pre-populated with the playlist's title, description, sorting, and selected songs

#### Scenario: Direct edit link
- **WHEN** the author opens `/editor/music/:id` directly
- **THEN** the editor SHALL load that playlist without visiting the listing first

#### Scenario: Unknown playlist
- **WHEN** the author opens an edit route for a playlist that does not exist
- **THEN** the editor SHALL show a not-found message and a link back to the listing

#### Scenario: Edit saved
- **WHEN** the author saves a valid playlist from its edit route
- **THEN** the editor SHALL report success and remain on the edit route

#### Scenario: Delete a playlist
- **WHEN** the author chooses Delete on a playlist
- **THEN** the editor SHALL ask for confirmation
- **WHEN** the author confirms
- **THEN** the playlist SHALL be deleted and the listing refreshed

#### Scenario: Title required
- **WHEN** the playlist title is blank
- **THEN** the save action SHALL be unavailable

#### Scenario: Empty listing
- **WHEN** no playlists exist
- **THEN** the editor SHALL show an empty state

#### Scenario: Failed operation
- **WHEN** loading, saving, or deleting a playlist fails
- **THEN** the editor SHALL report the failure without blocking the rest of the page

## ADDED Requirements

### Requirement: Playlist detail API

The system SHALL provide a public GET endpoint `/music/playlist/:id` that returns one playlist with its songs ordered by `sorting` descending, then creation time ascending.

#### Scenario: Read one playlist
- **WHEN** a client requests GET `/music/playlist/:id` for an existing playlist
- **THEN** the API SHALL return the playlist with its songs

#### Scenario: Unknown playlist
- **WHEN** a client requests a playlist identifier that does not exist or has been deleted
- **THEN** the API SHALL respond with not found
