package handler

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"grd0.net/api/utils"
)

type JwtCustomClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type LoginResponseBody struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (h *Handler) Login(c echo.Context) error {
	expires_at := time.Now().Add(time.Hour * 72)

	// Set custom claims
	claims := &JwtCustomClaims{
		"ken.lam@grd0.net",
		jwt.RegisteredClaims{
			Issuer:    "api.grd0.net",
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

	return c.JSON(http.StatusOK, LoginResponseBody{
		Token:     t,
		ExpiresAt: expires_at,
	})
}
