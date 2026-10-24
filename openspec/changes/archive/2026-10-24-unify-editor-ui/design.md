# Design

## Context

See `proposal.md` for motivation and `specs/` for the behaviour contract. The relevant current state:

- `EditorLayout.vue` provides only the shell. Every list and edit page composes its own header, spacing, cards, controls, buttons, loading and empty states with copy-pasted Tailwind recipes.
- The visual vocabulary in `src/assets/main.css` is `background`, `foreground`, `shade`, `secondary`, `accent`. `secondary` currently equals `foreground`, so `text-secondary` renders at full brightness; real muted text is written as `text-foreground/40…60` in some places.
- Strings live in gettext-style `src/locales/en.yaml`, flattened by `vite-plugin-gettext-yaml-i18n.ts`, with keys equal to `msgid`. Editor pages mix `$t`, `useI18n().t`, and a global non-reactive `t` import.
- Article editing is a bespoke split workspace with scroll sync and block editing; the rest of the editors are conventional forms.
- There is no unit-test suite. Verification is `pnpm type-check` and `pnpm build` inside `vue-web`.

## Goals / Non-Goals

**Goals:**

- One page skeleton and one header composition shared by all `/editor` list and edit pages.
- One set of form controls, buttons, loading/empty/error states.
- One vocabulary for create, save, delete, remove, failures, and confirmations.
- Small UX repairs that fall out of the consistency work: loading states everywhere, Ctrl/⌘+S on every editable form, dirty-form cancel confirmation, accurate playlist removal, success feedback for silent actions.

**Non-Goals:**

- No redesign of the article block model, scroll sync, pickers, or Excalidraw integration.
- No changes to public-site visuals, routing, authentication, APIs, storage, or data models.
- No new test framework, modal system, or global theme-token changes.
- The legacy `/blog/editor` (`BlogEditorView`, `PostTable`, `SideBySideEditor`) is out of scope; it keeps its current implementation.

## Decisions

### Shared primitives live in `src/components/editor/ui/`

New components:

| Component | Responsibility |
|---|---|
| `EditorPage.vue` | Page skeleton: vertical stack, standard width (`max-w-5xl`), optional `full` width for the article workspace. |
| `EditorPageHeader.vue` | Back link, `h1`, optional description, trailing action slot; supports a screen-reader-only title. |
| `EditorSection.vue` | Card section with optional heading, shared border/background/padding. |
| `EditorField.vue` | Label, required marker, help/error slot, `for` wiring. |
| `EditorInput.vue`, `EditorSelect.vue`, `EditorTextarea.vue` | Shared control visuals, `v-model`, `id`, `required`, `aria-required`. |
| `EditorButton.vue` | `primary`, `secondary`, `destructive`, `icon` variants; `label` required for icon-only. |
| `EditorLoading.vue` | Shared skeleton loading block built on the existing `Skeleton.vue`. |
| `EditorEmptyState.vue` | Shared muted message block. |
| `EditorError.vue` | Shared inline error banner, used for load/validation failures. |

A single `useSaveShortcut(save, { enabled })` composable centralizes Ctrl/⌘+S, replacing the four ad-hoc handlers. `EditorEntityRow` is refactored onto `EditorButton` rather than replaced.

**Why components over class constants:** components keep the accessible-name, required, and disabled contracts in one place, and Vue's `<script setup>` typing catches drift at build time. Class-string helpers would leave those contracts to each call site.

### Page skeleton and header

- Skeleton: `EditorPage` renders `<div class="mx-auto flex w-full max-w-5xl flex-col gap-6">`. Article uses `full` to keep the split workspace.
- Header: `header.flex flex-wrap items-start justify-between gap-4`, left column holds the back link (`← Back to {section}`), `h1.text-2xl.font-bold` with `CursorBlink`, and `p.text-sm.text-foreground/60`; right column is the action slot.
- Edit pages pass `back-to` and a section label; list pages omit the back link and put the New button in the action slot.
- The Articles list drops its inner "Articles" heading; the page header is the only heading.
- The article workspace uses `EditorPageHeader` with actions and a screen-reader-only `h1`, because the visible title is the editable title input in the canvas. This satisfies the one-`h1` contract without duplicating the title.

### Visual recipes

