package traveler_map

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"grd0.net/api/data"
)

func GetMapLocation(c echo.Context) error {
	map_loc, err := data.ReadMapLocation()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map_loc)
}
