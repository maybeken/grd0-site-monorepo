package handler

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-jose/go-jose/v4"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"grd0.net/api/utils"
)

var (
	oidcIssuer       = utils.GetEnv("OIDC_ISSUER")
	oidcClientID     = utils.GetEnv("OIDC_CLIENT_ID")
	oidcClientSecret = utils.GetEnv("OIDC_CLIENT_SECRET")
	authRedirectURL  = utils.GetEnv("AUTH_REDIRECT_URL")
	authCallbackURL  = utils.GetEnv("AUTH_CALLBACK_URL")
	authStateSecret  = utils.GetEnv("AUTH_STATE_SECRET")
	authJWTSecret    = utils.GetEnv("AUTH_JWT_SECRET")
	jwksCacheTTLStr  = utils.GetEnvWithFallback("JWKS_CACHE_TTL", "1h")
)

type JwtCustomClaims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	jwt.RegisteredClaims
}

type jwksCacheManager struct {
	mu        sync.RWMutex
	keys      map[string]*jose.JSONWebKey
	fetchedAt time.Time
	ttl       time.Duration
	jwksURI   string
}

func (c *jwksCacheManager) fetchKeys(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.jwksURI, nil)
	if err != nil {
		return fmt.Errorf("failed to create JWKS request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("JWKS fetch returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read JWKS response: %w", err)
	}

	var jwks jose.JSONWebKeySet
	if err := json.Unmarshal(body, &jwks); err != nil {
		return fmt.Errorf("failed to parse JWKS: %w", err)
	}

	newKeys := make(map[string]*jose.JSONWebKey)
	for i := range jwks.Keys {
		k := jwks.Keys[i]
		newKeys[k.KeyID] = &k
	}

	c.mu.Lock()
	c.keys = newKeys
	c.fetchedAt = time.Now()
	c.mu.Unlock()

	return nil
}

func (c *jwksCacheManager) getKey(ctx context.Context, kid string) (*jose.JSONWebKey, error) {
	c.mu.RLock()
	if time.Since(c.fetchedAt) > c.ttl {
		c.mu.RUnlock()
		if err := c.fetchKeys(ctx); err != nil {
			return nil, err
		}
		c.mu.RLock()
	}

	key, ok := c.keys[kid]
	c.mu.RUnlock()

	if ok {
		return key, nil
	}

	if err := c.fetchKeys(ctx); err != nil {
		return nil, err
	}

	c.mu.RLock()
	key, ok = c.keys[kid]
	c.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("key %s not found in JWKS", kid)
	}

	return key, nil
}

var (
	oidcProvider  *oidc.Provider
	jwksCacheInst *jwksCacheManager
	tokenEndpoint string

	refreshTokens = make(map[string]string)
	refreshMu     sync.RWMutex
)

func InitOIDC(ctx context.Context) error {
	var err error
	oidcProvider, err = oidc.NewProvider(ctx, oidcIssuer)
	if err != nil {
		return fmt.Errorf("OIDC discovery failed for issuer %s: %w", oidcIssuer, err)
	}

	var providerClaims struct {
		TokenEndpoint string `json:"token_endpoint"`
		JWKSURI       string `json:"jwks_uri"`
	}
	if err := oidcProvider.Claims(&providerClaims); err != nil {
		return fmt.Errorf("failed to parse OIDC provider claims: %w", err)
	}

	tokenEndpoint = providerClaims.TokenEndpoint

	ttl, err := time.ParseDuration(jwksCacheTTLStr)
	if err != nil {
		ttl = time.Hour
	}

	jwksCacheInst = &jwksCacheManager{
		keys:    make(map[string]*jose.JSONWebKey),
		ttl:     ttl,
		jwksURI: providerClaims.JWKSURI,
	}

	return nil
}

func generateCodeVerifier() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func computeCodeChallenge(verifier string) string {
	h := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h[:])
}

func encodeState(verifier string) string {
	currentUnix := time.Now().UnixMilli()
	duration := 2 * time.Hour.Milliseconds()
	timeWindow := currentUnix - currentUnix%duration + duration

	mac := hmac.New(sha256.New, []byte(authStateSecret))
	mac.Write([]byte(fmt.Sprintf("%s:%d", verifier, timeWindow)))
	sig := fmt.Sprintf("%x", mac.Sum(nil))

	return base64.RawURLEncoding.EncodeToString([]byte(verifier)) + "." + sig
}

