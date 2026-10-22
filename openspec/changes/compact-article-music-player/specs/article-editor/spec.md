# Spec Delta

## MODIFIED Requirements

### Requirement: Music block authoring

A music block SHALL be authored either by selecting a song from the existing music catalogue or by entering a YouTube link, and SHALL snapshot the song's video identifier, title, and artist into the block. The YouTube link form SHALL accept a bare video identifier or a link in the watch, `youtu.be`, `embed`, `shorts`, or `live` form, including `music.youtube.com` and `youtube-nocookie.com` hosts. The block's inline control SHALL render the compact player presentation, matching the article and the live preview, while keeping the change-song action available.

#### Scenario: Select a song
- **WHEN** the author opens the music picker
- **THEN** the picker SHALL list the songs available from the music catalogue
- **THEN** choosing a song SHALL store its video identifier, title, and artist in the block

#### Scenario: Author from a YouTube link
- **WHEN** the author pastes a YouTube link into the music picker and confirms it
- **THEN** the picker SHALL extract the video identifier from the link
- **THEN** the picker SHALL resolve the video's title and artist
- **THEN** confirming SHALL store the video identifier, title, and artist in the block

#### Scenario: Unsupported link
- **WHEN** the author confirms text that is not a recognised YouTube link
- **THEN** the picker SHALL report the error and leave the block unchanged

#### Scenario: Link cannot be resolved
- **WHEN** the video behind a valid YouTube link cannot be read
- **THEN** the picker SHALL report the error and leave the block unchanged

#### Scenario: Catalogue entry changes later
- **WHEN** a selected song's title changes in the music catalogue after the article is saved
- **THEN** the article SHALL continue displaying the title captured at authoring time

#### Scenario: Compact control
- **WHEN** an article music block has a selected song
- **THEN** its inline control SHALL render the compact horizontal player
- **THEN** the change-song action SHALL remain available below the control
