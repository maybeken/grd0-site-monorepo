package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"grd0.net/api/data"
)

func (h *Handler) GetMusic(c echo.Context) error {
	music, err := data.ReadMusic()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, music)
}
