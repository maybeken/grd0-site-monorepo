package main

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
	"grd0.net/api/database"
	"grd0.net/api/handler"
	"grd0.net/api/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
)

func main() {
	log := utils.InitiateLogger()

	db_path := utils.GetEnvWithFallback("DB_PATH", "data/api.db")

	db := database.OpenDatabase(db_path)
	database.AutoMigrate(db)

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "https://grd0.net"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(middleware.RequestID())

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:        true,
		LogStatus:     true,
		Skipper:       utils.EchoLogSkipper,
		LogValuesFunc: utils.EchoLogValueFunc,
	}))

	rateLimitConfig := middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: rate.Limit(20), Burst: 30, ExpiresIn: 3 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}
	e.Use(middleware.RateLimiterWithConfig(rateLimitConfig))
	e.Use(middleware.Decompress())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	h := &handler.Handler{DB: db, Logger: log}
	handler.RegisterRouter(e, h, echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(handler.JwtCustomClaims)
		},
		SigningKey: []byte(utils.GetEnv("JWT_SECRET")),
	}))

	e.Logger.Fatal(e.Start(":80"))
}
