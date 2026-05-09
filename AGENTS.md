# AGENTS.md

## Monorepo structure

- **pnpm workspace** (`pnpm-workspace.yaml`): `vue-web/`, `tldraw-web/`
- **Independent packages** (not in workspace, own lockfiles):
  - `api/` — Go backend (module `grd0.net/api`, Go 1.24)
  - `assets-helper/` — Node.js TOC generator (has its own `pnpm-lock.yaml`)

There is **no root `package.json`**. Run pnpm commands from workspace package dirs or use `pnpm --filter`.

## Commands

### vue-web (main site, Vue 3 + Vite)

```
pnpm dev          # vite dev server (default port 5173)
pnpm build        # type-check (vue-tsc) + vite build (run in parallel via run-p)
pnpm type-check   # vue-tsc --build --force
pnpm lint         # eslint --fix (⚠ currently broken — uses removed ESLint 9 flags)
pnpm format       # prettier --write src/
```

### tldraw-web (React + Vite)

```
pnpm dev          # vite dev server
pnpm build        # tsc -b && vite build (sequential)
pnpm lint         # eslint .
```

### api (Go)

```
go run server.go         # dev server (port :80)
docker build -t api .    # multi-stage build from api/Dockerfile
```

### assets-helper

```
pnpm run generate   # generates files.json into ../api/data/files/
```

### Build all web apps

```
bash build-web.sh   # builds tldraw-web then vue-web sequentially
```

## Key conventions

- **Prettier** (`vue-web`): `semi: false`, `singleQuote: true`, `tabWidth: 2`, `trailingComma: none`, `printWidth: 100`
- **ESLint**: `tldraw-web` uses ESLint 9 flat config (`eslint.config.js`); `vue-web` uses legacy `.eslintrc.cjs` with `@rushstack/eslint-patch`. The `pnpm lint` command in vue-web is broken — it passes `--ignore-path` and `--ext` flags that ESLint 9 removed. When linting manually, omit those flags.
- **TypeScript**: Both `vue-web` and `tldraw-web` use project references (`tsconfig.app.json` + `tsconfig.node.json`)
- **Tailwind CSS v4**: Uses `@tailwindcss/vite` plugin (not the v3 PostCSS approach). No `tailwind.config.js` — config is in the CSS file.
- **Component auto-import**: `unplugin-vue-components` auto-imports components from `src/components/`. Components like `Modal`, `Navbar`, `Icon` are usable in templates without explicit imports. The generated file is `vue-web/components.d.ts`.

## Tiptap editor (vue-web)

The blog editor uses Tiptap (`@tiptap/vue-3`). Key gotchas:

- **`BubbleMenu`** is exported from `@tiptap/vue-3/menus`, NOT the main `@tiptap/vue-3` package:
  ```ts
  import { BubbleMenu } from '@tiptap/vue-3/menus'
  ```
- **`@tiptap/core` is not directly resolvable** — import `Extension` and other core types from `@tiptap/vue-3` instead (it re-exports them):
  ```ts
  import { Extension } from '@tiptap/vue-3'
  ```
- **Type augmentation** — `src/types/tiptap-markdown.d.ts` augments `@tiptap/core`'s `Storage` interface so `editor.storage.markdown.getMarkdown()` type-checks. Do not delete this file.
- **Markdown round-trip** — `tiptap-markdown` intercepts `setContent()` to parse markdown, and `editor.storage.markdown.getMarkdown()` serializes back. Content is stored as markdown in the backend (no format change).

## API / data fetching

Uses `alova` (not axios/fetch). Services in `src/services/` wrap alova's `useRequest()`:

```ts
// Returns { loading: Ref<boolean>, data: Ref<T> }
const { loading, data } = useRequest(someMethod())
```

Two alova instances exist:
- `dataInstance` — public read-only API (with GET caching)
- `adminInstance` — authenticated API (attaches JWT from `sessionStorage.getItem("jwt_token")`)

## Environment variables (api)

Required env vars for the Go API server: `LOG_LEVEL`, `JWT_SECRET`, `SESSION_SECRET`, `NEXTCLOUD_URL`, `NEXTCLOUD_CLIENT_KEY`, `NEXTCLOUD_CLIENT_SECRET`, `S3_ENDPOINT`, `S3_ACCESS_KEY_ID`, `S3_SECRET_ACCESS_KEY`, `S3_BUCKET`

## Testing

No test runner or test suites are configured in this repo.
