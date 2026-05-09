## Why

The blog post editor (`/blog/editor`) uses OverType — a basic markdown editor — with raw browser `confirm()`/`alert()`/`prompt()` dialogs for all interactions. The editing experience feels amateur compared to modern CMS tools. Upgrading to a Notion/Medium-style block editor with proper modal/toast UX brings the editor in line with user expectations and improves authoring speed.

## What Changes

- **Replace OverType with Tiptap block editor** (`@tiptap/vue-3`). Tiptap provides a ProseMirror-powered block-based editing experience with slash commands, bubble menus, and native image handling. Content continues to be stored as markdown in the backend — Tiptap parses markdown on load and serializes back to markdown on save.
- **Add modal and toast utilities**. Replace all `confirm()`/`alert()`/`prompt()` calls across the blog editor (and reusable elsewhere) with proper styled modal dialogs and toast notifications.
- **Spike Tiptap markdown round-trip first**. Before committing fully, verify that parsing existing posts into Tiptap and serializing back produces semantically equivalent markdown. Minor formatting normalizations (e.g., `**` over `__`, whitespace) are acceptable. If fidelity is poor, alternate approaches will be evaluated.
- **One-time content migration**. Existing posts will be opened and re-saved through the new editor, accepting minor markdown normalizations from the serializer.

## Capabilities

### New Capabilities

- `tiptap-editor`: Block-based rich text editor with slash commands, bubble formatting menu, drag-and-drop image upload, and markdown import/export. Replaces OverType in SideBySideEditor.
- `modal-toast-notifications`: Reusable modal dialog and toast notification components. Used initially by the blog editor but designed as shared utilities for future use.

### Modified Capabilities

_None._ No existing specs to modify. Backend API and data model remain unchanged (markdown storage).

## Impact

- **vue-web**: `SideBySideEditor.vue`, `BlogEditorView.vue` — replaced editor and workflow logic. New dependencies: `@tiptap/vue-3`, `@tiptap/starter-kit`, `@tiptap/pm`, `@tiptap/extension-placeholder`, `@tiptap/extension-image`, `tiptap-markdown`. New shared utilities: modal and toast composables/components.
- **api**: No changes. Markdown content field, `PUT /blog/:uri`, and `PUT /blog/attachment/:key` endpoints unchanged.
- **No breaking changes**. Public blog rendering (`MarkdownDisplay.vue`, unified/remark/rehype pipeline) unchanged. Existing posts remain readable.
