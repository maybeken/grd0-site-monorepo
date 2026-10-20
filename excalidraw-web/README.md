# grd0 draw

React 19 + [Excalidraw](https://www.npmjs.com/package/@excalidraw/excalidraw) drawing tool, embedded by `vue-web` as a same-origin iframe at `/excalidraw/index.html` and surfaced at `/tool/draw`.

## Commands

```bash
pnpm dev      # copies fonts, then starts Vite
pnpm build    # copies fonts, runs tsc -b and vite build
pnpm lint     # eslint
pnpm preview  # preview the production build
```

`scripts/copy-fonts.mjs` copies Excalidraw's runtime fonts into `public/fonts`, and `src/main.tsx` points `window.EXCALIDRAW_ASSET_PATH` at the current base path so no fonts are fetched from a CDN.

## Parent bridge

- Parent → frame: `{ event: 'command', func: 'save' | 'load', payload }`
- Frame → parent: `{ save: <Excalidraw scene object> }` or `{ loadError: 'legacy' | 'invalid' }`

The scene is also autosaved to `localStorage` under `grd0.draw.scene`.
