package handler

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"grd0.net/api/utils"
)

func LocalAuthGuard() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := &JwtCustomClaims{
				Email: utils.GetEnvWithFallback("LOCAL_AUTHOR_EMAIL", "local@grd0.net"),
				Name:  utils.GetEnvWithFallback("LOCAL_AUTHOR_NAME", "Local Author"),
			}

			c.Set("user", jwt.NewWithClaims(jwt.SigningMethodHS256, claims))

			return next(c)
		}
	}
}
