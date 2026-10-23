# Tasks

## 1. Playlist API

- [x] 1.1 Add `Playlist` and `PlaylistSong` models and register the playlist migration — verify `go build ./...` succeeds.
- [x] 1.2 Add `GET /music/playlists` returning playlists with their songs ordered by `sorting` then creation time — verify the response nests songs and excludes soft-deleted songs.
- [x] 1.3 Add `PUT /music/playlist` creating or updating a playlist in one transaction, replacing membership and deduplicating repeated song IDs — verify a blank title is rejected and an unknown ID returns 404.
- [x] 1.4 Add `DELETE /music/playlist/:id` removing the playlist and its associations while keeping the songs — verify a repeated delete returns 404.
- [x] 1.5 Smoke-test the endpoints against a scratch SQLite database with two songs and several playlists — verify create, update, ordering, shared songs, delete, and 404 responses.

## 2. Public music page

- [x] 2.1 Add a playlist selector modelled on the gallery dropdown and place it in the right-hand playlist column above the song list.
- [x] 2.2 Select the first playlist by default and reselect the first playlist when the current selection disappears.
- [x] 2.3 Remove the combined all-songs option and filter the player and song list to the selected playlist.
- [x] 2.4 Show an empty state for a playlist with no songs and keep previous/next and automatic advancement inside the playlist.

## 3. Playlist editor

- [x] 3.1 Add the authenticated `/music/editor` route, hidden navigation entry, and translation keys.
- [x] 3.2 List playlists with title, description, and song count, with create, edit, and delete actions.
- [x] 3.3 Build the playlist form with title, description, sorting, and a catalogue song checklist; disable saving without a title.
- [x] 3.4 Confirm deletion, refresh the listing after writes, and report load, save, and delete failures through toasts.

## 4. Editor song management

- [x] 4.1 Add an Add song action that opens the shared music picker in link-only mode.
- [x] 4.2 Persist a confirmed song that is missing from the catalogue through `PUT /music`, refresh the catalogue, and select it.
- [x] 4.3 Preview a resolved link with the compact player and emit the selection only after confirmation; clear the preview when the link changes.
- [x] 4.4 Add a per-song Remove action that confirms, deletes the catalogue entry, and clears it from the pending selection.
- [x] 4.5 Add all new strings to `vue-web/src/locales/en.yaml` and remove the unused all-songs string.

## 5. Documentation and verification

- [x] 5.1 Document the playlist endpoints in `api/README.md`.
- [x] 5.2 Document the music routes and playlist editor in `vue-web/README.md`.
- [x] 5.3 Run `pnpm type-check` and `pnpm build` in `vue-web` and confirm both succeed.
- [x] 5.4 Run `gofmt`, `go vet ./...`, and `go build ./...` in `api` and confirm they succeed.

## Workflow follow-up

- Archive the change after review.
