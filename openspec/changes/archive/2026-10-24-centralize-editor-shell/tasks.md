# Tasks

## 1. API detail endpoints

- [x] 1.1 Add a `GetCoffeeBean` handler in `api/handler/coffee.go` that loads one non-deleted bean by `id` and returns 404 when absent; verify with `cd api && go build -o ../grd0_blog_api .`
- [x] 1.2 Add a `GetCoffeeEquipment` handler that loads one non-deleted equipment record by `id` and returns 404 when absent; verify with `cd api && go build -o ../grd0_blog_api .`
- [x] 1.3 Add a `GetCoffeeTasting` handler that preloads Bean, Equipment, and Grinder and returns 404 when absent; verify with `cd api && go build -o ../grd0_blog_api .`
- [x] 1.4 Add a `GetPlaylist` handler in `api/handler/music.go` that preloads songs with `orderSongs` and returns 404 when absent; verify with `cd api && go build -o ../grd0_blog_api .`
- [x] 1.5 Register `GET /coffee/bean/:id`, `GET /coffee/equipment/:id`, `GET /coffee/tasting/:id`, and `GET /music/playlist/:id` on the public `v1` group; verify each returns a JSON object for an existing id and 404 for an unknown id with the API running locally
- [x] 1.6 Add the four requests to the Restfox collection under `docs/api/`; verify the JSON still parses

## 2. Frontend services

- [x] 2.1 Add `getCoffeeBeanRaw(id)`, `getCoffeeEquipmentRaw(id)`, and `getCoffeeTastingRaw(id)` to `vue-web/src/services/coffee.ts` using the data instance; verify with `cd vue-web && pnpm type-check`
- [x] 2.2 Add `getPlaylistRaw(id)` to `vue-web/src/services/music.ts` using the data instance; verify with `cd vue-web && pnpm type-check`

## 3. Shared unsaved-change guard

- [x] 3.1 Create `vue-web/src/composables/useUnsavedGuard.ts` exposing `onBeforeRouteLeave` confirmation plus a `beforeunload` listener, and refactor `ArticleEditor.vue` to use it without changing its current prompts; verify a dirty article still prompts on sidebar/back navigation and on reload
- [x] 3.2 Add baseline snapshot tracking to `BeanForm.vue`, `EquipmentForm.vue`, `TastingForm.vue`, and `PlaylistForm.vue`, and apply the shared guard in each; verify editing a field then navigating away prompts, while an untouched form does not

## 4. Editor shell, auth, and routing

- [x] 4.1 Create `EditorLayout.vue` and `EditorSidebar.vue` rendering the Content, Coffee, and Music groups with a site link; verify the layout renders at `/editor` with no public navbar or footer
- [x] 4.2 Add the `meta.layout: 'editor'` check to `App.vue` so the public background, navbar, and footer are skipped only for matched editor routes; verify public pages render unchanged and `/editor` renders full-screen
- [x] 4.3 Add the `/editor` route subtree with the section list and edit routes, `''` and `/coffee` redirects, and the four legacy path redirects; verify each legacy URL lands on its new equivalent
- [x] 4.4 Enforce `requiresAuth` for matched editor routes in `router.beforeEach` with the existing local-environment bypass; verify an unauthenticated visit redirects to login and a local visit does not
- [x] 4.5 Derive the sidebar's active entry from the matched route and collapse the sidebar behind a toggle below the large breakpoint; verify list and edit pages highlight correctly and the narrow layout stays usable
- [x] 4.6 Add sidebar labels and route title keys to `vue-web/src/locales/*.yaml`; verify `cd vue-web && pnpm type-check` and that document titles resolve for each new route
- [x] 4.7 Update the hidden editor entries in `Navbar.vue` to the new paths; verify each opens the matching section

## 5. Articles section

- [x] 5.1 Move the listing to an `/editor/articles` view and keep `?page=` pagination and the create-draft-then-open-route action; verify paging, creation, and return-to-listing preserve the page
- [x] 5.2 Move the workspace to `/editor/articles/:uri`, update its back link and post-rename `router.replace` to the new paths; verify direct open, rename, and back navigation
- [x] 5.3 Remove the view's duplicated `onMounted` session check now that the shell enforces auth; verify no article route renders without the shell guard in place
- [x] 5.4 Verify the workspace's height and scrolling work in the full-screen shell, including independent editor/preview scrolling and preview alignment

## 6. Coffee sections

- [x] 6.1 Create the Beans listing view with inline create, row edit links, and delete confirmation; verify create refreshes the list and edit navigates to `/editor/coffee/beans/:id`
- [x] 6.2 Create the Bean edit view that loads by id, renders `BeanForm`, and shows a not-found state with a link back; verify direct open, save, and an unknown id
- [x] 6.3 Create the Equipment listing view with inline create, row edit links, and delete confirmation; verify the equivalent flows at `/editor/coffee/equipment`
- [x] 6.4 Create the Equipment edit view with the same load/not-found behavior; verify direct open, save, and an unknown id
- [x] 6.5 Create the Tastings listing view with inline create, row edit links, and delete confirmation; verify create and delete refresh the list
- [x] 6.6 Create the Tasting edit view that loads the tasting, beans, and equipment in parallel and passes them to `TastingForm`; verify direct open, populated selects, save, and an unknown id
- [x] 6.7 Remove `CoffeeEditorView.vue` and its duplicated session check; verify `pnpm type-check` reports no dangling imports

## 7. Music section

- [x] 7.1 Create the Playlists listing view at `/editor/music` with inline create, row edit links, and delete confirmation; verify create refreshes the list and edit navigates to `/editor/music/:id`
- [x] 7.2 Create the Playlist edit view that loads by id and renders `PlaylistForm`, including a not-found state with a link back; verify direct open, save, and an unknown id
- [x] 7.3 Remove `MusicEditorView.vue` and its duplicated session check; verify `pnpm type-check` reports no dangling imports

## 8. Integration verification

- [x] 8.1 Verify all four legacy editor URLs redirect to their new routes and that no public route behavior changed
- [x] 8.2 Verify the unsaved-change prompt across every section for sidebar navigation, back navigation, and reload, and confirm cancel keeps the edits
- [x] 8.3 Run `cd vue-web && pnpm type-check` and `cd vue-web && pnpm build`; confirm both succeed
- [x] 8.4 Run `cd api && go build -o ../grd0_blog_api .`; confirm it succeeds
- [x] 8.5 Smoke-test the full editor flow against the local API: create, edit, direct-link refresh, and delete in each section

## Workflow follow-up

- Archive the change once the project's review requirements are satisfied.
