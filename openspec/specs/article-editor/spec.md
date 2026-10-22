# article-editor Specification

## Purpose

Defines the authoring experience for block articles: creating an article first, then adding, editing, ordering, and removing its blocks.

## Requirements

### Requirement: Two-page editor navigation

The editor SHALL be addressed as two distinct pages: a listing page and an editing page. The listing MUST be reachable at a fixed route, and the editing page MUST identify its article in the route.

#### Scenario: Listing page
- **WHEN** the author opens the editor entry point
- **THEN** the listing SHALL be shown without an editing surface

#### Scenario: Editing page
- **WHEN** the author opens an editing route for an existing article
- **THEN** the editing workspace SHALL load that article

#### Scenario: Direct link
- **WHEN** an author opens an editing route directly
- **THEN** the workspace SHALL load the article without visiting the listing first

#### Scenario: Return to listing
- **WHEN** the author returns to the listing from an article
- **THEN** the listing SHALL restore the page the author was viewing

#### Scenario: Unknown article
- **WHEN** an editing route names an article that does not exist
- **THEN** the workspace SHALL show an error and offer a return to the listing

### Requirement: Listing pagination

The listing SHALL request articles one page at a time from the server and SHALL keep the current page addressable.

#### Scenario: Paged listing
- **WHEN** the listing loads
- **THEN** it SHALL show one page of articles together with the total number of articles and pages

#### Scenario: Change page
- **WHEN** the author selects another page
- **THEN** the listing SHALL request and show that page

#### Scenario: Page in the route
- **WHEN** the author is viewing page 2
- **THEN** the route SHALL identify page 2
- **WHEN** that route is opened
- **THEN** the listing SHALL show page 2

#### Scenario: Empty page after removal
- **WHEN** unpublishing or deleting leaves the current page empty and it is not the first page
- **THEN** the listing SHALL move to the previous page

### Requirement: Create before authoring

The editor SHALL require an existing article before any block can be added. Starting a new article MUST create an empty draft through the API and MUST open the editing route only after creation succeeds. The listing MUST NOT require a URI before creation; the URI is set afterwards in the metadata sidebar. The listing SHALL create the draft immediately after the author confirms, without a separate create form step.

#### Scenario: Start a new article
- **WHEN** an authenticated author starts a new article from the listing
- **THEN** the editor SHALL create an empty draft with a generated URI
- **THEN** the editor SHALL open the editing route for that URI
- **THEN** the block list SHALL be empty

#### Scenario: Single-step creation
- **WHEN** the author chooses New Article
- **THEN** the editor SHALL ask for confirmation
- **WHEN** the author confirms
- **THEN** the editor SHALL create the draft and open it
- **THEN** the editor SHALL NOT require a further create action

#### Scenario: Set the URI after creation
- **WHEN** the author enters a URI in the metadata sidebar and saves
- **THEN** the article SHALL be addressable under the new URI
- **THEN** the generated URI SHALL no longer resolve

#### Scenario: Duplicate URI on creation
- **WHEN** the supplied URI already exists
- **THEN** the editor SHALL report the conflict and SHALL NOT open an article
- **THEN** the author SHALL be able to choose a different URI

#### Scenario: Unauthenticated editor access
- **WHEN** an unauthenticated visitor opens the editor
- **THEN** article and block operations SHALL fail without exposing draft content

### Requirement: Article list actions

The editor list SHALL offer open, unpublish, and permanent delete actions.

#### Scenario: Open from the list
- **WHEN** the author chooses an article in the listing
- **THEN** the editor SHALL navigate to that article's editing route
- **THEN** the route SHALL identify the article by its URI

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

### Requirement: Editing workspace layout

The editing page SHALL present a shared title and subtitle header, an editor pane, a preview pane, and a metadata sidebar, with the editor and preview scrolling independently of the page.

#### Scenario: Shared header
- **WHEN** the workspace is shown on a wide screen
- **THEN** the title and subtitle fields SHALL span the editor and preview width above both panes
- **THEN** the header SHALL NOT extend over the metadata sidebar

