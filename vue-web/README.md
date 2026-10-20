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

## Drawing tool

The `/tool/draw` route embeds the Excalidraw bundle served from `public/excalidraw`, a symlink to `excalidraw-web/dist`. Build `excalidraw-web` before `vue-web`, or run the root `build-web.sh`, which does both in order.