func decodeState(state string) (string, error) {
	parts := strings.SplitN(state, ".", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid state format")
	}

	verifierBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", fmt.Errorf("failed to decode verifier from state: %w", err)
	}
	verifier := string(verifierBytes)

	currentUnix := time.Now().UnixMilli()
	duration := 2 * time.Hour.Milliseconds()
	timeWindow := currentUnix - currentUnix%duration + duration

	mac := hmac.New(sha256.New, []byte(authStateSecret))
	mac.Write([]byte(fmt.Sprintf("%s:%d", verifier, timeWindow)))
	expectedSig := fmt.Sprintf("%x", mac.Sum(nil))

	if !hmac.Equal([]byte(parts[1]), []byte(expectedSig)) {
		return "", fmt.Errorf("state HMAC verification failed")
	}

	return verifier, nil
}

func verifyIDToken(ctx context.Context, rawIDToken string) (*oidc.IDToken, error) {
	verifier := oidc.NewVerifier(oidcIssuer, &oidcKeySet{cache: jwksCacheInst}, &oidc.Config{
		ClientID: oidcClientID,
	})

	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("id_token verification failed: %w", err)
	}

	return idToken, nil
}

type oidcKeySet struct {
	cache *jwksCacheManager
}

func (ks *oidcKeySet) VerifySignature(ctx context.Context, jwtStr string) ([]byte, error) {
	jws, err := jose.ParseSigned(jwtStr, []jose.SignatureAlgorithm{jose.RS256})
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT: %w", err)
	}

	if len(jws.Signatures) == 0 {
		return nil, fmt.Errorf("JWT has no signatures")
	}

	kid := jws.Signatures[0].Header.KeyID
	if kid == "" {
		return nil, fmt.Errorf("JWT header missing kid")
	}

	key, err := ks.cache.getKey(ctx, kid)
	if err != nil {
		return nil, err
	}

	payload, err := jws.Verify(key)
	if err != nil {
		return nil, fmt.Errorf("signature verification failed: %w", err)
	}

	return payload, nil
}

type idTokenClaims struct {
	Sub  string `json:"sub"`
	Name string `json:"name"`
}

func extractClaims(idToken *oidc.IDToken) (*idTokenClaims, error) {
	var claims idTokenClaims
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to extract claims: %w", err)
	}

	if claims.Sub == "" {
		return nil, fmt.Errorf("id_token missing required sub claim")
	}

	return &claims, nil
}

