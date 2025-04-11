package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"

	"grd0.net/api/utils"
)

type JwtCustomClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func (h *Handler) Login(c echo.Context) error {
	q := c.Request().URL.Query()
	q.Add("provider", "nextcloud")
	c.Request().URL.RawQuery = q.Encode()

	req := c.Request()
	res := c.Response().Writer
	if gothUser, err := gothic.CompleteUserAuth(res, req); err == nil {
		return c.JSON(http.StatusOK, gothUser)
	}
	gothic.BeginAuthHandler(res, req)
	return nil
}

func (h *Handler) LoginCallback(c echo.Context) error {
	req := c.Request()
	res := c.Response().Writer
	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	email := user.Email
	issuer := user.Provider

	h.Logger.Infof("Authentication from %s for %s", issuer, email)

	expires_at := time.Now().Add(time.Hour * 72)

	// Set custom claims
	claims := &JwtCustomClaims{
		email,
		jwt.RegisteredClaims{
			Issuer:    issuer,
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expires_at),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(utils.GetEnv("JWT_SECRET")))
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s?token=%s&expires_at=%d", utils.GetEnv("AUTH_CALLBACK_URL"), t, expires_at.UnixMilli()))
}
