# Proposal

## Why

Gallery assets are already published from the `assets/` pipeline, but the metadata that makes them readable — collection titles, covers, per-image descriptions, and timezone corrections — can only be changed by calling the API by hand. The editor console covers articles, coffee, and music but has no gallery section, so curating a trip after uploading photos is a manual database exercise.

## What Changes

- Add a Gallery section to the editor console with a collections listing at `/editor/gallery` and a per-collection workspace at `/editor/gallery/:collection`.
- List existing collections with their title, cover, and asset count.
- Discover asset directories present in `files.json` that have no collection metadata and offer to publish them, so a collection is only created for a directory that actually contains assets.
- Edit collection metadata: title and cover. Cover selection SHALL reuse the article editor's gallery picker, restricted to the collection's assets.
- Expose the currently implicit `filename: "*"` detail row as collection defaults for timezone adjustment and fallback description.
- Edit per-image overrides against the assets that exist: description and timezone adjustment, with a preview of the resulting local time.
- Surface overrides whose file no longer exists so they can be removed.
- Add the frontend write services for gallery collections and details; no database migration is required.
- Add an admin endpoint that lists gallery asset directories from `files.json`, because the existing asset endpoints only return directories that already have a collection row.

## Capabilities

### New Capabilities

- `gallery-editor`: Collection publishing and metadata curation, including collection defaults, per-image overrides, and stale-override cleanup.

### Modified Capabilities

- `editor-console`: The sidebar gains a Gallery entry, and the shell's route and section-default requirements cover the gallery listing and workspace.

## Impact

- **Frontend**: new `views/editor/EditorGalleryListView.vue` and `EditorGalleryView.vue`; new gallery form components; gallery write services; `GalleryPicker` gains an optional collection restriction; router and `EditorSidebar` additions; new `gallery.editor.*` and `editor.gallery` locale strings.
- **API**: one new admin `GET` handler for asset directories, registered beside the existing gallery routes. Existing collection and detail write endpoints are reused unchanged.
- **Data**: no schema change and no migration. Collections remain keyed by the bare directory name; detail rows remain keyed by `/gallery/<name>` and filename.
- **Out of scope**: uploading, moving, or deleting image assets; changing the TOC generator or R2 upload pipeline.