func mintAppJWT(sub, name string) (string, time.Time, error) {
	expiresAt := time.Now().Add(15 * time.Minute)

	email := sub
	if u, err := url.Parse(oidcIssuer); err == nil {
		host := u.Hostname()
		parts := strings.Split(host, ".")
		if len(parts) >= 2 {
			domain := strings.Join(parts[len(parts)-2:], ".")
			email = sub + "@" + domain
		}
	}

	claims := &JwtCustomClaims{
		Email: email,
		Name:  name,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    oidcIssuer,
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(authJWTSecret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign app JWT: %w", err)
	}

	return signed, expiresAt, nil
}

func storeRefreshToken(appRT, oidcRT string) {
	refreshMu.Lock()
	refreshTokens[appRT] = oidcRT
	refreshMu.Unlock()
}

func lookupRefreshToken(appRT string) (string, bool) {
	refreshMu.RLock()
	oidcRT, ok := refreshTokens[appRT]
	refreshMu.RUnlock()
	return oidcRT, ok
}

func deleteRefreshToken(appRT string) {
	refreshMu.Lock()
	delete(refreshTokens, appRT)
	refreshMu.Unlock()
}

type tokenResponse struct {
	IDToken      string `json:"id_token"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

func exchangeCode(ctx context.Context, code, verifier string) (*tokenResponse, error) {
	data := fmt.Sprintf(
		"grant_type=authorization_code&code=%s&redirect_uri=%s&client_id=%s&client_secret=%s&code_verifier=%s",
		code, authRedirectURL, oidcClientID, oidcClientSecret, verifier,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenEndpoint, strings.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token endpoint request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token endpoint returned status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp tokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &tokenResp, nil
}

func refreshOIDCToken(ctx context.Context, oidcRT string) (*tokenResponse, error) {
	data := fmt.Sprintf(
		"grant_type=refresh_token&refresh_token=%s&client_id=%s&client_secret=%s",
		oidcRT, oidcClientID, oidcClientSecret,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenEndpoint, strings.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read refresh response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("refresh endpoint returned status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp tokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse refresh response: %w", err)
	}

	return &tokenResp, nil
}

func (h *Handler) Login(c echo.Context) error {
	verifier, err := generateCodeVerifier()
	if err != nil {
		h.Logger.Errorf("PKCE verifier generation failed: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate PKCE verifier")
	}

	challenge := computeCodeChallenge(verifier)
	state := encodeState(verifier)

	params := fmt.Sprintf(
		"response_type=code&scope=%s&client_id=%s&redirect_uri=%s&state=%s&code_challenge=%s&code_challenge_method=S256",
		"openid profile offline_access", oidcClientID, authRedirectURL, state, challenge,
	)

	var providerEndpoints struct {
		AuthEndpoint string `json:"authorization_endpoint"`
	}
	if err := oidcProvider.Claims(&providerEndpoints); err != nil {
		h.Logger.Errorf("failed to get authorization endpoint: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get OIDC endpoints")
	}

	authURL := providerEndpoints.AuthEndpoint + "?" + params
	return c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func (h *Handler) LoginCallback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	if code == "" {
		h.Logger.Errorf("callback missing authorization code")
		return echo.NewHTTPError(http.StatusBadRequest, "missing authorization code")
	}

	if state == "" {
		h.Logger.Errorf("callback missing state parameter")
		return echo.NewHTTPError(http.StatusBadRequest, "missing state parameter")
	}

	verifier, err := decodeState(state)
	if err != nil {
		h.Logger.Errorf("state verification failed: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid state parameter")
	}

	tResp, err := exchangeCode(c.Request().Context(), code, verifier)
	if err != nil {
		h.Logger.Errorf("code exchange failed: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to exchange authorization code")
	}

	if tResp.IDToken == "" {
		h.Logger.Errorf("token response missing id_token")
		return echo.NewHTTPError(http.StatusInternalServerError, "missing id_token in token response")
	}

	idToken, err := verifyIDToken(c.Request().Context(), tResp.IDToken)
	if err != nil {
		h.Logger.Errorf("id_token verification failed: %v", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "id_token verification failed")
	}

	var claimsMap map[string]interface{}
	if err := idToken.Claims(&claimsMap); err == nil {
		h.Logger.Debugf("id_token claims: %+v", claimsMap)
	}

	claims, err := extractClaims(idToken)
	if err != nil {
		h.Logger.Errorf("claim extraction failed: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to extract user claims")
	}

	h.Logger.Infof("Authentication for %s", claims.Sub)

	appJWT, expiresAt, err := mintAppJWT(claims.Sub, claims.Name)
	if err != nil {
		h.Logger.Errorf("JWT minting failed: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create session token")
	}

	var appRT string
	if tResp.RefreshToken != "" {
		appRT = uuid.New().String()
		storeRefreshToken(appRT, tResp.RefreshToken)
	}

	redirectURL := fmt.Sprintf("%s?token=%s&refresh=%s&expires_at=%d",
		authCallbackURL, appJWT, appRT, expiresAt.UnixMilli())

	return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

type refreshRequest struct {
	Refresh string `json:"refresh"`
}

type refreshResponse struct {
	Token     string `json:"token"`
	Refresh   string `json:"refresh"`
	ExpiresAt int64  `json:"expires_at"`
}

func (h *Handler) Refresh(c echo.Context) error {
	var req refreshRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if req.Refresh == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing refresh token")
	}

	oidcRT, ok := lookupRefreshToken(req.Refresh)
	if !ok {
		h.Logger.Errorf("refresh token not found")
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid refresh token")
	}

	tResp, err := refreshOIDCToken(c.Request().Context(), oidcRT)
	if err != nil {
		h.Logger.Errorf("OIDC refresh failed: %v", err)
		deleteRefreshToken(req.Refresh)
		return echo.NewHTTPError(http.StatusUnauthorized, "refresh failed")
	}

	if tResp.IDToken == "" {
		h.Logger.Errorf("refresh response missing id_token")
		deleteRefreshToken(req.Refresh)
		return echo.NewHTTPError(http.StatusUnauthorized, "missing id_token in refresh response")
	}

	idToken, err := verifyIDToken(c.Request().Context(), tResp.IDToken)
	if err != nil {
		h.Logger.Errorf("id_token verification failed on refresh: %v", err)
		deleteRefreshToken(req.Refresh)
		return echo.NewHTTPError(http.StatusUnauthorized, "id_token verification failed")
	}

	claims, err := extractClaims(idToken)
	if err != nil {
		h.Logger.Errorf("claim extraction failed on refresh: %v", err)
		deleteRefreshToken(req.Refresh)
		return echo.NewHTTPError(http.StatusUnauthorized, "failed to extract claims")
	}

	appJWT, expiresAt, err := mintAppJWT(claims.Sub, claims.Name)
	if err != nil {
		h.Logger.Errorf("JWT minting failed on refresh: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create session token")
	}

	newAppRT := req.Refresh
	if tResp.RefreshToken != "" {
		deleteRefreshToken(req.Refresh)
		newAppRT = uuid.New().String()
		storeRefreshToken(newAppRT, tResp.RefreshToken)
	}

	return c.JSON(http.StatusOK, refreshResponse{
		Token:     appJWT,
		Refresh:   newAppRT,
		ExpiresAt: expiresAt.UnixMilli(),
	})
}
