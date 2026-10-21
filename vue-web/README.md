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
| `/gallery`, `/travel/map`, `/music`, `/coffee` | Existing sections |

Article pages render an ordered block list (`markdown`, `image`, `music`, `excalidraw`) through `components/article/BlockRenderer.vue`. The editor lives under `components/editor/article/` and is split into a paginated listing at `/article/editor` and a workspace at `/article/editor/:uri`. The workspace keeps the editor and preview in independent scroll panes and aligns the preview to the editor by block; the metadata sidebar carries the URI, publish state, publishing actions, and delete. Creating an article happens on the listing and opens its workspace. Drawing blocks use the authenticated draft scene endpoints, so previews work without publishing. Both editor routes redirect to the login flow when no session token is present.

## Drawing tool

The `/tool/draw` route embeds the Excalidraw bundle served from `public/excalidraw`, a symlink to `excalidraw-web/dist`. Build `excalidraw-web` before `vue-web`, or run the root `build-web.sh`, which does both in order.
