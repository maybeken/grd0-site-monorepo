package music

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"grd0.net/api/data"
)

func GetMusic(c echo.Context) error {
	music, err := data.ReadMusic()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, music)
}
