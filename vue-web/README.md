> Partly generated with ChatGPT. Treat the details with care and verify against the code.

# grd0.net web

Vue 3 + Vite 8 + Tailwind v4 + Pinia single-page app for [grd0.net](https://grd0.net).

## Commands

```bash
pnpm dev          # Vite dev server on :5173
pnpm build        # type-check and production build
pnpm type-check   # vue-tsc --build --force
pnpm format       # prettier --write src/
pnpm preview      # preview the production build
```

`pnpm lint` is currently broken (ESLint 9 flat-config migration); use `pnpm type-check` and `pnpm build` for verification.

## Environment

Copy `.env.example` to `.env` and adjust as needed:

```bash
VITE_API_URL=//api.grd0.net/api
VITE_ASSETS_URL=//assets.grd0.net
VITE_DEFAULT_TITLE=IDEA GROUND ZERO
```

`pnpm dev` is treated as a local environment and skips the frontend auth redirect. To do the same in a local production build, set `VITE_APP_ENV=LOCAL`; pair it with the API's `APP_ENV=LOCAL`. Never use either value in a deployed environment.

## Routes

| Route | Purpose |
| --- | --- |
| `/article` | Published article index (main navigation entry) |
| `/article/:slug` | Published article reader (block renderer) |
| `/article/editor` | Authenticated article listing, paginated (hidden nav entry) |
| `/article/editor/:uri` | Authenticated article editing workspace |
| `/blog`, `/blog/:slug`, `/blog/editor` | Legacy markdown blog, kept for existing links; no longer linked in navigation or the sitemap |
| `/gallery`, `/travel/map`, `/coffee` | Existing sections |
| `/music` | Music page: playlist selector, player, and selected playlist |
| `/music/editor` | Authenticated playlist editor (hidden nav entry) |

Article pages render an ordered block list (`markdown`, `image`, `music`, `excalidraw`) through `components/article/BlockRenderer.vue`. The editor lives under `components/editor/article/` and is split into a paginated listing at `/article/editor` and a workspace at `/article/editor/:uri`. The workspace keeps the editor and preview in independent scroll panes and aligns the preview to the editor by block; the metadata sidebar carries the URI, publish state, publishing actions, and delete. Creating an article happens on the listing and opens its workspace. Drawing blocks use the authenticated draft scene endpoints, so previews work without publishing. Both editor routes redirect to the login flow when no session token is present.

## Music

The music page groups the catalogue into named playlists. A `PlaylistSelector` in the right-hand column switches between playlists and the first playlist is selected by default; there is no combined all-songs view. Each playlist drives the full `MusicPlayer` and the song list beneath the selector.

The editor at `/music/editor` lists playlists and creates, edits, or deletes them. `PlaylistForm` edits the title, description, sorting, and the catalogue song checklist. The add-song action reuses `components/editor/article/MusicPicker.vue` with `:link_only="true"`, so the picker hides the catalogue rows already shown by the form. A resolved YouTube link is previewed with the compact `MusicPlayer` before confirmation, and a confirmed song that is not yet in the catalogue is persisted through `PUT /music`. The per-song Remove action deletes the catalogue entry with confirmation.

Playlist membership is many-to-many: a song can appear in several playlists. Playlist writes replace the complete membership in one transaction and never delete songs. Deleting a song from the catalogue removes it from playlists, while article music blocks keep rendering from the identifier, title, and artist snapshotted into the block.

## Drawing tool

The `/tool/draw` route embeds the Excalidraw bundle served from `public/excalidraw`, a symlink to `excalidraw-web/dist`. Build `excalidraw-web` before `vue-web`, or run the root `build-web.sh`, which does both in order.
