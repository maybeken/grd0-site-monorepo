# Purpose

OIDC authentication capability for the grd0.net API and frontend. Enables single sign-on via an external OpenID Connect provider with PKCE, automatic token refresh, and backward-compatible JWT middleware.

# Requirements

### Requirement: OIDC Discovery
The system SHALL discover OIDC endpoints from the issuer URL using the `.well-known/openid-configuration` endpoint.

#### Scenario: Successful discovery
- **WHEN** the system initializes with `OIDC_ISSUER=https://cloud.grd0.net`
- **THEN** the system SHALL fetch `https://cloud.grd0.net/index.php/.well-known/openid-configuration`
- **THEN** the system SHALL extract `authorization_endpoint`, `token_endpoint`, `jwks_uri`, and `userinfo_endpoint`
- **THEN** the system SHALL cache these endpoints for the lifetime of the application

#### Scenario: Discovery failure
- **WHEN** the OIDC issuer is unreachable or returns an error
- **THEN** the system SHALL log the error and fail to start
- **THEN** the system SHALL NOT accept authentication requests

### Requirement: PKCE Code Challenge
The system SHALL use PKCE (Proof Key for Code Exchange) with S256 challenge method for all authorization requests.

#### Scenario: Login request with PKCE
- **WHEN** a user initiates login via `/auth/login`
- **THEN** the system SHALL generate a cryptographically random `code_verifier` (43-128 characters)
- **THEN** the system SHALL compute `code_challenge = base64url(SHA256(code_verifier))`
- **THEN** the system SHALL redirect to the OIDC authorization endpoint with `code_challenge` and `code_challenge_method=S256`
- **THEN** the system SHALL encode the `code_verifier` in the `state` parameter

#### Scenario: State parameter format
- **WHEN** the system generates the state parameter
- **THEN** the state SHALL be formatted as `base64url(code_verifier) + "." + HMAC(code_verifier + timestamp, AUTH_STATE_SECRET)`
- **THEN** the HMAC SHALL use a time-windowed timestamp (2-hour rolling window) for CSRF protection

### Requirement: Authorization Request Parameters
The system SHALL include all required parameters in the authorization request to the OIDC provider.

#### Scenario: Authorization redirect
- **WHEN** the system redirects to the OIDC authorization endpoint
- **THEN** the request SHALL include `response_type=code`
- **THEN** the request SHALL include `scope=openid profile offline_access`
- **THEN** the request SHALL include `client_id` from `OIDC_CLIENT_ID` environment variable
- **THEN** the request SHALL include `redirect_uri` from `AUTH_REDIRECT_URL` environment variable
- **THEN** the request SHALL include `state` parameter with embedded PKCE verifier

### Requirement: GET Callback Handling
The system SHALL accept GET requests from the OIDC provider at `/auth/callback` with the authorization code in query parameters.

#### Scenario: Successful GET callback
- **WHEN** the OIDC provider redirects to `/auth/callback` with `code` and `state` as query parameters
- **THEN** the system SHALL extract the `code` and `state` from the query string
- **THEN** the system SHALL proceed with code exchange

#### Scenario: Invalid state parameter
- **WHEN** the `state` parameter is missing or fails HMAC verification
- **THEN** the system SHALL reject the request with HTTP 400 Bad Request
- **THEN** the system SHALL log the invalid state attempt

### Requirement: Authorization Code Exchange
The system SHALL exchange the authorization code for tokens using the PKCE code_verifier.

#### Scenario: Successful code exchange
- **WHEN** the system has a valid authorization code and code_verifier
- **THEN** the system SHALL POST to the OIDC token endpoint with `grant_type=authorization_code`
- **THEN** the request SHALL include `code`, `redirect_uri`, `client_id`, `client_secret`, and `code_verifier`
- **THEN** the system SHALL receive `id_token`, `access_token`, and optionally `refresh_token` in the response

#### Scenario: Code exchange failure
- **WHEN** the token endpoint returns an error
- **THEN** the system SHALL log the error
- **THEN** the system SHALL return HTTP 500 Internal Server Error
- **THEN** the system SHALL NOT redirect to the frontend

### Requirement: ID Token Verification
The system SHALL verify the id_token signature using the OIDC provider's JWKS (JSON Web Key Set).

#### Scenario: Valid id_token
- **WHEN** the system receives an id_token from the token endpoint
- **THEN** the system SHALL extract the `kid` (key ID) from the token header
- **THEN** the system SHALL fetch the corresponding public key from the JWKS endpoint
- **THEN** the system SHALL verify the token signature using RS256
- **THEN** the system SHALL verify the `iss` claim matches `OIDC_ISSUER`
- **THEN** the system SHALL verify the `aud` claim matches `OIDC_CLIENT_ID`
- **THEN** the system SHALL verify the token is not expired (with 30-second leeway)
- **THEN** the system SHALL verify the `nbf` (not before) claim
- **THEN** the system SHALL extract the `sub` claim as the user identifier

