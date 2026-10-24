# Tasks

## 1. Foundations

- [x] 1.1 Create `src/components/editor/ui/EditorPage.vue`, `EditorPageHeader.vue`, and `EditorSection.vue` with the shared skeleton, header composition, and card recipe from design.md; verify `pnpm type-check` passes in `vue-web`.
- [x] 1.2 Create `EditorField.vue`, `EditorInput.vue`, `EditorSelect.vue`, and `EditorTextarea.vue` with label/required/error wiring and shared control styles; verify `pnpm type-check` passes.
- [x] 1.3 Create `EditorButton.vue` with primary, secondary, destructive, and icon variants, requiring an accessible name for icon-only buttons; verify `pnpm type-check` passes and the icon variant renders `aria-label` plus `title`.
- [x] 1.4 Create `EditorLoading.vue`, `EditorEmptyState.vue`, and `EditorError.vue`; verify `pnpm type-check` passes and each renders the shared class recipe.
- [x] 1.5 Add `useSaveShortcut.ts` and a shared `confirmDiscard()` helper; verify `pnpm type-check` passes and the composable removes its listener on unmount.
- [x] 1.6 Add the canonical `common.*` and `editor.ui.*` keys, convert one count label (gallery assets) to a vue-i18n plural pipe, and verify `pnpm build` succeeds and the label renders "1 asset" and "2 assets".
- [x] 1.7 Refactor `EditorEntityRow.vue` onto `EditorButton` and the shared muted-text recipe; verify bean/equipment/tasting/playlist/map rows still render Edit and Delete actions.
- [x] 1.8 Apply the readability rules to the primitives: page-title size that beats the global unlayered `h1`, 18px body/label/control/button text at weight 500, 14px metadata, `leading-relaxed`, `placeholder:text-foreground/60`, visible focus outline, 3:1 control boundaries, and light text on the accent primary button; verify the computed `h1` size and the muted, border, and primary-button contrast in DevTools.
- [x] 1.9 Ensure list rows and headers wrap or truncate long values with a full-value `title`; verify a long article title and map slug do not overlap adjacent content.

## 2. List pages

- [x] 2.1 Migrate `EditorArticleListView.vue` and `ArticleList.vue` to `EditorPage`/`EditorPageHeader`, removing the duplicate "Articles" heading and moving New Article into the header action slot; verify the page renders one `h1` and draft creation still opens the edit route.
- [x] 2.2 Migrate `EditorGalleryListView.vue` to the shared header and states, including the shared empty state for collections and directories; verify the page renders one `h1`, the loading state precedes content, and publishing still updates the listing.
- [x] 2.3 Migrate `EditorBeanListView.vue`, `EditorEquipmentListView.vue`, and `EditorTastingListView.vue` to the shared header, add the missing loading states, and add success feedback for deletes; verify each list shows loading then rows or the shared empty state.
- [x] 2.4 Migrate `EditorPlaylistListView.vue` and `EditorSongListView.vue` to the shared header and states, add the missing playlist loading state, and keep inline song editing; verify each list renders headers, loading, and empty states consistently.
- [x] 2.5 Migrate `EditorMapLocationListView.vue` to the shared header and states and render loading above the list; verify the create form opens below the header, slug collision still reports an error, and loading no longer appears after the rows.

## 3. Edit pages

- [x] 3.1 Update the article workspace header and action area to use `EditorPageHeader` with a screen-reader-only `h1` and the shared button variants; verify the visible editable title, save state, conflict banner, and Ctrl/⌘+S still work.
- [x] 3.2 Migrate `EditorGalleryView.vue` to the shared header/actions and consolidate "Save metadata" and "Save changes" into one Save action, keeping the pending-count affordance; verify the collection saves, publishes/unpublishes, and shows the shared not-found/error state.
- [x] 3.3 Migrate `BeanForm.vue`, `EquipmentForm.vue`, and `TastingForm.vue` to the shared header, fields, and buttons, and correct the `idkFields.grider` typo; verify each form saves and the IDK toggle clears the grinder value.
- [x] 3.4 Migrate `EditorBeanView.vue`, `EditorEquipmentView.vue`, and `EditorTastingView.vue` to the shared back link and not-found state; verify direct edit routes load items and unknown ids show the shared error with a list link.
- [x] 3.5 Migrate `PlaylistForm.vue` and `MusicEntryForm.vue` to the shared controls and terminology, and change playlist removal to unselect the song while leaving the catalogue entry intact; verify removing a song from a playlist does not delete the song from the catalogue.
- [x] 3.6 Migrate `EditorPlaylistView.vue` and the inline song editor to the shared back link, loading, and not-found states; verify direct routes load and unknown ids show the shared error.
- [x] 3.7 Migrate `EditorMapLocationView.vue` to the shared header, sections, fields, and buttons; verify save, publish-state changes on save, slug display, and coordinate validation still behave as specified.

## 4. UX repairs

- [x] 4.1 Apply `useSaveShortcut` to every editable editor form, including the coffee forms; verify Ctrl/⌘+S saves in each form and does nothing when the form is unchanged or invalid.
- [x] 4.2 Wire `confirmDiscard()` into Cancel and close handlers on dirty forms and pickers; verify cancelling with edits prompts, confirming discards, and declining keeps the edits without a second route-guard prompt.
- [x] 4.3 Make each failure report exactly once and give silent mutations success feedback; verify a failed load shows only the inline error, a failed mutation shows only a toast, and gallery publish plus coffee deletes show a toast.
- [x] 4.4 Render the gallery collection not-found state instead of falling through to the generic branch; verify an unknown collection directory shows the shared error and a back link.
- [x] 4.5 Make every text input and textarea span the full width of its form container, including the article URI field, the YouTube link field, and the brew-time field; verify no text field carries a width cap or a `max-w-prose` limit.

## 5. Wording and translation cleanup

- [x] 5.1 Rewrite editor `msgstr` values to the shared terminology (New/Save/Saved/Delete/Remove from playlist, one error family, sentence case, `…`, punctuation); verify the affected pages show the new wording.
- [x] 5.2 Convert all counted editor labels to plural pipes; verify each label renders singular at one and plural otherwise.
- [x] 5.3 Extract the remaining hardcoded editor strings (`IDK`, `mm:ss`, separators, `★`, `—`, `/article/` prefixes) and remove the global `t` import from editor components; verify type-check passes.
- [x] 5.4 Merge duplicate editor keys, remove dead editor keys, and replace `text-secondary` with `text-foreground/60` in editor files; verify a search finds no `text-secondary` in `src/components/editor` or `src/views/editor` and type-check passes.
- [x] 5.5 Replace remaining low-contrast editor text (`text-foreground/30` and `/40`), `text-[10px]` pills, and sub-14px labels; verify muted text and placeholders meet 4.5:1, body/controls render at 18px/500, and no editor text renders below 14px.

## 6. Verification

- [x] 6.1 Run `cd vue-web && pnpm type-check` and confirm it passes with no errors.
- [x] 6.2 Run `cd vue-web && pnpm build` and confirm the production build succeeds.
- [x] 6.3 Manually check every `/editor` list and edit page at a wide and a narrow viewport for one visible `h1`, consistent header/action placement, loading/empty states, and no duplicate headings.
- [x] 6.4 Confirm the `/blog/editor` page and all public routes are unchanged.
- [x] 6.5 Run a manual readability audit on one list and one edit page per section: check body/muted/placeholder contrast, keyboard focus visibility, the 18px/500 floor, line length, and long-value truncation at both viewports.

## Workflow follow-up

- Archive the change after review and verification are complete.
