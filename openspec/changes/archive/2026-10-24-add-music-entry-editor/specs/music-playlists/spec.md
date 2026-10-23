# Spec Delta

## MODIFIED Requirements

### Requirement: Add a catalogue song from the editor

The playlist editor SHALL offer an action to add a catalogue song, and SHALL reuse the article editor's music picker without its catalogue list. The picker SHALL offer title and artist overrides, and adding an already-catalogued song MUST NOT rewrite that entry.

#### Scenario: Open the add-song picker
- **WHEN** the author activates Add song
- **THEN** the music picker SHALL open with the YouTube link form only
- **THEN** the picker SHALL NOT repeat the catalogue song list

#### Scenario: Add an existing catalogue song
- **WHEN** the confirmed song already exists in the catalogue
- **THEN** the song SHALL be selected in the playlist form
- **THEN** no new catalogue entry SHALL be created
- **THEN** the existing catalogue entry SHALL remain unchanged

#### Scenario: Add a new catalogue song
- **WHEN** the author confirms a resolved link whose song is not in the catalogue
- **THEN** the editor SHALL persist the song through the music catalogue API
- **THEN** the song SHALL be selected in the playlist form

#### Scenario: Apply overrides to a new entry
- **WHEN** the author overrides the resolved title or artist and confirms a song that is not yet catalogued
- **THEN** the created catalogue entry SHALL use the overridden values
- **THEN** the song SHALL be selected in the playlist form

#### Scenario: Overrides for an existing entry
- **WHEN** the author overrides the title or artist and confirms a song that is already catalogued
- **THEN** the existing catalogue entry SHALL NOT be rewritten
- **THEN** the editor SHALL indicate that the entry is corrected from the songs page

### Requirement: Link preview before adding a song

When a YouTube link is resolved in the music picker, the picker SHALL preview it with the shared compact player and editable title and artist overrides, and MUST NOT emit the selection until the author confirms.

#### Scenario: Preview a resolved link
- **WHEN** the author submits a valid YouTube link
- **THEN** the picker SHALL resolve the video's identifier, title, and artist
- **THEN** the picker SHALL show the compact player for that video
- **THEN** the title and artist inputs SHALL be pre-filled with the resolved values
- **THEN** the selection SHALL NOT be emitted yet

#### Scenario: Confirm the preview
- **WHEN** the author confirms the preview
- **THEN** the picker SHALL emit the effective title and artist
- **THEN** the preview SHALL be cleared

#### Scenario: Override before confirming
- **WHEN** the author changes the title or artist input before confirming
- **THEN** the emitted value SHALL be the entered override
- **WHEN** an input is left blank
- **THEN** the emitted value SHALL be the resolved value

#### Scenario: Edit the link
- **WHEN** the author changes the link after a preview is shown
- **THEN** the preview and any overrides SHALL be cleared

#### Scenario: Invalid or unreadable link
- **WHEN** the link is not recognised or the video cannot be read
- **THEN** the picker SHALL report the error
- **THEN** no preview SHALL be shown and no selection emitted

## ADDED Requirements

### Requirement: Playlist save shortcut

The playlist form SHALL save the playlist when the author presses Ctrl+S or Cmd+S, using the same validation and feedback as the Save action, and MUST prevent the browser's save dialog.

#### Scenario: Keyboard save
- **WHEN** the author presses Ctrl+S or Cmd+S while the playlist form is open
- **THEN** the editor SHALL prevent the browser default and save the playlist
- **THEN** a blank title SHALL NOT be saved
