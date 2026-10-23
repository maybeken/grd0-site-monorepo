# Design

## Context

See `proposal.md` for motivation. `MapLocation` stores `slug` (unique), `title`, `subtitle`, `icon` (an Iconify name from the `mynaui` set), `color`, `text_color`, `pos` (`[latitude, longitude]`), the resolution thresholds `display_at` and `hide_at`, and a `published` flag. The API exposes `GET /travel/map` (published locations), `GET /travel/map/all` (authenticated, all locations), `PUT /travel/map` (upsert keyed by slug), and `DELETE /travel/map/:slug`.

The public renderer (`TravelersMap.vue`) uses the shared coordinate conversion in `helpers/map.ts` and the shared palette in `helpers/mapPalette.ts`. Live data already contains `darkblue`, which the shared palette now supplies. Locations are HTML overlays positioned at their coordinate; the map has no marker layer. Icon rendering uses `@iconify/vue/offline` backed by locally served `mynaui` data, so neither the editor nor the public map requests icons from a CDN.

The editor shell already provides routing, unsaved-change guarding, toasts, and a section sidebar. `vue3-openlayers` v12 provides `OlVectorLayer`, `OlSourceVector`, `OlFeature`, `OlGeomPoint`, and an `OlMap` `singleclick` event carrying the projected coordinate.

## Goals / Non-Goals

**Goals:**

- A list and workspace that make every stored location visible and correctable.
- Visual coordinate selection that matches the public map's placement exactly.
- One shared palette so the editor and renderer cannot drift again.
- A draft/published state so incomplete locations can be prepared without appearing publicly.
- An icon selector that covers the full icon set and renders offline from locally served data.
- Keep the public map's rendered output and the write payload compatible.

**Non-Goals:**

- Renaming or re-slugging existing locations.
- Editing icon sets other than `mynaui`.
- Scheduled publication or per-location publish dates.
- Changing map tiles, projection, or the public page layout.

## Decisions

**1. A dedicated picker component with shared helpers**

`components/map/LocationMapPicker.vue` renders one OSM map, one marker at the working coordinate, and a label preview using the selected colour and icon. It emits a coordinate on `singleclick` (converted to EPSG:4326). Conversion and palette move out of `TravelersMap.vue` into `helpers/map.ts` and `helpers/mapPalette.ts`; both components import them.

Alternatives considered:

- *Overload `TravelersMap` with pickable/selected props*: deliberately rejected in favour of a separate editor component; the public map keeps its own data flow.
- *Copy conversion and palette into the picker*: rejected because the two would drift again, which is the bug being fixed.

**2. Slug generated on creation, immutable afterwards**

The create form accepts only a title. The editor derives a kebab-case slug (lowercase, punctuation removed, whitespace collapsed to hyphens) and sends it explicitly in the `PUT` payload, so the API's `url.PathEscape(title)` fallback is never used. Before writing, the editor compares the slug against the loaded locations and rejects a collision without writing.

Alternatives considered:

- *Let the API derive the slug*: rejected because `PathEscape` yields `%20` separators rather than kebab case.
- *Server-side slug generation*: a larger API change that would also alter existing behaviour for other callers.
- *Editable slug*: rejected because the upsert is keyed by slug, so a rename would create a second row and leave the old location behind.

**3. Preview the map while editing, save explicitly**

The picker keeps its own view state (centre and zoom) and always renders the working coordinate, so numeric entry, map clicks, and colour changes update the preview immediately. Nothing is written until Save. A snapshot of the form (including untouched threshold resolutions) drives `useUnsavedGuard`.

Alternatives considered:

- *Save on every change*: rejected; it conflicts with the editor console's established save and unsaved-guard behaviour.

**4. Thresholds are edited as zoom levels, stored as resolutions**

`display_at` and `hide_at` are metres-per-pixel thresholds. The workspace converts them with `zoom = log2(156543.03392 / resolution)` and back with `resolution = 156543.03392 / 2^zoom`, showing zoom levels to two decimals. Empty input means unset. Untouched fields are written back with their original resolution so opening and saving cannot drift stored values; only an edited field is recomputed and rounded to the nearest integer.

Semantics mapped to labels: `display_at` is "visible from zoom" (the location appears at or above that zoom), `hide_at` is "visible until zoom" (it disappears past that zoom), matching the public renderer's comparisons.

Alternatives considered:

- *Expose raw resolutions*: rejected as unreadable and easy to get wrong.
- *Round-trip through zoom for every save*: rejected because repeated rounding would slowly change stored data.

**5. Coordinate validation is explicit and non-destructive**

Latitude must be within -90 to 90 and longitude within -180 to 180. Stored out-of-range values (the existing `o-canada` longitude of `236.88`) load as-is with an inline error, and saving stays disabled until they are corrected. This makes existing bad data fixable rather than hiding it.

**6. Separate public and editor reads; reuse the write endpoint**

The public map reads `GET /travel/map`, which returns published locations only. The editor listing and workspace read the authenticated `GET /travel/map/all`, so drafts remain editable; the workspace finds its row in that response by slug, so no single-location endpoint is needed for 18 rows. Both handlers format `pos` the same way. Writes use `PUT /travel/map` with the full form payload.

