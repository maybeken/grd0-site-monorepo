package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type (
	Handler struct {
		DB     *gorm.DB
		Logger *logrus.Logger
	}
)

func RegisterRouter(e *echo.Echo, h *Handler, auth_guard echo.MiddlewareFunc) {
	e.GET("/health", h.GetHealthcheck)

	r := e.Group("")
	r.Use(auth_guard)

	e.POST("/login", h.Login)

	e.GET("/blog", h.GetBlog)
	e.GET("/blog/:uri", h.GetBlog)

	e.GET("/gallery/collection", h.GetGalleryCollection)
	e.GET("/gallery/collection/:path", h.GetGalleryDetail)
	e.GET("/gallery/:collection", h.GetAsset)

	e.GET("/travel/map", h.GetMapLocation)

	e.GET("/music", h.GetMusic)
}

type ErrorResponseBody struct {
	TraceID string `json:"trace_id"`
	Error   string `json:"error"`
}

type HealthcheckResponseBody struct {
	Status        string `json:"status"`
	SQLiteVersion string `json:"sqlite_version"`
}

func (h *Handler) GetHealthcheck(c echo.Context) error {
	db := h.DB

	var version string
	if err := db.Raw("SELECT sqlite_version()").Scan(&version).Error; err != nil {
		return ErrorResponseConstructor(c, http.StatusInternalServerError, fmt.Sprintf("%s", err))
	}

	return c.JSON(http.StatusOK, HealthcheckResponseBody{
		Status:        "Healthy",
		SQLiteVersion: version,
	})
}

func ErrorResponseConstructor(c echo.Context, status int, message string) error {
	return c.JSON(status, ErrorResponseBody{
		TraceID: c.Response().Header().Get(echo.HeaderXRequestID),
		Error:   fmt.Sprintf("%s", message),
	})
}
