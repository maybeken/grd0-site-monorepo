# Proposal

## Why

The `/editor` pages were built one section at a time, and each re-invented its own page header, spacing, form controls, buttons, loading and empty states, and vocabulary. Authors now see "Save" next to "Save changes" and "Create", "Delete" next to "Remove", "Failed to …" next to "Could not …", and a different heading treatment on nearly every page. The inconsistency slows authoring and makes the console feel unfinished.

## What Changes

- Introduce a small set of shared editor UI primitives under `src/components/editor/ui/`: page header (title, description, actions, back link), section card, form field/label, text input, select, textarea, button variants (primary, secondary, destructive, icon), loading state, empty state, and inline error banner.
- Adopt one page skeleton for every `/editor` list page and edit page: consistent wrapper spacing, header composition, action placement, content width, and back navigation.
- Replace the current three loading treatments (`Skeleton`, `Loading...` paragraph, nothing) with a single loading state, and the duplicated empty-state markup with one component.
- Unify terminology across all editor pages: "New" for create actions, "Save" for persistence ("Saved" is status text, not a button), "Delete" for destructive actions, "Remove from playlist" for playlist membership, one error family, consistent sentence casing, punctuation, required markers, ellipses, and pluralized counts.
- Raise readability: normal text meets a 4.5:1 contrast ratio, muted and placeholder text stop using the failing `text-foreground/40` recipe, the editor stays on Doto but raises body and control text to 18px at weight 500 with 14px metadata, multi-line text gets a readable line height and measure, controls get 3:1 boundaries, and every control gets a visible keyboard focus indicator.
- Fix small UX gaps surfaced by the inconsistency review: missing loading states on Beans/Equipment/Tastings/Playlists, Ctrl/⌘+S on every editable form, accurate confirmation text for catalogue deletion, and success/failure feedback for actions that are currently silent.
- Extend the translation catalogue for the new shared strings, remove duplicate/dead editor keys, and add plural forms for count-based labels.

## Capabilities

### New Capabilities

- `editor-ui`: the shared visual and verbal language of the editor console — page skeleton, shared controls, loading/empty/error states, button and form conventions, terminology, feedback, and keyboard conventions.

### Modified Capabilities

- `editor-console`: list and edit pages must render through the shared page header and page skeleton, so navigation, headings, and primary actions are in the same place in every section.
- `ui-string-extraction`: editor UI strings must follow the shared terminology, use plural forms for counts, and leave no hardcoded editor strings (including `IDK`, `mm:ss`, and separator glyphs).

## Impact

- `vue-web/src/components/editor/` — new `ui/` primitives; existing `EditorLayout`, `EditorEntityRow`, `ArticleList`, `ArticleSidebar`, `BlockItem`, `BlockPalette`, pickers, and block editors updated to use them.
- `vue-web/src/views/editor/*` — all list and edit views refactored onto the shared primitives and wording.
- `vue-web/src/components/coffee/*`, `music/*`, `map/*` — forms and pickers updated for shared controls, wording, pluralization, and save shortcuts.
- `vue-web/src/locales/en.yaml` — new `editor.ui.*`/`common.*` keys, terminology rewrites, plural forms, removal of duplicate and dead editor keys.
- `openspec/specs/editor-ui/spec.md` (new), `openspec/specs/editor-console/spec.md`, `openspec/specs/ui-string-extraction/spec.md`.
- No API, dependency, routing, authentication, persistence, or data-model changes. The public site is untouched.
