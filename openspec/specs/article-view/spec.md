# article-view Specification

## Purpose

Defines the public article experience: the article index, the article page, and how each block type is presented to readers.

## Requirements

### Requirement: Article index

The public article index SHALL list published articles with their title, subtitle, author, timestamps, summary, and cover.

#### Scenario: Published articles listed
- **WHEN** a reader opens the article index
- **THEN** every published article SHALL be shown, most recent first
- **THEN** unpublished and future-dated articles SHALL NOT be shown

#### Scenario: Cover shown when present
- **WHEN** a listed article has a cover
- **THEN** the index SHALL display the cover image
- **WHEN** the cover is a relative asset path
- **THEN** the index SHALL request it through the CDN resizing helper
- **WHEN** the cover is an absolute URL
- **THEN** the index SHALL render it directly without applying the CDN transform
- **WHEN** a listed article has no cover
- **THEN** the index SHALL render the entry without an image

#### Scenario: Summary-card presentation
- **WHEN** the index renders a listed article
- **THEN** it SHALL use the same summary-card presentation as the blog listing
- **THEN** the card SHALL show the author, title, subtitle, optional cover, rendered summary, a continuation link, and timestamps

### Requirement: Public navigation

The article index SHALL be the public article entry point. The legacy blog routes MAY remain available for existing links but MUST NOT be linked from the main navigation.

#### Scenario: Main navigation
- **WHEN** a reader views the main navigation
- **THEN** the article entry SHALL link to `/article`
- **THEN** the legacy `/blog` index SHALL NOT be linked

#### Scenario: Legacy routes remain
- **WHEN** a reader opens an existing `/blog` or `/blog/:slug` link
- **THEN** the legacy blog SHALL still render

### Requirement: Article sitemap

The sitemap SHALL include published article URLs and SHALL NOT include legacy blog URLs.

#### Scenario: Sitemap contents
- **WHEN** the sitemap is generated
- **THEN** it SHALL include `/article` and each published article under `/article/:uri`
- **THEN** it SHALL NOT include `/blog` or `/blog/:uri`
- **THEN** drafts and future-dated articles SHALL NOT appear

### Requirement: Article page

The public article page SHALL resolve an article by URI regardless of publish state and render its header and blocks in stored order.

#### Scenario: Published article resolves
- **WHEN** a reader opens a URI whose article is published
- **THEN** the page SHALL show the title, subtitle, author, and timestamps
- **THEN** the blocks SHALL be rendered in their stored order

#### Scenario: Draft article resolves
- **WHEN** a reader opens the URI of an article with no publish timestamp
- **THEN** the page SHALL render the article and its blocks instead of the not-found state

#### Scenario: Scheduled article resolves
- **WHEN** a reader opens the URI of an article whose publish timestamp is in the future
- **THEN** the page SHALL render the article and its blocks instead of the not-found state

#### Scenario: Unknown or unpublished URI
- **WHEN** a reader opens a URI that no article uses
- **THEN** the page SHALL show the not-found state
- **THEN** a URI whose article is unpublished SHALL NOT show the not-found state

### Requirement: Unpublished article notice

The public article page SHALL display a notice when the article is not published yet, meaning its publish timestamp is absent or in the future. The notice MUST state that the article is unpublished.

#### Scenario: Draft article notice
- **WHEN** a reader opens the URI of an article with no publish timestamp
- **THEN** the page SHALL show a notice that the article is not published yet

#### Scenario: Scheduled article notice
- **WHEN** a reader opens the URI of an article whose publish timestamp is in the future
- **THEN** the page SHALL show a notice that the article is not published yet

#### Scenario: Published article has no notice
- **WHEN** a reader opens the URI of an article whose publish timestamp is in the past
- **THEN** the page SHALL NOT show the not-published notice

### Requirement: Markdown block rendering

Markdown blocks SHALL be rendered through the existing sanitised markdown pipeline, with relative image references resolved through the asset CDN.

#### Scenario: Safe rendering
- **WHEN** a markdown block contains embedded HTML or unsafe link protocols
- **THEN** the rendered output SHALL be sanitised
- **THEN** external links SHALL open with the established external-link attributes

#### Scenario: Resized images
- **WHEN** a markdown block contains a relative image reference
- **THEN** the image SHALL be served through the CDN resizing helper with a responsive source set

### Requirement: Image block rendering

An image block SHALL display its gallery asset at a readable size and SHALL show its caption, when present, using the same presentation as a markdown blockquote. A lone image block SHALL occupy the full content width; consecutive image blocks SHALL be presented as one image grid rather than stacked full-width blocks.

#### Scenario: Gallery asset displayed
- **WHEN** an article contains an image block
- **THEN** the image SHALL be loaded from the gallery collection and filename in the block
- **THEN** the caption SHALL be shown when present

