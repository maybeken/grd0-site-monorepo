package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"golang.org/x/time/rate"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"grd0.net/api/aws_lambda"
	"grd0.net/api/blog"
	"grd0.net/api/gallery"
	"grd0.net/api/traveler_map"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "https://grd0.net"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(middleware.RequestID())

	skipper := func(c echo.Context) bool {
		// Skip health check endpoint
		return c.Request().URL.Path == "/health"
	}
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		Skipper:   skipper,
		BeforeNextFunc: func(c echo.Context) {
			c.Set("customValueFromContext", 42)
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			fmt.Printf("REQUEST: uri: %v, status: %v\n", v.URI, v.Status)
			return nil
		},
	}))

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "healthy")
	})

	e.GET("/blog", blog.GetBlog)
	e.GET("/blog/:uri", blog.GetBlog)

	e.GET("/gallery/category", gallery.GetGalleryCategory)
	e.GET("/gallery/category/:path", gallery.GetGalleryDetail)
	e.GET("/gallery/:category", gallery.GetAsset)

	e.GET("/travel/map", traveler_map.GetMapLocation)

	isLambda := os.Getenv("LAMBDA")

	if isLambda == "TRUE" {
		lambdaAdapter := &aws_lambda.LambdaAdapter{Echo: e}
		lambda.Start(lambdaAdapter.Handler)
	} else {
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

		e.Logger.Fatal(e.Start(":80"))
	}
}
