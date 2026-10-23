# Proposal

## Why

Music catalogue entries can only be created or deleted as a side effect of the playlist form, and the catalogue's `description` field — already rendered publicly on `/music` as the song's meaning — has no editor at all. Every entry therefore keeps an empty note forever. Authors need a first-class place to see the catalogue, add entries, correct the title/artist resolved from YouTube, write notes, and delete entries.

## What Changes

- Add a catalogue editor at `/editor/music/songs`:
  - Listing of every catalogue entry with title, artist, notes preview, and how many playlists include it.
  - Inline create form on the listing that resolves a YouTube link through the shared music picker, then saves title, artist, and notes.
  - Inline editing on the listing: choosing Edit on a row exposes title, artist, and notes fields in place, replacing the separate song edit route; notes write the existing `description` field.
  - Delete with confirmation from both the songs page and the existing playlist form (both entry points stay).
- Add title and artist override inputs to the shared music picker's resolved-link preview. The picker emits effective values: an entered override wins, and a blank field falls back to the resolved oEmbed value. Notes stay in the entry form, not the picker.
- Update `UpsertMusic` so selected fields are written even when empty (notes and artist can be cleared) and so song `sorting` is preserved when the editor does not supply it. Client-supplied identity and timestamps are not trusted on update.
- Add Ctrl+S/Cmd+S save shortcuts to the catalogue entry form and the playlist form.
- Add a **Songs** entry beside **Playlists** in the Music sidebar group. The public `/music` page needs no change: it already renders `description`.
- Existing catalogue songs selected from a playlist are not rewritten by overrides; correcting catalogue metadata is the songs page's job.
- No schema or migration change: `title`, `artist`, `description`, and `sorting` already exist on the music table.

## Capabilities

### New Capabilities

- `music-catalogue`: Owns the catalogue entry itself — identity by YouTube video id, public read, authenticated upsert and delete, ordering, and the catalogue editor that lists, creates, edits notes/title/artist, and deletes entries.

### Modified Capabilities

- `music-playlists`: The add-song picker in the playlist editor allows title and artist overrides, adding an already-catalogued song leaves the catalogue entry unchanged, and the playlist form gains a keyboard save shortcut.
- `article-editor`: The music block picker allows title and artist overrides, and the block snapshots the values confirmed in the picker.
- `article`: A music block authored from a YouTube link carries the title and artist confirmed at authoring time, which may be author overrides rather than the raw oEmbed values.
- `editor-console`: The Music sidebar group lists Songs and Playlists, and the console route requirement permits a section to edit entries inline on its listing instead of requiring an edit route.

## Impact

- `api/handler/music.go` — `UpsertMusic` update path (explicit field selection, zero-value writes, preserve `sorting`, ignore client identity/timestamps).
- `vue-web/src/components/editor/article/MusicPicker.vue` — shared by the playlist form and the article music block editor; override fields change both flows.
- `vue-web/src/views/editor/` — new songs listing with inline editing; shared entry form component.
- `vue-web/src/router/index.ts`, `vue-web/src/components/editor/EditorSidebar.vue`, `vue-web/src/locales/en.yaml` — routes, navigation, and strings.
- `vue-web/src/services/music.ts` and `vue-web/src/interfaces/Music.d.ts` — reuse existing catalogue calls; no new endpoint expected.
- Public `/music` page and article rendering — unchanged, but now show author-written notes.
