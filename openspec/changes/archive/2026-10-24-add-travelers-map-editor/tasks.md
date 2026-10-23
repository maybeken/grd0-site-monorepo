# Tasks

## 1. Shared map helpers and palette parity

- [x] 1.1 Extract the EPSG:4326 to EPSG:3857 conversion in both directions into `vue-web/src/helpers/map.ts`, plus zoom/resolution conversion and formatting; verify `pnpm type-check` passes and the public map still renders locations at the same places.
- [x] 1.2 Extract the colour and text colour class maps into `vue-web/src/helpers/mapPalette.ts`, add `darkblue`, and switch `TravelersMap.vue` to the shared lookups; verify the existing `sharks-and-rabbits` location now renders with a dark blue background and the other colours are unchanged.

## 2. Frontend services and shell wiring

- [x] 2.1 Add map write services (`upsertMapLocation`, `deleteMapLocation`) using `adminInstance`, a kebab-case slug generator, and a location lookup from the existing read service; verify `pnpm type-check` passes.
- [x] 2.2 Add the Travel route pair (`/editor/travel`, `/editor/travel/:slug`) inside the editor shell and the Map Locations sidebar entry under a Travel heading; verify navigating to `/editor/travel` renders the listing and the sidebar marks it active.
- [x] 2.3 Add `editor.mapLocations` and `map.editor.*` translation keys to `vue-web/src/locales/en.yaml`; verify no raw keys appear on the listing, workspace, or picker.

## 3. Listing and creation

- [x] 3.1 Build `EditorMapLocationListView.vue` listing every location with title, optional subtitle, slug, and a delete action with confirmation; verify deleting a location removes it from the public map.
- [x] 3.2 Implement the title-only inline create form, generating a kebab-case slug and rejecting a slug that already exists without writing; verify creating with a duplicate title reports the collision and leaves the existing row unchanged.
- [x] 3.3 Navigate to the new location's workspace after creation and show an empty state when no locations exist; verify the inline form closes and the route matches the generated slug.
- [x] 3.4 Report load, create, and delete failures through the shared toast without blocking the page; verify by pointing the dev server at an API that returns an error.

## 4. Location workspace

- [x] 4.1 Build `EditorMapLocationView.vue` loading the location by slug with title, subtitle, icon, colour, and text colour fields, a not-found state, and a link back to the listing; verify direct navigation works and an unknown slug shows not-found.
- [x] 4.2 Display the slug as read-only information with a link to the public map and no rename control; verify no input can change the slug.
- [x] 4.3 Add live label preview using the selected icon, colour, and text colour from the shared palette; verify the preview matches the public map label after saving.
- [x] 4.4 Wire an explicit save with snapshot-based dirty state and `useUnsavedGuard`; verify navigating away with edits prompts, cancelling keeps them, and saving reports success and stays on the route.
- [x] 4.5 Disable saving when the title is blank and identify invalid coordinate fields inline; verify a stored out-of-range longitude loads in an error state and can be corrected and saved.

## 5. Map picker and coordinates

- [x] 5.1 Build `components/map/LocationMapPicker.vue` with an OSM base layer, a vector marker and label preview at the working coordinate, `singleclick` coordinate selection, and independent centre/zoom; verify clicking updates the numeric fields and the marker follows numeric edits.
- [x] 5.2 Verify the picker uses the shared conversion so a stored position renders at the same place as on the public map, and that it neither requests the public location list nor emits navigation.
- [x] 5.3 Round-trip an existing location through load, preview, and save without touching coordinates, and confirm the stored position and label are unchanged.

## 6. Visibility thresholds

- [x] 6.1 Add "visible from zoom" and "visible until zoom" inputs converting to and from resolution, with empty meaning unset; verify existing 4000/3000 values display as approximately 5.29/5.70.
- [x] 6.2 Preserve untouched threshold resolutions on save and recompute only edited fields; verify loading a location and saving without changes produces no value drift.
- [x] 6.3 Verify the public map hides and shows a location at the expected zoom levels after setting and clearing each threshold.

## 7. Delete handler hardening

- [x] 7.1 Remove the redundant `url.PathEscape` in `DeleteMapLocation` so it matches the decoded slug parameter; verify `cd api && go build ./...` passes and an authenticated delete for a normal slug still returns 202.

## 8. Integration verification

- [x] 8.1 Run `cd vue-web && pnpm type-check` and `cd vue-web && pnpm build`; both must pass (the repo's lint script is known broken, so it is not part of verification).
- [x] 8.2 Walk the full flow against the running API: create a location, place it on the map, set a colour and visibility threshold, save, confirm the public map rendering, then delete it.
- [x] 8.3 Run `openspec validate add-travelers-map-editor --strict` and confirm it passes with no warnings.

## 9. Publishing, icon selection, and keyboard saving

- [x] 9.1 Add `published` to the API schema and migration, backfill existing rows to published, filter `GET /travel/map`, and add the authenticated `GET /travel/map/all` for the editor.
- [x] 9.2 Show publish state in the listing and workspace, create new locations as drafts, and persist publish/unpublish through the existing save flow.
- [x] 9.3 Replace the free-text icon field with a searchable selector covering the full `mynaui` set with paging.
- [x] 9.4 Serve the full icon data locally and register it with `@iconify/vue/offline` so listing, searching, and rendering make no icon CDN requests.
- [x] 9.5 Save with Ctrl+S or Cmd+S, respecting the workspace's validity checks.
- [x] 9.6 Verify the publish/unpublish flow, icon search and paging, keyboard save, `pnpm type-check`, `pnpm build`, and `go build ./...`.

## Workflow follow-up

- Archive the change after review and verification are complete.
