# Proposal

## Why

The site has a home for coffee, music, travel and photography, but nothing for the vehicles the author has driven. The Grand Tour needs a showroom-style page where each vehicle is presented by an on-device cutout rather than the raw photograph.

Producing those covers with hosted image generation did not work: free tiers have no usable image models, and the edit step distorts each vehicle into a generic car. The author already has photographs of the vehicles, so the cover should be derived on the user's own machine from an uploaded photo, which also keeps the feature free, private, and offline-capable.

## What Changes

- New `Vehicle` entity with make, model, year, body type, country code and label, city, drive date, technical details, notes, rating, ordering, and cover asset keys.
- Public read API for vehicles plus authenticated create, update and delete, following the existing coffee and gallery handler pattern.
- New public page `/drive` ("Grand Tour") presenting vehicles as a full-height, horizontally snapped carousel of cards — cutout cover, flag and country, city, and vehicle details. No detail route.
- New authenticated editor section `/editor/vehicles` with a listing and a form, including cover preview and upload.
- New in-browser cutout tool that segments an uploaded photo on-device (WebGPU, WebAssembly fallback) and produces two transparent PNG outputs: a photographic cutout and a flat silhouette. No hosted inference service and no per-image cost.
- A manual erase step in the cutout tool so leftovers such as ground shadows can be removed before saving, applied to the cutout and both silhouettes.
- Cover assets upload to R2 through a presigned PUT URL, using the existing attachment pattern with a `vehicles/` prefix.
- Bundled country-flag icon collection so flags render offline and identically across platforms.
- An explicit "Unknown" placeholder so a vehicle without a generated cover still renders a complete slide.
- Route titles and translation keys; the "Grand Tour" navigation entry is deferred until the page is ready to publish, so `/drive` stays reachable by direct URL only.

## Capabilities

### New Capabilities
- `vehicle-catalogue`: the vehicle entity, its public read endpoints, authenticated mutations, and cover asset keys.
- `grand-tour-view`: the `/drive` page — full-height carousel, card composition, ordering, flag rendering, and loading and empty states.
- `vehicle-editor`: the authenticated listing and form at `/editor/vehicles`, including cover preview, upload and reassignment.
- `vehicle-cutout-tool`: on-device photo segmentation producing transparent PNG cutouts and flat silhouettes, including model loading, execution-provider fallback and upload.

### Modified Capabilities

None. Existing capabilities keep their current requirements; the new route titles and translation keys extend the shared UI without changing specified behaviour.

## Impact

- **API (`api/`)**: new `schema/vehicle.go`, `handler/vehicle.go`, routes in `handler/handler.go`, table registration in `database/`, and a vehicle attachment endpoint for presigned uploads.
- **Web (`vue-web/`)**: `views/DriveView.vue`, `components/drive/*`, `views/editor/EditorVehicleListView.vue`, `views/editor/EditorVehicleView.vue`, `components/vehicle/VehicleForm.vue`, `services/vehicle.ts`, `interfaces/Vehicle.d.ts`, router, `Navbar.vue`, `EditorSidebar.vue`, `App.vue` layout handling, `locales/en.yaml`.
- **Assets**: cutout and silhouette PNGs under the vehicle asset prefix in R2; a country-flag icon JSON in `vue-web/public/`; no bundled silhouette artwork is shipped.
- **Dependencies**: an in-browser inference runtime (`@huggingface/transformers` or `onnxruntime-web`) and a quantized segmentation model (~42 MB) served from R2. Web bundle growth is deferred until the editor tool is opened.
- **Storage**: a new `vehicles/` prefix for covers, plus `models/` for the segmentation model and `onnxruntime/` for the runtime backend in the assets bucket; no change to existing prefixes.
- **Ordering**: the feature ships in two phases — entity, API, editor and page first; the cutout tool second — so vehicles can be entered before the tool exists, showing the Unknown placeholder until a cover is generated.