#### Scenario: Wide layout
- **WHEN** the workspace is shown on a wide screen
- **THEN** the editor, preview, and sidebar SHALL be visible together
- **THEN** the editor and preview SHALL be separated by a single vertical divider

#### Scenario: Narrow layout
- **WHEN** the workspace is shown on a narrow screen
- **THEN** the author SHALL be able to switch between the editor and the preview

#### Scenario: Independent scrolling
- **WHEN** the author scrolls the editor pane
- **THEN** the page SHALL NOT scroll with it

#### Scenario: Preview scrollbar
- **WHEN** the preview pane scrolls
- **THEN** its scrollbar SHALL be hidden

### Requirement: Metadata sidebar

The editing page SHALL show the article's URI, publication state, and publishing controls in a sidebar.

#### Scenario: URI
- **WHEN** the workspace is open
- **THEN** the sidebar SHALL show the article's current URI
- **THEN** the author SHALL be able to edit the URI

#### Scenario: Publication state
- **WHEN** the workspace is open
- **THEN** the sidebar SHALL show the publish timestamp and whether the article is a draft, scheduled, or published

#### Scenario: Public link
- **WHEN** the workspace is open
- **THEN** the sidebar SHALL offer a link to the article's public page regardless of its publication state

#### Scenario: Copy link
- **WHEN** the article has a publish date
- **THEN** the sidebar SHALL offer a copy control beside the URI
- **THEN** activating it SHALL copy the article's public URL

#### Scenario: Publication actions
- **WHEN** the article is published
- **THEN** Publish now SHALL be hidden
- **WHEN** the article is a draft
- **THEN** Unpublish SHALL be hidden
- **WHEN** the article is scheduled
- **THEN** both Publish now and Unpublish SHALL be available

#### Scenario: Delete
- **WHEN** the author chooses to delete the article from the sidebar
- **THEN** the editor SHALL require confirmation before permanently deleting it

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

A markdown block SHALL be edited as a full-width markdown surface that grows with its content. The article preview SHALL provide the live rendering; a separate per-block preview is not required.

#### Scenario: Live preview
- **WHEN** the author types in a markdown block
- **THEN** the editor SHALL stage the text immediately
- **THEN** the article preview SHALL reflect the content within one second without requiring a save

#### Scenario: Block grows with content
- **WHEN** the text in a markdown block exceeds the visible area
- **THEN** the block SHALL grow instead of scrolling internally

#### Scenario: Upload an image into markdown
- **WHEN** the author drops a file onto a markdown block or chooses a file
- **THEN** the editor SHALL upload it through the article attachment flow
- **THEN** the editor SHALL insert a relative image reference into the markdown at the cursor

#### Scenario: Failed upload
- **WHEN** an upload fails
- **THEN** the editor SHALL report the failure and SHALL NOT insert a broken reference

### Requirement: Always-visible preview

The editor SHALL present the public block renderer beside the editing surface on wide screens and SHALL keep it synchronised with the editor's scroll position. The preview SHALL render blocks only; the title and subtitle are edited in the shared header.

#### Scenario: Preview layout
- **WHEN** an article is open in the editor on a wide screen
- **THEN** the title and subtitle fields SHALL appear above the editor and preview
- **THEN** the editing controls SHALL appear on the left, the preview beside them, and the metadata sidebar on the right
- **THEN** the preview SHALL be present without a toggle
- **WHEN** the workspace is shown on a narrow screen
- **THEN** a toggle SHALL switch between the editor and the preview

#### Scenario: No duplicated header
- **WHEN** the preview renders
- **THEN** it SHALL NOT repeat the article title or subtitle

#### Scenario: Editor-driven scrolling
- **WHEN** the author scrolls the editor pane
- **THEN** the preview SHALL scroll so that the corresponding block is aligned
- **WHEN** the author scrolls the preview pane alone
- **THEN** the editor's scroll position SHALL NOT change

#### Scenario: Block alignment
- **WHEN** the preview follows the editor
- **THEN** the alignment SHALL use block identity rather than proportional scroll position
- **THEN** blocks with different rendered heights SHALL remain corresponding

