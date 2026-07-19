## Why

The current OAuth2 implementation uses Nextcloud's custom OAuth2 endpoints (`/apps/oauth2/authorize`, `/apps/oauth2/api/v1/token`) and a separate OCS API call to fetch user identity. This approach is non-standard, tightly coupled to Nextcloud's specific implementation, and has a bug where the JWT email claim is populated from `identity.Address` instead of `identity.EMail`. Additionally, there's no token refresh mechanism—users must re-authenticate every 72 hours. Migrating to OpenID Connect provides standardized endpoints, proper identity verification via id_token, PKCE security, and refresh token support.

## What Changes

- Replace Nextcloud OAuth2 endpoints with OIDC discovery (`/.well-known/openid-configuration`)
- Use standard OIDC authorization code flow with GET callback (authorization code in query params)
- Implement PKCE (Proof Key for Code Exchange) with code_verifier embedded in the state parameter
- Verify the OIDC id_token signature using JWKS (JSON Web Key Set) with configurable cache TTL
- Extract user identity from `sub` claim in id_token instead of making a separate OCS API call
- Construct app JWT email claim as `sub@issuer_domain` (e.g., `ken.lam@grd0.net`)
- Request `offline_access` scope to obtain OIDC refresh tokens
- Backend mints its own app JWT (HS256, 15-minute expiry) after verifying the id_token
- Implement app refresh tokens (1-day expiry) that wrap OIDC refresh tokens
- Add `/auth/refresh` endpoint for token refresh without re-authentication
- Store OIDC refresh tokens in-memory on the backend, mapped to app refresh tokens
- Frontend stores app JWT and app refresh token in sessionStorage
- Frontend implements automatic token refresh when JWT is near expiry
- Environment variables change from `NEXTCLOUD_*` to `OIDC_*` prefix

## Capabilities

### New Capabilities
- `oidc-auth`: OpenID Connect authentication flow with PKCE, id_token verification, app JWT minting, and refresh token management

### Modified Capabilities
<!-- No existing auth-related specs to modify -->

## Impact

**Backend (Go API):**
- `api/handler/authenticate.go`: Complete rewrite of OAuth2 flow to OIDC
- `api/server.go`: No changes to JWT middleware (still HS256 with AUTH_JWT_SECRET)
- New dependencies: `github.com/coreos/go-oidc/v3` for OIDC verification
- New in-memory stores: JWKS cache, OIDC refresh token mapping
- Environment variables: Remove `NEXTCLOUD_URL`, `NEXTCLOUD_CLIENT_KEY`, `NEXTCLOUD_CLIENT_SECRET`, `AUTH_STATE_SECRET`; Add `OIDC_ISSUER`, `OIDC_CLIENT_ID`, `OIDC_CLIENT_SECRET`, `JWKS_CACHE_TTL`, `AUTH_REFRESH_SECRET`

**Frontend (Vue SPA):**
- `vue-web/src/router/index.ts`: Callback route stores refresh_token in addition to existing fields
- `vue-web/src/services/api.ts`: Add token refresh logic (check expiry before requests, call `/auth/refresh`)
- New auth service or utility for token management

**Deployment:**
- Environment variable updates required in production
- API restart invalidates active sessions (in-memory refresh token store)

**Security:**
- PKCE prevents authorization code interception
- RS256 signature verification of id_token (asymmetric, no shared secret)
- Short-lived app JWT (15min) reduces token theft window
- Refresh tokens stored server-side, never exposed to frontend