- Muted text: `text-foreground/60`. Do not use `text-secondary` or `text-foreground/40` in new editor code. On `neutral-950`, 40% foreground is roughly 3.6:1 and fails normal-text contrast; 60% is roughly 6:1.
- Control: `w-full rounded border border-foreground/40 bg-background px-2.5 py-1.5 text-lg font-medium placeholder:text-foreground/60 focus:border-foreground/60 focus:outline-none focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-foreground/60`.
  - `accent` is `neutral-600`, which is about 2.6:1 against `neutral-950` even at full opacity; `border-accent/40` is effectively invisible and fails the 3:1 control-boundary threshold. Controls therefore use `foreground/40` (about 3.55:1).
- Primary button: `rounded bg-accent px-3 py-1.5 text-lg font-semibold text-foreground hover:brightness-110 disabled:opacity-50`. The current `text-background` on `accent` is only about 2.5:1; light text on the same fill is about 7.2:1 and keeps the existing grey-fill look.
- Secondary: `rounded border border-accent/40 px-3 py-1.5 text-lg font-medium hover:bg-accent/10 disabled:opacity-50`.
- Destructive: `rounded border border-red-500/40 px-3 py-1.5 text-lg font-medium text-red-400 hover:bg-red-500/10 disabled:opacity-50`.
- Icon-only: `rounded-full border border-dotted p-1 hover:bg-accent/10`.
- Section card: `rounded-xl border border-accent/20 bg-background p-4`. Section borders are decorative, not control boundaries, so the subtle accent tint stays.
- Standard widths collapse the current `max-w-2xl` / `max-w-4xl` / `max-w-5xl` split into `max-w-5xl` for every non-article page.

### Readability rules

| Role | Size / treatment |
|---|---|
| Page title (`h1`) | `text-3xl` (30px) at `md`, `text-2xl` below, `font-bold` |
| Section heading (`h2`) | `text-xl` (20px) `font-bold` |
| Body, labels, hints | `text-lg` (18px), `font-medium` (500) |
| Controls | `text-lg` (18px), `font-medium` (500) |
| Metadata, timestamps, status pills | `text-sm` (14px) minimum, `font-medium` (500) |
| Text inputs and textareas | `text-lg` (18px), `font-medium` (500), full width of the form container |

- Doto stays. It is a dot-matrix display face, so size and weight carry the legibility: editor body, labels, controls, and button text use 18px at weight 500, and no editor text falls below 14px. Avoid weights below 500 in body and controls.
- Metadata, timestamps, and status pills use 14px. The map status pill's `text-[10px]` becomes `text-sm`; row action buttons move from `text-xs` to `text-lg` like other controls.
- Multi-line text uses `leading-relaxed` (1.625). Text inputs and textareas span the full width of their form container; measure comes from the form/card width rather than a character cap.
- Placeholders use `placeholder:text-foreground/60`; browser default placeholder colour is not relied on in the dark theme.
- Focus is visible on keyboard interaction via `focus-visible:outline-2 outline-offset-2 outline-foreground/60` on buttons, links, inputs, selects, and textareas.
- Status pills always include a text label ("Published", "Draft"); colour is never the only signal.
- Long titles and paths in list rows use `min-w-0` + `truncate` (or `break-words`) and carry a `:title` with the full value so truncation never hides information.
- Global `h1` in `main.css` is unlayered (`text-2xl md:text-5xl`), so it takes precedence over utility classes. `EditorPageHeader` must define its own explicit title size that wins (for example by using a scoped class block or `text-2xl! md:text-3xl!`), and the change must verify the computed size on desktop rather than assuming the utility applies.
- Density is the trade-off: 18px body text makes rows and forms taller than today. That is accepted in exchange for Doto legibility; do not shrink labels or controls to recover density.

### States and feedback

- Loading: `EditorLoading` (three skeleton rows) replaces both the `Loading...` paragraph and the missing states. It renders before content and suppresses the empty state.
- Empty: one `EditorEmptyState` with sentence-case, sentence-final messages, e.g. "No beans yet." and "No songs in the catalogue yet."
- Failure: loads use inline `EditorError`; mutations use toasts. The same failure is never reported twice. `Failed to …` and `Could not …` collapse into one family: "Could not … .".
- Success: mutations toast unless the page visibly reflects the result; gallery publish and bean/equipment/tasting delete gain feedback.

### Terminology and catalogue changes

