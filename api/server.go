package main

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/time/rate"
	"grd0.net/api/database"
	"grd0.net/api/handler"
	"grd0.net/api/storage"
	"grd0.net/api/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
)

func main() {
	log := utils.InitiateLogger()

	local_environment := utils.IsLocalEnvironment()

	if local_environment {
		log.Warn("APP_ENV=LOCAL: authentication bypass is enabled; never use this in production")
	} else if err := handler.InitOIDC(context.Background()); err != nil {
		log.Fatalf("Failed to initialize OIDC: %v", err)
	}

	db_path := utils.GetEnvWithFallback("DB_PATH", "data/api.db")

	db := database.OpenDatabase(db_path)
	database.AutoMigrate(db)

	var s3client *storage.S3
	if !local_environment {
		client, err := storage.InitiateClient(utils.GetEnv("S3_ENDPOINT"), utils.GetEnv("S3_ACCESS_KEY_ID"), utils.GetEnv("S3_SECRET_ACCESS_KEY"))
		if err != nil {
			panic(err)
		}
		s3client = client
	}

	blob_backend := utils.GetEnvWithFallback("BLOB_BACKEND", "local")
	blob_local_path := utils.GetEnvWithFallback("BLOB_LOCAL_PATH", "data/scenes")
	blobs, err := storage.NewBlobStore(
		blob_backend,
		blob_local_path,
		s3client,
		utils.GetEnvWithFallback("S3_BUCKET", "assets-grd0-net"),
	)
	if err != nil {
		panic(err)
	}

	scene_max_bytes, err := strconv.ParseInt(
		utils.GetEnvWithFallback("ARTICLE_SCENE_MAX_BYTES", strconv.FormatInt(handler.DEFAULT_ARTICLE_SCENE_MAX_BYTE, 10)),
		10,
		64,
	)
	if err != nil {
		log.Fatalf("Invalid ARTICLE_SCENE_MAX_BYTES: %v", err)
	}

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "https://grd0.net"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
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

	h := &handler.Handler{
		DB:                   db,
		Logger:               log,
		S3:                   s3client,
		Blobs:                blobs,
		ArticleSceneMaxBytes: scene_max_bytes,
	}
	var auth_guard echo.MiddlewareFunc
	if local_environment {
		auth_guard = handler.LocalAuthGuard()
	} else {
		auth_guard = echojwt.WithConfig(echojwt.Config{
			NewClaimsFunc: func(c echo.Context) jwt.Claims {
				return new(handler.JwtCustomClaims)
			},
			SigningKey: []byte(utils.GetEnv("AUTH_JWT_SECRET")),
		})
	}

	handler.RegisterRouter(e, h, auth_guard)

	e.Logger.Fatal(e.Start(utils.GetEnvWithFallback("PORT", ":80")))
}
