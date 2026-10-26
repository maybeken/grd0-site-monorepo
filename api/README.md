> Partly generated with ChatGPT. Treat the details with care and verify against the code.

# grd0.net API

Go 1.26 BFF built with Echo v4, GORM/SQLite and the MinIO client for S3. Module path: `grd0.net/api`.

## Commands

```bash
# From this directory

# Development loop: rebuilds and restarts on Go changes, listens on :8080
go tool air

# Build and run once (requires the environment variables below; listens on :80 by default)
go run .

go build -o ../grd0_blog_api .
```

`go tool air` is pinned in `go.mod` (`tool github.com/air-verse/air`) and reads `api/.air.toml`. It loads `api/.env.dev`, which sets `APP_ENV=LOCAL`, `PORT=:8080`, SQLite at `data/api.db` and the local blob backend at `data/scenes`. Air watches only `.go` files under `api/` and ignores `data/` and `.git`; its build output lives in `/tmp/grd0-api-air`, so the checkout stays clean.

Point the frontend at it by setting `VITE_API_URL=http://localhost:8080` in `vue-web/.env`. Keep the dev server on `http://localhost:5173`, because the API CORS allowlist is fixed.

### Docker development loop

```bash
# From this directory
docker compose -f compose.dev.yaml up --build
```

The container bind-mounts `api/` and runs the same `.air.toml` and `.env.dev`, publishing `:8080`. Air's tmp directory is backed by a named volume at `/tmp/grd0-api-air`, so the container never creates root-owned paths in the checkout. `./local-test.sh` runs this command, and `./local-test.sh --no-docker` runs the native `go tool air` loop instead.

File watching uses fsnotify, which works with Linux bind mounts. On Docker Desktop or another host that does not deliver inotify events into the container, set `poll = true` under `[build]` in `api/.air.toml` and restart the service.

The production Docker image is built from the repository root: `cd api; go build -o ../grd0_blog_api`.

## Environment variables

Required:

```bash
OIDC_ISSUER=
OIDC_CLIENT_ID=
OIDC_CLIENT_SECRET=
AUTH_JWT_SECRET=
AUTH_STATE_SECRET=
AUTH_REDIRECT_URL=
AUTH_CALLBACK_URL=
S3_ENDPOINT=
S3_ACCESS_KEY_ID=
S3_SECRET_ACCESS_KEY=
```

Optional:

```bash
PORT=:80                        # listen address
DB_PATH=data/api.db             # SQLite database path
JWKS_CACHE_TTL=1h
LOG_LEVEL=WARN                  # DEBUG, INFO, WARN or ERROR
S3_BUCKET=assets-grd0-net
LAMBDA=TRUE                     # enables in-memory JSON caching
BLOB_BACKEND=local              # local or s3
BLOB_LOCAL_PATH=data/scenes     # root for the local blob backend
ARTICLE_SCENE_MAX_BYTES=5242880 # maximum Excalidraw scene write size
APP_ENV=LOCAL                   # local only: bypass authentication (never in production)
LOCAL_AUTHOR_EMAIL=local@grd0.net
LOCAL_AUTHOR_NAME=Local Author
```

### Local authentication bypass

Set `APP_ENV=LOCAL` to skip OIDC initialization and authenticate every restricted route with a synthetic author (`LOCAL_AUTHOR_EMAIL` / `LOCAL_AUTHOR_NAME`, defaulting to `local@grd0.net` / `Local Author`). The bypass is intended only for local development and logs a warning at startup. Never set `APP_ENV=LOCAL` in a deployed environment. `/auth/login`, `/auth/callback`, and `/auth/refresh` return `404` while it is active. `local-test.sh` runs the Docker development loop, which loads `api/.env.dev` with `APP_ENV=LOCAL`.

Local mode also skips the S3 client, so `S3_ENDPOINT`, `S3_ACCESS_KEY_ID` and `S3_SECRET_ACCESS_KEY` are not required. The presigned attachment endpoints are unavailable and return `501`:

