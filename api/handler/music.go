package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"grd0.net/api/schema"
)

func (h *Handler) GetMusic(c echo.Context) error {
	db := h.DB

	var music []schema.Music
	db.Order("sorting DESC").Order("created_at ASC").Find(&music)

	return c.JSON(http.StatusOK, music)
}

func (h *Handler) UpsertMusic(c echo.Context) error {
	db := h.DB

	// var music schema.Music
	var music schema.Music
	if err := c.Bind(&music); err != nil {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "Unable to parse input.")
	}

	if err := db.Create(&music).Error; errors.Is(err, gorm.ErrDuplicatedKey) {
		if err := db.Unscoped().Model(schema.Music{}).Where("v = ?", music.V).Updates(&music).Update("deleted_at", nil).Error; err != nil {
			return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
		}
	} else if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, music)
}

func (h *Handler) DeleteMusic(c echo.Context) error {
	vid := c.Param("v")

	db := h.DB

	res := db.Where("v = ?", vid).Delete(&schema.Music{})
	row_count := res.RowsAffected
	err := res.Error

	if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	} else if row_count <= 0 {
		return h.ErrorResponseConstructor(c, http.StatusNotFound, "")
	}

	return c.JSON(http.StatusAccepted, nil)
}
