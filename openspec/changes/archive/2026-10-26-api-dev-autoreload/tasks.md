# Tasks

## 1. Add Air as a pinned Go tool

- [x] 1.1 Run `go get -tool github.com/air-verse/air@v1.67.4` from `api/` and verify `go tool air -v` prints `v1.67.4` and `api/go.mod` contains both the `tool github.com/air-verse/air` line and a `require` entry.
- [x] 1.2 Confirm the existing module versions are not downgraded by the new graph: `go mod graph | grep -E 'golang.org/x/(sys|text)'` still resolves `x/sys v0.49.0` and `x/text v0.43.0`, and `go build ./...` in `api/` succeeds.
- [x] 1.3 Add `api/.air.toml` with `root = "."`, `tmp_dir = "/tmp/grd0-api-air"`, `env_files = [".env.dev"]`, `include_ext = ["go"]`, `exclude_dir = ["data", ".git"]`, `cmd = "go build -o /tmp/grd0-api-air/main ."`, `entrypoint = ["/tmp/grd0-api-air/main"]`, `stop_on_error = false`, `exclude_unchanged = true`, `delay = 300`, and `[misc] clean_on_exit = true`; verify `cd api && go tool air` logs a successful build and creates `/tmp/grd0-api-air/main`.
- [x] 1.4 Add tracked `api/.env.dev` with `APP_ENV=LOCAL`, `PORT=:8080`, `DB_PATH=data/api.db`, `BLOB_BACKEND=local`, `BLOB_LOCAL_PATH=data/scenes`, `LOG_LEVEL=INFO`; verify `git check-ignore -v api/.env.dev` reports nothing while `git check-ignore -v api/api.env` still matches `*.env`.
- [x] 1.5 Update `AGENTS.md` so the `:80` note says the deployed default is `:80` and development uses `:8080` via `PORT`; verify the note matches `api/.env.dev` and `server.go`.

## 2. Start without S3 in local mode

- [x] 2.1 In `api/server.go`, skip `storage.InitiateClient` when `utils.IsLocalEnvironment()` is true and pass a nil client to the handler; verify `cd api && env -u S3_ENDPOINT -u S3_ACCESS_KEY_ID -u S3_SECRET_ACCESS_KEY APP_ENV=LOCAL PORT=:8080 go run .` starts and serves `GET /article` without panicking.
- [x] 2.2 Return `501 Not Implemented` from `GeneratePresignedUrl` (`handler/blog.go`) and `GenerateArticlePresignedUrl` (`handler/article.go`) when `h.S3 == nil`; verify `curl -i -X PUT localhost:8080/blog/attachment/test` and `curl -i -X PUT localhost:8080/article/attachment/test` both return `501` under `APP_ENV=LOCAL`.
- [x] 2.3 Confirm the non-local path is unchanged: with `APP_ENV` unset and `S3_ENDPOINT` unset, the process still halts on the missing environment variable, so production cannot silently run without S3.
- [x] 2.4 Document the local `501` behaviour and the "no S3 in local mode" rule in the `api/README.md` local-authentication section; verify the documented curl commands return what the README states.

## 3. Verify the native loop

- [x] 3.1 Start `cd api && go tool air`, verify `curl -s localhost:8080/article` returns the article list JSON, then edit a Go file (for example add a log line in `server.go`) and verify the air log shows a rebuild and the server answers again without manual intervention.
- [x] 3.2 Verify only Go sources trigger rebuilds: `touch api/data/files/files.json` and `touch api/data/api.db` while air runs, and confirm no rebuild line appears in the air log.
- [x] 3.3 Verify `/tmp/grd0-api-air` is removed after `Ctrl-C`, and that `git status` shows no new build artefacts.

## 4. Add the Docker development path

- [x] 4.1 Add `api/Dockerfile.dev` (`golang:1.26-alpine`, `git`, `go mod download`, `RUN go tool air -v`) and verify `docker build -f Dockerfile.dev -t grd0/api-dev .` succeeds from `api/`.
- [x] 4.2 Add `api/compose.dev.yaml` with a bind mount of `./:/app`, `ports: 8080:8080`, and a named volume at `/tmp/grd0-api-air`; verify `docker compose -f compose.dev.yaml up --build` serves `GET /article` on `localhost:8080`.
- [x] 4.3 Verify the container reloads on a host edit: change a Go file while the compose service runs and confirm a rebuild in the container logs and a live response afterwards.
- [x] 4.4 Verify the named volume keeps the working tree clean: stop the container and confirm `api/tmp` does not exist and no root-owned files appear under `api/`.
- [x] 4.5 Document the container loop, the `poll = true` fallback for non-Linux hosts, and the `vue-web/.env` port override in `api/README.md`; verify each documented command runs as written.

## 5. Integration checks

- [x] 5.1 Run the frontend against the local API: set `VITE_API_URL=http://localhost:8080` in `vue-web/.env`, start `pnpm dev` and `cd api && go tool air`, and verify the article list loads and an editor save succeeds while `APP_ENV=LOCAL` bypasses auth.
- [x] 5.2 Verify the production image is unaffected: `cd api && go build -o ../grd0_blog_api .` succeeds and `go list -deps . | grep -E 'air|gohugoio'` prints nothing.
- [x] 5.3 Run `openspec validate api-dev-autoreload --strict` and verify it passes.

## Workflow follow-up

- Archive the change after the loop has been used for one real API edit and the review requirements are satisfied, using `openspec archive api-dev-autoreload --skip-specs` because the change declares no spec deltas.
- Verify the archive move succeeds and `openspec list` no longer reports the change as active.
