# Spec Delta

## MODIFIED Requirements

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

## ADDED Requirements

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
