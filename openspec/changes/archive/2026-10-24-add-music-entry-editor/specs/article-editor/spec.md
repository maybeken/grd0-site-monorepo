# Spec Delta

## MODIFIED Requirements

### Requirement: Music block authoring

A music block SHALL be authored by selecting a song from the music catalogue or entering a YouTube link, and SHALL snapshot the video identifier and the title and artist confirmed in the picker. The link form SHALL accept a bare video identifier or a watch, `youtu.be`, `embed`, `shorts`, or `live` link, including `music.youtube.com` and `youtube-nocookie.com` hosts.

#### Scenario: Select a song
- **WHEN** the author opens the music picker
- **THEN** the picker SHALL list the songs available from the music catalogue
- **THEN** choosing a song SHALL store its video identifier, title, and artist in the block

#### Scenario: Author from a YouTube link
- **WHEN** the author pastes a YouTube link into the music picker and confirms it
- **THEN** the picker SHALL extract the video identifier from the link
- **THEN** the picker SHALL resolve the video's title and artist
- **THEN** confirming SHALL store the video identifier and the effective title and artist in the block

#### Scenario: Override resolved metadata
- **WHEN** the author changes the title or artist override after a link resolves
- **THEN** confirming SHALL store the entered values in the block
- **WHEN** an override input is left blank
- **THEN** the block SHALL store the resolved value

#### Scenario: Unsupported link
- **WHEN** the author confirms text that is not a recognised YouTube link
- **THEN** the picker SHALL report the error and leave the block unchanged

#### Scenario: Link cannot be resolved
- **WHEN** the video behind a valid YouTube link cannot be read
- **THEN** the picker SHALL report the error and leave the block unchanged

#### Scenario: Catalogue entry changes later
- **WHEN** a selected song's title changes in the music catalogue after the article is saved
- **THEN** the article SHALL continue displaying the title captured at authoring time