- Create: `New Article`, `New Bean`, `New Equipment`, `New Tasting`, `New Playlist`, `New Song`, `New Location`. The `+` glyph becomes the leading icon (`mynaui:plus`), not part of the text.
- Save: `common.save` = "Save", `common.saving` = "Saving…", `common.saved` = "Saved" (status only). Gallery's "Save metadata" and "Save changes" both become "Save".
- Delete: destructive item actions use `common.delete`. Playlist membership uses `music.editor.removeFromPlaylist` = "Remove from playlist". In `PlaylistForm`, that action removes the selected id instead of deleting the catalogue entry; catalogue deletion stays in the Songs editor.
- Confirmations: one key per concept, naming the item and stating reversibility. `article.editor.confirmDiscard` merges into `common.confirmDiscard`. Legacy "Are you sure to …" phrasings are rewritten.
- Spelling: British `colour`, consistent with the existing map editor labels.
- Required fields: `*` appended to the label and `aria-required` on the control (Title, Name, and the coffee quick-add name fields).
- Units: parentheses with lower-case units, e.g. `Dose (g)`, `Water in (ml)`, `Temp (°C)`.
- Punctuation: no trailing period on labels/buttons; full sentences end with `.`; ellipses use `…`.

### Translation mechanics

- Keep existing key names where they are stable; rewrite `msgstr` values. Merge duplicate keys (`common.confirmDiscard`, the bean/equipment delete texts) onto canonical keys and delete dead editor keys.
- Plural forms use vue-i18n pipes in `msgstr`, e.g. `'{count} asset | {count} assets'`, called as `t(key, { count }, count)` so both the interpolation value and the plural branch are supplied. The gettext plugin passes `msgstr` through verbatim, so pipe syntax survives the transform.
- Editor components use `$t` in templates and `useI18n().t` in scripts. The global `t` import is removed from editor components.
- Add keys for `IDK`, `mm:ss`, `·`/`★` separators, `/article/` prefixes, and the `—` fallback, per `ui-string-extraction`.

### Small UX repairs included

- `useSaveShortcut` is applied to `BeanForm`, `EquipmentForm`, `TastingForm` (new) and keeps existing behaviour in the article, gallery, map, playlist, and song forms.
- Cancel handlers call a shared `confirmDiscard()` when the form is dirty; route-level `useUnsavedGuard` remains unchanged, so navigation does not double-prompt.
- Missing loading states are added to Beans, Equipment, Tastings, and Playlists.
- `EditorMapLocationListView` renders its loading state above the list rather than after the rows.
- `TastingForm`'s `idkFields.grider` typo is corrected while touching the IDK labels.

## Risks / Trade-offs

- **Visual regression across many pages** → migrate section by section (articles, gallery, coffee, music, travel), running `pnpm type-check` and `pnpm build` after each, and check one wide and one narrow viewport per section.
- **Doto is a dot-matrix display face** → Doto stays by decision. The 18px/500 floor is the mitigation; verify it in the readability lab and the manual pass. If it still reads poorly, revisit the typeface decision rather than silently shrinking text.
- **Global unlayered `h1` rule can override the shared header size** → verify the computed size in DevTools at `md` and above, and make `EditorPageHeader`'s title rule explicit if the utility loses.
- **Article workspace is bespoke** → limit its change to the header, action placement, states, and wording; do not touch scroll sync, block model, or save logic.
- **`max-w-5xl` everywhere reduces usable width for the map's two-column layout** → verify at `lg`; if cramped, the map keeps a two-column grid that wraps below the breakpoint rather than reintroducing a page-width variant.
- **Plural pipe syntax through the custom plugin** → validate early with one count label in gallery and music before converting the rest.
- **Dirty-cancel confirmation can double with the router guard** → one `confirmDiscard()` helper, called only from Cancel/close handlers; the router guard remains the only navigation guard.
- **Behaviour change in playlist removal (no longer deletes the catalogue entry)** → call it out in the change summary; catalogue deletion stays available in the Songs editor.
- **Large `en.yaml` diff** → change wording and structure in one pass per namespace, and let `openspec validate` plus `pnpm type-check` catch missing keys.

## Migration Plan

1. Add `editor/ui/` primitives, `useSaveShortcut`, and the new/renamed translation keys without touching pages; build to confirm the primitives compile.
2. Migrate list pages (Articles, Gallery, Beans, Equipment, Tastings, Playlists, Songs, Map Locations) onto `EditorPage` / `EditorPageHeader` / states.
3. Migrate edit pages (Article workspace header only, Gallery collection, Bean, Equipment, Tasting, Playlist, Song inline form, Map Location).
4. Apply the UX repairs and remove duplicate/dead keys.
5. Run `pnpm type-check` and `pnpm build`; manually verify each section at wide and narrow viewports.

Rollback: the primitives are additive until step 2; after that, each section's migration is an independent revert. No data or routing migration is involved.
