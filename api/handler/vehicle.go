package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"grd0.net/api/schema"
	"grd0.net/api/storage"
	"grd0.net/api/utils"
)

func (h *Handler) GetVehicles(c echo.Context) error {
	db := h.DB

	var vehicles []schema.Vehicle
	if err := db.
		Order("CASE WHEN sort_order IS NULL THEN 1 ELSE 0 END").
		Order("sort_order ASC").
		Order("driven_at ASC").
		Find(&vehicles).Error; err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, vehicles)
}

func (h *Handler) GetVehicle(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusNotFound, "")
	}

	var vehicle schema.Vehicle
	if err := h.DB.First(&vehicle, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return h.ErrorResponseConstructor(c, http.StatusNotFound, "")
		}
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, vehicle)
}

func (h *Handler) UpsertVehicle(c echo.Context) error {
	db := h.DB

	var vehicle schema.Vehicle
	if err := c.Bind(&vehicle); err != nil {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "Unable to parse input.")
	}

	vehicle.Make = strings.TrimSpace(vehicle.Make)
	vehicle.Model = strings.TrimSpace(vehicle.Model)

	if vehicle.Make == "" {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "Vehicle make is required.")
	}

	if vehicle.Model == "" {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "Vehicle model is required.")
	}

	if vehicle.Year == 0 {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "Vehicle year is required.")
	}

	if vehicle.ID != uuid.Nil {
		res := db.Model(&schema.Vehicle{}).Where("id = ?", vehicle.ID).Select("*").Updates(&vehicle)
		if res.Error != nil {
			return h.ErrorResponseConstructor(c, http.StatusInternalServerError, res.Error.Error())
		}
		if res.RowsAffected == 0 {
			return h.ErrorResponseConstructor(c, http.StatusNotFound, "")
		}
	} else {
		if err := db.Create(&vehicle).Error; err != nil {
			return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, vehicle)
}

func (h *Handler) DeleteVehicle(c echo.Context) error {
	id := c.Param("id")

	res := h.DB.Where("id = ?", id).Delete(&schema.Vehicle{})
	if res.Error != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, res.Error.Error())
	} else if res.RowsAffected <= 0 {
		return h.ErrorResponseConstructor(c, http.StatusNotFound, "")
	}

	return c.JSON(http.StatusAccepted, nil)
}

func (h *Handler) GenerateVehiclePresignedUrl(c echo.Context) error {
	key := c.Param("key")
	if key == "" || strings.Contains(key, "/") || strings.Contains(key, "\\") || strings.Contains(key, "..") {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "Invalid attachment key.")
	}

	if h.S3 == nil {
		return h.ErrorResponseConstructor(c, http.StatusNotImplemented, "Attachment uploads are unavailable in local mode.")
	}

	s3 := h.S3
	year := time.Now().UTC().Year()
	month := time.Now().UTC().Format("01")
	object_path := fmt.Sprintf("vehicles/%d/%s/%s", year, month, key)

	presigned_url, err := s3.GeneratePresignedUrl(
		storage.S3PUT,
		utils.GetEnvWithFallback("S3_BUCKET", "assets-grd0-net"),
		object_path,
		60,
	)

	if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusAccepted, PresignedUrlResponseBody{
		Url:       presigned_url.String(),
		ObjectKey: object_path,
	})
}
