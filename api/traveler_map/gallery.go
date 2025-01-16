package traveler_map

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"grd0.net/api/data"
)

func GetMapLocation(c echo.Context) error {
	blog_posts, err := data.ReadMapLocation()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, blog_posts)
}
