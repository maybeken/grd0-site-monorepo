# Tasks

## 1. Vehicle data model and API

- [x] 1.1 Add a `Vehicle` GORM schema with `BaseColumns`, required make/model/year, nullable technical fields, `country_code`, `country`, `city`, `driven_at`, `sort_order`, `cutout_key` and `silhouette_key`; verify `cd api && go build ./...` succeeds
- [x] 1.2 Register the schema in the database table list so startup migration creates it; verify the table exists afterwards with `sqlite3 api/data/api.db '.schema vehicles'`
- [x] 1.3 Implement the public collection endpoint returning vehicles ordered by `sort_order` then `driven_at`; verify with `curl localhost:8080/vehicles` on an empty and a populated database
- [x] 1.4 Implement public single-vehicle retrieval for editor loading, returning 404 for unknown identifiers; verify both with `curl`
- [x] 1.5 Implement authenticated upsert with validation of make, model and year, returning 400 on a missing required field and 200 on update; verify each case with `curl` against the local auth guard
- [x] 1.6 Implement authenticated delete returning success and removing the vehicle from the public listing; verify with `curl` before and after
- [x] 1.7 Add a vehicle attachment endpoint that mints a presigned PUT URL under a `vehicles/{year}/{month}/` prefix and reports uploads as unavailable when storage is unconfigured; verify the local-mode response is a clear not-implemented status
- [x] 1.8 Add the new endpoints to the Restfox collection under `docs/api/`; verify the collection imports without errors

## 2. Shared frontend data and assets

- [x] 2.1 Add the `Vehicle` TypeScript interface and a `services/vehicle.ts` module exposing public listing, admin upsert, delete and presign calls through the existing alova instances; verify `pnpm type-check` passes
- [x] 2.2 Add a country-flag icon collection JSON to `vue-web/public/` and register it with the existing icon loader alongside the `mynaui` collection; verify a flag renders while the network panel shows no external request
- [x] 2.3 Add bundled body-type silhouette assets and a mapping helper with a generic default; verify every body-type value resolves to an asset and an unknown value falls back
- [x] 2.4 Add translation keys for the navigation entry, route title, page copy, editor labels and cutout tool messages to `locales/en.yaml`; verify `pnpm build` succeeds and the new strings render

## 3. Public Grand Tour page

- [x] 3.1 Add the `/drive` route with a translated title and a new immersive layout marker; verify direct navigation renders the page and the document title updates
- [x] 3.2 Add the "Grand Tour" entry to the public navigation; verify it appears and navigates to `/drive`
- [x] 3.3 Extend the application shell so the immersive layout keeps the header, drops the container padding and hides the footer; verify `/drive` fits the viewport with no vertical page scroll
- [x] 3.4 Implement the horizontally snapping carousel with full-height slides, arrow-key navigation and reduced-motion handling; verify one slide per step, keyboard movement, and no intercepted vertical scrolling
- [x] 3.5 Implement the vehicle slide: silhouette cover, make/model/year, flag and country label, city, compact details and optional notes; verify with a fully populated vehicle and with a vehicle that has only the required fields
- [x] 3.6 Implement loading placeholders and the empty-collection message; verify by throttling the network and by running against an empty database

## 4. Vehicle editor

- [x] 4.1 Add an editor navigation group and routes for the listing and workspace; verify both routes are reachable and highlighted in the sidebar
- [x] 4.2 Implement the listing with identity, cover state, loading, retryable error and confirmed deletion; verify deletion removes the entry and an error state offers retry
- [x] 4.3 Implement the workspace and form with required-field validation, optional fields, cancel and save; verify a blocked save on missing model and a successful save with only required fields
- [x] 4.4 Implement cover preview that distinguishes a stored silhouette from the bundled fallback; verify both states by editing vehicles with and without a silhouette key
- [x] 4.5 Verify direct workspace navigation for an existing identifier loads the vehicle and an unknown identifier shows a recoverable error

## 5. In-browser cutout tool

- [x] 5.1 Add the in-browser inference runtime dependency and publish the quantized segmentation model to R2, referencing it from configuration; verify the model URL responds with permissive CORS headers and long-lived caching
- [x] 5.2 Implement a worker that loads the model once behind a provider-selection helper, preferring WebGPU and falling back to WebAssembly; verify the WebGPU path in a supporting browser and a forced WebAssembly path
- [x] 5.3 Implement preprocessing with an exact square resize and plain value scaling, and postprocessing that maps the matte back to the source frame; verify a non-square photo yields outputs matching the source aspect ratio
- [x] 5.4 Implement the photographic cutout PNG and both silhouette colour variants, filling enclosed regions and removing small disconnected fragments; verify the silhouette is a single continuous shape with windows filled
- [x] 5.5 Implement first-load download progress with caching on later runs, plus a clear unsupported-browser state; verify the progress appears once and that reopening the tool performs no download
- [x] 5.6 Integrate the tool with the vehicle form: supplying a photo generates covers, and saving uploads them through the presigned flow and stores the returned keys; verify discarding the form leaves the vehicle without cover keys
- [x] 5.7 Keep the runtime out of the public path and the PWA precache manifest; verify the Grand Tour page requests no model files and the built manifest lists none

## 6. Integration verification

