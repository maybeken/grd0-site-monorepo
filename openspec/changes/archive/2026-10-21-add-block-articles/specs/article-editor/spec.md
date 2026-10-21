# Spec Delta

## Purpose

Defines the authoring experience for block articles: creating an article first, then adding, editing, ordering, and removing its blocks.

## ADDED Requirements

### Requirement: Create before authoring

The editor SHALL require an existing article before any block can be added. Starting a new article MUST create an empty draft through the API and MUST open only after creation succeeds.

#### Scenario: Start a new article
- **WHEN** an authenticated author starts a new article and supplies an unused URI
- **THEN** the editor SHALL create an empty draft and open it for editing
- **THEN** the block list SHALL be empty

#### Scenario: Duplicate URI on creation
- **WHEN** the supplied URI already exists
- **THEN** the editor SHALL report the conflict and SHALL NOT open an article
- **THEN** the author SHALL be able to choose a different URI

#### Scenario: Unauthenticated editor access
- **WHEN** an unauthenticated visitor opens the editor
- **THEN** article and block operations SHALL fail without exposing draft content

### Requirement: Article list actions

The editor list SHALL offer open, unpublish, and permanent delete actions.

#### Scenario: Unpublish from the list
- **WHEN** the author unpublishes a published article from the list
- **THEN** the article SHALL remain editable and its blocks and assets SHALL be retained
- **THEN** the list SHALL show the article as unpublished

#### Scenario: Delete from the list
- **WHEN** the author confirms permanent deletion from the list
- **THEN** the article and its scene files SHALL be removed
- **THEN** the row SHALL disappear and the URI SHALL become available

### Requirement: Article metadata editing

The editor SHALL allow the author to edit the title, subtitle, publish timestamp, and URI of an article.

#### Scenario: Rename an article
- **WHEN** the author changes the URI to an unused value and saves
- **THEN** the editor SHALL continue editing the same article under its new URI
- **THEN** the previous URI SHALL no longer resolve

#### Scenario: Duplicate rename rejected
- **WHEN** the author changes the URI to one already in use
- **THEN** the editor SHALL report the conflict and retain the previous URI

### Requirement: Block palette

The editor SHALL offer markdown, image, music, and excalidraw blocks as additions to the article.

#### Scenario: Add a block
- **WHEN** the author chooses a block type
- **THEN** an empty block of that type SHALL be appended to the article
- **THEN** the block SHALL receive a stable identifier before its first save

#### Scenario: Controls follow the block list
- **WHEN** the author views an article with blocks
- **THEN** the block palette and save action SHALL appear after the last block

### Requirement: Markdown block editing

A markdown block SHALL be edited as a full-width markdown surface. The article preview SHALL provide the live rendering; a separate per-block preview is not required.

#### Scenario: Live preview
- **WHEN** the author types in a markdown block
- **THEN** the editor SHALL stage the text immediately
- **THEN** the article preview SHALL reflect the content within a short debounce (3 seconds) without requiring a save

#### Scenario: Upload an image into markdown
- **WHEN** the author drops a file onto a markdown block or chooses a file
- **THEN** the editor SHALL upload it through the article attachment flow
- **THEN** the editor SHALL insert a relative image reference into the markdown at the cursor

#### Scenario: Failed upload
- **WHEN** an upload fails
- **THEN** the editor SHALL report the failure and SHALL NOT insert a broken reference

### Requirement: Always-visible preview

The editor SHALL keep the public block renderer visible beside the editing surface at all times.

#### Scenario: Preview layout
- **WHEN** an article is open in the editor
- **THEN** the editing controls SHALL appear on the left and the preview on the right on wide screens
- **THEN** the preview SHALL be present without a toggle

#### Scenario: Debounced refresh
- **WHEN** the author changes metadata or blocks
- **THEN** the preview SHALL refresh 3 seconds after the last change
- **WHEN** an article is opened or saved
- **THEN** the preview SHALL refresh immediately

#### Scenario: Preview uses the public renderer
- **WHEN** the preview renders blocks
- **THEN** it SHALL use the same renderer as the public article view
- **THEN** drawing scenes SHALL be read through the authenticated draft scene endpoint

### Requirement: Image block authoring

An image block SHALL be authored by selecting an existing gallery asset and MAY include a caption.

#### Scenario: Select a gallery image
- **WHEN** the author opens the gallery picker
- **THEN** the picker SHALL list gallery collections and their assets with thumbnails
- **THEN** choosing an asset SHALL record its collection and filename in the block

#### Scenario: Caption
- **WHEN** the author enters caption text for an image block
- **THEN** the caption SHALL be stored with the block and displayed with the image

### Requirement: Music block authoring

