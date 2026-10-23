# Spec Delta

## Purpose

Lets an authenticated author adopt asset directories as gallery collections and curate the metadata that the public gallery renders, without editing the database directly.

## ADDED Requirements

### Requirement: Gallery collections listing

The gallery editor SHALL provide a listing at `/editor/gallery` that shows every collection with its display title, cover, and asset count, and SHALL show the directories available for adoption separately from existing collections.

#### Scenario: Existing collections
- **WHEN** the author opens the gallery listing
- **THEN** every collection SHALL be listed with its display title, its cover when set, and the number of assets in its directory

#### Scenario: Title fallback
- **WHEN** a collection has no display title
- **THEN** the listing SHALL show the collection's directory name

#### Scenario: Empty gallery
- **WHEN** no collections exist and no asset directories are available
- **THEN** the listing SHALL show an empty state

#### Scenario: Open a collection
- **WHEN** the author selects a collection
- **THEN** the editor SHALL navigate to that collection's workspace route

#### Scenario: Failed load
- **WHEN** collections or directories cannot be loaded
- **THEN** the editor SHALL report the failure without blocking the rest of the page

### Requirement: Asset directory discovery

The system SHALL expose the gallery asset directories found in `files.json` to authenticated editors, including each directory's asset count, so directories without collection metadata can be discovered.

#### Scenario: List directories
- **WHEN** an authenticated client requests the gallery directories
- **THEN** the API SHALL return every directory under `/gallery` found in the asset index with its asset count

#### Scenario: Unauthenticated access
- **WHEN** a client without an application token requests the gallery directories
- **THEN** the API SHALL reject the request as unauthorised

#### Scenario: Directories already adopted
- **WHEN** a directory already has collection metadata
- **THEN** it SHALL appear as a collection rather than as an adoptable directory

### Requirement: Collection adoption

The gallery listing SHALL offer an adopt action for each directory that has no collection metadata. Adopting SHALL create the collection, keep the author on the listing, and refresh it.

#### Scenario: Adopt a directory
- **WHEN** the author adopts an available directory
- **THEN** the collection SHALL be created for that directory
- **THEN** the directory SHALL move from the available list into the collections list
- **THEN** no asset SHALL be created, moved, or deleted

#### Scenario: Adopt an unavailable directory
- **WHEN** the author adopts a directory that already has collection metadata
- **THEN** the system SHALL NOT create a second collection for it

#### Scenario: Adoption fails
- **WHEN** creating the collection fails
- **THEN** the editor SHALL report the failure and leave the directory available

### Requirement: Collection metadata editing

Each collection SHALL have a workspace at `/editor/gallery/:collection` that edits its display title and cover. Opening the route directly SHALL load the collection without visiting the listing first.

#### Scenario: Edit title
- **WHEN** the author changes the collection title and saves
- **THEN** the new title SHALL be persisted
- **THEN** the listing SHALL show the new title on return

#### Scenario: Clear title
- **WHEN** the author clears the collection title and saves
- **THEN** the public gallery SHALL fall back to the directory name

#### Scenario: Choose a cover
- **WHEN** the author opens the cover picker
- **THEN** the picker SHALL offer only assets from this collection
- **WHEN** the author selects an asset and saves
- **THEN** that filename SHALL become the collection's cover

#### Scenario: Direct link
- **WHEN** the author opens a collection workspace route directly
- **THEN** the collection and its assets SHALL load

#### Scenario: Unknown collection
- **WHEN** the route names a collection that does not exist
- **THEN** the editor SHALL show a not-found message and a link back to the listing

#### Scenario: Directory identity is fixed
- **WHEN** the author edits a collection
- **THEN** the directory the collection points at SHALL NOT be changeable

#### Scenario: Delete a collection
- **WHEN** the author deletes a collection and confirms
- **THEN** the collection metadata SHALL be removed
- **THEN** the assets and per-image overrides SHALL remain on disk and in storage
- **THEN** the editor SHALL return to the listing

