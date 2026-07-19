## 1. Backend Dependencies and Configuration

- [x] 1.1 Add `github.com/coreos/go-oidc/v3` dependency to `api/go.mod`
- [x] 1.2 Update environment variable reading in `api/utils/utils.go` to support new OIDC variables
- [x] 1.3 Add new environment variables to deployment configuration: `OIDC_ISSUER`, `OIDC_CLIENT_ID`, `OIDC_CLIENT_SECRET`, `JWKS_CACHE_TTL`, `AUTH_REFRESH_SECRET`
- [x] 1.4 Remove old environment variables from deployment configuration: `NEXTCLOUD_URL`, `NEXTCLOUD_CLIENT_KEY`, `NEXTCLOUD_CLIENT_SECRET`

## 2. Backend OIDC Provider Setup

- [x] 2.1 Create OIDC provider initialization in `api/handler/authenticate.go` using `oidc.NewProvider(ctx, OIDC_ISSUER)`
- [x] 2.2 Implement JWKS cache manager with in-memory storage and configurable TTL
- [x] 2.3 Add thread-safe JWKS cache with `sync.RWMutex` for concurrent access
- [x] 2.4 Implement JWKS fetch function with fallback to cache on failure
- [x] 2.5 Implement "fetch on kid miss" logic to refresh JWKS when key not found
- [x] 2.6 Add TTL-based cache expiry check before using cached keys

## 3. Backend PKCE Implementation

- [x] 3.1 Implement PKCE code_verifier generation (cryptographically random, 43-128 characters)
- [x] 3.2 Implement code_challenge computation: `base64url(SHA256(code_verifier))`
- [x] 3.3 Implement state parameter encoding: `base64url(code_verifier) + "." + HMAC(code_verifier + timestamp, AUTH_STATE_SECRET)`
- [x] 3.4 Implement state parameter decoding and HMAC verification
- [x] 3.5 Implement time-windowed timestamp (2-hour rolling window) for CSRF protection

## 4. Backend Login Flow

- [x] 4.1 Update `Login()` handler to generate PKCE code_verifier and code_challenge
- [x] 4.2 Update `Login()` handler to encode code_verifier in state parameter
- [x] 4.3 Update `Login()` handler to redirect to OIDC authorization endpoint with standard GET callback
- [x] 4.4 Add required parameters to authorization redirect: `response_type=code`, `scope=openid profile offline_access`, `client_id`, `redirect_uri`, `state`, `code_challenge`, `code_challenge_method=S256`

## 5. Backend Callback Flow

- [x] 5.1 Keep `LoginCallback()` handler as GET method
- [x] 5.2 Update `LoginCallback()` to extract `code` and `state` from query parameters
- [x] 5.3 Update `LoginCallback()` to decode state parameter and recover code_verifier
- [x] 5.4 Update `LoginCallback()` to verify state HMAC and timestamp
- [x] 5.5 Implement authorization code exchange: POST to token endpoint with `grant_type=authorization_code`, `code`, `redirect_uri`, `client_id`, `client_secret`, `code_verifier`
- [x] 5.6 Parse token response to extract `id_token`, `access_token`, and `refresh_token`

## 6. Backend ID Token Verification

- [x] 6.1 Implement id_token signature verification using JWKS cache
- [x] 6.2 Extract `kid` from id_token header and lookup corresponding public key
- [x] 6.3 Verify id_token claims: `iss` matches `OIDC_ISSUER`, `aud` matches `OIDC_CLIENT_ID`
- [x] 6.4 Verify id_token expiry with 30-second leeway for clock skew
- [x] 6.5 Verify `nbf` (not before) claim
- [x] 6.6 Extract `sub` claim from id_token as user identifier
- [x] 6.7 Validate that `sub` claim is present

## 7. Backend App JWT Minting

- [x] 7.1 Update JWT minting logic to create app JWT after id_token verification
- [x] 7.2 Set app JWT custom claims: `email` (constructed as `sub@issuer_domain`), `name` from id_token
- [x] 7.3 Set app JWT registered claims: `iss` (OIDC_ISSUER), `nbf` (current time), `exp` (current time + 15 minutes)
- [x] 7.4 Sign app JWT using HS256 with `AUTH_JWT_SECRET`
- [x] 7.5 Remove old Nextcloud OCS API call and `NextcloudUser` struct
- [x] 7.6 Remove old `userFromReader()` function

## 8. Backend Refresh Token Storage