```bash
curl -i -X PUT localhost:8080/blog/attachment/example     # 501 Not Implemented
curl -i -X PUT localhost:8080/article/attachment/example  # 501 Not Implemented
```

Attachment uploads therefore only work outside local mode.

### Blob storage backends

Article assets (Excalidraw scenes) go through a `BlobStore`. The local backend writes under `BLOB_LOCAL_PATH` (`data/scenes` by default); keys are rejected if they are absolute or contain `..`. Set `BLOB_BACKEND=s3` to use the configured `S3_ENDPOINT`, `S3_ACCESS_KEY_ID`, `S3_SECRET_ACCESS_KEY`, and `S3_BUCKET` instead. The S3 backend uses the same key layout and streamed reads, so article code is unchanged. Local storage is not durable on Lambda; use the S3 backend there.

## Article API

All write endpoints require `Authorization: Bearer <token>`. Article URIs must match `[a-z0-9-]+`.

| Method | Path | Auth | Purpose |
| --- | --- | --- | --- |
| `POST` | `/article` | yes | Create an empty, unpublished article. Duplicate URI returns `409`. |
| `GET` | `/article` | no | List published articles (no block bodies). |
| `GET` | `/article/:uri` | no | Read a published article with its ordered blocks. Drafts return `404`. |
| `GET` | `/article/all` | yes | List every article, including drafts and unpublished ones. |
| `GET` | `/article/:uri/draft` | yes | Read any article with its blocks, including drafts. |
| `PUT` | `/article/:uri` | yes | Update metadata and blocks; a differing `body.uri` renames the article. |
| `DELETE` | `/article/:uri` | yes | Unpublish the article. Blocks and assets are retained. |
| `DELETE` | `/article/:uri/purge` | yes | Permanently delete the article and every scene file under its UUID prefix. The URI becomes reusable. |
| `PUT` | `/article/:uri/scene/:sid` | yes | Store the scene JSON and return `{ "scene_path": "scenes/<article-uuid>/<sid>.json" }`. |
| `GET` | `/article/:uri/scene/:sid` | no | Read a scene belonging to a published article. |
| `GET` | `/article/:uri/draft/scene/:sid` | yes | Read a scene for a draft. |
| `DELETE` | `/article/:uri/scene/:sid` | yes | Delete a scene. |
| `PUT` | `/article/attachment/:key` | yes | Request a presigned upload under `article/YYYY/MM/`. |

Blocks are an ordered array. Each block has `id` and `type`; payloads are `markdown` (`md`), `image` (`collection`, `filename`, optional `caption`), `music` (`v`, `title`, optional `artist`), and `excalidraw` (`scene_path`). The server validates types, required fields, and unique block IDs, and derives the list `summary` and `cover` on save.

## Music API

Playlists group catalogue songs; a song can belong to more than one playlist. Playlist writes replace the full song list in a transaction.

| Method | Path | Auth | Purpose |
| --- | --- | --- | --- |
| `GET` | `/music` | no | List every song in the catalogue, ordered by `sorting` then creation time. |
| `PUT` | `/music` | yes | Upsert a song by its YouTube video ID (`v`). |
| `DELETE` | `/music/:v` | yes | Soft-delete a song; it disappears from playlists, while article blocks keep their snapshotted values. |
| `GET` | `/music/playlists` | no | List playlists with their songs, ordered by `sorting` then creation time. |
| `PUT` | `/music/playlist` | yes | Create or update a playlist. Send `id` to update; `songs` is the full replacement list and may contain only song IDs. A blank title is rejected and repeated IDs are stored once. |
| `DELETE` | `/music/playlist/:id` | yes | Delete a playlist and its song associations. Songs themselves are kept. |

## Data

The API reads its content from `data/files/*.json`, which is not committed. Generate it with `assets-helper` or from the `grd0-site-data` sub-repo.
