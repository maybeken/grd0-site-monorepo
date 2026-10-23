# music-playlists Specification

## Purpose

Groups the music catalogue into named, ordered playlists and provides an editor for curating each playlist's songs.

## Requirements

### Requirement: Playlist grouping

The music module SHALL group catalogue songs into named playlists. A playlist MUST have a non-empty title, and songs MAY belong to more than one playlist. Playlists and their songs SHALL be ordered by `sorting` descending, then by creation time ascending.

#### Scenario: Create a playlist
- **WHEN** an authenticated author creates a playlist with a title
- **THEN** the playlist SHALL be stored and returned with its songs
- **THEN** a new playlist MAY be created with no songs

#### Scenario: Reject an empty title
- **WHEN** a playlist write has a blank title
- **THEN** the API SHALL reject the write with a bad-request response
- **THEN** no playlist SHALL be created or changed

#### Scenario: A song in several playlists
- **WHEN** a song is added to more than one playlist
- **THEN** each playlist SHALL return that song
- **THEN** editing one playlist SHALL NOT change the songs of another

#### Scenario: Ordering
- **WHEN** playlists or songs are listed
- **THEN** they SHALL be ordered by `sorting` descending, then by creation time ascending

#### Scenario: Deleting a playlist keeps its songs
- **WHEN** an author deletes a playlist
- **THEN** the playlist and its song associations SHALL be removed
- **THEN** the songs SHALL remain in the catalogue

### Requirement: Playlist persistence

Playlist reads SHALL be public and playlist writes SHALL require authentication. A playlist write SHALL replace the playlist's complete song membership inside a transaction, SHALL store a repeated song identifier only once, and SHALL leave the previous membership intact when the write fails.

#### Scenario: Read playlists
- **WHEN** a client requests the playlists
- **THEN** the API SHALL return every playlist with its songs
- **THEN** soft-deleted songs SHALL NOT be returned

#### Scenario: Replace membership
- **WHEN** an author saves a playlist with a new set of song identifiers
- **THEN** the stored membership SHALL equal the submitted set

#### Scenario: Repeated identifier
- **WHEN** a submitted membership repeats a song identifier
- **THEN** the identifier SHALL be stored once

#### Scenario: Unknown playlist
- **WHEN** an author updates or deletes a playlist that does not exist
- **THEN** the API SHALL respond with not found

#### Scenario: Failed write
- **WHEN** a playlist write fails while storing its membership
- **THEN** the playlist metadata and membership SHALL be rolled back

### Requirement: Playlist detail API

The system SHALL provide a public GET endpoint `/music/playlist/:id` that returns one playlist with its songs ordered by `sorting` descending, then creation time ascending.

#### Scenario: Read one playlist
- **WHEN** a client requests GET `/music/playlist/:id` for an existing playlist
- **THEN** the API SHALL return the playlist with its songs

#### Scenario: Unknown playlist
- **WHEN** a client requests a playlist identifier that does not exist or has been deleted
- **THEN** the API SHALL respond with not found

### Requirement: Public playlist browsing

The `/music` page SHALL present a playlist selector and SHALL display the first playlist by default. The page MUST NOT offer an all-songs option, and selecting a playlist SHALL change both the player and the song list.

#### Scenario: Default playlist
- **WHEN** a reader opens the music page
- **THEN** the first playlist returned by the API SHALL be selected
- **THEN** its songs SHALL be shown in the song list
- **THEN** the player SHALL be ready to play the first song

#### Scenario: Switch playlists
- **WHEN** the reader selects another playlist
- **THEN** the song list SHALL show that playlist's songs
- **THEN** playback SHALL reset to the first song of that playlist

#### Scenario: No combined view
- **WHEN** the reader opens the playlist selector
- **THEN** only the playlists SHALL be offered
- **THEN** no option SHALL show every catalogue song together

#### Scenario: Selector placement
- **WHEN** the music page renders on a wide screen
- **THEN** the playlist selector SHALL appear in the right-hand playlist column above the song list

