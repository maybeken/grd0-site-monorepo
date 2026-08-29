## ADDED Requirements

### Requirement: Feed editor route
The system SHALL provide a `/feed/editor` route that renders the feed editor page. The route SHALL require authentication.

#### Scenario: Navigate to feed editor
- **WHEN** an authenticated user navigates to `/feed/editor`
- **THEN** the system SHALL render the feed editor page with an empty feed entry form

#### Scenario: Navigate to feed editor for existing entry
- **WHEN** an authenticated user navigates to `/feed/editor/summer-update`
- **THEN** the system SHALL fetch the existing feed entry and populate the editor with its current state

#### Scenario: Unauthenticated access
- **WHEN** an unauthenticated user navigates to `/feed/editor`
- **THEN** the system SHALL redirect to the login page

### Requirement: Feed entry title and subtitle input
The system SHALL provide text inputs for the feed entry title and subtitle.

#### Scenario: Enter title
- **WHEN** a user types in the title input
- **THEN** the system SHALL update the feed entry title in the editor state

#### Scenario: Enter subtitle
- **WHEN** a user types in the subtitle input
- **THEN** the system SHALL update the feed entry subtitle in the editor state

#### Scenario: Load existing title and subtitle
- **WHEN** editing an existing feed entry
- **THEN** the system SHALL populate the title and subtitle inputs with the existing entry's values

### Requirement: Add blog block
The system SHALL allow users to add a blog block to the feed entry, with a markdown editor for content creation.

#### Scenario: Add new blog block
- **WHEN** a user clicks "Add Block" and selects "Blog"
- **THEN** the system SHALL add a new blog block to the editor with a markdown editor (reusing SideBySideEditor component)

#### Scenario: Edit blog block content
- **WHEN** a user types in the blog block's markdown editor
- **THEN** the system SHALL update the block's content (title, content) in the editor state

#### Scenario: Blog block fields
- **WHEN** a blog block is added
- **THEN** the editor SHALL provide inputs for title and content (markdown). Subtitle is not available at the block level (it is set at the feed entry level).

### Requirement: Add music block
The system SHALL allow users to add a music block to the feed entry, with inputs for YouTube video ID and metadata.

#### Scenario: Add new music block
- **WHEN** a user clicks "Add Block" and selects "Music"
- **THEN** the system SHALL add a new music block to the editor with inputs for YouTube video ID, title, and artist

#### Scenario: Edit music block metadata
- **WHEN** a user enters a YouTube video ID, title, and artist
- **THEN** the system SHALL update the block's data in the editor state

### Requirement: Add reference block
The system SHALL allow users to add a reference block pointing to an existing Blog or Music entry via a searchable dropdown.

#### Scenario: Add reference to existing blog
- **WHEN** a user clicks "Add Block" and selects "Reference Existing", then searches for and selects a blog post
- **THEN** the system SHALL add a blog reference block pointing to the selected blog post

#### Scenario: Add reference to existing music
- **WHEN** a user clicks "Add Block" and selects "Reference Existing", then searches for and selects a music entry
- **THEN** the system SHALL add a music reference block pointing to the selected music entry

#### Scenario: Searchable dropdown
- **WHEN** a user types in the reference block's search input
- **THEN** the system SHALL query the blog or music API and display matching results in a dropdown

#### Scenario: Select from dropdown
- **WHEN** a user clicks on a search result
- **THEN** the system SHALL populate the reference block with the selected entry's ID

### Requirement: Edit referenced blog
The system SHALL allow users to edit the blog post referenced by a selected blog block inline within the feed editor. Clicking a block selects it; when a blog block is selected, the full blog editor (SideBySideEditor) SHALL be revealed below the block list, allowing the user to edit the referenced blog post's title, subtitle, publish date, and markdown content.

#### Scenario: Select blog block to reveal editor
- **WHEN** a user clicks on a blog block in the block list
- **THEN** the system SHALL highlight the selected block and render the SideBySideEditor below the block list, loaded with the referenced blog post's content

#### Scenario: Edit referenced blog content
- **WHEN** a user modifies the title, subtitle, date, or markdown content in the SideBySideEditor
- **THEN** the system SHALL track the edits locally; the user SHALL save the blog post via the SideBySideEditor's save button to persist changes to the blog post

#### Scenario: Switch selected blog block
- **WHEN** a user clicks a different blog block while another is selected
- **THEN** the system SHALL replace the SideBySideEditor with the newly selected blog post's content

#### Scenario: Deselect blog block
- **WHEN** a user clicks the currently selected blog block again
- **THEN** the system SHALL deselect the block and hide the SideBySideEditor

#### Scenario: Remove selected block
- **WHEN** a user removes the currently selected block
- **THEN** the system SHALL deselect the block, hide the SideBySideEditor, and remove the block from the editor

### Requirement: Remove block
The system SHALL allow users to remove a block from the feed entry.

#### Scenario: Remove block
- **WHEN** a user clicks the remove button on a block
- **THEN** the system SHALL remove the block from the editor state and re-order the remaining blocks

### Requirement: Reorder blocks
The system SHALL allow users to reorder blocks using up/down buttons.

#### Scenario: Move block up
- **WHEN** a user clicks the up button on a block
- **THEN** the system SHALL swap the block with the one above it, updating the order field

#### Scenario: Move block down
- **WHEN** a user clicks the down button on a block
- **THEN** the system SHALL swap the block with the one below it, updating the order field

#### Scenario: First block cannot move up
- **WHEN** a block is at the top of the list
- **THEN** the up button SHALL be disabled

#### Scenario: Last block cannot move down
- **WHEN** a block is at the bottom of the list
- **THEN** the down button SHALL be disabled

### Requirement: Save draft
The system SHALL allow users to save the feed entry as a draft (published_at set to null).

#### Scenario: Save new draft
- **WHEN** a user clicks "Save Draft" on a new feed entry
- **THEN** the system SHALL send `PUT /feed/:slug` with published_at set to null, creating the entry and its blocks, and display a success message

#### Scenario: Update existing draft
- **WHEN** a user clicks "Save Draft" on an existing draft entry
- **THEN** the system SHALL send `PUT /feed/:slug` with the updated content, replacing all blocks

### Requirement: Publish
The system SHALL allow users to publish the feed entry (published_at set to current timestamp or a specified future timestamp).

#### Scenario: Publish immediately
- **WHEN** a user clicks "Publish" on a feed entry
- **THEN** the system SHALL send `PUT /feed/:slug` with published_at set to the current timestamp, making the entry visible on the feed listing

#### Scenario: Publish with error
- **WHEN** the publish request fails
- **THEN** the system SHALL display an error message and retain the editor state

### Requirement: Validation
The system SHALL validate the feed entry before saving.

#### Scenario: Title required
- **WHEN** a user attempts to save a feed entry with an empty title
- **THEN** the system SHALL display a validation error and prevent the save

#### Scenario: At least one block required
- **WHEN** a user attempts to save a feed entry with no blocks
- **THEN** the system SHALL display a validation error and prevent the save

#### Scenario: Blog block content required
- **WHEN** a user attempts to save a feed entry with a blog block that has no content
- **THEN** the system SHALL display a validation error and prevent the save

#### Scenario: Music block YouTube ID required
- **WHEN** a user attempts to save a feed entry with a music block that has no YouTube video ID
- **THEN** the system SHALL display a validation error and prevent the save
