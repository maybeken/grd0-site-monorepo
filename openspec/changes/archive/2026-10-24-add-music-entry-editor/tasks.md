# Tasks

## 1. Catalogue write integrity

- [x] 1.1 Update `UpsertMusic` in `api/handler/music.go` to drop client-supplied identity and timestamps, select `title`, `artist`, and `description` explicitly on update so empty values persist, leave `sorting` untouched, and re-read the row before responding; verify `cd api && go build -o ../grd0_blog_api .` succeeds.
- [x] 1.2 Exercise the catalogue API to confirm an entry can be created, its notes and artist can be cleared, its sorting survives an update, and a deleted entry stops appearing in `GET /music` and `GET /music/playlists`.

## 2. Shared picker overrides

- [x] 2.1 Add a resolved baseline plus title and artist override inputs to `MusicPicker.vue`'s link preview, emit the effective values (`override || resolved`), and clear preview and overrides when the link changes; verify an article music block stores a typed override and stores the resolved value when the input is blank.
- [x] 2.2 Update `PlaylistForm.vue` so a new catalogue entry uses the effective title and artist, while a song already in the catalogue is selected without being rewritten and shows the "correct it on the songs page" hint; verify both cases against the catalogue after adding from a playlist.
- [x] 2.3 Add a Ctrl+S/Cmd+S save shortcut to the playlist form using the same validation and feedback as the Save action; verify the build passes and a blank title still prevents saving.

## 3. Songs editor

- [x] 3.1 Add `MusicEntryForm.vue` supporting create and edit modes with title, artist, and notes fields, a Resolve link action that opens the picker in create mode, a blank-title save guard, success feedback, and the unsaved-change guard; verify creating and editing both persist the expected values.
- [x] 3.2 Add the songs listing view with title, artist, notes preview, playlist membership counts loaded in parallel with the catalogue, an inline create form, per-row Edit and Delete actions, and empty and load-failure states; verify the listing renders and the create form returns to the refreshed list.
- [x] 3.3 Allow editing an entry inline on the songs listing by replacing its row with the edit form, keeping the other rows visible and their actions unavailable while editing; verify save and cancel both return the row to its read-only state.
- [x] 3.4 Register the `/editor/music/songs` route with a page title and remove the separate song edit route and view; verify the listing renders and its document title updates.
- [x] 3.5 Add a Songs entry under the Music sidebar group and make active highlighting choose the longest matching route, so `/editor/music/songs` marks only Songs while `/editor/music/:id` marks Playlists; verify both states in the sidebar.
- [x] 3.6 Add the new `music.editor.*`, `editor.songs`, and router title strings to `vue-web/src/locales/en.yaml`; verify no raw translation keys appear on the songs pages.
- [x] 3.7 Add a Ctrl+S/Cmd+S save shortcut to the entry form using the same validation and feedback as the Save action; verify the build passes and a blank title still prevents saving.

## 4. Deletion consistency

- [x] 4.1 Add deletion from the songs listing and from the inline edit form, with confirmation naming the title and playlist membership count, and refresh the listing after success; verify cancelling leaves the entry and failing reports an error without removing the row.
- [x] 4.2 Keep deletion in the playlist form, confirming it removes the entry from the catalogue and from the pending selection; verify the removed song disappears from every playlist and from the playlist form.

## 5. Verification

- [x] 5.1 Run `cd vue-web && pnpm type-check` and confirm it passes.
- [x] 5.2 Run `cd vue-web && pnpm build` and confirm it passes.
- [x] 5.3 Run `cd api && go build -o ../grd0_blog_api .` and confirm it passes.
- [x] 5.4 Run `openspec validate add-music-entry-editor --strict` and confirm the change is valid.
- [ ] 5.5 Walk the full flow: create an entry with overrides and notes, confirm `/music` shows the notes, edit the title and artist inline on the songs listing, clear the notes, delete the entry, then add the same video from a playlist (new) and again (existing); confirm the catalogue behaves as specified at each step.

## Workflow follow-up

- Archive the change after review and verification are complete.
