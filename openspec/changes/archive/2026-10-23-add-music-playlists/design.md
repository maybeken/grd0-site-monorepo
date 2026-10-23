# Design

## Context

The music catalogue is stored in the `musics` table and served from `/music`, ordered by `sorting` then creation time. The public music page renders a full `MusicPlayer` with previous/next controls and lists every catalogue song in a right-hand column under a "Playlist" heading. There is no music editor: catalogue entries are written through `PUT /music`.

The article editor already owns a reusable `MusicPicker` that lists catalogue songs and can resolve a pasted YouTube link with YouTube oEmbed. The shared `MusicPlayer` has a compact variant used by article music blocks. The gallery module provides the closest precedent for grouped content, but its collection membership is stored per asset path rather than as a many-to-many relation.

## Goals / Non-Goals

**Goals:**

- Group catalogue songs into named, ordered playlists.
- Keep the existing catalogue and article music blocks working unchanged.
- Default the public music page to the first playlist, with no combined view.
- Give the music module an editor that can create playlists and manage catalogue songs.
- Reuse the existing picker and players rather than introducing new media components.

**Non-Goals:**

- Changing the song shape or the article music block schema.
- Multiple tracks per article music block.
- Per-playlist song ordering distinct from the catalogue `sorting`.
- Playlist covers, descriptions rendered publicly, or playlist URLs.
- Reordering songs inside a playlist by drag and drop.

## Decisions

**Many-to-many membership.** A playlist references songs through a `playlist_songs` join table, so one song can appear in several playlists. A one-to-many column on `musics` was rejected because it would prevent reuse across playlists.

**An explicit join model.** GORM's association `Replace` was rejected after testing: the `BaseColumns.BeforeCreate` hook overwrites supplied IDs when saving associations, which made GORM create new song rows and fail the foreign key. The handler instead deletes and inserts `PlaylistSong` rows explicitly inside the same transaction as the playlist write. This also deduplicates repeated song IDs.

**Ordering.** Playlists and songs both use `sorting DESC, created_at ASC`. Per-playlist ordering was deferred; the catalogue order is already the site's curation mechanism, so a second ordering would duplicate it.

**First playlist selected.** The public page selects the first playlist returned by the API and only changes the selection when the current one no longer exists. This keeps the page deterministic after edits and removes the "all songs" option.

**Selector placement.** The selector sits at the top of the right-hand playlist column rather than above the whole page, so the full player remains the visual focus and the control sits with the list it filters.

**Reuse the article picker in link-only mode.** The playlist editor passes `link_only` to `MusicPicker`, which hides the catalogue rows because the form already lists the catalogue with checkboxes. The article editor keeps the default mode with catalogue rows. A separate link form was rejected as duplication.

**Preview before persistence.** Resolving a link sets a preview instead of emitting the selection. The picker renders the compact `MusicPlayer` and waits for confirmation. Only the confirmed selection reaches the playlist editor, which then persists a new catalogue entry through `PUT /music` and selects it. This prevents a mistyped or unwanted link from creating a catalogue entry.

**Catalogue deletion from the editor.** The remove button calls the existing `DELETE /music/:v`. Because deletion is soft, article blocks that snapshot the song's identifier, title, and artist keep rendering, while playlists stop returning the row. The confirmation copy names the catalogue as the target so the wider effect is explicit.

## Risks / Trade-offs

- [A song deleted from the catalogue disappears from every playlist] → the removal confirmation names the catalogue, and article blocks keep rendering from their snapshotted values.
- [Playlist writes replace the whole membership] → the editor always submits the complete selection, and the handler performs the replacement atomically so a failed write leaves the previous membership intact.
- [The selector is capped by the right column width on narrow screens] → the selector uses the full column width and the layout stacks the columns below the medium breakpoint.
- [The picker now has two modes] → the mode is a single boolean prop and only gates the catalogue rows; link resolution and preview behaviour are shared.
