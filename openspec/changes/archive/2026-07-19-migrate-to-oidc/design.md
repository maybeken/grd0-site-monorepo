## Context

The current authentication system uses Nextcloud's custom OAuth2 implementation with hardcoded endpoints (`/apps/oauth2/authorize`, `/apps/oauth2/api/v1/token`) and a separate OCS API call to fetch user identity. The backend acts as a middleman: it performs the OAuth2 dance, fetches user info, then mints its own JWT (HS256, 72-hour expiry) for the frontend to use on subsequent API calls.

**Current flow:**
1. Frontend redirects to API `/auth/login`
2. API redirects to Nextcloud OAuth2 authorize endpoint
3. User authenticates on Nextcloud
4. Nextcloud redirects to API `/auth/callback` with authorization code
5. API exchanges code for access token, calls OCS API for user info
6. API mints its own JWT and redirects to frontend with token in URL
7. Frontend stores JWT in sessionStorage, uses it for API calls

**Constraints:**
- Single-instance API deployment (no distributed session store)
- SQLite database with GORM (no dedicated user/session tables)
- Frontend is a Vue 3 SPA with sessionStorage-based token storage
- Must maintain backward compatibility with existing protected routes
- Personal CMS with low user count (single admin)

## Goals / Non-Goals

**Goals:**
- Migrate to OpenID Connect for standardized authentication
- Use PKCE (Proof Key for Code Exchange) to prevent authorization code interception
- Verify id_token signature using JWKS (RS256) instead of relying on shared secrets
- Implement token refresh to avoid 72-hour re-authentication cycles
- Use `sub` claim as user identifier (Nextcloud OIDC doesn't release email claim)
- Construct email as `sub@issuer_domain` for app JWT
- Maintain the current app JWT pattern (HS256, backend-minted) for API authorization
- Fix the email claim bug (currently uses `identity.Address` instead of `identity.EMail`)

**Non-Goals:**
- Implement distributed session storage or database-backed refresh tokens
- Add role-based access control or user management
- Support multiple OIDC providers (Nextcloud-only)
- Implement logout flow (can be added later)
- Change the frontend's sessionStorage-based token storage pattern
- Modify protected route logic or JWT middleware signature verification

## Decisions

### 1. PKCE Code Verifier Storage: Embed in State Parameter

**Decision:** Encode the PKCE code_verifier in the OAuth state parameter alongside a time-windowed HMAC for CSRF protection.

**Rationale:** 
- Eliminates need for server-side session storage between `/auth/login` and `/auth/callback`
- Stateless across API restarts (no lost login attempts)
- Simpler than cookie-based or session-based storage
- The state parameter is already required for CSRF protection; extending it is natural

**Alternatives considered:**
- **Server-side session store:** Requires session management, cleanup, lost on restart
- **Signed cookie:** Adds complexity, browser cookie handling, SameSite considerations
- **Separate code_verifier storage:** Overkill for single-instance deployment

**Format:**
```
state = base64url(code_verifier) + "." + HMAC(code_verifier + timestamp, AUTH_STATE_SECRET)
```

On callback: split state, verify HMAC, decode verifier, use for token exchange.

### 2. JWKS Cache: In-Memory with Configurable TTL

**Decision:** Cache JWKS (JSON Web Key Set) in memory with a TTL controlled by the `JWKS_CACHE_TTL` environment variable (default: 1 hour).

**Rationale:**
- JWKS keys change rarely (key rotation events)
- Fetching on every token verification would be wasteful
- In-memory cache is sufficient for single-instance deployment
- Configurable TTL allows tuning without code changes

**Strategy:**
- On token verify: look up `kid` from token header
- Cache hit: use cached public key
- Cache miss: fetch JWKS from provider, update cache, retry
- TTL expiry: next request triggers fresh fetch
- Thread-safe with `sync.RWMutex`

**Alternatives considered:**
- **No caching:** Too many HTTP requests to provider
- **Persistent cache (Redis/DB):** Overkill for single instance, adds dependency
- **Fixed TTL:** Less flexible, requires code changes to tune

### 3. OIDC Refresh Token Storage: In-Memory Map

**Decision:** Store OIDC refresh tokens in an in-memory map, keyed by a generated app refresh token. App refresh tokens are opaque strings returned to the frontend.

**Rationale:**
- OIDC refresh tokens must be kept secret (cannot be exposed to browser)
- In-memory map is simple and sufficient for single-instance deployment
- No database migration required
- On API restart, active sessions are invalidated (acceptable for personal CMS)

**Flow:**
1. On login: generate app refresh token (random UUID), store `{appRT: oidcRT}` in map
2. On refresh: lookup OIDC RT by app RT, exchange with provider, update map with new tokens
3. On logout/expiry: delete from map

**Alternatives considered:**
- **Encrypted in app refresh token:** More complex, requires encryption key management
- **Database storage:** Requires schema migration, overkill for single admin
- **No refresh tokens:** Users re-authenticate every 15 minutes (unacceptable UX)

### 4. App JWT Minting: Backend-Issued HS256

**Decision:** Continue minting app-specific JWTs (HS256, 15-minute expiry) on the backend after verifying the OIDC id_token. Do not pass the id_token directly to the frontend.

**Rationale:**
- Maintains current authorization pattern (no changes to JWT middleware)
- Backend controls claims and expiry independently of provider
- Can add custom claims in the future without provider dependency
- Frontend doesn't need to understand OIDC token format
- HS256 verification is simpler and faster than RS256 on every request

**Alternatives considered:**
- **Trust id_token directly:** Requires RS256 verification on every API request, ties authorization to provider
- **Pass id_token to frontend:** Frontend must parse OIDC format, less control over claims

### 5. Token Refresh: Frontend-Initiated with Expiry Check

**Decision:** Frontend checks JWT expiry before each API request. If JWT is near expiry (< 2 minutes remaining), call `/auth/refresh` with the app refresh token to get a new JWT and refresh token.

**Rationale:**
- Proactive refresh prevents failed requests due to expired tokens
- 2-minute buffer accounts for network latency and clock skew
- Frontend-controlled refresh is simpler than backend middleware interception
- Matches common SPA authentication patterns

**Flow:**
1. Before API request: check `jwt_expires` from sessionStorage
2. If `Date.now() >= expires - 120000` (2 min buffer): call `POST /auth/refresh`
3. Update sessionStorage with new `jwt_token`, `jwt_expires`, `refresh_token`
4. Proceed with original request using new JWT
5. If refresh fails: redirect to login

### 6. Callback Method: GET Redirect

**Decision:** Use standard OIDC authorization code flow with GET callback. The OIDC provider redirects back to `/auth/callback` with `code` and `state` as query parameters.

**Rationale:**
- Simpler implementation that works with all OIDC providers
- Nextcloud OIDC app doesn't support `response_mode=form_post`
- Authorization code is short-lived and single-use, reducing risk from URL exposure
- PKCE provides additional security against code interception
- The API immediately exchanges the code for tokens server-side

**Alternatives considered:**
- **POST callback (form_post):** Nextcloud OIDC doesn't support this response mode
- **One-time exchange code:** More complex, requires additional endpoint and in-memory store

### 7. User Identity: `sub` Claim with Constructed Email

**Decision:** Use the `sub` claim from the id_token as the user identifier. Construct the app JWT email claim as `sub@issuer_domain` (e.g., `ken.lam@grd0.net` by extracting the domain from `OIDC_ISSUER`).

**Rationale:**
- Nextcloud's OIDC implementation does not release the `email` claim in the id_token
- The `sub` claim is always present and contains the username
- Constructing email from `sub` + issuer domain provides a usable email-like identifier for the app JWT
- Avoids requiring provider-side configuration changes to release the email claim

**Alternatives considered:**
- **Require email claim from provider:** Requires Nextcloud OIDC app configuration changes that may not be supported
- **Use only `sub` without domain:** Less useful as an identifier in the app JWT, harder to distinguish from other sources

## Risks / Trade-offs

### Risk: API Restart Invalidates Active Sessions

**Impact:** In-memory refresh token store is lost on restart. Users must re-authenticate.

**Mitigation:** Acceptable for personal CMS with single admin. Can migrate to database-backed storage later if needed.

### Risk: Token Exposure in Redirect URL

**Impact:** App JWT and refresh token are passed to frontend via redirect URL query params. These may be logged by browsers or proxy servers.

**Mitigation:** App JWT expires in 15 minutes (short window). App refresh token expires in 1 day (bounded). This is a personal CMS, not high-security. If this becomes a concern, implement one-time exchange code pattern.

### Risk: JWKS Cache Staleness

**Impact:** If JWKS keys are rotated and cache TTL is long, token verification may fail until cache refreshes.

**Mitigation:** "Fetch on kid miss" pattern ensures immediate refresh when a new key is encountered. Default TTL of 1 hour is conservative. Can reduce TTL if needed.

### Risk: Clock Skew Between API and Nextcloud

**Impact:** id_token expiry (15 minutes) may be rejected if server clocks drift significantly.

**Mitigation:** Add leeway (30 seconds) to token expiry validation. Ensure NTP is configured on both servers.

### Risk: Refresh Storm from Multiple Tabs

**Impact:** If multiple browser tabs detect expired JWT simultaneously, all will call `/auth/refresh`, causing redundant token exchanges.

**Mitigation:** For single-user CMS, unlikely to be a problem. Can add single-flight pattern (mutex per refresh token) if needed.

### Trade-off: Stateless Login vs. Server-Side Session

**Decision:** PKCE verifier in state (stateless) vs. server-side session storage.

**Trade-off:** Stateless is simpler and survives restarts, but state parameter is larger. Acceptable for this use case.

### Trade-off: In-Memory Refresh Tokens vs. Database

**Decision:** In-memory map vs. database-backed storage.

**Trade-off:** In-memory is simpler and requires no migration, but sessions are lost on restart. Acceptable for single-instance personal CMS.

## Migration Plan

### Pre-Deployment

1. **Update environment variables:**
   - Remove: `NEXTCLOUD_URL`, `NEXTCLOUD_CLIENT_KEY`, `NEXTCLOUD_CLIENT_SECRET`, `AUTH_STATE_SECRET`
   - Add: `OIDC_ISSUER`, `OIDC_CLIENT_ID`, `OIDC_CLIENT_SECRET`, `JWKS_CACHE_TTL`, `AUTH_REFRESH_SECRET`

2. **Update Nextcloud OIDC app configuration:**
   - Set redirect URI to `https://api.grd0.net/auth/callback`
   - Note the OIDC client ID and secret
   - Ensure `offline_access` scope is allowed (for refresh tokens)

3. **Deploy new API version:**
   - New binary with OIDC implementation
   - New dependencies: `github.com/coreos/go-oidc/v3`

### Deployment

1. **Deploy API:** Push new version with OIDC support
2. **Update environment variables:** Set new OIDC_* variables, remove NEXTCLOUD_* variables
3. **Restart API:** Pick up new configuration
4. **Test login flow:** Verify end-to-end authentication works
5. **Test refresh flow:** Wait for JWT to near expiry, verify automatic refresh

### Rollback

If issues arise:
1. Revert to previous API binary
2. Restore NEXTCLOUD_* environment variables
3. Restart API

**Note:** Active sessions will be invalidated on rollback (in-memory refresh tokens lost).

### Post-Deployment

1. **Monitor logs:** Watch for authentication errors, JWKS fetch failures
2. **Verify refresh tokens:** Check that refresh flow works as expected
3. **Update documentation:** Document new environment variables and flow

## Open Questions

**None at this time.** All major decisions have been made during exploration.