- [x] 6.1 Run `pnpm type-check` and `pnpm build` in `vue-web/` and fix any failures
- [x] 6.2 Run `go build ./...` in `api/` and fix any failures
- [x] 6.3 Walk the manual checklist end to end: create a vehicle, generate covers from a photo, save, confirm the page renders the new slide with its flag, and confirm deletion removes it
- [x] 6.4 Record the tool's known limitations (side-on photos preferred, shadows may survive thresholding) in the editor's help text and verify the copy is visible where the photo is supplied

## 7. Country selection refinement

- [x] 7.1 Replace the separate country and country code inputs with a single autocompleted field that suggests names and codes from the bundled flag collection and shows the matching flag; verify suggestions filter by name and code
- [x] 7.2 Store the selected country label and upper-case ISO code together, clear both when the field is cleared, and revert unmatched text on blur; verify a saved vehicle renders its flag and label on the page
- [x] 7.3 Run `pnpm type-check`, `pnpm build` and the editor tests; verify the country field round-trips an existing vehicle

## 8. Launch gating

- [x] 8.1 Keep the Grand Tour unlisted: remove its entry from the public navigation and verify no menu link points to `/drive` while direct visits still render

## 9. Missing-cover placeholder

- [x] 9.1 Remove the bundled body-type silhouettes and their mapping, and render an "Unknown" placeholder wherever a vehicle has no silhouette asset
- [x] 9.2 Update the public card, editor preview, specs and tests to assert the placeholder, then run `pnpm type-check` and the page/editor specs

## 10. Spinning carousel

- [x] 10.1 Replace the single-slide scroller with a three-dimensional deck where the active vehicle sits in front and its neighbours are visible behind it, with wrap-around navigation, arrow keys, swipe and reduced-motion handling
- [x] 10.2 Update the page specs for the stacked deck and wrap-around, then run `pnpm type-check` and the page specs

## 11. Deployment size

- [x] 11.1 Drop the emitted ONNX Runtime wasm binary from the bundle so no deployed file exceeds the 25 MiB Pages limit, and publish the asyncify and plain variants to the asset CDN
- [x] 11.2 Load the runtime wasm from the asset CDN via configuration, choosing the variant per browser, and verify the tool runs and caches both the model and the runtime on later visits

## 12. Coloured cutout covers

- [x] 12.1 Render the photographic cutout as the card cover, fall back to the silhouette when only that exists, and prefer the cutout in the editor preview and listing cover state
- [x] 12.2 Update the specs, design and tests, then run `pnpm type-check`, `pnpm build` and the page/editor specs

## 13. Trimmed outputs

- [x] 13.1 Trim the cutout and both silhouette variants to the subject's alpha bounds with a small margin, sharing the same dimensions, and report when no subject is found
- [x] 13.2 Update the output-geometry spec/design and verify the delivered dimensions through the end-to-end flow

## 14. Manual erase

- [x] 14.1 Add a manual erase canvas to the cover tool that removes painted regions from the cutout and both silhouettes, with brush size, undo and clear, and uploads the edited blobs
- [x] 14.2 Update the cutout-tool spec/design and verify the erase and undo flow end to end

## 15. Precise erase controls

- [x] 15.1 Add zoom controls (fit to 6×) over a scrollable canvas, and draw a brush outline at the pointer scaled to the brush size
- [x] 15.2 Update the manual-erase spec/design and verify zoom, brush scaling and erasing while zoomed

## 16. Full-frame brushing

- [x] 16.1 Accept erase strokes anywhere in the canvas frame, mapping them to image coordinates and letting the image clip the paint
- [x] 16.2 Update the manual-erase spec/design and verify a stroke started outside the image erases where it crosses

## 17. Frame margins

- [x] 17.1 Reserve a scrollable margin around the erase canvas and subtract it from the fit scale so image edges can be brought inward at any zoom
- [x] 17.2 Update the erase spec/design and verify scroll room and edge erasing at high zoom

## 18. Canvas panning

- [x] 18.1 Give the zoomed canvas a working surface larger than the image and re-centre the view on every zoom change
- [x] 18.2 Pan the frame by dragging with the right mouse button, suppressing the context menu
- [x] 18.3 Update the erase spec/design and verify scrolling, right-drag panning and erasing while zoomed
- [x] 18.4 Capture Ctrl/⌘-wheel over the frame to zoom around the pointer without scrolling the page

## 19. Trim after erase

- [x] 19.1 Re-trim the edited cutout and both silhouettes to the remaining content bounds after every erase change
- [x] 19.2 Update the spec/design and verify the saved dimensions shrink after an edge erasure

## 20. Carousel pointer isolation

- [x] 20.1 Make the cover ignore pointer input and prevent dragging or selection so drags only drive the deck
- [x] 20.2 Update the carousel spec/design and verify drag navigation plus non-draggable, unselectable covers

## 21. Editor affordances

- [x] 21.1 Use a grab cursor on the carousel, support Ctrl/⌘+Z for eraser undo outside form fields, and show each vehicle's cover thumbnail in the editor listing
- [x] 21.2 Update the editor listing spec/design and verify the thumbnails, cursor and keyboard undo

## 22. Failure reporting and documentation refresh

- [x] 22.1 Surface the worker's error detail, distinguish unreadable files and empty results, and retry on WebAssembly when a WebGPU session fails
- [x] 22.2 Refresh the spec/design for shipped behaviour: failure reporting, provider retry, thumbnail placement, absolute asset origins, R2 runtime hosting, and stale migration/risk/open-question text

## Workflow follow-up

- Archive the change after review and a successful manual pass.
- Verify the archived capability specs render under `openspec/specs/`.
