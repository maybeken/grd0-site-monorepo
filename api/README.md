> Partly generated with ChatGPT. Treat the details with care and verify against the code.

# grd0.net API

Go 1.26 BFF built with Echo v4, GORM/SQLite and the MinIO client for S3. Module path: `grd0.net/api`.

## Commands

```bash
# From this directory
go run .                              # listens on :80 by default
go build -o ../grd0_blog_api .
```

The Docker image is built from the repository root: `cd api; go build -o ../grd0_blog_api`.

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
```

## Data

The API reads its content from `data/files/*.json`, which is not committed. Generate it with `assets-helper` or from the `grd0-site-data` sub-repo.
