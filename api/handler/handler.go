package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"grd0.net/api/storage"
)

type (
	Handler struct {
		DB     *gorm.DB
		Logger *logrus.Logger
		S3     *storage.S3
	}
)

func RegisterRouter(e *echo.Echo, h *Handler, auth_guard echo.MiddlewareFunc) {
	e.Any("/*", func(c echo.Context) error {
		return c.NoContent(http.StatusNotFound)
	})

	e.GET("/health", h.GetHealthcheck)

	restrictV1 := e.Group("")
	restrictV1.Use(auth_guard)

	v1 := e.Group("")
	v2 := e.Group("/v2")

	v1.GET("/auth/login", h.Login)
	v1.GET("/auth/callback", h.LoginCallback)

	v1.GET("/blog", h.GetBlog)
	v1.GET("/blog/:uri", h.GetBlog)
	restrictV1.GET("/blog/all", h.GetBlogByAdmin)
	restrictV1.PUT("/blog/:uri", h.UpsertBlog)
	restrictV1.DELETE("/blog/:uri", h.DeleteBlog)
	restrictV1.PUT("/blog/attachment/:key", h.GeneratePresignedUrl)

	v1.GET("/gallery/details/:path", h.GetGalleryDetail)
	v2.GET("/gallery/details/:path", h.GetGalleryDetailV2)
	restrictV1.PUT("/gallery/details/:path", h.UpsertGalleryDetail)
	restrictV1.DELETE("/gallery/details/:path", h.DeleteGalleryDetail)

	v1.GET("/gallery/collection", h.GetGalleryCollection)
	restrictV1.PUT("/gallery/collection", h.UpsertGalleryCollection)
	restrictV1.DELETE("/gallery/collection/:path", h.DeleteGalleryCollection)

	v1.GET("/gallery/:collection", h.GetAsset)
	v2.GET("/gallery/:collection", h.GetAssetV2)

	v1.GET("/travel/map", h.GetMapLocation)
	restrictV1.PUT("/travel/map", h.UpsertMapLocation)
	restrictV1.DELETE("/travel/map/:slug", h.DeleteMapLocation)

	v1.GET("/music", h.GetMusic)
	restrictV1.PUT("/music", h.UpsertMusic)
	restrictV1.DELETE("/music/:v", h.DeleteMusic)

	v1.GET("/coffee/beans", h.GetCoffeeBeans)
	restrictV1.PUT("/coffee/bean", h.UpsertCoffeeBean)
	restrictV1.DELETE("/coffee/bean/:id", h.DeleteCoffeeBean)

	v1.GET("/coffee/equipment", h.GetCoffeeEquipment)
	restrictV1.PUT("/coffee/equipment", h.UpsertCoffeeEquipment)
	restrictV1.DELETE("/coffee/equipment/:id", h.DeleteCoffeeEquipment)

	v1.GET("/coffee/tastings", h.GetCoffeeTastings)
	restrictV1.PUT("/coffee/tasting", h.UpsertCoffeeTasting)
	restrictV1.DELETE("/coffee/tasting/:id", h.DeleteCoffeeTasting)
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
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, fmt.Sprintf("%s", err))
	}

	return c.JSON(http.StatusOK, HealthcheckResponseBody{
		Status:        "Healthy",
		SQLiteVersion: version,
	})
}

func (h *Handler) ErrorResponseConstructor(c echo.Context, status int, message string) error {
	return c.JSON(status, ErrorResponseBody{
		TraceID: c.Response().Header().Get(echo.HeaderXRequestID),
		Error:   fmt.Sprintf("%s", message),
	})
}
