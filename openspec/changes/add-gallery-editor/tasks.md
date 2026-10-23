# Tasks

## 1. Asset directory endpoint

- [ ] 1.1 Add `GET /gallery/directories` under `restrictV1`, returning each `/gallery/*` key from `data.ReadAsset()` with its bare directory name and asset count; verify with `cd api && go build ./...` and an authenticated `curl` against the running API that the response includes `new zealand` and excludes `/blog/*` keys.
- [ ] 1.2 Verify the endpoint requires authentication by calling it without a bearer token and confirming a 401 response.

## 2. Frontend services and shell wiring

- [ ] 2.1 Add gallery write and admin read services: directory listing, collection upsert/delete, and detail upsert/delete, with a single `galleryDetailPath(collection)` helper producing the stored `/gallery/<name>` path; verify `cd vue-web && pnpm type-check` passes.
- [ ] 2.2 Add the Gallery route pair (`/editor/gallery`, `/editor/gallery/:collection`) inside the editor shell and the Gallery sidebar entry under Content; verify navigating to `/editor/gallery` in the dev server renders the listing and the sidebar marks Gallery active.
- [ ] 2.3 Add `editor.gallery` and `gallery.editor.*` translation keys to `vue-web/src/locales/en.yaml`; verify no raw keys appear on the listing or workspace and `pnpm type-check` passes.

## 3. Collections listing and adoption

- [ ] 3.1 Build `EditorGalleryListView.vue` merging collections from `/gallery/collection` with directories from the new endpoint, showing title, cover thumbnail, asset count, and separating adoptable directories; verify in the dev server that directories with metadata appear only as collections.
- [ ] 3.2 Implement adoption: create the collection row for a directory, keep the author on the listing, refresh both lists, and surface failures through the shared toast; verify by adopting an orphan directory and confirming it moves lists without touching `files.json`.
- [ ] 3.3 Show an empty state when no collections or directories exist and a not-blocking error when loading fails; verify by pointing the dev server at an API without gallery data.

## 4. Collection workspace

- [ ] 4.1 Build `EditorGalleryView.vue` loading the collection by route name, its assets, and its override rows, with a not-found state and a link back to the listing; verify direct navigation to an existing collection loads without visiting the listing, and an unknown name shows not-found.
- [ ] 4.2 Add the metadata form for title and cover, with cover selection through `GalleryPicker` restricted by a new optional `collection` prop; verify the picker offers only that collection's assets and the saved cover renders on the public gallery page.
- [ ] 4.3 Add the collection defaults row for the wildcard `*` override, explaining that it applies to images without their own row; verify saving a timezone default changes an uncaptioned image's displayed time on the public gallery.
- [ ] 4.4 Implement collection deletion with confirmation and a message that assets and overrides remain; verify the row disappears from the listing and the assets remain visible in `files.json`.

## 5. Per-image overrides

- [ ] 5.1 Render the paginated asset list with thumbnail, filename, effective description, timezone adjustment, and computed local time, falling back to collection defaults for assets with no override row; verify against a collection containing both override styles.
- [ ] 5.2 Implement the pending-change map so edits survive pagination, pre-fill new overrides from the collection defaults, and wire `useUnsavedGuard`; verify editing on page one, moving to page two, and saving persists the page-one edit, and navigating away prompts first.
- [ ] 5.3 Implement save for pending rows, including deleting a row when both description and timezone are cleared, and reporting failures without losing pending state; verify a cleared image returns to the collection default description on the public gallery.
- [ ] 5.4 List stale overrides whose filename is absent from the asset index, excluding `*`, and allow deleting them; verify by adding a detail row for a non-existent filename directly through the API and removing it in the workspace.

## 6. Integration verification

- [ ] 6.1 Run `cd vue-web && pnpm type-check` and `cd vue-web && pnpm build`; both must pass (the repo's lint script is known broken, so it is not part of verification).
- [ ] 6.2 Walk the full flow against the running API: adopt a directory, rename it, choose a cover, set defaults, override one image, clear that override, and delete the collection; confirm the public gallery reflects each step.
- [ ] 6.3 Run `openspec validate add-gallery-editor --strict` and confirm it passes with no warnings.

## Workflow follow-up

- Archive the change after review and verification are complete.