`DeleteMapLocation` currently calls `url.PathEscape` on `c.Param("slug")`, which Echo has already decoded, so a slug containing characters that escape would be double-encoded and never match. The handler will use the decoded parameter directly.

Alternatives considered:

- *Return drafts from the public endpoint and filter in the client*: rejected because drafts would leak through the public API.
- *Add a per-location read endpoint*: rejected as unnecessary for the current row count; one list read serves both screens.
- *Create a location with no fields and navigate*: rejected because a row at `[0, 0]` would appear on the public map immediately; the draft state now prevents that.

**7. Palette parity through one exported map**

`helpers/mapPalette.ts` exports the background and text colour maps plus class lookup functions. The picker renders previews from it, and the editor's selects are built from its keys. `darkblue` is added as `bg-blue-900/75`; the existing entries keep their current classes so public rendering of existing rows does not change.

**8. Publishing is a boolean, and drafts are hidden from the public read**

`published` is a boolean on `MapLocation`. New locations are created as drafts so the title-only create flow cannot place a `[0, 0]` marker on the public map. `GET /travel/map` filters to published rows, while `GET /travel/map/all` returns everything for the authenticated editor. The upsert includes `published` in its conflict update. The workspace shows the state and toggles it as a pending form change, so it follows the same explicit-save and unsaved-guard rules as every other field. Existing rows are backfilled to published by an idempotent migration, preserving public visibility.

Alternatives considered:

- *Publish immediately when the toggle is clicked*: rejected because it would break the workspace's explicit-save model and the unsaved-change guard.
- *Default new locations to published*: rejected because it reintroduces the accidental public marker that drafts solve.
- *Reuse the public list in the editor*: rejected because drafts would be missing from the listing.

**9. A searchable icon selector over locally served data**

`components/editor/IconSelector.vue` loads the full `mynaui` icon data from `public/mynaui-icons.json`, registers it with `@iconify/vue/offline`, and offers a searchable grid. Search filters the full set locally; results render in pages of 120 with scroll and “Load more”, so every icon is reachable without rendering thousands at once. The trigger shows the icon name, and the rendered icon preview stays beside it. `Icon.vue` uses the offline entry, and the collection is registered before the app mounts, so no icon request leaves the site.

Alternatives considered:

- *Use the Iconify search API*: rejected because it requires an external service and keyword search adds nothing over filtering the local names.
- *Bundle the icon JSON into the main chunk*: rejected to keep the initial JavaScript smaller; the JSON is a cacheable static asset.
- *Cap results without paging*: rejected after confirming the visible list appeared incomplete.

**10. Ctrl+S and Cmd+S mirror the article editor**

A window keydown handler intercepts Ctrl/Cmd+S, prevents the browser save dialog, ignores key repeats, and calls the same save routine as the Save button. The routine's existing validity checks prevent a save when the title is blank or a coordinate is invalid.

## Risks / Trade-offs

- **Slug collision checking is client-side only** → The listing is loaded before creation, the check happens immediately before the write, and the editor is single-author. A server-side create endpoint is the follow-up if concurrent editing ever matters.
- **Zoom conversion rounding** → Untouched thresholds are written back verbatim, and edited values are rounded once; the displayed zoom is informative rather than treated as exact precision.
- **A separate picker duplicates map setup (tiles, view, overlays)** → Acceptable: the duplication is view markup, while the logic that matters (conversion, palette) is shared.
- **Adding `darkblue` changes how one existing location renders** → Intended fix; the public label for that location becomes dark blue instead of the default background.
- **Invalid stored coordinates block saving until corrected** → The error names the field, so the author can fix it immediately; no other field is lost.
- **Public map still has no marker layer** → Unchanged by design; the picker's marker is an editor affordance only.
- **New locations are drafts** → A newly created row no longer appears publicly until the author publishes it. This is intentional, but the workspace must make the state obvious.
- **The public read no longer returns drafts** → Any client that expected every stored row from `GET /travel/map` now sees only published rows; the editor uses `GET /travel/map/all`.
- **The local icon file is about 1.5 MB** → It is fetched once and registered before mount so offline icons render on first paint. The file is cacheable, but it adds to application startup compared with a small list.
- **Only the `mynaui` set is available locally** → The selector offers only those icons. Existing locations use that set; adding another set requires shipping its JSON.
- **Toggling publish is a pending change** → The public state changes only after Save or Ctrl+S, matching the rest of the workspace.

## Migration Plan

`AutoMigrate` adds a nullable `published` column and an idempotent backfill sets NULL rows to true, so existing locations keep their public visibility. The backfill uses `UpdateColumn` to leave `updated_at` untouched. New rows always write an explicit boolean, with new locations created as drafts. Deploy the API and web build together. Existing coordinates, thresholds, and colours are unchanged. The locally served icon JSON ships with the web build. Correct the known invalid longitude through the editor after deploy.

## Open Questions

None.
