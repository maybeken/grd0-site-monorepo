package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"grd0.net/api/schema"
)

func (h *Handler) GetMusic(c echo.Context) error {
	db := h.DB

	var music []schema.Music
	db.Order("sorting").Order("created_at ASC").Find(&music)

	return c.JSON(http.StatusOK, music)
}

func (h *Handler) AddMusic(c echo.Context) error {
	db := h.DB

	var music schema.Music
	if err := c.Bind(&music); err != nil {
		return ErrorResponseConstructor(c, http.StatusBadRequest, "Unable to parse input.")
	}

	db.Create(&music)

	return c.JSON(http.StatusOK, music)
}
