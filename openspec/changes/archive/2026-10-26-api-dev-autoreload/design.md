# Design

## Context

See `proposal.md` for motivation. The approach is shaped by constraints found in the current code and on this machine:

- `utils.GetEnv` panics when a variable is missing and has no fallback, so `S3_ENDPOINT`, `S3_ACCESS_KEY_ID` and `S3_SECRET_ACCESS_KEY` are mandatory today even with `APP_ENV=LOCAL`. `minio.New("")` would also fail.
- S3 is only used by `GeneratePresignedUrl` (`handler/blog.go`) and `GenerateArticlePresignedUrl` (`handler/article.go`). Blob reads and writes already switch on `BLOB_BACKEND`.
- `PORT` defaults to `:80`. The host has `net.ipv4.ip_unprivileged_port_start = 1024`, so a non-root process cannot bind it.
- The frontend already treats `import.meta.env.DEV` as local: it skips the router auth guard and does not attach a JWT, so the API must run `APP_ENV=LOCAL` for writes to work.
- `api/data/` contains the SQLite database and generated JSON, both of which change at runtime and must not trigger rebuilds.
- Air v1.67.4 requires Go 1.26, which matches `api/go.mod`. It only auto-discovers `.air.toml`/`.config/air.toml`, and it has no `[build.env]` table; environment comes from `env_files` or `full_bin`.
- Go 1.27 is installed, so tool directives are available. `~/go/bin` is not on `PATH`, so `go install` would not be directly runnable.
- Docker 29.9 and Compose v5 are installed.

## Goals / Non-Goals

**Goals:**

- Rebuild and restart the API automatically when Go sources under `api/` change, with no manual step and no global tool install.
- Provide an identical loop natively and in a container, configured from one file each.
- Start locally with no S3 configuration and no OIDC configuration.
- Keep deployed behaviour byte-for-byte equivalent: `:80`, OIDC, S3 and blob backends unchanged.

**Non-Goals:**

- Orchestrating Vite or Excalidraw from the same command; their own dev servers are unchanged.
- Restarting on `data/*.json` edits. Those are re-read per request unless `LAMBDA=TRUE`, so a rebuild would be wasted.
- In-process hot reload, graceful shutdown on restart, or zero-downtime dev restarts.
- Making S3-backed attachment uploads work locally. They return `501` by design.
- Windows development support.

## Decisions

### Air, pinned as a Go tool

Use `go get -tool github.com/air-verse/air@v1.67.4` and run it as `go tool air`. Alternatives considered:

- **`go install`**: rejected because `~/go/bin` is not on `PATH`; every developer would need a PATH edit and the version would drift per machine.
- **`go run github.com/air-verse/air@v1.67.4`**: keeps `go.mod`/`go.sum` clean, but the version lives in a script instead of the module graph and is re-resolved on use.
- **`watchexec`/`entr`**: rejected because the host is an atomic Fedora derivative; installing a system watcher means layering or rpm-ostree, and neither offers Air's ignore rules.
- **`reflex`**: viable, but it has weaker directory exclusions, which matters because `data/` must be ignored.
- **`nodemon`**: rejected as an unnecessary Node layer for a pure Go loop.
- **`docker compose watch`**: rejected as the primary loop because image rebuilds are much slower than Air's incremental `go build`.

The cost is that Air imports Hugo's `watcher/filenotify`, so Hugo v0.164.0 and ~20 transitive modules enter `go.sum` (~52 lines in a clean probe module, 25 MB once in the module cache). Tool dependencies are a separate graph: `go build` in the production Dockerfile ignores them and none are linked into the API binary. Existing higher versions in `api/go.mod` (`golang.org/x/sys v0.49.0`, `golang.org/x/text v0.43.0`) win under MVS, so no dependency is downgraded.

### One `.air.toml` at `api/`, watching Go only

`root = "."`, `tmp_dir = "/tmp/grd0-api-air"`, `include_ext = ["go"]`, `exclude_dir = ["data", ".git"]`. The command is `go build -o /tmp/grd0-api-air/main .` with `entrypoint = ["/tmp/grd0-api-air/main"]` rather than `go run`, so Air owns the child process and can replace it cleanly.

The tmp directory is deliberately outside the source tree. Docker creates missing paths for nested mounts as root, so a repo-local `tmp/` inside the bind mount would become root-owned and block the next native Air run once `clean_on_exit` removes it. Keeping it under `/tmp` means the repository is never touched and both loops can alternate safely.

`stop_on_error = false` keeps the last working server running while a compile error is fixed, which is the main DX benefit of a watcher. `exclude_unchanged = true` and `delay = 300` avoid duplicate builds from editor save bursts. `clean_on_exit = true` removes `/tmp/grd0-api-air` on shutdown; nothing build-related lands in the repository.