- [x] 8.1 Create in-memory refresh token store: `map[string]string` (app RT → OIDC RT)
- [x] 8.2 Add `sync.RWMutex` for thread-safe access to refresh token store
- [x] 8.3 Implement app refresh token generation (random UUID)
- [x] 8.4 Store OIDC refresh token mapped to app refresh token after successful login
- [x] 8.5 Implement refresh token lookup function by app refresh token
- [x] 8.6 Implement refresh token deletion function for cleanup

## 9. Backend Refresh Endpoint

- [x] 9.1 Create new `/auth/refresh` POST endpoint in `api/handler/handler.go`
- [x] 9.2 Implement `Refresh()` handler to accept `{refresh: app_refresh_token}` in request body
- [x] 9.3 Lookup OIDC refresh token by app refresh token
- [x] 9.4 POST to OIDC token endpoint with `grant_type=refresh_token` and OIDC refresh token
- [x] 9.5 Parse refresh response to extract new `id_token` and optionally new `refresh_token`
- [x] 9.6 Verify new id_token using JWKS cache
- [x] 9.7 Mint new app JWT (15-minute expiry) from new id_token
- [x] 9.8 If new OIDC refresh token received, update in-memory store
- [x] 9.9 Return new app JWT and app refresh token to frontend
- [x] 9.10 Handle refresh failures: delete invalid token from store, return HTTP 401

## 10. Backend Callback Redirect

- [x] 10.1 Update `LoginCallback()` to redirect to `AUTH_CALLBACK_URL` with query parameters: `token` (app JWT), `refresh` (app RT), `expires_at` (expiry timestamp)
- [x] 10.2 Ensure redirect uses HTTP 302 Temporary Redirect
- [x] 10.3 Add error handling for redirect failures

## 11. Frontend Token Storage

- [x] 11.1 Update `/redirect/auth/callback` route in `vue-web/src/router/index.ts` to extract `refresh` query parameter
- [x] 11.2 Store `refresh_token` in sessionStorage alongside existing `jwt_token` and `jwt_expires`
- [x] 11.3 Update callback route to handle missing refresh token (redirect to login)

## 12. Frontend Token Refresh Logic

- [x] 12.1 Create auth service or utility in `vue-web/src/services/` for token management
- [x] 12.2 Implement function to check if JWT is near expiry (< 2 minutes remaining)
- [x] 12.3 Implement `POST /auth/refresh` call with refresh token
- [x] 12.4 Update sessionStorage with new `jwt_token`, `refresh_token`, and `jwt_expires` after refresh
- [x] 12.5 Integrate refresh check into `adminInstance` in `vue-web/src/services/api.ts` before each request
- [x] 12.6 Handle refresh failure: clear sessionStorage, redirect to login

## 13. Backend Error Handling and Logging

- [x] 13.1 Add error logging for OIDC discovery failures
- [x] 13.2 Add error logging for JWKS fetch failures
- [x] 13.3 Add error logging for id_token verification failures
- [x] 13.4 Add error logging for code exchange failures
- [x] 13.5 Add error logging for refresh token failures
- [x] 13.6 Ensure sensitive information (tokens, secrets) is not logged

## 14. Backend Route Registration

- [x] 14.1 Register `/auth/refresh` as a public route in `api/handler/handler.go`
- [x] 14.2 Ensure `/auth/login` and `/auth/callback` remain public routes
- [x] 14.3 Verify all protected routes still use JWT middleware correctly

## 15. Testing and Verification

- [ ] 15.1 Test OIDC discovery endpoint fetch and parsing
- [ ] 15.2 Test PKCE code_verifier generation and code_challenge computation
- [ ] 15.3 Test state parameter encoding and decoding
- [ ] 15.4 Test login flow end-to-end: redirect to Nextcloud, GET callback, token exchange
- [ ] 15.5 Test id_token verification with valid and invalid signatures
- [ ] 15.6 Test JWKS caching: cache hit, cache miss, TTL expiry
- [ ] 15.7 Test app JWT minting with correct claims and expiry
- [ ] 15.8 Test refresh token storage and lookup
- [ ] 15.9 Test refresh endpoint: successful refresh, invalid token, OIDC failure
- [ ] 15.10 Test frontend token storage and retrieval
- [ ] 15.11 Test frontend automatic refresh when JWT near expiry
- [ ] 15.12 Test protected route access with valid and expired JWTs

## 16. Documentation and Deployment

- [x] 16.1 Update environment variable documentation with new OIDC variables
- [ ] 16.2 Document the new authentication flow in README or deployment guide
- [ ] 16.3 Create deployment checklist for environment variable migration
- [ ] 16.4 Test deployment with new environment variables
- [ ] 16.5 Verify rollback procedure: revert binary, restore old environment variables