#### Scenario: Invalid signature
- **WHEN** the id_token signature verification fails
- **THEN** the system SHALL reject the token
- **THEN** the system SHALL log the verification failure
- **THEN** the system SHALL return HTTP 401 Unauthorized

### Requirement: JWKS Caching
The system SHALL cache JWKS keys in memory with a configurable TTL to avoid fetching on every token verification.

#### Scenario: Cache hit
- **WHEN** the system needs to verify a token with a `kid` that is in the cache
- **THEN** the system SHALL use the cached public key
- **THEN** the system SHALL NOT make an HTTP request to the JWKS endpoint

#### Scenario: Cache miss
- **WHEN** the system needs to verify a token with a `kid` that is not in the cache
- **THEN** the system SHALL fetch the JWKS from the provider's `jwks_uri`
- **THEN** the system SHALL update the cache with the new keys
- **THEN** the system SHALL retry the verification with the updated cache

#### Scenario: Cache TTL expiry
- **WHEN** the cache age exceeds `JWKS_CACHE_TTL` (default: 1 hour)
- **THEN** the next token verification SHALL trigger a fresh JWKS fetch
- **THEN** the system SHALL update the cache with the new keys

#### Scenario: JWKS fetch failure
- **WHEN** the JWKS endpoint is unreachable
- **THEN** the system SHALL log the error
- **THEN** the system SHALL continue using the cached keys (if any)
- **THEN** the system SHALL retry fetching on the next cache miss

### Requirement: App JWT Minting
The system SHALL mint its own app JWT after successfully verifying the OIDC id_token.

#### Scenario: Successful JWT minting
- **WHEN** the id_token is successfully verified
- **THEN** the system SHALL extract the `sub` claim from the id_token as the user identifier
- **THEN** the system SHALL construct the email claim as `sub@issuer_domain` (e.g., `ken.lam@grd0.net`)
- **THEN** the system SHALL create a new JWT with custom claims: `email` (constructed), `name`
- **THEN** the system SHALL set registered claims: `iss` (OIDC_ISSUER), `nbf` (current time), `exp` (current time + 15 minutes)
- **THEN** the system SHALL sign the JWT using HS256 with `AUTH_JWT_SECRET`
- **THEN** the system SHALL return the signed JWT to the frontend

#### Scenario: Missing sub claim
- **WHEN** the id_token does not contain a `sub` claim
- **THEN** the system SHALL reject the token
- **THEN** the system SHALL log the missing claim
- **THEN** the system SHALL return HTTP 400 Bad Request

### Requirement: Refresh Token Storage
The system SHALL store OIDC refresh tokens in memory, mapped to generated app refresh tokens.

#### Scenario: Storing refresh tokens
- **WHEN** the OIDC token response includes a `refresh_token` (obtained via `offline_access` scope)
- **THEN** the system SHALL generate a random app refresh token (UUID)
- **THEN** the system SHALL store the mapping `{app_refresh_token: oidc_refresh_token}` in memory
- **THEN** the system SHALL return the app refresh token to the frontend
- **THEN** the OIDC refresh token SHALL NOT be exposed to the frontend

#### Scenario: Refresh token lookup
- **WHEN** the system receives a refresh request with an app refresh token
- **THEN** the system SHALL look up the corresponding OIDC refresh token
- **THEN** if found, the system SHALL proceed with token refresh
- **THEN** if not found, the system SHALL return HTTP 401 Unauthorized

#### Scenario: Refresh token cleanup
- **WHEN** the API restarts
- **THEN** all in-memory refresh token mappings SHALL be lost
- **THEN** users SHALL need to re-authenticate

### Requirement: Token Refresh Endpoint
The system SHALL provide a `/auth/refresh` endpoint for refreshing expired app JWTs using app refresh tokens.

#### Scenario: Successful refresh
- **WHEN** the frontend POSTs to `/auth/refresh` with a valid app refresh token
- **THEN** the system SHALL look up the corresponding OIDC refresh token
- **THEN** the system SHALL POST to the OIDC token endpoint with `grant_type=refresh_token` and the OIDC refresh token
- **THEN** the system SHALL receive a new `id_token` and optionally a new `refresh_token`
- **THEN** the system SHALL verify the new id_token using JWKS
- **THEN** the system SHALL mint a new app JWT (15-minute expiry)
- **THEN** if a new OIDC refresh token is received, the system SHALL update the in-memory mapping
- **THEN** the system SHALL return the new app JWT and app refresh token to the frontend

