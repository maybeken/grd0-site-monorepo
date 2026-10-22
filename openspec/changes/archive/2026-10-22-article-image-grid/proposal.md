# Proposal

## Why

Images render inconsistently depending on how the author inserted them. Gallery image blocks show their caption as plain centered text and always occupy a full row, so a run of gallery images pushes everything else off the screen, while markdown images flow into a compact grid. The same article looks different depending on which mechanism the author used.

## What Changes

- Gallery image block captions adopt the markdown blockquote presentation, so a caption reads the same whether the image came from the gallery or from markdown.
- Consecutive gallery image blocks render as one image grid using the same geometry as markdown images; a lone gallery image keeps its full width.
- Consecutive image-only paragraphs inside a markdown block merge into that same grid, so blank-line-separated images group just like consecutive-line images; any intervening text, heading, or other content ends the run.
- Grouping is same-type only: gallery blocks group with gallery blocks and markdown images group with markdown images. A gallery block adjacent to a markdown image does not merge with it.
- A lone image in either path keeps its current full-width presentation.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `article-view`: The **Image block rendering** requirement changes so the caption uses the markdown blockquote presentation and consecutive image blocks form a grid, and a new **Image sequence grouping** requirement describes runs of markdown images and gallery images.

## Impact

- `vue-web`: `BlockRenderer` groups consecutive image blocks into a new grid component; `ImageBlock` caption styling changes; `MarkdownDisplay` gains a post-processing pass over sanitised output that merges image-only paragraphs; `main.css` gains shared image-grid and caption rules.
- Editor preview: it renders through `BlockRenderer`, so it inherits the grouping; each grouped image retains its own block identity for scroll alignment.
- No API, worker, storage, or database changes.
