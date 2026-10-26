# Proposal

## Why

The Go API has no development reload loop. Every change needs `cd api && go run .`, and `go run .` cannot even start without exported OIDC and S3 variables, so local work defaults to `local-test.sh`, which rebuilds a Docker image per iteration. The listener also defaults to `:80`, which an unprivileged user cannot bind on this machine (`net.ipv4.ip_unprivileged_port_start = 1024`). The result is a slow, error-prone inner loop for the code that changes most often.

## What Changes

- Pin Air as a Go tool dependency and run it with `go tool air`, so no global install or PATH change is required.
- Add `api/.air.toml` that watches only Go sources under `api/` and excludes `data/`, `tmp/` and `.git`, rebuilding `./tmp/main` on change.
- Add a tracked `api/.env.dev` for local development: `APP_ENV=LOCAL`, `PORT=:8080`, SQLite at `data/api.db`, and the local blob backend under `data/scenes`.
- Make the API start without S3 configuration when `APP_ENV=LOCAL`: skip the MinIO client initialisation and return `501` from the two presigned-upload handlers instead of dereferencing a nil client.
- Add a Docker development path (`api/Dockerfile.dev` and `api/compose.dev.yaml`) that runs the same Air config against a bind-mounted `api/` for developers who prefer containers.
- Document the loop in `api/README.md` and correct the `:80` note in `AGENTS.md`; document the local `vue-web/.env` override to `http://localhost:8080`.
- Leave the deployed listener default at `:80` and production authentication, blob storage and S3 behaviour unchanged.

## Capabilities

### New Capabilities

None. This change only affects the development workflow, not system behavior.

### Modified Capabilities

None.

## Impact

- `api/go.mod`, `api/go.sum`: Air is added as a tool dependency, which pulls the Hugo module graph into `go.sum` through Air's `watcher/filenotify` import. No runtime package is added to the API binary.
- `api/.air.toml`, `api/.env.dev`: new development configuration.
- `api/server.go`: conditionally skip `storage.InitiateClient` when `APP_ENV=LOCAL`.
- `api/handler/blog.go`, `api/handler/article.go`: return `501` from presigned-URL handlers when the S3 client is absent.
- `api/Dockerfile.dev`, `api/compose.dev.yaml`: new container development path.
- `api/README.md`, `AGENTS.md`: document the native and container loops and the `:8080` dev port.
- No production configuration, endpoint, schema or frontend code changes; `vue-web/.env` remains untracked and is adjusted locally.
