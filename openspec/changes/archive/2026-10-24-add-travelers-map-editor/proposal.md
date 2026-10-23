# Proposal

## Why

Traveler's map locations can only be created by right-clicking the public map, copying a JSON fragment to the clipboard, and sending it to the API by hand. There is no listing, no edit form, and no way to correct a coordinate — which is how an invalid longitude such as `236.88` survives on the public map. There is also no way to prepare a location without exposing it publicly, and the icon field is free text. The editor console has no travel section, and the map's colour palette exists only inside the renderer, so a stored colour like `darkblue` silently falls back to the default.

## What Changes

- Add a Travel section to the editor console with a locations listing at `/editor/travel` and a per-location workspace at `/editor/travel/:slug`.
- Show every location with its title, subtitle, slug, and publish state, and allow deleting a location with confirmation.
- Create a location from the listing using a title-only inline form. The slug SHALL be generated in kebab case from the title and SHALL be immutable afterwards. New locations SHALL start as drafts.
- Edit a location's title, subtitle, slug display, icon, colour, text colour, publish state, visibility thresholds, and coordinates.
- Publish or unpublish a location. The public map SHALL show published locations only, and existing locations SHALL remain published.
- Choose an icon through a searchable selector covering the full icon set, with the icon list and rendered icons served locally rather than fetched from an icon CDN.
- Set coordinates either numerically or by clicking a map preview that shows the location's marker and label.
- Present visibility thresholds as zoom levels rather than raw resolution values, while storing the existing `display_at` and `hide_at` fields unchanged.
- Reject coordinates outside valid latitude and longitude ranges.
- Extract the coordinate conversion and colour palette from `TravelersMap.vue` into shared helpers, add the missing `darkblue` colour, and use the same helpers in the editor picker and the public map so stored colours always render.
- Fix the delete handler so it does not escape an already-decoded slug a second time.
- Save from the keyboard with Ctrl+S or Cmd+S, using the same validation as the Save button.

## Capabilities

### New Capabilities

- `travelers-map-editor`: Location listing, creation, editing, deletion, publish state, searchable icon selection, coordinate picking, visibility thresholds, slug lifecycle, keyboard saving, and palette parity between the editor and the public map.

### Modified Capabilities

- `editor-console`: The sidebar gains a Travel group with a Map Locations entry, and the shell's route requirement covers the locations listing and workspace.

## Impact

- **Frontend**: new `views/editor/EditorMapLocationListView.vue` and `EditorMapLocationView.vue`; new map picker and icon selector components; new shared map, palette, and icon helpers; `TravelersMap.vue` switches to the shared helpers; `Icon.vue` renders from the offline Iconify entry; router and `EditorSidebar` additions; new `map.editor.*` and `editor.mapLocations` locale strings; new map write services.
- **API**: `MapLocation` gains a `published` flag with a backfill for existing rows. `GET /travel/map` returns published locations only, and the authenticated `GET /travel/map/all` serves the editor. `DeleteMapLocation` no longer double-escapes its slug parameter. The existing `PUT /travel/map` handles writes.
- **Assets**: `vue-web/public/mynaui-icons.json` ships the full icon set so the editor can list, search, and render icons without external requests.
- **Data**: existing rows keep their stored resolution thresholds, coordinates, and colours and are marked published by the migration. New locations start as drafts.
- **Out of scope**: map tiles, projection changes, the public map's layout, editing icon sets other than `mynaui`, and scheduled publication.
