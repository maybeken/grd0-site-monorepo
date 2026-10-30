# Design

## Context

See `proposal.md` for motivation and `specs/` for requirements.

The project is a pnpm monorepo: a Vue 3 single-page app (`vue-web/`) built with Vite and Tailwind v4, and a Go service (`api/`) using Echo, GORM over SQLite, and the MinIO client for R2 storage. The frontend talks to the API through alova instances: `dataInstance` for public reads with a five-minute cache, `adminInstance` for authenticated writes.

Two existing patterns carry most of this work:

- The coffee feature is the shape to follow: a GORM entity with `BaseColumns`, an Echo handler with public reads and authenticated mutations, a service module, a public view, and an editor list plus form.
- Attachments already upload through presigned PUT URLs (`/blog/attachment/:key`, `/article/attachment/:key`), which return both the URL and the object key to store.

The application shell is `App.vue`: a sticky header, a `main` container, and a footer, with one escape hatch — routes whose metadata sets `layout: 'editor'` render bare; the immersive layout is a second value on the same switch. The lint script is broken, so verification is `pnpm type-check`, `pnpm build`, and the Playwright specs under `vue-web/tests/`.

Spike results that constrain the design:

- A quantized `ormbg` segmentation model (Apache-2.0, 42 MB) produced a clean vehicle matte from a cluttered photo in ~680 ms on CPU at 1024².
- The model declares `do_normalize: false`; input is plain `pixel / 255`. Applying the conventional mean-0.5 normalisation produces garbage.
- Its preprocessor resizes to an exact 1024×1024 square. Letterboxing to preserve aspect was tested with black, gray and white padding and made masks worse in every case, because padding turns background detail into salience. Distort-then-map-back is correct.
- Silhouettes derive cleanly from the matte: threshold, flood-fill from the corners to fill enclosed window holes, drop small disconnected specks. Colour keying could never do this because it had no notion of an object.
- WebGPU is supported by Chrome/Edge 113+, Safari 26, Firefox 141+ on Windows and 145+ on macOS, and Chrome for Android 121+. WebAssembly with threads is a usable fallback at roughly 2 s per image.

## Goals / Non-Goals

**Goals:**

- Ship the public page, the data model, the API and the editor independently of the cutout tool.
- Keep all cover generation on the user's device and free of recurring cost.
- Produce deterministic, small cover assets that work in the dark theme.
- Reuse existing conventions rather than introduce new architectural layers.

**Non-Goals:**

- A vehicle detail route; cards are self-contained.
- Multiple drive locations per vehicle, or a vehicle-versus-drive split.
- Vector (SVG) covers, server-side inference, or hosted image APIs.
- Additional locales beyond the existing English catalogue.

## Decisions

### Single vehicle record, not vehicle plus drive

One row per vehicle, with `country_code`, `country` and `city` on the row. A vehicle-plus-drive join would model repeat drives and border crossings, but the page is one card per vehicle and the author's history is small. The join doubles the editor surface for a case that does not exist yet.

Alternative rejected: normalising into `Vehicle` and `Drive` with a relation. Revisit only if a vehicle genuinely needs two simultaneous locations.

### Ordering by explicit sort order, then drive date

`sort_order` ascending when set, `driven_at` ascending otherwise, so the tour reads chronologically while allowing a favourite to be pulled forward. The coffee `pinned` flag was considered and rejected as a second mechanism for the same need.

### Immersive layout through route metadata

A new `layout: 'immersive'` value extends the switch already in `App.vue`: the header stays, the container padding is dropped and the footer is suppressed. This avoids `calc(100dvh - 4rem)` arithmetic that breaks whenever the navbar height changes.

Alternative rejected: absolutely positioning the carousel over the shell, which breaks scroll restoration and the sticky header.

### Three-dimensional spinning carousel

Slides are absolutely stacked in one stage with a `perspective`. The active vehicle sits at the front; the adjacent vehicles are pushed sideways, receded on the Z axis and rotated about the Y axis, so the neighbouring cars stay visible behind it and the deck reads as a rotating rack. Advancing moves the deck one position and the offsets wrap around, so the tour spins rather than stopping at the ends.

Slides are `select-none` and the cover is `pointer-events-none` with `draggable="false"`, so a drag always moves the deck instead of dragging the image or selecting the copy. The deck uses a grab cursor to signal that it can be dragged.

