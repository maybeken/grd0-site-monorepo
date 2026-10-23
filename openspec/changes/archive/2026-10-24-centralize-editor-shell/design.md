# Design

## Context

See `proposal.md` for motivation. The constraints that shape this design come from the current code:

- `App.vue` unconditionally renders `ParticleBackground`, a sticky `Navbar`, `<main><RouterView/></main>`, and a footer. There is no layout abstraction.
- Editor authentication is implemented in each editor view's `onMounted`, while the router declares `requiresAuth` metadata that nothing reads.
- The article editor already implements the target shape: a listing route, a `:uri` workspace route, a dirty snapshot, `onBeforeRouteLeave`, and a `beforeunload` listener.
- Coffee and music editors hold the selected item in a component `ref` and render the form inline; neither form tracks dirty state.
- `TastingForm` requires `beans` and `equipmentList` props; `PlaylistForm` fetches the song catalogue itself; `BeanForm` and `EquipmentForm` need only the item.
- The Go API has list, upsert, and delete endpoints for coffee and music, but no get-by-id endpoints. Upsert handlers already return 404 for missing rows.

## Goals / Non-Goals

**Goals:**

- One persistent editor shell with a sidebar, so every editor shares layout, navigation, and authentication.
- Every editable item has an addressable URL that survives refresh and direct opening.
- One implementation each of the authentication redirect and the unsaved-change guard.
- Preserve all existing public pages, API behavior, and article editing semantics.

**Non-Goals:**

- Redesigning the forms or list rows beyond the prop and dirty-state changes needed for routing.
- Changing the legacy blog editor or unifying it with the new shell.
- Introducing a shared client cache or state manager for editor data.
- Changing upload, scene storage, or article block behavior.

## Decisions

### Route subtree under `/editor`

A single top-level layout route owns the subtree, with child routes for each list and edit page:

```
/editor                          EditorLayout   meta: requiresAuth, layout: 'editor'
  ''            -> redirect /editor/articles
  /articles                    list (?page=)
  /articles/:uri               article workspace
  /coffee       -> redirect /editor/coffee/beans
  /coffee/beans                list
  /coffee/beans/:id            edit
  /coffee/equipment            list
  /coffee/equipment/:id        edit
  /coffee/tastings             list
  /coffee/tastings/:id         edit
  /music                       list
  /music/:id                   edit
```

Legacy paths are top-level redirect routes, outside the layout: `/article/editor`, `/article/editor/:uri`, `/coffee/editor`, `/music/editor`.

**Why over preserving existing paths:** a shared parent gives the sidebar, auth guard, and chrome switch one place to live, and it removes the current requirement that `/article/editor` be registered before `/article/:slug` to avoid being swallowed by the public article route. Keeping the old paths would force a path-less layout route with absolute-path children and leave the shell straddling unrelated public paths.

**Alternative considered:** keeping `/article/editor`, `/coffee/editor`, and `/music/editor` and wrapping them in a parent with absolute child paths. Rejected for the ambiguous namespace and because the editor would still be spread across three top-level prefixes.

### Chrome switched by inherited route meta

`EditorLayout` sets `meta.layout: 'editor'` (and `meta.requiresAuth`), which Vue Router exposes on every matched child. `App.vue` reads `route.matched` and renders either the public chrome or the bare `<RouterView/>`. No new layout system is introduced.

**Alternative considered:** a `v-if` on `route.path.startsWith('/editor')`. Rejected because meta survives future path changes and expresses intent.

### New get-by-id endpoints on the public group

Add `GET /coffee/bean/:id`, `/coffee/equipment/:id`, `/coffee/tasting/:id`, and `/music/playlist/:id` to the `v1` group, mirroring the public list endpoints. Tasting preloads Bean, Equipment, and Grinder; playlist preloads songs with the existing `orderSongs`.

**Why over finding the item in the list response:** edit routes are deep-linkable by design, so a cold load must resolve one item, and a 404 needs to be distinguishable from an empty list. These resources have no draft or unpublished state, so item-by-id exposes nothing the public list does not already expose.

**Alternative considered:** fetching the section list and searching client-side. Works today because collections are small, but ties every edit page to the list payload and makes not-found detection a client-side guess.

**Alternative considered:** a shared Pinia cache primed by list pages. Rejected because no cache invalidation story exists and cold loads still need a fetch.

### Creation remains on list pages

List pages own create: an inline form, a save that refreshes the list, and no route change. Edit routes never create. Articles are the one exception required by their spec: the listing creates a draft and then navigates to the workspace route.

**Why over routable `/new`:** it matches the current implementation of all four forms, avoids an empty URL whose contents depend on whether a draft is created on mount, and keeps "create" and "delete" in one place.

### Shared composable for unsaved changes

Extract the article editor's pattern into `useUnsavedGuard(dirty)` implementing `onBeforeRouteLeave` plus a `beforeunload` listener. Each edit container computes `dirty` from a snapshot captured after load and after save. The sidebar becomes safe because it is ordinary in-app navigation and therefore covered by the same guard.

**Alternative considered:** autosave on navigation. Rejected because the article editor's explicit save and conflict handling would need a separate model.

### Authentication enforced once

Move the session check into `router.beforeEach` for routes whose matched metadata has `requiresAuth`, preserving the `isLocalEnvironment()` bypass. Delete the five per-view `onMounted` redirects.

### Param naming

Articles use `:uri` because the API keys articles by URI and the URI can change during a save (the workspace already calls `router.replace` on rename). Coffee and music use `:id` because those records are keyed by UUID. The shell spec calls this out so the difference is intentional rather than accidental.

### View decomposition

- Articles: the existing listing and workspace views move under the shell; the workspace keeps its own layout and guards but adopts the shared composable.
- Coffee: one list view and one edit view per section, since the three forms already take different props. The tasting edit view fetches the tasting, the bean list, and the equipment list in parallel.
- Music: one list view and one edit view; the playlist form keeps fetching the catalogue itself.
- A shared sidebar plus a small section list-row component cover the repeated row markup currently inlined in `CoffeeEditorView` and `MusicEditorView`.

## Risks / Trade-offs

- **Existing editor bookmarks break** -> permanent redirects for all four legacy paths, and update the hidden navbar entries in the same change.
- **Frontend and API must ship together** -> the endpoints are additive and safe to deploy first; the frontend can then be released, and a frontend rollback leaves unused endpoints behind.
- **Dirty detection could produce false prompts** -> capture the baseline after load and after each successful save, and treat a pristine form as clean; article behavior is unchanged because it already does this.
- **Full-screen layout breaks the article editor's height math** -> replace `calc(100vh - 11rem)` assumptions with a flex column sized to the shell viewport and verify the independent scroll panes and preview alignment still hold.
- **Sidebar active-state drift across nested routes** -> derive the active entry from the matched route rather than the raw path, and cover list and edit pages in both checks.
- **Tasting edit needs three requests** -> issue them in parallel and show the existing loading state; the collections are small and public.
- **Blog editor remains unguarded** -> explicitly out of scope; record it as a follow-up rather than silently changing it here.

## Migration Plan

1. Add the four API handlers and routes, and the corresponding service functions. Deploy the API; this is additive and changes no existing response.
2. Add the `/editor` layout and child routes, move the article views, split the coffee and music editors, and add the legacy redirects in one frontend release.
3. Verify each section's list, create, edit, delete, direct-link refresh, and unsaved-change prompt before removing nothing else.
4. Rollback: revert the frontend release to restore the previous editors; the added API endpoints remain harmless.

## Open Questions

- Should login return the author to the editor route they originally requested instead of the home page? Independent of this change's specs and safe to decide during implementation.
