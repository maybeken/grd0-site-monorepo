package handler

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"grd0.net/api/schema"
)

func (h *Handler) GetMapLocation(c echo.Context) error {
	db := h.DB

	var map_locations []schema.MapLocation
	db.Find(&map_locations)

	locations_formatted := make([]schema.MapLocation, 0, len(map_locations))
	for _, loc := range map_locations {
		loc.Position = []float64{
			loc.Latitude,
			loc.Longitude,
		}
		locations_formatted = append(locations_formatted, loc)
	}

	return c.JSON(http.StatusOK, locations_formatted)
}

func (h *Handler) UpsertMapLocation(c echo.Context) error {
	db := h.DB

	var map_loc schema.MapLocation
	if err := c.Bind(&map_loc); err != nil {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "Unable to parse input.")
	}

	if len(map_loc.Position) != 2 {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "Position must be an array of two numbers")
	}
	map_loc.Latitude = map_loc.Position[0]
	map_loc.Longitude = map_loc.Position[1]

	if map_loc.Slug == "" {
		map_loc.Slug = url.PathEscape(map_loc.Title)
	} else {
		map_loc.Slug = url.PathEscape(map_loc.Slug)
	}

	if err := db.Create(&map_loc).Error; errors.Is(err, gorm.ErrDuplicatedKey) {
		if err := db.Unscoped().Model(schema.MapLocation{}).Where("slug = ?", map_loc.Slug).Updates(&map_loc).Update("deleted_at", nil).Error; err != nil {
			return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
		}
	} else if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map_loc)
}

func (h *Handler) DeleteMapLocation(c echo.Context) error {
	slug := url.PathEscape(c.Param("slug"))

	db := h.DB

	res := db.Where("slug = ?", slug).Delete(&schema.MapLocation{})
	row_count := res.RowsAffected
	err := res.Error

	if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	} else if row_count <= 0 {
		return h.ErrorResponseConstructor(c, http.StatusNotFound, "")
	}

	return c.JSON(http.StatusAccepted, nil)
}