Slides stay transparent, so the particle background shows through the deck and the neighbouring cards layer behind the active one. Transforms and opacity transitions are disabled under `prefers-reduced-motion`. The project has no carousel dependency and this does not justify one.

Alternative rejected: the initial scroll-snap implementation, which showed exactly one vehicle per viewport step and gave no hint that the others existed.

Use `100dvh` for the viewport unit minus the header, and never intercept vertical scrolling.

### Flags from a bundled icon collection

Country flags render through the existing offline Iconify component, using a `circle-flags` collection served as JSON from `public/` and registered with `addCollection` — the same mechanism as `mynaui-icons.json`. Flags pass `dotify="false"` so they stay legible.

Alternatives rejected: emoji flags render as letter pairs on Linux and Windows, including the author's own machine; an external flag CDN would violate the established no-external-icon-services constraint.

### One country field, two stored values

The editor presents the country label and ISO code as a single autocompleted field rather than two inputs that could disagree. Suggestions are derived from the bundled `circle-flags` collection filtered to two-letter codes and named with `Intl.DisplayNames`, so the suggestion list, the stored label, the stored code and the rendered flag cannot drift apart, and no country dataset or network call is added.

Selecting a suggestion writes both the label and the upper-case ISO code. Clearing the field clears both. Text that matches no suggestion is reverted on blur instead of being saved as an unmatched label, so a flag always has the country it belongs to.

The record still keeps `country` and `country_code` as separate columns: the card needs the code to choose a flag and the label to print beside it, and neither is reliably recoverable from the other once stored.

Alternatives rejected: separate label and code inputs, which is the inconsistency this replaces; storing only the code and resolving the label at render time, which would require the public page to load a naming dataset it otherwise does not need.

### On-device segmentation with `ormbg`, hosted from R2

`@huggingface/transformers` runs the quantized `ormbg` model in a Web Worker, preferring the WebGPU provider and falling back to WebAssembly. The 42 MB model is copied to R2 and referenced from the app origin rather than the Hugging Face CDN, so CORS, versioning and offline behaviour stay under the project's control, and Transformer.js caches it after first use.

Alternatives rejected: RMBG-1.4 and RMBG-2.0 are source-available for non-commercial use only; the `onnx-community/ISNet-ONNX` repo is AGPL-3.0; BiRefNet is MIT but 109 MB at fp16 and 467 MB at fp32; MODNet at 6 MB is tuned for portraits.

The runtime's own WebAssembly backend is not bundled. Its SIMD binary is roughly 27 MB, which exceeds Cloudflare Pages' 25 MiB per-file limit, so the emitted copy is dropped from the build. The asyncify and plain variants are published to the asset CDN under `onnxruntime/` and referenced from configuration; the worker picks the variant the way Transformer.js does — asyncify everywhere except Safari below 26 without WebGPU — and the browser plus the runtime's wasm cache keep it for later runs.

Asset origins must be absolute `https://` URLs. Transformer.js builds its file-metadata probe from the model host, and a protocol-relative value such as `//assets.grd0.net` is silently rejected there, which leaves the image processor unloaded and produces a `this.processor is not a function` failure at inference time.

### Exact square resize, mapped back

Inference input is the source resized to exactly 1024×1024, ignoring aspect. The matte is then resized back to the source frame before any compositing. This matches the model's declared preprocessing and was measurably better than letterboxing.

The delivered assets are then trimmed to the subject's alpha bounds with a two-pixel margin, so the slide is not dominated by the photograph's empty margins and the cutout and silhouettes share one canvas. The inference frame is only an intermediate.

### PNG outputs, no vectorizer

Both assets are PNG with an alpha channel: a photographic cutout that keeps colour and shading, and a two-tone silhouette in a light and a dark variant. Silhouettes are 7–9 KB at 1600 px wide. No tracer ships, and the soft photographic cutout — which cannot be vectorised — stays available.

The page uses the photographic cutout as the card cover, because a coloured car reads better against the dark stage than a flat white shape; the silhouette is rendered when a vehicle has no cutout, and the editor previews the cutout first for the same reason.

Alternative rejected: tracing to SVG with potrace, which worked but added a WASM tracer for no size benefit.

### Cover assets keyed, uploads deferred

The form holds the generated assets in memory and only requests a presigned URL and uploads when the author saves. Object keys are stored on the vehicle record under a `vehicles/{year}/{month}/{name}` prefix. Discarding the form discards the assets.

### Cover thumbnails in the editor listing

