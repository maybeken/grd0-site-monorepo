# Proposal

## Why

The three content editors are unrelated pages that each fuse a listing and an editing form into one view. Only the article editor addresses an individual item in the URL; the coffee and music editors hold the selected item in component state, so a specific edit cannot be linked, refreshed, or reached with the browser's back button. Each editor also duplicates its own authentication redirect, and only the article editor protects unsaved changes.

## What Changes

- Add an authenticated `/editor` shell: a full-screen layout (no public navbar, footer, or particle background) with a left sidebar navigating to every editor list page.
- Move the article editor into the shell: `/editor/articles` for the listing (keeping `?page=`) and `/editor/articles/:uri` for the workspace.
- Split the coffee tabbed editor into routed sections: `/editor/coffee/beans`, `/editor/coffee/equipment`, and `/editor/coffee/tastings` listings, each with a `/:id` edit route.
- Move the music editor to `/editor/music` with `/editor/music/:id` for editing a playlist.
- Keep creation on the list page for every section: the list renders the create form and refreshes on save. Articles keep their existing create-draft-then-open-route flow.
- Add public get-by-id endpoints `GET /coffee/bean/:id`, `GET /coffee/equipment/:id`, `GET /coffee/tasting/:id`, and `GET /music/playlist/:id` so an edit route can resolve its item on a cold load.
- Add a shared unsaved-changes guard used by every edit route, confirming in-app navigation and prompting on reload or tab close.
- Enforce authentication once for the `/editor` subtree, replacing the per-view session checks.
- Redirect the previous editor paths: `/article/editor` to `/editor/articles`, `/article/editor/:uri` to `/editor/articles/:uri`, `/coffee/editor` to `/editor/coffee/beans`, and `/music/editor` to `/editor/music`.
- **BREAKING** for saved editor links: the editor URLs change. Redirects preserve them, and no public page URL changes.
- The legacy blog editor is out of scope and stays where it is.

## Capabilities

### New Capabilities

- `editor-console`: the authenticated editor shell. Owns the full-screen layout, sidebar section navigation, the list/edit route conventions, centralized auth redirect, and shared unsaved-change protection.

### Modified Capabilities

- `coffee-tasting-editor`: the tabbed single-page editor is replaced by routed section listings and per-item edit routes.
- `music-playlists`: the editor moves into the shell and playlist editing becomes an item route while creation stays on the listing; a public playlist detail endpoint is added.
- `coffee-bean`: a public get-by-id endpoint is added.
- `brew-equipment`: a public get-by-id endpoint is added.
- `coffee-tasting`: a public get-by-id endpoint returning nested bean and equipment is added.

## Impact

- **Frontend routing and chrome**: `vue-web/src/router/index.ts`, `App.vue` (editor layout conditional), a new `EditorLayout`/sidebar, and per-section list and edit views replacing `ArticleEditorListView`, `ArticleEditorView`, `CoffeeEditorView`, and `MusicEditorView`.
- **Shared behavior**: a new `useUnsavedGuard` composable extracted from `ArticleEditor`; dirty tracking added to the coffee and music forms.
- **Services**: `getCoffeeBeanRaw`, `getCoffeeEquipmentRaw`, `getCoffeeTastingRaw`, and `getPlaylistRaw` helpers.
- **API**: four handlers and routes in `api/handler/coffee.go`, `api/handler/music.go`, and `api/handler/handler.go`.
- **i18n**: new sidebar labels and route title keys in `src/locales/*.yaml`.
- **Unaffected**: public pages, the API's existing list/upsert/delete behavior, article editing semantics, and the blog editor.
