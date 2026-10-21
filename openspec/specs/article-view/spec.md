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

The public article page SHALL resolve a published article by URI and render its header and blocks in stored order.

#### Scenario: Published article resolves
- **WHEN** a reader opens a URI whose article is published
- **THEN** the page SHALL show the title, subtitle, author, and timestamps
- **THEN** the blocks SHALL be rendered in their stored order

#### Scenario: Unknown or unpublished URI
- **WHEN** a reader opens a URI with no published article
- **THEN** the page SHALL show the not-found state

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

An image block SHALL display its gallery asset at a readable size, with its caption when present.

#### Scenario: Gallery asset displayed
- **WHEN** an article contains an image block
- **THEN** the image SHALL be loaded from the gallery collection and filename in the block
- **THEN** the caption SHALL be shown when present

#### Scenario: Asset no longer available
- **WHEN** the referenced gallery asset cannot be loaded
- **THEN** the article SHALL continue rendering the remaining blocks

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
