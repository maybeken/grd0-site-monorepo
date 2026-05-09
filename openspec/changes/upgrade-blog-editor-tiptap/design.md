## Context

The blog editor at `/blog/editor` consists of `BlogEditorView.vue` (post table + editor host) and `SideBySideEditor.vue` (OverType markdown editor + live preview). All user interactions use browser-native `confirm()`/`alert()`/`prompt()` dialogs. No UI component library is in use — just Tailwind CSS and `@iconify/vue` icons. The backend stores content as raw markdown in SQLite and serves images via S3 presigned URLs through Cloudflare CDN.

Phase 1 scope is narrow: replace the editor component and add basic modal/toast UX. No backend changes, no layout restructuring, no taxonomy or media library.

## Goals / Non-Goals

**Goals:**
- Replace OverType with Tiptap (`@tiptap/vue-3`) providing block-based editing, slash commands, bubble formatting menu, and native image drag-and-drop
- Serialize editor content back to markdown for backend compatibility
- Replace all `confirm()`/`alert()`/`prompt()` calls with proper modal dialogs and toast notifications
- Spike the markdown round-trip first to verify fidelity before committing fully

**Non-Goals:**
- Admin sidebar/navigation layout
- Media library page (browse previously uploaded images)
- Categories, tags, or any taxonomy system
- SEO metadata fields
- Multiple author support
- Rich-text (HTML) storage format
- Backend API changes
- Public-facing blog rendering changes (MarkdownDisplay.vue stays)

## Decisions

### 1. Tiptap over Milkdown / Editor.js / Toast UI Editor

**Decision:** Use Tiptap (`@tiptap/vue-3`).

**Rationale:**
- **Tiptap**: First-class Vue 3 support, ProseMirror foundation, block-based model, extensible via extensions, active ecosystem, slash commands and bubble menu are built-in extensions, supports markdown import/export via `prosemirror-markdown`.
- **Milkdown**: Also ProseMirror-based with native markdown awareness, but Vue support is community-maintained and less mature. React-first library.
- **Editor.js**: Clean block model and is framework-agnostic, but outputs JSON blocks — would require a custom markdown serializer, adding complexity and risk.
- **Toast UI Editor**: Has markdown + WYSIWYG modes with a Vue wrapper, but lacks the Notion/Medium block-editor feel (no slash commands, no block reordering).