#### Scenario: Caption presentation
- **WHEN** an image block has a caption
- **THEN** the caption SHALL use the markdown blockquote presentation, including its filled background, rounded corners, and centered secondary text

#### Scenario: Caption absent
- **WHEN** an image block has no caption
- **THEN** the block SHALL render without caption space

#### Scenario: Lone image keeps full width
- **WHEN** an image block is not adjacent to another image block
- **THEN** the image SHALL occupy the full content width

#### Scenario: Asset no longer available
- **WHEN** the referenced gallery asset cannot be loaded
- **THEN** the article SHALL continue rendering the remaining blocks

### Requirement: Image sequence grouping

The article page SHALL present consecutive images from the same source as one grid. Markdown images group only with markdown images; gallery image blocks group only with gallery image blocks. Adjacent images of different sources MUST NOT share a grid, and a run of one image SHALL keep its full-width presentation. Non-image content, including text, headings, and horizontal rules, SHALL end a markdown run.

#### Scenario: Consecutive gallery image blocks
- **WHEN** two or more image blocks are adjacent
- **THEN** they SHALL render as one image grid
- **THEN** each image SHALL load from its own gallery asset

#### Scenario: Grid geometry matches markdown
- **WHEN** images render in a grid
- **THEN** the grid SHALL use the same layout as markdown image sequences, including the same cell sizing and the same handling of the final row

#### Scenario: Consecutive markdown image paragraphs
- **WHEN** a markdown block contains two or more image-only paragraphs in sequence, including when they are separated by blank lines
- **THEN** their images SHALL render as one image grid

#### Scenario: Markdown images in one paragraph
- **WHEN** a markdown paragraph contains multiple images
- **THEN** those images SHALL render as one image grid

#### Scenario: Text ends a markdown run
- **WHEN** text, a heading, or a horizontal rule appears between markdown images
- **THEN** the images before and after that content SHALL NOT share a grid

#### Scenario: Mixed sources stay separate
- **WHEN** a gallery image block is adjacent to a markdown image
- **THEN** the two images SHALL NOT share a grid

#### Scenario: Failure inside a grid
- **WHEN** one image in a grid cannot be loaded
- **THEN** the remaining images SHALL stay in the grid and the article SHALL continue rendering

### Requirement: Music block rendering

A music block SHALL render a playable music player using the shared player component.

#### Scenario: Play a song
- **WHEN** a reader presses play on a music block
- **THEN** the referenced song SHALL play through the shared player
- **THEN** the displayed title and artist SHALL come from the block's stored values

#### Scenario: Several music blocks
- **WHEN** an article contains more than one music block
- **THEN** each block SHALL control only its own playback

#### Scenario: Player readiness
- **WHEN** a music block's embedded player signals readiness
- **THEN** the play control SHALL become enabled
- **THEN** play, pause, and seek commands SHALL be delivered to that player

### Requirement: Read-only drawing rendering

An excalidraw block SHALL render its stored scene in the embedded drawing surface in view mode. Readers MUST NOT be able to modify the drawing or affect the author's locally stored drawing.

#### Scenario: Scene displayed read-only
- **WHEN** a reader views an article with a drawing block
- **THEN** the scene SHALL be loaded into the drawing surface
- **THEN** editing tools SHALL be unavailable and changes SHALL NOT be persisted

#### Scenario: Read-only from first paint
- **WHEN** a reader opens an article with a drawing block
- **THEN** the drawing surface SHALL be in view mode from its initial render and remain read-only after the scene loads

#### Scenario: Save shortcut not captured
- **WHEN** a reader presses Ctrl+S or Cmd+S while a public drawing has focus
- **THEN** the drawing surface SHALL NOT persist or capture the scene

#### Scenario: Reader's local drawing untouched
- **WHEN** a reader views any drawing block
- **THEN** the reader's locally stored drawing SHALL NOT be read, replaced, or autosaved

#### Scenario: Scene load failure
- **WHEN** the scene file is missing or invalid
- **THEN** the block SHALL show a failure state
- **THEN** the rest of the article SHALL continue rendering

#### Scenario: Several drawing blocks
- **WHEN** an article contains more than one drawing block
- **THEN** each drawing surface SHALL load only its own scene

### Requirement: Embedded drawing scroll safety

Embedded drawings SHALL NOT capture wheel or trackpad scrolling until the reader activates them.

#### Scenario: Page scroll over a drawing
- **WHEN** the reader scrolls with the pointer over a drawing that has not been activated
- **THEN** the page SHALL scroll and the drawing SHALL NOT pan or zoom

#### Scenario: Activate a drawing
- **WHEN** the reader clicks the drawing
- **THEN** the drawing SHALL become interactive until the pointer leaves it
