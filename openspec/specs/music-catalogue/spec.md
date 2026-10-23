# music-catalogue Specification

## Purpose

Owns the music catalogue: each song entry's video identity, metadata and notes, ordering, and the authenticated editor that lists, creates, corrects, and deletes entries.

## Requirements

### Requirement: Catalogue identity

A catalogue entry SHALL be identified by its YouTube video identifier `v`, which MUST be unique across the catalogue. An entry SHALL carry a title, an optional artist, an optional description used as the song's public note, and a sorting value. Public catalogue reads SHALL return every non-deleted entry ordered by sorting descending, then creation time ascending.

#### Scenario: Public read of the catalogue
- **WHEN** a client requests the music catalogue
- **THEN** the API SHALL return every entry that has not been deleted
- **THEN** entries SHALL be ordered by sorting descending, then creation time ascending

#### Scenario: Saving an unseen video
- **WHEN** an authenticated author saves an entry whose video identifier is not in the catalogue
- **THEN** the entry SHALL be created and returned with its identifier

#### Scenario: Saving a known video
- **WHEN** an authenticated author saves an entry whose video identifier already exists
- **THEN** the existing entry SHALL be updated
- **THEN** no second entry SHALL be created

### Requirement: Catalogue write integrity

Catalogue writes SHALL require authentication. An upsert SHALL persist title, artist, and description even when a field is empty, SHALL preserve fields the editor does not supply such as sorting, and SHALL NOT accept client-supplied identity or creation timestamps. A soft-deleted entry SHALL no longer appear in catalogue or playlist reads.

#### Scenario: Clear the notes
- **WHEN** an author clears an entry's notes and saves
- **THEN** the stored description SHALL be empty
- **THEN** later reads SHALL return the empty description

#### Scenario: Clear the artist
- **WHEN** an author clears an entry's artist and saves
- **THEN** the stored artist SHALL be empty
- **THEN** the player SHALL fall back to the video's own metadata when rendering

#### Scenario: Preserve sorting
- **WHEN** an author saves an entry without supplying a sorting value
- **THEN** the entry's stored sorting SHALL be unchanged

#### Scenario: Unauthenticated write
- **WHEN** a client without an authenticated session attempts to write an entry
- **THEN** the API SHALL reject the request
- **THEN** no entry SHALL be created or changed

#### Scenario: Deleted entry disappears
- **WHEN** an author deletes an entry
- **THEN** catalogue and playlist reads SHALL no longer return it
- **THEN** article blocks that captured the song SHALL keep rendering

### Requirement: Link resolution overrides

The shared music picker SHALL offer title and artist override inputs for a resolved YouTube link. On confirmation it SHALL emit the effective values: the entered override when present, otherwise the value resolved from the link. It MUST NOT emit a selection before the author confirms.

#### Scenario: Override resolved metadata
- **WHEN** the author changes the title or artist input after a link resolves
- **THEN** confirming SHALL emit the changed values to the caller

#### Scenario: Blank override
- **WHEN** an override input is left blank
- **THEN** the emitted value SHALL be the value resolved from the link

#### Scenario: Change the link
- **WHEN** the author edits the link after a preview is shown
- **THEN** the preview and any overrides SHALL be cleared

### Requirement: Catalogue editor listing

The editor console SHALL provide an authenticated songs listing at `/editor/music/songs`. Each row SHALL show the entry's title, artist, notes preview, and the number of playlists that include it, and SHALL offer Edit and Delete actions. The listing SHALL show an empty state when the catalogue is empty, and SHALL report load failures without blocking the rest of the page.

#### Scenario: List the catalogue
- **WHEN** the author opens the songs listing
- **THEN** every catalogue entry SHALL be listed with its title, artist, and notes preview

#### Scenario: Playlist membership
- **WHEN** an entry belongs to one or more playlists
- **THEN** its row SHALL show how many playlists include it
- **WHEN** an entry belongs to no playlist
- **THEN** its row SHALL show a zero membership indicator

#### Scenario: Empty listing
- **WHEN** the catalogue has no entries
- **THEN** the listing SHALL show an empty state

