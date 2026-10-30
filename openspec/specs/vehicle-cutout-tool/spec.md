# vehicle-cutout-tool Specification

## Purpose

Turns a photograph the author already has into the two cover assets the Grand Tour needs — a transparent photographic cutout and a flat silhouette — using segmentation that runs entirely in the browser.

## Requirements

### Requirement: On-device segmentation

The tool SHALL segment the supplied photograph locally in the browser. Image data MUST NOT be sent to a remote inference service, and the feature MUST NOT require an API key or create per-image cost.

#### Scenario: Photograph processed locally
- **WHEN** the author supplies a photograph
- **THEN** the segmentation runs in the browser and no image content is transmitted off the device

#### Scenario: Works without a network after first use
- **WHEN** the model has already been downloaded and the device is offline
- **THEN** segmentation still runs

### Requirement: Execution provider selection

The tool SHALL use the WebGPU execution provider when the browser supports it, SHALL fall back to a WebAssembly provider when WebGPU is unavailable or a GPU session fails, and SHALL report clearly when neither can run.

#### Scenario: WebGPU available
- **WHEN** the browser exposes WebGPU
- **THEN** inference uses the GPU provider

#### Scenario: Fallback
- **WHEN** WebGPU is unavailable but WebAssembly is supported
- **THEN** inference completes on the CPU provider with reduced speed

#### Scenario: Provider failure
- **WHEN** the WebGPU provider is selected but the session cannot start or run
- **THEN** the tool retries on WebAssembly before reporting a failure

#### Scenario: Unsupported browser
- **WHEN** neither provider can run
- **THEN** the tool explains that the feature is unavailable on this browser and does not attempt inference

### Requirement: Cutout output

The tool SHALL produce a cutout image in PNG format with a transparent background that retains the subject's colour and detail.

#### Scenario: Background removed
- **WHEN** segmentation completes on a photograph with a background
- **THEN** the cutout contains only the vehicle with everything else transparent

### Requirement: Silhouette output

The tool SHALL produce a flat silhouette as a PNG with a transparent background, filling enclosed regions such as windows and discarding small disconnected fragments, and SHALL offer the silhouette in a colour suited to a dark page and one suited to a light page.

#### Scenario: Solid silhouette
- **WHEN** the silhouette is generated from a matte containing window holes and small specks
- **THEN** the windows are filled, the specks are absent, and the shape is a single continuous region

#### Scenario: Both colour variants
- **WHEN** the author generates a silhouette
- **THEN** a light-on-transparent and a dark-on-transparent variant are available

### Requirement: Manual erase

After segmentation, the author SHALL be able to erase regions from the generated cover, applied to the cutout and both silhouettes. The tool SHALL offer undo and a way to clear every erasure, and the erased result SHALL be what the form saves.

#### Scenario: Erasing leftovers
- **WHEN** the author drags across the generated image
- **THEN** the painted region is removed from the cutout and from both silhouettes

#### Scenario: Undo
- **WHEN** the author undoes an erasure
- **THEN** the erased region is restored in all three assets

#### Scenario: Saved result
- **WHEN** the author saves after erasing
- **THEN** the uploaded assets reflect the erased regions

### Requirement: Trim after erasing

The edited assets SHALL be re-trimmed to the remaining content on every erasure.

#### Scenario: Erasing an edge
- **WHEN** an erasure removes content at the edge of the cover
- **THEN** the saved cutout and silhouettes are cropped to what remains

#### Scenario: Everything erased
- **WHEN** an erasure leaves no content
- **THEN** the tool reports that everything was erased and keeps the previous assets

### Requirement: Erase canvas controls

The erase canvas SHALL zoom for precise control and SHALL accept strokes anywhere in its frame, including outside the image bounds. A brush outline matching the current brush size SHALL follow the pointer. Ctrl or ⌘ with the wheel SHALL zoom around the pointer without scrolling the page.

#### Scenario: Stroke from the empty frame
- **WHEN** the author starts a stroke outside the image but inside the canvas frame and drags into the image
- **THEN** the stroke erases where it crosses the image instead of being ignored

#### Scenario: Zoom
- **WHEN** the author zooms the canvas
- **THEN** the image enlarges inside a scrollable frame and erasing still removes the region under the pointer

#### Scenario: Wheel zoom
- **WHEN** the author holds Ctrl or ⌘ and scrolls over the frame
- **THEN** the canvas zooms around the pointer instead of scrolling the page

#### Scenario: Brush outline
- **WHEN** the pointer is over the canvas frame
- **THEN** an outline matching the current brush size follows it

### Requirement: Canvas panning

When zoomed in, the erase frame SHALL provide a working surface larger than the image, SHALL scroll over that surface, and SHALL pan when the author drags with the right mouse button without opening a context menu.

#### Scenario: Scrolling the canvas
- **WHEN** the author zooms until the image exceeds the frame
- **THEN** the frame scrolls over a surface that extends beyond the image's edges on every side

#### Scenario: Right-drag panning
- **WHEN** the author drags with the right mouse button inside the frame
- **THEN** the view moves with the pointer and the context menu does not open

### Requirement: Output geometry

Both outputs SHALL be trimmed to the subject. Transparent margins MUST be removed so the canvas fits the visible vehicle, the cutout and both silhouette variants SHALL share the same dimensions, and a small margin SHALL remain so the subject's edges are not clipped. The full source frame MUST NOT be preserved.

#### Scenario: Empty margins removed
- **WHEN** a photograph has empty space around the vehicle
- **THEN** both the cutout and the silhouette are delivered cropped to the vehicle's bounds

#### Scenario: Shared dimensions
- **WHEN** the cutout and the silhouette variants are generated
- **THEN** all three images have the same width and height

#### Scenario: No subject
- **WHEN** segmentation finds no subject in the photograph
- **THEN** the tool reports that no subject was found instead of returning a blank canvas

### Requirement: Model delivery and caching

The model SHALL be served from the application's own asset origin. On first use the tool SHALL show download progress, and on subsequent uses SHALL load from browser cache without redownloading. The model MUST NOT be loaded on public pages.

#### Scenario: First use
- **WHEN** the tool is opened for the first time
- **THEN** download progress is visible and the tool becomes usable when the model is ready

#### Scenario: Later use
- **WHEN** the tool is reopened after a successful run
- **THEN** the model loads from cache and no download progress is shown

#### Scenario: Public page unaffected
- **WHEN** a visitor browses the Grand Tour page
- **THEN** no segmentation model is requested

### Requirement: Cover upload

The tool SHALL upload generated assets to storage through the authenticated presigned upload flow and SHALL return the object keys that the vehicle record stores. Upload failure MUST NOT discard the generated assets locally.

#### Scenario: Successful upload
- **WHEN** the author accepts a generated cover and authentication is valid
- **THEN** the asset is uploaded and its key is available to the vehicle form

#### Scenario: Upload fails
- **WHEN** the upload request fails
- **THEN** the generated cover remains available in the tool and the failure is reported

### Requirement: Failure handling

When segmentation or upload fails, the tool SHALL report the failure without losing the vehicle form's unsaved values.

#### Scenario: Segmentation error
- **WHEN** inference fails for an image
- **THEN** the error is shown and the form retains its current values

#### Scenario: Unreadable image
- **WHEN** the supplied file cannot be decoded as an image
- **THEN** the tool reports that the file could not be read and the form retains its current values
