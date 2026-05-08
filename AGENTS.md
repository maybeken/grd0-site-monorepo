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
pnpm dev          # vite dev server
pnpm build        # type-check (vue-tsc) + vite build (run in parallel via run-p)
pnpm type-check   # vue-tsc --build --force
pnpm lint         # eslint --fix
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

- **Prettier** (`vue-web`): `semi: false`, `singleQuote: true`, `tabWidth: 2`, `trailingComma: none`
- **ESLint**: `tldraw-web` uses ESLint 9 flat config (`eslint.config.js`); `vue-web` uses legacy `.eslintrc.cjs`
- **TypeScript**: Both `vue-web` and `tldraw-web` use project references (`tsconfig.app.json` + `tsconfig.node.json`)

## Environment variables (api)

Required env vars for the Go API server: `LOG_LEVEL`, `JWT_SECRET`, `SESSION_SECRET`, `NEXTCLOUD_URL`, `NEXTCLOUD_CLIENT_KEY`, `NEXTCLOUD_CLIENT_SECRET`, `S3_ENDPOINT`, `S3_ACCESS_KEY_ID`, `S3_SECRET_ACCESS_KEY`, `S3_BUCKET`

## Testing

No test runner or test suites are configured in this repo.
