# Design

## Context

The catalogue model already supports everything the editor needs: `schema.Music` carries `Title`, `Artist`, `V` (unique YouTube id), `Description`, and `Sorting`, and `api/handler/music.go` already exposes public `GET /music` plus authenticated `PUT /music` and `DELETE /music/:v`. The public `/music` page already renders `Description` as the current song's meaning. No DB migration is required.

Three constraints shape the approach:

- The editor console spec currently expects every section to have a list route and an edit route (`openspec/specs/editor-console/spec.md`). This change amends that expectation so a section may edit inline instead; the music catalogue takes that option.
- `MusicPicker.vue` is shared by `PlaylistForm.vue` and `MusicBlockEditor.vue`; changes to its contract reach both flows.
- `UpsertMusic` updates with a struct (`Updates(&music)`), so GORM skips zero values. It also binds whatever the client sends.

## Goals / Non-Goals

**Goals:**

- Give the catalogue a first-class listing and entry editor that follows console conventions.
- Make notes (the existing `description`) writable and clearable without changing the data model.
- Let the picker override resolved title and artist, with predictable, non-destructive propagation.
- Keep deletion available from both the songs page and the playlist form.

**Non-Goals:**

- No per-entry sorting UI, pagination, search, or bulk operations; the catalogue is small and listed in full.
- No new API endpoints, tables, or columns.
- No change to the public `/music` page or article rendering beyond showing authored notes.
- No move of the existing playlist routes or legacy redirect.

## Decisions

### Inline editing on the songs listing

Use `/editor/music/songs` as the only catalogue route. Creation and editing both use `MusicEntryForm`: the create form appears above the listing, and choosing Edit on a row swaps that row for the same form in edit mode. Saving or cancelling returns the row to its read-only state without navigation, and the listing reloads after a save so the membership counts stay current.

Alternatives considered: a dedicated `/editor/music/songs/:v` edit route matching the coffee and gallery sections. Rejected because a catalogue entry has only three fields and keeping the author in the list is faster; the console route requirement is amended by this change to allow either pattern. A sticky-save variant over the whole list, like the gallery overrides screen, was also rejected because per-row editing makes deletion and validation simpler.

The catalogue API is keyed by `v`, but no route identifier is needed. Duplicate detection compares the resolved `v` against the catalogue already loaded for the listing. When it matches, the create form shows a notice whose action closes the form and opens that row for inline editing.

### Save shortcuts use the same guarded save path

`MusicEntryForm` and `PlaylistForm` listen for Ctrl+S/Cmd+S while mounted, prevent the browser's save dialog, and call the same guarded save function as the Save button. This mirrors the article editor and gallery overrides screens.

Alternative considered: extracting a shared `useSaveShortcut` composable. Not extracted for two short handlers; revisit if another editor adopts the shortcut.

### One form component, picker only for link resolution

`MusicEntryForm` owns `title`, `artist`, and `notes` fields. In create mode it offers a Resolve link action that opens `MusicPicker` with `link_only`; the confirmed selection fills `v`, `title`, and `artist`, and the author can keep editing any of them before saving. In edit mode `v` is fixed display text and the picker is not offered — changing the video would mean a different entry.

Alternative considered: let the picker's confirm save the entry directly, then refine notes in the listing. Rejected because it creates empty-notes entries on cancel and splits creation across two surfaces.

### Picker overrides are effective values, applied at creation only

The picker preview keeps a resolved baseline and two override inputs prefilled from it. On confirm it emits `title: override.trim() || resolved.title` and the same for artist. This keeps the emitted shape `{ v, title, artist }` unchanged for all callers.

- Article music blocks store the emitted values, which is already their snapshot semantics.
- The playlist form creates a new catalogue entry with the emitted values when `v` is unseen.
- When `v` already exists, the playlist form selects it and does not rewrite it; the author is pointed at the songs page to correct catalogue metadata.

Alternative considered: apply overrides to existing catalogue entries from the playlist form. Rejected because it silently mutates shared catalogue data from a playlist-scoped screen, and the request was to keep both delete entry points, not both edit entry points.

### Notes are `Description`; the UI label is Notes

No field is added. The entry form labels the field Notes and persists it to `description`, which the public page already displays. The existing public copy ("What's the meaning of this?") is unchanged.

Alternative considered: rename the column or add a separate `notes` column. Rejected as a needless migration that would duplicate the same concept.

### Explicit field selection on catalogue updates

`UpsertMusic` will select `title`, `artist`, and `description` explicitly on the update path so empty strings are written. `sorting` is deliberately not selected, preserving any existing value the editor does not edit. Identity and timestamps supplied by the client are dropped before create and update, and the handler re-reads the row before responding.

Alternative considered: switch the handler to a `map` update. Equivalent outcome, but explicit `Select` keeps the struct-based handler readable and makes the editable field set obvious.

The existing soft-delete/resurrect behaviour of the upsert is kept: saving a previously deleted `v` restores it.

### Membership counts derived client-side

The songs listing loads `GET /music` and `GET /music/playlists` in parallel and counts, per `v`, how many playlists include it. No aggregate endpoint is added. If the playlist request fails while the catalogue request succeeds, the listing still renders entries without counts and shows a warning; membership is advisory, not load-blocking.

Alternative considered: add a membership count to the catalogue response. Rejected because it couples the public music read to playlist aggregation for a value only the editor uses.

### Longest-match sidebar highlighting

The current `isActive` check marks every entry whose `match` is a path prefix, which would light up both Playlists (`/editor/music`) and Songs (`/editor/music/songs`) on the songs page. Active highlighting will choose the longest matching `match` within the sidebar, so `/editor/music/songs` resolves to Songs while `/editor/music/<id>` still resolves to Playlists.

## Risks / Trade-offs

- [Zero-value updates remain broken if only the UI changes] -> The handler change lands with the editor; verify by saving an entry with notes, clearing them, and reading the catalogue again.
- [Shared picker changes article and playlist flows] -> The emitted shape is unchanged and overrides default to the resolved value; verify the article music block and playlist add-song paths after the change.
- [Extra playlist request on the songs listing] -> Load both in parallel, tolerate playlist failure, and keep counts advisory.
- [Inline editing loses edits if another row is opened] -> Row actions are unavailable while an entry is being edited, and the form carries the unsaved-change guard.
- [Sidebar prefix matching regresses with future nested routes] -> Longest-match rule is deterministic; covered by the editor-console scenario for the songs route.
- [No automated tests exist in this repo] -> Verify with `pnpm type-check` and `pnpm build` for `vue-web`, `go build` for `api`, and manual checks of clear-notes, create, delete, and both override flows.

## Migration Plan

No data migration. The API and web build deploy together; the editor degrades to today's behaviour if the API change is not deployed (notes can be added but not cleared). Rollback is a straight revert of the handler and web changes; the schema is untouched.