#### Scenario: Empty playlist
- **WHEN** the selected playlist has no songs
- **THEN** the song list SHALL show an empty state
- **THEN** the play and skip controls SHALL be unavailable

#### Scenario: Selection no longer exists
- **WHEN** the selected playlist is no longer returned after a refresh
- **THEN** the page SHALL select the first available playlist

#### Scenario: Navigation stays inside the playlist
- **WHEN** a reader uses previous, next, or automatic advancement
- **THEN** playback SHALL stay within the selected playlist

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

### Requirement: Playlist song selection

The playlist form SHALL list the catalogue songs with a selection control and SHALL save the selected songs as the playlist's membership.

#### Scenario: Select songs
- **WHEN** the author toggles a catalogue song
- **THEN** the form SHALL add or remove that song from the pending selection
- **THEN** the displayed selected count SHALL update

#### Scenario: Save membership
- **WHEN** the author saves the playlist
- **THEN** the selected songs SHALL become the playlist's complete membership

#### Scenario: Catalogue unavailable
- **WHEN** the catalogue cannot be loaded and the author saves an existing playlist
- **THEN** the playlist's existing songs SHALL be preserved

### Requirement: Add a catalogue song from the editor

The playlist editor SHALL offer an action to add a catalogue song, and SHALL reuse the article editor's music picker without its catalogue list.

#### Scenario: Open the add-song picker
- **WHEN** the author activates Add song
- **THEN** the music picker SHALL open with the YouTube link form only
- **THEN** the picker SHALL NOT repeat the catalogue song list

#### Scenario: Add an existing catalogue song
- **WHEN** the confirmed song already exists in the catalogue
- **THEN** the song SHALL be selected in the playlist form
- **THEN** no new catalogue entry SHALL be created

#### Scenario: Add a new catalogue song
- **WHEN** the author confirms a resolved link whose song is not in the catalogue
- **THEN** the editor SHALL persist the song through the music catalogue API
- **THEN** the song SHALL be selected in the playlist form

### Requirement: Link preview before adding a song

When a YouTube link is resolved in the music picker, the picker SHALL preview it with the shared compact player and MUST NOT emit the selection until the author confirms.

#### Scenario: Preview a resolved link
- **WHEN** the author submits a valid YouTube link
- **THEN** the picker SHALL resolve the video's identifier, title, and artist
- **THEN** the picker SHALL show the compact player for that video
- **THEN** the selection SHALL NOT be emitted yet

#### Scenario: Confirm the preview
- **WHEN** the author confirms the preview
- **THEN** the picker SHALL emit the resolved song
- **THEN** the preview SHALL be cleared

#### Scenario: Edit the link
- **WHEN** the author changes the link after a preview is shown
- **THEN** the preview SHALL be cleared

#### Scenario: Invalid or unreadable link
- **WHEN** the link is not recognised or the video cannot be read
- **THEN** the picker SHALL report the error
- **THEN** no preview SHALL be shown and no selection emitted

### Requirement: Remove a catalogue song from the editor

The playlist form SHALL offer a remove action for each catalogue song. Removing SHALL require confirmation, delete the song through the catalogue API, and clear it from the pending selection.

#### Scenario: Remove a song
- **WHEN** the author chooses Remove for a song and confirms
- **THEN** the song SHALL be deleted from the catalogue
- **THEN** the song SHALL disappear from the list and from the pending selection

#### Scenario: Cancelled removal
- **WHEN** the author dismisses the confirmation
- **THEN** the catalogue and the selection SHALL remain unchanged

#### Scenario: Removed song in playlists
- **WHEN** a removed song belonged to one or more playlists
- **THEN** the playlists SHALL no longer return that song
- **THEN** article blocks that snapshot the song SHALL keep rendering

#### Scenario: Removal fails
- **WHEN** the catalogue deletion fails
- **THEN** the editor SHALL report the failure and keep the song
