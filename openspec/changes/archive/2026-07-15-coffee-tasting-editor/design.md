## Context

The coffee feature already has full CRUD API endpoints for beans, equipment, and tastings. The frontend has a `TastingForm` component for creating/editing tastings, but no UI for managing beans or equipment. The editor view (`CoffeeEditorView.vue`) was a minimal wrapper that only rendered `TastingForm`.

The existing `TastingForm` already follows a consistent pattern: card-style container with rounded borders, label/input pairs, and save/cancel buttons. New components should match this styling.

## Goals / Non-Goals

**Goals:**
- Provide a single-page admin interface for all coffee entity management
- Allow creating beans and equipment without leaving the tasting form
- Maintain visual consistency with existing `TastingForm` styling
- Reuse existing API service functions without backend changes
- Auto-select newly created items in the tasting form

**Non-Goals:**
- Bulk operations (import/export, multi-delete)
- Search/filter within lists
- Pagination (dataset remains small for personal use)
- Drag-and-drop reordering
- Image upload for beans or equipment
- Validation beyond required name field

## Decisions

### 1. Tabbed interface vs. separate routes

**Decision**: Use a single view with tab buttons to switch between Beans, Equipment, and Tasting sections.

**Rationale**: Keeps the editor as a single destination. Tabs are simpler than nested routes for this use case. The three sections are closely related and users may switch between them frequently when setting up new tastings.

**Alternative considered**: Separate routes (`/coffee/editor/beans`, `/coffee/editor/equipment`). Rejected — adds navigation overhead for a small admin tool.

### 2. Inline quick-add vs. modal dialogs

**Decision**: Implement inline expandable forms within `TastingForm` for quick-add of beans and equipment.

**Rationale**: Inline forms keep context visible — the user sees the tasting form while adding a bean. No modal overlay to manage. Simpler interaction: click "+ new", fill fields, click "Add", done.

**Alternative considered**: Modal dialogs. Rejected — adds complexity (focus trapping, backdrop, z-index) and obscures the tasting form context.

### 3. Form components as separate files

**Decision**: Create `BeanForm.vue` and `EquipmentForm.vue` as standalone components rather than inlining the forms in `CoffeeEditorView.vue`.

**Rationale**: Reusability — the same forms are used in both the tab list view and could be used elsewhere. Separation of concerns — each form handles its own validation, save logic, and edit-mode population.

### 4. List refresh strategy

**Decision**: After any save or delete operation, re-fetch the entire list from the API.

**Rationale**: Simple and correct. The dataset is small (personal use), so re-fetching is fast. Avoids optimistic update complexity and ensures consistency with the backend.

**Alternative considered**: Optimistic updates (add/remove from local array). Rejected — adds complexity for negligible UX gain at this scale.

### 5. Quick-add auto-select behavior

**Decision**: After successfully creating a bean or equipment via quick-add, automatically select it in the tasting form.

**Rationale**: The user just created the item because they wanted to use it. Auto-select saves a click and confirms the item was created successfully.

### 6. Inline form field subset

**Decision**: Quick-add forms expose only the most essential fields (name + roaster/origin for beans, name + type for equipment).

**Rationale**: Quick-add is for speed. Users who need to set all fields can use the full form in the respective tab. Keeps the inline form compact and non-intrusive.

## Risks / Trade-offs

### [Risk] Inline form state persistence
If the user fills a quick-add form, then cancels and re-opens it, the fields are cleared. → **Mitigation**: Acceptable behavior — cancel means "discard". Users can re-enter quickly.

### [Trade-off] No confirmation for inline quick-add
Quick-add creates items immediately without a separate confirmation step. → **Mitigation**: The "Add" button is the confirmation. Items can be deleted from the list view if created by mistake.

### [Trade-off] Full list re-fetch after changes
Re-fetching the entire list after each operation is not optimal for large datasets. → **Mitigation**: Acceptable for personal use with small datasets. Can optimize later if needed.

### [Risk] Quick-add form validation
Minimal validation (only name required) could lead to incomplete records. → **Mitigation**: Users can edit the item later from the list view to add missing fields.
