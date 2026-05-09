## ADDED Requirements

### Requirement: Block-based markdown editing
The blog editor SHALL provide a block-based editing experience where each content element (paragraph, heading, list, image, code block, blockquote, horizontal rule) is a distinct editable block.

#### Scenario: Create heading via slash command
- **WHEN** the user types `/heading` and selects the heading option from the command menu
- **THEN** a new heading block is created at the cursor position

#### Scenario: Type paragraph text
- **WHEN** the user types text in a paragraph block and presses Enter
- **THEN** a new paragraph block is created below

#### Scenario: Reorder blocks
- **WHEN** the user drags a block handle to a new position
- **THEN** the block is moved to the target position and surrounding content reflows

### Requirement: Slash command menu
The editor SHALL display a command menu when the user types `/` at the start of an empty block, listing available block types (heading, image, code block, bullet list, ordered list, blockquote, horizontal rule).

#### Scenario: Open slash menu
- **WHEN** the user types `/` at the beginning of an empty paragraph block
- **THEN** a command menu appears showing available block types

#### Scenario: Select block type
- **WHEN** the user selects a block type from the slash menu
- **THEN** the current block is converted to the selected type

#### Scenario: Close slash menu
- **WHEN** the user presses Escape or clicks outside the slash menu
- **THEN** the command menu is dismissed

### Requirement: Inline formatting bubble menu
The editor SHALL display a floating bubble menu when the user selects text, providing bold, italic, strikethrough, inline code, and link formatting options.

#### Scenario: Format selected text
- **WHEN** the user selects text and clicks "Bold" in the bubble menu
- **THEN** the selected text is wrapped in bold formatting

#### Scenario: Add link via bubble menu
- **WHEN** the user selects text, clicks the link button, enters a URL, and confirms
- **THEN** the selected text becomes a hyperlink to the specified URL

#### Scenario: Bubble menu dismissal
- **WHEN** the user deselects text or clicks outside the bubble menu
- **THEN** the bubble menu is hidden

### Requirement: Markdown import on load
The editor SHALL parse existing markdown content into its internal block model when a post is loaded for editing.

#### Scenario: Load existing markdown post
- **WHEN** a post with markdown content is opened in the editor
- **THEN** the content is parsed and rendered as editable blocks in the editor

#### Scenario: Load empty new post
- **WHEN** a new post with no content is created
- **THEN** the editor displays a single empty paragraph block with placeholder text

### Requirement: Markdown export on save
The editor SHALL serialize its internal block model back to a markdown string when the post is saved.

#### Scenario: Save edited post
- **WHEN** the user clicks "Save" on an edited post
- **THEN** the editor content is serialized to a markdown string and sent to the backend API

#### Scenario: Round-trip fidelity
- **WHEN** a post is loaded into the editor and saved without modifications
- **THEN** the output markdown is semantically equivalent to the input (minor whitespace or formatting convention differences such as `**bold**` vs `__bold__` are acceptable)

### Requirement: Drag-and-drop image upload
The editor SHALL support dragging image files directly into the editing area, uploading them via the existing S3 presigned URL flow, and inserting them as image blocks.

#### Scenario: Upload image via drag and drop
- **WHEN** the user drags an image file onto the editor
- **THEN** the image is uploaded to S3 and inserted as an image block at the drop position

#### Scenario: Upload progress indication
- **WHEN** an image is being uploaded
- **THEN** a loading state is shown on the image placeholder

#### Scenario: Upload failure handling
- **WHEN** an image upload fails
- **THEN** a toast notification displays the error and the image is not inserted

### Requirement: Real-time preview
The editor SHALL provide a live preview of the rendered markdown content in the right pane, using the same unified/remark/rehype rendering pipeline as the public blog.

#### Scenario: Preview updates on edit
- **WHEN** the user types or modifies content in the editor
- **THEN** the preview pane updates within 500ms to reflect the current markdown output

#### Scenario: Preview matches public render
- **WHEN** the preview renders markdown content
- **THEN** the rendered HTML matches what would appear on the public blog page for the same markdown