A music block SHALL be authored by selecting a song from the existing music catalogue, and SHALL snapshot the song's title and artist into the block.

#### Scenario: Select a song
- **WHEN** the author opens the music picker
- **THEN** the picker SHALL list the songs available from the music catalogue
- **THEN** choosing a song SHALL store its video identifier, title, and artist in the block

#### Scenario: Catalogue entry changes later
- **WHEN** a selected song's title changes in the music catalogue after the article is saved
- **THEN** the article SHALL continue displaying the title captured at authoring time

### Requirement: Drawing block authoring

An excalidraw block SHALL be authored in an embedded drawing surface, and saving the drawing SHALL store its scene in article asset storage.

#### Scenario: Author a drawing
- **WHEN** the author draws and saves within the block
- **THEN** the editor SHALL store the scene and record its path in the block

#### Scenario: Reopen a drawing
- **WHEN** the author reopens an article containing a drawing block
- **THEN** the editor SHALL load the stored scene into the drawing surface

#### Scenario: Scene save failure
- **WHEN** storing the scene fails
- **THEN** the editor SHALL report the failure and SHALL NOT record an unusable path

#### Scenario: Scroll safety
- **WHEN** the pointer is over an embedded drawing that has not been activated
- **THEN** wheel and trackpad scrolling SHALL scroll the page instead of the drawing
- **WHEN** the author clicks the drawing
- **THEN** the drawing SHALL become interactive until the pointer leaves it

### Requirement: Block ordering and removal

The editor SHALL allow blocks to be reordered and removed.

#### Scenario: Reorder
- **WHEN** the author moves a block up or down
- **THEN** the new order SHALL be persisted with the article

#### Scenario: Viewport follows the moved block
- **WHEN** the author moves a block up or down
- **THEN** the editor SHALL scroll the moved block into view at its new position

#### Scenario: Remove a block
- **WHEN** the author removes a block and saves
- **THEN** the block SHALL no longer appear in the article
- **THEN** any stored scene belonging only to that block SHALL be deleted on a best-effort basis

### Requirement: Saving

The editor SHALL persist article metadata and the full ordered block list together, and SHALL report success or failure.

#### Scenario: Successful save
- **WHEN** the author saves
- **THEN** the article and its ordered blocks SHALL be persisted
- **THEN** the editor SHALL report success

#### Scenario: Failed save
- **WHEN** persistence fails
- **THEN** the editor SHALL report the failure and SHALL retain the author's unsaved edits

### Requirement: Publish now

The editor SHALL provide a Publish now action beside the publish date.

#### Scenario: Publish immediately
- **WHEN** the author activates Publish now
- **THEN** the publish timestamp SHALL be set to the current time
- **THEN** the article and its current blocks SHALL be saved
- **THEN** the editor SHALL report success or failure

### Requirement: Save shortcut

The editor SHALL save the article when the author presses Ctrl+S or Cmd+S.

#### Scenario: Keyboard save
- **WHEN** the author presses Ctrl+S or Cmd+S in the article editor
- **THEN** the editor SHALL prevent the browser's default behavior and save the article
- **THEN** repeated or concurrent save requests SHALL NOT create duplicate saves

#### Scenario: Save from a drawing
- **WHEN** the author presses Ctrl+S while a drawing canvas has focus
- **THEN** the request SHALL be forwarded to the editor rather than handled by the drawing library

### Requirement: Unsaved changes protection

The editor SHALL warn the author before discarding unsaved changes.

#### Scenario: Leave the editor
- **WHEN** the author changes any metadata or block and attempts to navigate away
- **THEN** the editor SHALL ask for confirmation
- **THEN** cancelling SHALL keep the editor and the changes

#### Scenario: Reload or close
- **WHEN** the author has unsaved changes and reloads or closes the page
- **THEN** the browser SHALL prompt for confirmation

#### Scenario: Open another article
- **WHEN** the author has unsaved changes and selects a different article in the editor list
- **THEN** the editor SHALL ask for confirmation
- **THEN** cancelling SHALL keep the current article and its changes

### Requirement: Non-blocking feedback

The editor SHALL report save, conflict, upload, and scene failures without blocking the rest of the interface.

#### Scenario: Save feedback
- **WHEN** a save succeeds or fails
- **THEN** the editor SHALL show a dismissible notification
- **THEN** the author SHALL be able to continue editing

### Requirement: Draft preview

The editor SHALL preview unpublished articles through the authenticated draft endpoint using the public renderer.

#### Scenario: Preview a draft
- **WHEN** the author opens the preview of an unpublished article
- **THEN** the editor SHALL render it with the same block renderer used publicly
- **THEN** the preview SHALL be available without publishing