Each row in `/editor/vehicles` groups the vehicle's cover thumbnail with its name. The thumbnail shows the cutout when present, otherwise the silhouette, otherwise the "Unknown" placeholder, with the year, location and cover state beneath the name. Keeping the thumbnail and name together makes the row scannable instead of spreading the items across the full width.

### Failure reporting

The worker reports the underlying error rather than only a generic failure. The tool shows the translated message with the worker's detail beneath it, unreadable files are called out as unreadable images, and an empty matte reports that no subject was found. When a WebGPU session fails after selection, the worker retries once on WebAssembly before the error reaches the author. The form keeps its values throughout, and generated assets survive a failed upload so the author can save again.

### Manual erase after segmentation

The cover tool keeps the generated assets and an erase layer of brush strokes. Strokes are applied as `destination-out` over the photographic cutout and over both silhouettes, so the three deliverables stay in step, and the edited blobs are what the form uploads. Strokes are kept as a list so undo pops and replays them and clearing removes them all without re-running inference; Ctrl or ⌘ with Z mirrors the undo button outside form fields. Each edit is re-trimmed to the remaining alpha bounds before it is emitted, so erasing a shadow or a corner does not leave a transparent strip in the saved file.

The canvas fits the frame by default and zooms to 6× inside a scrollable viewport, which keeps the pointer-to-pixel mapping intact. A brush outline is drawn at the pointer, scaled to the brush size, so the stroke width is visible before the author commits to it. Pointer events are handled by the frame rather than the image, so a stroke can start in the empty margin and continue into the image; out-of-bounds points are simply clipped by the canvas. When zoomed, the image sits on a working surface padded by half a frame on each side and the view re-centres on every zoom change, so the scrollbars pan over the whole canvas rather than stopping at the image edges. Dragging with the right mouse button pans the frame as well, with the context menu suppressed so the gesture is reliable. Ctrl or ⌘ with the wheel zooms around the pointer — the wheel delta is divided into tenths of a step and the image ratio under the cursor is preserved by adjusting the scroll after the layout change — and the page does not scroll.

Alternative rejected: editing the matte and re-deriving everything on each stroke, which would re-run hole filling and speck removal and make undo expensive.

### Unknown placeholder instead of bundled silhouettes

A vehicle without a silhouette asset renders an "Unknown" placeholder where the cover would sit. A hand-drawn silhouette per body type was implemented first and then removed: it implied a vehicle shape the author had not chosen, added seven assets to maintain, and only covers a state that disappears as soon as a cover is generated. A placeholder is honest about the missing asset and still keeps the slide from rendering blank.

`body_type` stays on the record as metadata and remains a form option; it no longer selects an image.

Alternative rejected: the body-type silhouette map (sedan, hatchback, SUV, van, pickup, motorbike plus a generic default).

## Risks / Trade-offs

- **42 MB first-load cost** → Show download progress, lazy-import the runtime only when the tool opens, cache aggressively, and keep the model out of the PWA precache manifest.
- **Ground shadows survive segmentation** → The cover tool ships a manual erase brush; threshold tuning covers the common case.
- **Three-quarter photographs produce blobs** → State the side-on expectation in the tool's guidance; the cutout is only as good as the viewpoint.
- **WebGPU fragmentation** → Provider detection with an explicit CPU fallback, a retry on WebAssembly when a GPU session fails, and a clear unsupported-browser state.
- **Memory pressure on mobile** → Cap input resolution, dispose the inference session when the tool closes.
- **Model licensing changes upstream** → Pin the model revision and self-host the file; the licence is recorded in the design so a future swap is a conscious decision.
- **Heavy inference is hard to test automatically** → The page and editor are covered by Playwright specs; the cutout flow is verified with a throwaway end-to-end run against the hosted model.
- **Presigned URL expiry during a slow upload** → Request the URL only at save time, after generation completes.

## Migration Plan

1. Deploy the API with the additive `Vehicle` table migration and routes. No backfill; existing tables are untouched.
2. Upload the model to `models/ormbg/` and the ONNX Runtime wasm backend to `onnxruntime/` in the assets bucket, then deploy the web app with the page and editor. Vehicles can be entered immediately; their cards show the "Unknown" placeholder until a cover is generated.
3. The cutout tool runs only when the author opens it in the editor workspace; it has no effect until used.
4. Rollback: remove the navigation entry and routes; the table and assets can remain without affecting existing features.
