# Proposal

## Why

The music page presents every song in the catalogue as a single flat "Playlist", so the site has no way to group songs or to curate more than one listening list. The music catalogue is also only editable through the API, which makes maintaining playlists impractical.

## What Changes

- Add a playlist entity that groups catalogue songs. A song may belong to more than one playlist, and playlists are ordered explicitly with a sorting value.
- Expose playlist read and write endpoints alongside the existing music catalogue endpoints. Playlist writes replace the full song membership in a transaction and never delete the songs themselves.
- Replace the flat song list on the music page with a playlist selector that displays the first playlist by default. The combined "all songs" view is removed.
- Place the selector in the right-hand playlist column, above the song list, so the player stays the focus of the left column.
- Add an authenticated playlist editor at `/music/editor` that lists playlists, creates and edits them, and deletes them with confirmation.
- Let the editor choose playlist songs from the catalogue, add a catalogue entry through the shared music picker, and remove a catalogue entry with confirmation.
- Reuse the article editor's music picker in a link-only mode for adding songs, and preview a resolved YouTube link with the shared compact player before the song is persisted.
- Keep the existing music catalogue endpoints and article music blocks unchanged.

## Capabilities

### New Capabilities

- `music-playlists`: playlist grouping and ordering, public playlist browsing, the playlist editor, and catalogue song management from the editor.

### Modified Capabilities

None.

## Impact

- `api/schema/music.go`: playlist and playlist-song models.
- `api/database/migration.go`: playlist migration.
- `api/handler/music.go`, `api/handler/handler.go`: playlist read and write endpoints.
- `vue-web/src/services/music.ts`, `vue-web/src/interfaces/Music.d.ts`: playlist types and API calls.
- `vue-web/src/views/MusicView.vue`, `vue-web/src/components/music/PlaylistSelector.vue`: public playlist browsing.
- `vue-web/src/views/MusicEditorView.vue`, `vue-web/src/components/music/PlaylistForm.vue`: playlist editor and song management.
- `vue-web/src/components/editor/article/MusicPicker.vue`: link-only mode and link preview.
- `vue-web/src/router/index.ts`, `vue-web/src/components/layout/Navbar.vue`, `vue-web/src/locales/en.yaml`: route, hidden navigation entry, and strings.