#### Scenario: Load failure
- **WHEN** the catalogue or playlist data cannot be loaded
- **THEN** the listing SHALL report the failure and keep the page usable

### Requirement: Create a catalogue entry

The songs listing SHALL own creation through an inline form. The author SHALL resolve a YouTube link with the shared music picker, MAY override the resolved title and artist, and SHALL be able to enter notes before saving. A successful save SHALL create the entry, refresh the listing, and keep the author on the listing. A blank title SHALL prevent saving.

#### Scenario: Create an entry
- **WHEN** the author resolves a link, supplies a title, and saves the form
- **THEN** the entry SHALL be created with the effective title, artist, and notes
- **THEN** the listing SHALL refresh and the author SHALL remain on it

#### Scenario: Already in the catalogue
- **WHEN** the resolved video identifier already exists
- **THEN** the editor SHALL NOT create a second entry
- **THEN** the create form SHALL stay open with a notice offering the existing entry
- **WHEN** the author activates that notice
- **THEN** the create form SHALL close and the existing entry SHALL open for inline editing

#### Scenario: Blank title
- **WHEN** the title is blank
- **THEN** the save action SHALL be unavailable

#### Scenario: Failed create
- **WHEN** saving the entry fails
- **THEN** the editor SHALL report the failure and retain the author's input

### Requirement: Edit a catalogue entry inline

The songs listing SHALL allow an entry to be opened for editing in place. The inline form SHALL expose the entry's title, artist, and notes, SHALL save without leaving the listing, and SHALL return the row to its read-only state. Cancelling SHALL discard the edits, and a second entry MUST NOT be opened while one is being edited.

#### Scenario: Open an entry for editing
- **WHEN** the author chooses Edit on a row
- **THEN** that row SHALL become an editable form with its current title, artist, and notes
- **THEN** the other rows SHALL remain visible

#### Scenario: Save inline edits
- **WHEN** the author changes the title, artist, or notes and saves
- **THEN** the entry SHALL be updated and the row SHALL return to its read-only state
- **THEN** the listing SHALL show the new values

#### Scenario: Correct resolved metadata
- **WHEN** the author replaces the artist resolved from YouTube with a corrected value and saves
- **THEN** subsequent catalogue and playlist reads SHALL use the corrected artist

#### Scenario: Cancel inline edits
- **WHEN** the author cancels editing
- **THEN** the entry SHALL remain unchanged and the row SHALL return to its read-only state

#### Scenario: One entry at a time
- **WHEN** an entry is being edited
- **THEN** the actions on the other rows SHALL be unavailable until editing ends

#### Scenario: Failed inline save
- **WHEN** saving the inline edit fails
- **THEN** the editor SHALL report the failure and keep the form open with the edits

### Requirement: Catalogue save shortcut

The catalogue entry form SHALL save the entry when the author presses Ctrl+S or Cmd+S, using the same validation and feedback as the Save action, and MUST prevent the browser's save dialog.

#### Scenario: Keyboard save
- **WHEN** the author presses Ctrl+S or Cmd+S while the entry form is open
- **THEN** the editor SHALL prevent the browser default and save the entry
- **THEN** an invalid form SHALL NOT be saved

### Requirement: Delete a catalogue entry

The songs listing and playlist editor SHALL offer deletion. Confirmation on the songs page SHALL name the entry and the number of playlists that include it; the playlist editor SHALL continue to confirm deletion from the catalogue. A cancelled or failed delete SHALL leave the entry unchanged.

#### Scenario: Delete from the songs page
- **WHEN** the author confirms deletion on the songs page
- **THEN** the entry SHALL be deleted and the listing SHALL refresh

#### Scenario: Cancel deletion
- **WHEN** the author dismisses the confirmation
- **THEN** the entry SHALL remain in the catalogue

#### Scenario: Delete failure
- **WHEN** the catalogue deletion fails
- **THEN** the editor SHALL report the failure and keep the entry

#### Scenario: Delete from the playlist editor
- **WHEN** the author confirms deletion from the playlist form
- **THEN** the entry SHALL be deleted from the catalogue
- **THEN** it SHALL disappear from the playlist's pending selection
