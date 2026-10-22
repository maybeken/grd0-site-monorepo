# Design

## Context

See proposal.md for motivation and the spec delta for required behavior.

Two independent renderers produce images today. `BlockRenderer` walks the block list and renders each `image` block through `ImageBlock` inside its own wrapper. Markdown blocks render through `MarkdownDisplay`, which runs `remark-parse` → `remark-rehype` → `remarkBreaks` → `inspectUrls` → `rehype-sanitize` → `rehype-stringify` and assigns the result to `v-html`. The sanitizer strips the `break` nodes `remarkBreaks` emits, so consecutive-line images end up as adjacent `<img>` elements in one `<p>`.

Image layout lives entirely in `main.css`: `.markdown p:has(img)` is a wrapping flex row, and `.markdown p > img` is `grow object-cover w-[32%]`. A single image grows to full width; runs of two and three fill their row, and a fourth image in a three-per-row run stretches across the full final row. `.markdown blockquote` is the caption-like presentation the gallery caption must adopt.

`BlockRenderer` is also used by the editor's live preview, and `useScrollSync` aligns the editor and preview by measuring every `[data-block-id]` element inside the preview pane.

## Goals / Non-Goals

**Goals:**

- Make gallery and markdown images share one visual geometry and one caption presentation.
- Group image runs inside each source type without changing block storage or the API.
- Preserve editor preview fidelity and block-alignment behavior.
- Keep the existing markdown pipeline, sanitisation, and CDN image resolution intact.

**Non-Goals:**

- Merging gallery and markdown images into a single grid.
- Changing image resolution, compression, or CDN parameters.
- Changing how the editor stores or orders blocks.
- Introducing a lightbox, zoom, or caption editing in markdown.

## Decisions

**Group consecutive image blocks in `BlockRenderer`.** Compute display units before rendering: walk the block list and accumulate runs where `block.type === 'image'`. A run of two or more becomes one grid; a run of one renders as today. Alternative considered: give the images a shared parent with CSS sibling selectors. Rejected because CSS cannot count or bound an arbitrary run, and the renderer already owns block sequencing.

**Merge markdown image paragraphs in `MarkdownDisplay` after sanitisation.** Parse the sanitised HTML with `DOMParser`, walk the top-level children, and replace each run of two or more image-only `<p>` elements with one `<div class="image-grid">` holding their `<img>` nodes in order. Image-only means a `<p>` whose element children are all `<img>`, with whitespace text allowed. A lone image paragraph is left untouched so it keeps full width. Alternative considered: pure CSS with inline-flex image paragraphs. Rejected because CSS cannot distinguish a lone image from a run, so single images would shrink to 32%. Alternative considered: a remark/rehype plugin. Rejected because the merge is a presentation concern over already-sanitised output, and operating after sanitisation keeps the security boundary unchanged.

**Use one grid DOM shape per source.** Markdown grids contain `<img>` directly; gallery grids contain `<figure>` cells so each caption can sit under its image. Both use the same `image-grid` class and the same flex-wrap, growth, and cell sizing, which is what makes the two paths visually identical without sharing a component. Alternative considered: extract each markdown image into a component. Rejected because the markdown path is an HTML string rendered by `v-html`.

**Keep the literal markdown geometry.** Grid cells use `grow w-[32%]`, and images use `object-cover`, reproducing the existing three-per-row layout and the stretched final row. Alternative considered: fixed-aspect cells. Rejected because it would visibly diverge from the markdown path this change is meant to match.

**Consolidate the caption and grid styles in `main.css`.** Place the caption rule under a shared selector list (`.markdown blockquote, .article-caption`) and the grid rules adjacent to the markdown image rules, so the two paths read as one system and drift is visible in review. The gallery caption omits `mb-4`; the figure's own gap supplies the spacing.

**Keep block identity on each grid cell.** The gallery grid's wrapper carries no `data-block-id`; each `<figure>` carries its own. This preserves `useScrollSync`, which measures anchors in document order and does not care about nesting.

**Treat a linked image as not image-only.** `[![alt](src)](url)` renders an `<a>` around the `<img>` and is not merged. Alternative considered: unwrapping the anchor. Rejected because it would silently drop the author's link; leaving it unmerged preserves the link and only costs grid membership.

## Risks / Trade-offs

- **CSS specificity:** `.markdown img { w-full }` and `.image-grid > img { w-[32%] }` have equal specificity, and the merged div sits inside `.markdown`; if the grid rule comes first, images stay full width → Place the grid rules after the markdown block or scope them as `.markdown .image-grid > img`.
- **Mixed captions in a gallery row:** a captioned cell has a shorter image than an uncaptioned neighbor in the same row, because the caption consumes height → Accepted with the literal geometry; rows that are uniformly captioned or uncaptioned align correctly.
- **Cropping in gallery grids:** `object-cover` crops photographs to equalize row heights, unlike the current full-width `object-contain` → Accepted as the cost of matching markdown; a lone image is still uncropped.
- **Final-row stretch:** a fourth image in a three-per-row run spans the full row → Matches the existing markdown behavior, which this change explicitly preserves.
- **Extra DOM parse per markdown render:** `MarkdownDisplay` parses its own sanitised HTML once more → Small, local to markdown blocks, and includes no network or pipeline work.
- **Anchor elements break image-only detection:** linked images stay outside the grid and render at the old single-image size → Deliberate; documented in the spec decision above.
- **Preview regressions:** grouping changes measured heights in the preview pane → Each figure keeps its block id and the sync measures live DOM geometry, so alignment recomputes on the next measure.

## Migration Plan

Frontend-only deploy. No API, data, or configuration change. Roll back by redeploying the previous web build; stored blocks are untouched, so old and new clients render the same data differently but remain compatible. The editor preview and public page ship together because both consume `BlockRenderer`.

## Open Questions

None.
