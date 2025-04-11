package handler

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"

	"grd0.net/api/utils"
)

var (
	nextcloudUrl = utils.GetEnv("NEXTCLOUD_URL")

	oauthConfig = &oauth2.Config{
		ClientID:     utils.GetEnv("NEXTCLOUD_CLIENT_KEY"),
		ClientSecret: utils.GetEnv("NEXTCLOUD_CLIENT_SECRET"),
		RedirectURL:  utils.GetEnv("AUTH_REDIRECT_URL"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  nextcloudUrl + "/apps/oauth2/authorize",
			TokenURL: nextcloudUrl + "/apps/oauth2/api/v1/token",
		},
		Scopes: []string{},
	}
)

type NextcloudUser struct {
	EMail       string `json:"email"`
	DisplayName string `json:"display-name"`
	ID          string `json:"id"`
	Address     string `json:"address"`
}

type JwtCustomClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func (h *Handler) Login(c echo.Context) error {
	url := oauthConfig.AuthCodeURL(stateTokenGenerator())
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) LoginCallback(c echo.Context) error {
	state := c.QueryParam("state")
	if state != stateTokenGenerator() {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid oauth state")
	}

	code := c.QueryParam("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		h.Logger.Error(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get token")
	}

	client := oauthConfig.Client(context.Background(), token)
	resp, err := client.Get(nextcloudUrl + "/ocs/v2.php/cloud/user?format=json")
	if err != nil || resp.StatusCode != http.StatusOK {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get user info")
	}
	defer resp.Body.Close()

	bits, err := io.ReadAll(resp.Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not parse user info")
	}

	identity, err := userFromReader(bytes.NewReader(bits))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not parse user info")
	}

	h.Logger.Infof("Authentication for %s", identity.EMail)

	expires_at := time.Now().Add(time.Hour * 72)

	// Set custom claims
	claims := &JwtCustomClaims{
		identity.Address,
		jwt.RegisteredClaims{
			Issuer:    nextcloudUrl,
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expires_at),
		},
	}

	// Create token with claims
	jwt_token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := jwt_token.SignedString([]byte(utils.GetEnv("AUTH_JWT_SECRET")))
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s?token=%s&expires_at=%d", utils.GetEnv("AUTH_CALLBACK_URL"), t, expires_at.UnixMilli()))
}

func stateTokenGenerator() string {
	current_unix := time.Now().UnixMilli()
	duration := 2 * time.Hour.Milliseconds()
	expiration := current_unix - current_unix%duration + duration
	token_raw := []byte(fmt.Sprintf("%s:%d", utils.GetEnv("AUTH_STATE_SECRET"), expiration))

	hash := sha256.New()
	hash.Write(token_raw)
	hashedData := hash.Sum(nil)
	hexString := hex.EncodeToString(hashedData)

	return hexString
}

func userFromReader(r io.Reader) (*NextcloudUser, error) {
	u := struct {
		Ocs struct {
			Data struct {
				NextcloudUser
			}
		} `json:"ocs"`
	}{}
	if err := json.NewDecoder(r).Decode(&u); err != nil {
		return nil, err
	}

	return &NextcloudUser{
		EMail:       u.Ocs.Data.EMail,
		DisplayName: u.Ocs.Data.DisplayName,
		ID:          u.Ocs.Data.ID,
		Address:     u.Ocs.Data.Address,
	}, nil
}