#### Scenario: Position after refresh
- **WHEN** the preview refreshes after a debounced change
- **THEN** it SHALL realign to the block currently shown in the editor

#### Scenario: Debounced refresh
- **WHEN** the author changes a block
- **THEN** the preview SHALL refresh within one second of the last change
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

A music block SHALL be authored either by selecting a song from the existing music catalogue or by entering a YouTube link, and SHALL snapshot the song's video identifier, title, and artist into the block. The YouTube link form SHALL accept a bare video identifier or a link in the watch, `youtu.be`, `embed`, `shorts`, or `live` form, including `music.youtube.com` and `youtube-nocookie.com` hosts.

#### Scenario: Select a song
- **WHEN** the author opens the music picker
- **THEN** the picker SHALL list the songs available from the music catalogue
- **THEN** choosing a song SHALL store its video identifier, title, and artist in the block

#### Scenario: Author from a YouTube link
- **WHEN** the author pastes a YouTube link into the music picker and confirms it
- **THEN** the picker SHALL extract the video identifier from the link
- **THEN** the picker SHALL resolve the video's title and artist
- **THEN** confirming SHALL store the video identifier, title, and artist in the block

#### Scenario: Unsupported link
- **WHEN** the author confirms text that is not a recognised YouTube link
- **THEN** the picker SHALL report the error and leave the block unchanged

#### Scenario: Link cannot be resolved
- **WHEN** the video behind a valid YouTube link cannot be read
- **THEN** the picker SHALL report the error and leave the block unchanged

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
- **THEN** the preview SHALL reflect the new order

#### Scenario: Viewport follows the moved block
- **WHEN** the author moves a block up or down
- **THEN** the editor SHALL scroll the moved block into view at its new position

#### Scenario: Remove a block
- **WHEN** the author removes a block and saves
- **THEN** the block SHALL no longer appear in the article
- **THEN** any stored scene belonging only to that block SHALL be deleted on a best-effort basis

### Requirement: Block duplication

The editor SHALL allow the author to duplicate an existing block. A duplicate MUST receive a new stable identifier, MUST carry the source block's content, and MUST be inserted immediately after the source block.

#### Scenario: Duplicate a block
- **WHEN** the author duplicates a block
- **THEN** a copy SHALL be inserted immediately after it
- **THEN** the copy SHALL have a new stable identifier
- **THEN** the copy SHALL carry the source block's content

#### Scenario: Duplicate a drawing
- **WHEN** the author duplicates a drawing block
- **THEN** the copy SHALL receive its own stored scene
- **THEN** editing one copy SHALL NOT change the other

#### Scenario: Persist the duplicate
- **WHEN** the author saves after duplicating
- **THEN** the article SHALL contain both blocks in the new order

### Requirement: Saving

The editor SHALL persist article metadata and the full ordered block list together, and SHALL report success or failure. The save control SHALL indicate whether the article is saved, unsaved, or saving.

#### Scenario: Successful save
- **WHEN** the author saves
- **THEN** the article and its ordered blocks SHALL be persisted
- **THEN** the editor SHALL report success

#### Scenario: Failed save
- **WHEN** persistence fails
- **THEN** the editor SHALL report the failure and SHALL retain the author's unsaved edits

#### Scenario: Save state
- **WHEN** the article has unsaved changes
- **THEN** the save control SHALL show the save action
- **WHEN** the article has no unsaved changes
- **THEN** the save control SHALL show the saved state
- **WHEN** a save is running
- **THEN** the save control SHALL show the saving state

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

#### Scenario: Return to the listing
- **WHEN** the author has unsaved changes and returns to the listing
- **THEN** the editor SHALL ask for confirmation
- **THEN** cancelling SHALL keep the current article and its changes

#### Scenario: Reload or close
- **WHEN** the author has unsaved changes and reloads or closes the page
- **THEN** the browser SHALL prompt for confirmation

#### Scenario: Open another article
- **WHEN** the author has unsaved changes and navigates to another article's editing route
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