### Requirement: Collection defaults

The collection workspace SHALL expose the wildcard detail row (`filename` equal to `*`) as editable collection defaults for timezone adjustment and fallback description, and SHALL explain that these apply to images without their own override.

#### Scenario: Defaults applied
- **WHEN** an image has no per-image override
- **THEN** the public gallery SHALL use the collection defaults for description and timezone adjustment

#### Scenario: Per-image override wins
- **WHEN** an image has its own override
- **THEN** the public gallery SHALL use the per-image values instead of the collection defaults

#### Scenario: Clear defaults
- **WHEN** the author clears the default description or timezone adjustment and saves
- **THEN** the wildcard row SHALL stop supplying that value

### Requirement: Per-image overrides

The collection workspace SHALL list the assets that exist in the collection's directory and SHALL allow each asset's description and timezone adjustment to be overridden. Changing an override MUST NOT modify the asset file or its EXIF metadata.

#### Scenario: List assets
- **WHEN** the author opens a collection workspace
- **THEN** every asset in the collection SHALL be listed with its thumbnail, filename, and existing override values

#### Scenario: Show inherited defaults
- **WHEN** an asset has no stored override row
- **THEN** the workspace SHALL show the collection defaults as that asset's inherited values

#### Scenario: Start from defaults
- **WHEN** the author begins editing an asset that has no stored override row
- **THEN** the form SHALL start from the collection defaults rather than from blank values

#### Scenario: Set an override
- **WHEN** the author sets a description or timezone adjustment for an asset and saves
- **THEN** an override row SHALL be stored for that collection and filename

#### Scenario: Override replaces defaults
- **WHEN** an asset has a stored override row
- **THEN** the public gallery SHALL use that row's values and SHALL NOT merge individual fields from the collection defaults

#### Scenario: Effective time preview
- **WHEN** an asset has a capture time and an effective timezone adjustment
- **THEN** the workspace SHALL show the resulting local time

#### Scenario: Clear an override
- **WHEN** the author clears both the description and the timezone adjustment for an asset and saves
- **THEN** the stored override row SHALL be removed
- **THEN** the asset SHALL fall back to the collection defaults

#### Scenario: Large collection
- **WHEN** a collection has more assets than one page
- **THEN** the workspace SHALL paginate the assets
- **THEN** overrides edited on an earlier page SHALL remain pending when the author changes page

### Requirement: Save and discard behaviour

The collection workspace SHALL apply override changes only when the author saves, and SHALL warn before discarding pending changes through in-application navigation, reload, or closing the tab.

#### Scenario: Pending changes
- **WHEN** the author edits overrides without saving
- **THEN** the stored overrides SHALL remain unchanged

#### Scenario: Save pending changes
- **WHEN** the author saves with valid pending changes
- **THEN** every pending override SHALL be persisted
- **THEN** the workspace SHALL report success

#### Scenario: Navigate away with pending changes
- **WHEN** the author navigates away with pending changes
- **THEN** the editor SHALL ask for confirmation
- **THEN** cancelling SHALL keep the workspace and the pending changes

#### Scenario: Save fails
- **WHEN** persisting an override fails
- **THEN** the workspace SHALL report the failure and keep the pending changes

### Requirement: Stale override cleanup

The collection workspace SHALL identify stored overrides whose filename no longer exists in the collection's assets, excluding the wildcard row, and SHALL allow them to be deleted.

#### Scenario: Show stale overrides
- **WHEN** a stored override names a file that is not in the asset index
- **THEN** the workspace SHALL list it separately from the existing assets

#### Scenario: Delete a stale override
- **WHEN** the author deletes a stale override and confirms
- **THEN** the stored override SHALL be removed
- **THEN** the asset files SHALL remain unchanged

#### Scenario: Do not apply stale overrides
- **WHEN** an override is stale
- **THEN** it SHALL NOT be applied to any asset
