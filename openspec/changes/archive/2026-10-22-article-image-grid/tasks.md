# Tasks

## 1. Shared image presentation

- [x] 1.1 Consolidate the caption presentation in `vue-web/src/assets/main.css` under a shared selector list for `.markdown blockquote` and `.article-caption`, keeping `mb-4` only for the blockquote context; verify a markdown blockquote still renders filled and rounded and `cd vue-web && pnpm type-check` passes
- [x] 1.2 Apply the shared caption class to the gallery `figcaption` in `vue-web/src/components/article/blocks/ImageBlock.vue`; verify a gallery image with a caption shows the same filled, rounded, centered treatment as a markdown blockquote and a captionless image renders without extra space
- [x] 1.3 Add the `.image-grid` rules to `main.css` after the markdown rules, covering direct `<img>` cells and `<figure>` cells, and verify in the browser that a merged cell renders at 32% width rather than full width (specificity check)

## 2. Gallery image grouping

- [x] 2.1 Add an `ImageGrid` component under `vue-web/src/components/article/blocks/` that renders one `figure` per image with the shared caption treatment; verify a gallery image inside it displays from its collection and filename and shows the unavailable placeholder when the asset cannot load
- [x] 2.2 Compute display units in `vue-web/src/components/article/BlockRenderer.vue`, grouping runs of consecutive `image` blocks and rendering a run of one as a lone `ImageBlock`; verify three adjacent image blocks render as one three-per-row grid while a separated image stays full width
- [x] 2.3 Put each grid figure's block id on the figure and none on the grid wrapper; verify the editor preview still aligns to grouped images while scrolling the editor, and that mixed caption/no-caption rows render without layout errors
- [x] 2.4 Verify grid geometry across run lengths; confirm two images split the row, three images sit three-per-row, and a fourth stretches across the final row, matching the markdown path

## 3. Markdown image merging

- [x] 3.1 Add the post-sanitisation merge pass in `vue-web/src/components/shared/MarkdownDisplay.vue`, replacing runs of two or more image-only paragraphs with one `image-grid` div holding their images in order; verify images separated by blank lines render as one grid while a run of one stays full width
- [x] 3.2 Verify the run boundaries; confirm text, a heading, and a horizontal rule each split a run, and confirm a linked image (`[![alt](src)](url)`) stays outside the grid and keeps its link
- [x] 3.3 Verify image attributes survive the pass; confirm `src`, `srcset`, `alt`, `title`, and the lazy-loading attribute are unchanged on merged images, and that the sanitised output still has no injected markup
- [x] 3.4 Verify the two sources stay separate; confirm a gallery image block adjacent to a markdown image renders as two independent units, never one grid

## 4. Integration verification

- [x] 4.1 Run `cd vue-web && pnpm type-check` and confirm it passes
- [x] 4.2 Run `cd vue-web && pnpm build` and confirm it succeeds
- [x] 4.3 Open an article containing markdown image runs, consecutive gallery blocks, a lone gallery image, and a mixed adjacency; confirm each spec scenario in `specs/article-view/spec.md` is observable, and check the same content in the editor preview

## Workflow follow-up

- Archive the change after review.
- Verify the archived spec contains both the modified image block requirement and the new image sequence grouping requirement.
