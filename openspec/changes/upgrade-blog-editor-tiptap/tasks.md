## 1. Spike: Tiptap markdown round-trip validation

- [x] 1.1 Install Tiptap dependencies: `@tiptap/vue-3`, `@tiptap/starter-kit`, `@tiptap/pm`, `@tiptap/extension-placeholder`, `@tiptap/extension-image`, `tiptap-markdown`
- [x] 1.2 Create spike component at `src/components/editor/TiptapSpike.vue` with a Tiptap editor instance, markdown parse on load, and serialized output display
- [x] 1.3 Load a real existing blog post (fetch via API or hardcode) and parse into Tiptap
- [x] 1.4 Display side-by-side diff: original markdown vs serialized output (with no edits)
- [x] 1.5 Edit the content and verify serialized output remains valid markdown
- [x] 1.6 Test edge cases: tables, nested lists, code blocks, images, links, blockquotes, GFM extensions
- [x] 1.7 Document findings: acceptable normalizations vs deal-breakers. If fidelity is poor, customize the serializer or re-evaluate approach
- [x] 1.8 Remove spike component after validation

## 2. Shared modal and toast utilities

- [x] 2.1 Create `src/composables/useModal.ts` composable with `open()`, `close()`, `isOpen` ref, promise-based confirm
- [x] 2.2 Create `src/components/Modal.vue` component with `<Teleport to="body">`, backdrop overlay, focus trap, Escape to close, ARIA dialog attributes
- [x] 2.3 Create `src/composables/useToast.ts` composable with `success()`, `error()`, `info()` methods, auto-dismiss timers
- [x] 2.4 Create `src/components/ToastContainer.vue` component with positioned toast stack, enter/leave transitions, manual dismiss
- [x] 2.5 Mount `ToastContainer` once at app root (e.g., in `App.vue`)

## 3. Replace OverType with Tiptap in SideBySideEditor

- [x] 3.1 Remove `overtype` import and `new OverType(...)` initialization from `SideBySideEditor.vue`
- [x] 3.2 Add Tiptap editor component with extensions: StarterKit, Placeholder, Image, and markdown parse/serialize via `tiptap-markdown`
- [x] 3.3 Configure slash command menu using `@tiptap/suggestion` extension
- [x] 3.4 Configure bubble formatting menu using `@tiptap/extension-bubble-menu`
- [x] 3.5 Initialize editor content from loaded post markdown on mount and on post change (via watch or `:key`)
- [x] 3.6 Extract serialized markdown on content change (debounced 300ms) to feed the `MarkdownDisplay` preview pane
- [x] 3.7 Remove `ot_editor` ref and `watch([editable_content])` sync logic
- [x] 3.8 Remove OverType CSS references (if any) from the component or global styles

## 4. Wire up drag-and-drop image upload to Tiptap

- [x] 4.1 In Tiptap's `handleDrop` callback, intercept file drops and call `requestBlogAttachmentRaw(file.name)`
- [x] 4.2 Upload the file to S3 via `s3PresignedUrlUpload(res.url, file)`
- [x] 4.3 Insert the uploaded image into the editor at cursor/drop position using Tiptap's image node
- [x] 4.4 Show image placeholder during upload progress
- [x] 4.5 Handle upload errors with a toast notification
- [x] 4.6 Remove the old `@drop.prevent="fileUpload"` handler on the outer wrapper div

## 5. Replace browser dialogs with modals and toasts

- [x] 5.1 Replace `confirm('Are you sure to create/update the post?')` in `savePost()` with modal confirm
- [x] 5.2 Replace `alert("Done.")` / `alert("Failed.")` in `savePost()` with success/error toasts
- [x] 5.3 Replace `confirm('Are you sure to create new post?')` + `prompt('What is the new post uri?')` in `newPost()` with a modal containing a URI input field
- [x] 5.4 Replace `confirm("Are you sure to unpublish the post?")` in `BlogEditorView.vue` `unpublishPost()` with modal confirm
- [x] 5.5 Replace `alert("Done.")` / `alert("Failed.")` in `unpublishPost()` with success/error toasts
- [x] 5.6 Ensure save button shows loading state during API call (existing `loading` ref)

## 6. Verify and clean up

- [x] 6.1 Open each existing blog post in the editor, verify content renders correctly and preview matches public view
- [x] 6.2 Re-save each post through the new editor and verify the public blog still renders correctly
- [x] 6.3 Test create-new-post flow: click New, enter URI in modal, write content, save, verify appears in PostTable and public blog
- [x] 6.4 Test unpublish flow: click Delete, confirm in modal, verify post status changes to Unpublished
- [x] 6.5 Run `pnpm type-check` in vue-web and fix any type errors
- [x] 6.6 Run `pnpm lint` in vue-web and fix any lint errors
- [x] 6.7 Run `pnpm build` in vue-web and verify no build errors