Air must be launched from `api/` (that is also where `go tool` resolves the module), so `root = "."` never accidentally watches the monorepo's `node_modules`.

### Environment from a tracked `api/.env.dev`

Air's `env_files = [".env.dev"]` loads values into the process environment and the child inherits them. The file is tracked because `.env.dev` does not match the root `.env` or `*.env` ignore rules; `api/api.env` does match `*.env` and stays untracked and untouched.

Contents: `APP_ENV=LOCAL`, `PORT=:8080`, `DB_PATH=data/api.db`, `BLOB_BACKEND=local`, `BLOB_LOCAL_PATH=data/scenes`, `LOG_LEVEL=INFO`. No secrets. Because Air uses `godotenv` semantics, a variable already exported in the shell wins over the file, so real credentials can still be injected for a one-off session. Editing `.env.dev` does not itself trigger a rebuild; the next Go change applies it.

### Development port `:8080`

The dev listener moves to `:8080`; the deployed default stays `:80`. The frontend's untracked `vue-web/.env` points at `http://localhost:8080` locally. CORS needs no change because it is origin-based and already allows `http://localhost:5173`. The dev port is an override of `PORT`, not a code change.

### No S3 in local mode

In `server.go`, only call `storage.InitiateClient` when the environment is not local. The two presigned-URL handlers return `501 Not Implemented` when `h.S3 == nil`. Production behaviour is unchanged because `APP_ENV != LOCAL` always initialises the client; the nil guard is unreachable there.

Rejected alternative: keep the API untouched and put dummy `S3_ENDPOINT`/keys in `.env.dev`. It is zero code, but it spreads fake configuration, and the editor's upload path would receive a presigned URL to a dead host instead of a clear `501`. Skipping initialisation better matches "no S3" and fails honestly.

### Persistent dev data under `data/`

Use `data/api.db` and `data/scenes`, which already exist and are gitignored, so schema and content survive restarts. Pointing either at Air's tmp directory would be deleted by `clean_on_exit`. `data/` is excluded from watching, so SQLite WAL writes and `assets-helper` regenerating `files.json` never cause rebuilds. `LAMBDA` stays unset so JSON is read from disk per request.

### Container path with the same configuration

`api/Dockerfile.dev` builds on `golang:1.26-alpine`, installs `git`, runs `go mod download`, and prebuilds Air with `RUN go tool air -v` so the first start is fast. `api/compose.dev.yaml` bind-mounts `./:/app`, publishes `8080:8080`, and puts a named volume at `/tmp/grd0-api-air`.

The named volume mounts Air's tmp directory outside the bind-mounted source, so the container never creates a root-owned `api/tmp` in the checkout. No module or build cache volumes are mounted: they would hide the image's warmed caches from the `RUN` layers. Air uses fsnotify, which works on Linux bind mounts; `poll = true` remains a documented fallback for non-Linux hosts. Command: `docker compose -f compose.dev.yaml up --build`.

### Documentation

`api/README.md` documents both loop commands and the dev port; `AGENTS.md`'s "API listens on `:80`" note is corrected to say the deployed default is `:80` and dev uses `:8080`.

## Risks / Trade-offs

- **[`go.sum` growth from Hugo]** → Tool dependencies are not compiled into the API; if the graph is unacceptable, switch to `go run github.com/air-verse/air@v1.67.4` and drop the tool directive. Verified that `go tool air -v` builds in ~3 s and no existing dependency is downgraded.
- **[Restarts kill in-flight requests]** → Acceptable for development. The server has no signal handling today, so `send_interrupt` behaves like termination; adding graceful shutdown is out of scope.
- **[SQLite during rapid restarts]** → Air stops the old process before starting the new one, so only one process holds the database; SQLite recovers its WAL automatically.
- **[`PORT=:8080` hard-coded in config]** → Shell exports override `.env.dev`, and `PORT` can be changed there; the frontend `.env` must match.
- **[Vite on a non-5173 port]** → CORS only allows `localhost:5173`; note this in the docs rather than widening origins.
- **[`vue-web/.env` is untracked]** → Each developer sets `VITE_API_URL` locally; `.env.example` keeps the production URL, so the README must call out the local override.
- **[Root-owned container artefacts]** → Air's tmp directory sits outside the bind mount at `/tmp/grd0-api-air` and is backed by a named volume, so a fresh checkout never gains a root-owned directory and a hard-killed container leaves data only in the volume.
- **[Air does not watch `.env.dev`]** → Changing only the env file needs a touch of a Go file or a restart; documented.
