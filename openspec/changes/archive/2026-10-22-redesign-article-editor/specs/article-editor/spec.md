# Spec Delta

## ADDED Requirements

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
- **WHEN** the article is published
- **THEN** the sidebar SHALL offer a link to its public page

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

## MODIFIED Requirements

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