**Alternatives considered:** Milkdown (weaker Vue support), Editor.js (no native markdown output), Toast UI (not block-based), sticking with OverType (doesn't solve the core UX problem).

### 2. Markdown serialization via `tiptap-markdown` / `prosemirror-markdown`

**Decision:** Use the `tiptap-markdown` package (which wraps `prosemirror-markdown`) for both parsing markdown into Tiptap's ProseMirror JSON model and serializing back to markdown.

**Rationale:**
- `prosemirror-markdown` is the authoritative ProseMirror markdown parser/serializer, maintained by the ProseMirror author.
- `tiptap-markdown` wraps it as Tiptap extensions, providing a clean integration.
- Both parse and serialize are configurable — we can control serializer preferences (e.g., always use `**bold**` not `__bold__`).
- The spike will validate round-trip fidelity. If issues arise, the serializer can be customized.

### 3. Modal and Toast as shared composables + components

**Decision:** Build lightweight, dependency-free modal and toast utilities using Vue 3 `<Teleport>` + composables + Tailwind CSS, rather than pulling in a UI library.

**Rationale:**
- No UI library currently in the project. Adding one (e.g., `@headlessui/vue`, `radix-vue`) for just 2 components adds unnecessary weight.
- The modal pattern is: `useModal()` composable returns `{ open, close, isOpen }`, wrapped in a `<Modal>` component with `<Teleport>`.
- The toast pattern is: `useToast()` composable with a global `<ToastContainer>` mounted once, exposing `toast.success()`, `toast.error()`, `toast.info()`.
- Both are small (~100 lines each), self-contained, and reusable across the entire app (Map.vue and DotMatrix.vue also have `prompt()`/`alert()` calls).

**Alternatives considered:** `vue-sonner` (extra dependency for simple needs), `@headlessui/vue` (heavy for just modals), custom `confirm()` wrapper (doesn't improve UX).

### 4. Image upload integration

**Decision:** Wire Tiptap's native image extension drag-and-drop handler to the existing `requestBlogAttachmentRaw()` + `s3PresignedUrlUpload()` flow. Tiptap's `handleDrop` callback receives the file, uploads via S3 presigned URL, then inserts the image node.

**Rationale:**
- Current image flow works well: drag file → request presigned URL from API → PUT file to S3 → insert `![](/path)` markdown.
- Tiptap natively supports images via `@tiptap/extension-image`. The existing `requestBlogAttachmentRaw` service function returns `{ url: presignedUrl, key: path }`.
- No changes needed to the upload API or S3 bucket — same `PUT /blog/attachment/:key` endpoint.

### 5. Side-by-side preview preserved

**Decision:** Keep the 50/50 split layout with `MarkdownDisplay.vue` in the right pane. Tiptap serializes content to markdown on every change (debounced), which feeds into the existing preview pipeline.

**Rationale:**
- `MarkdownDisplay.vue` uses unified/remark/rehype for rendering — the exact same pipeline as the public blog. This provides an accurate preview of how the post will appear to readers.
- Tiptap's real-time serialization (debounced ~300ms) is efficient enough for live preview.
- Removing the preview would lose a valuable authoring feature and isn't warranted by Phase 1 scope.

### 6. Editor initialization and watch sync

**Decision:** On mount, parse the post's markdown content into Tiptap's internal model using the markdown parser. When the post data changes (user selects a different post in PostTable), re-initialize the editor content using a watch.

**Rationale:**
- Current pattern: OverType gets mounted to a ref div, and a `watch` syncs `editable_content` to the editor via `.setValue()`.
- Tiptap equivalent: the `<TiptapEditor>` component accepts `content` as a prop/model. A `:key` on the editor forces re-creation when switching posts (already done in `BlogEditorView.vue` via `:key="selected"`).

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| **Markdown round-trip fidelity**: Serializer may normalize formatting (e.g., `__bold__` → `**bold**`, list whitespace changes, reference links → inline). | Spike first with real post content. Accept minor normalizations. If a severe issue is found, customize the serializer or fall back to a different approach. |
| **Tiptap bundle size**: Tiptap + starter-kit + extensions can add ~150KB gzipped to the bundle. | Monitor with `pnpm run size` (vite-bundle-visualizer). Acceptable for an admin-only page; code-split if needed. |
| **ProseMirror markdown parsing edge cases**: Tables with alignment, nested blockquotes, HTML mixed in markdown may not parse cleanly. | Test with existing posts during the spike. The unified/remark pipeline for the preview is robust — the editor only needs to parse what the author writes, not arbitrary markdown. |
| **State sync between editor and preview**: If the serializer is too slow, preview lag becomes noticeable. | Debounce serialization at 300ms. ProseMirror serialization is synchronous and fast for typical post sizes (<50KB markdown). |
| **Modal/toast accessibility**: Custom-built modals may miss ARIA attributes, focus trapping, or keyboard navigation. | Implement focus trapping, Escape key dismissal, `role="dialog"`, `aria-modal`, and focus restoration per WAI-ARIA dialog pattern. |

## Migration Plan

1. **Spike first** (Task 1): Install Tiptap in a spike component, load real posts, test round-trip fidelity.
2. **Create shared utilities**: modal composable/component and toast composable/component.
3. **Replace OverType with Tiptap** in `SideBySideEditor.vue`.
4. **Replace browser dialogs** with modals/toasts in `SideBySideEditor.vue` and `BlogEditorView.vue`.
5. **Test with existing posts**: Open each post in the new editor, verify preview matches public render, re-save.

**Rollback**: Since the backend and storage format are unchanged, rolling back is as simple as reverting to the previous `SideBySideEditor.vue` and removing Tiptap deps. Posts saved through the new editor remain valid markdown.

## Open Questions

- **Spike result**: Does the markdown round-trip produce acceptable fidelity? If not, what customization is needed?
- **Custom blocks**: Should Tiptap include custom node types (e.g., callout blocks) in Phase 1, or stick to standard markdown elements?
- **Editor chrome**: Should the OverType toolbar (bold, italic, heading, etc.) be replaced with Tiptap's bubble menu only, or keep a fixed toolbar?