#### Scenario: Invalid refresh token
- **WHEN** the app refresh token is not found in the in-memory store
- **THEN** the system SHALL return HTTP 401 Unauthorized
- **THEN** the frontend SHALL redirect to login

#### Scenario: OIDC refresh failure
- **WHEN** the OIDC token endpoint rejects the refresh request
- **THEN** the system SHALL log the error
- **THEN** the system SHALL delete the invalid refresh token from the in-memory store
- **THEN** the system SHALL return HTTP 401 Unauthorized
- **THEN** the frontend SHALL redirect to login

### Requirement: Frontend Token Storage
The frontend SHALL store the app JWT, app refresh token, and expiry timestamp in sessionStorage.

#### Scenario: Storing tokens after login
- **WHEN** the frontend receives the authentication callback with `token`, `refresh`, and `expires` query parameters
- **THEN** the frontend SHALL store `jwt_token` in sessionStorage
- **THEN** the frontend SHALL store `refresh_token` in sessionStorage
- **THEN** the frontend SHALL store `jwt_expires` in sessionStorage
- **THEN** the frontend SHALL navigate to the home page

#### Scenario: Token retrieval
- **WHEN** the frontend makes an authenticated API request
- **THEN** the frontend SHALL retrieve `jwt_token` from sessionStorage
- **THEN** the frontend SHALL include it in the `Authorization: Bearer <token>` header

### Requirement: Automatic Token Refresh
The frontend SHALL automatically refresh the app JWT when it is near expiry.

#### Scenario: Proactive refresh
- **WHEN** the frontend is about to make an API request
- **THEN** the frontend SHALL check if `Date.now() >= jwt_expires - 120000` (2-minute buffer)
- **THEN** if true, the frontend SHALL call `POST /auth/refresh` with the `refresh_token`
- **THEN** the frontend SHALL update sessionStorage with the new `jwt_token`, `refresh_token`, and `jwt_expires`
- **THEN** the frontend SHALL proceed with the original API request using the new JWT

#### Scenario: Refresh failure
- **WHEN** the `/auth/refresh` endpoint returns an error
- **THEN** the frontend SHALL clear sessionStorage
- **THEN** the frontend SHALL redirect to the login page

#### Scenario: No refresh needed
- **WHEN** the JWT has more than 2 minutes until expiry
- **THEN** the frontend SHALL proceed with the API request using the current JWT

### Requirement: Environment Variable Configuration
The system SHALL use environment variables for all OIDC configuration.

#### Scenario: Required environment variables
- **WHEN** the API starts
- **THEN** the system SHALL require `OIDC_ISSUER` (e.g., `https://cloud.grd0.net`)
- **THEN** the system SHALL require `OIDC_CLIENT_ID`
- **THEN** the system SHALL require `OIDC_CLIENT_SECRET`
- **THEN** the system SHALL require `AUTH_JWT_SECRET` (for app JWT signing)
- **THEN** the system SHALL require `AUTH_REDIRECT_URL` (OIDC callback URL)
- **THEN** the system SHALL require `AUTH_CALLBACK_URL` (frontend callback URL)
- **THEN** the system SHALL require `AUTH_STATE_SECRET` (for PKCE state HMAC)

#### Scenario: Optional environment variables
- **WHEN** the API starts
- **THEN** the system SHALL accept `JWKS_CACHE_TTL` (default: `1h`)
- **THEN** the system SHALL parse the TTL as a duration string (e.g., `1h`, `30m`, `2h`)

#### Scenario: Missing required variable
- **WHEN** a required environment variable is not set
- **THEN** the system SHALL panic with a clear error message
- **THEN** the system SHALL NOT start

### Requirement: Backward Compatibility
The system SHALL maintain backward compatibility with existing protected routes and JWT middleware.

#### Scenario: Protected route access
- **WHEN** a user makes a request to a protected route with a valid app JWT
- **THEN** the system SHALL validate the JWT using HS256 with `AUTH_JWT_SECRET`
- **THEN** the system SHALL extract claims from the JWT
- **THEN** the system SHALL allow access to the protected route

#### Scenario: Invalid JWT
- **WHEN** a user makes a request with an invalid or expired JWT
- **THEN** the system SHALL return HTTP 401 Unauthorized
- **THEN** the frontend SHALL redirect to login or attempt refresh

### Requirement: Error Logging
The system SHALL log authentication errors for debugging and monitoring.

#### Scenario: Error logging
- **WHEN** an authentication error occurs (invalid token, failed verification, missing claims)
- **THEN** the system SHALL log the error with sufficient context (user email if available, error type)
- **THEN** the system SHALL NOT log sensitive information (tokens, secrets, refresh tokens)
