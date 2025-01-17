package main

import (
	"github.com/labstack/echo/v4"

	"github.com/labstack/echo/v4/middleware"

	"grd0.net/api/blog"
	"grd0.net/api/gallery"
	"grd0.net/api/traveler_map"

	"fmt"
	"net/http"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "https://grd0.net"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(middleware.Decompress())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
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

	e.Logger.Fatal(e.Start(":80"))
}
