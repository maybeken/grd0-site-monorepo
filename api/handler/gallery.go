package handler

import (
	"errors"
	"net/http"
	"sort"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/thoas/go-funk"

	"grd0.net/api/data"

	"grd0.net/api/schema"
)

func (h *Handler) GetGalleryDetail(c echo.Context) error {
	path := c.Param("path")

	if path == "" {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "")
	}

	db := h.DB

	var gallery_detail []schema.GalleryDetail
	if err := db.Where("path = ?", path).Find(&gallery_detail).Error; err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, gallery_detail)
}

func (h *Handler) UpsertGalleryDetail(c echo.Context) error {
	path := c.Param("path")

	if path == "" {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "")
	}

	db := h.DB

	var gallery_detail schema.GalleryDetail
	if err := c.Bind(&gallery_detail); err != nil {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "Unable to parse input.")
	}
	gallery_detail.Path = path

	if err := db.Create(&gallery_detail).Error; errors.Is(err, gorm.ErrDuplicatedKey) {
		if err := db.Unscoped().Model(schema.GalleryDetail{}).Where("path = ? AND filename = ?", path, gallery_detail.Filename).Updates(&gallery_detail).Update("deleted_at", nil).Error; err != nil {
			return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
		}
	} else if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, &gallery_detail)
}

func (h *Handler) DeleteGalleryDetail(c echo.Context) error {
	path := c.Param("path")

	var gallery_detail schema.GalleryDetail
	if err := c.Bind(&gallery_detail); err != nil {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "Unable to parse input.")
	}

	if path == "" || gallery_detail.Filename == "" {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "")
	}

	db := h.DB

	res := db.Where("path = ? AND filename = ?", path, gallery_detail.Filename).Delete(&schema.GalleryDetail{})
	row_count := res.RowsAffected
	err := res.Error

	if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	} else if row_count <= 0 {
		return h.ErrorResponseConstructor(c, http.StatusNotFound, "")
	}

	return c.JSON(http.StatusAccepted, nil)
}

func (h *Handler) GetGalleryCollection(c echo.Context) error {
	db := h.DB

	var collections []schema.GalleryCollectionDetail

	if err := db.Order("updated_at DESC").Find(&collections).Error; err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	collection_by_key := make(schema.GalleryCollection, len(collections))

	for _, collection := range collections {
		modified_collection := collection
		modified_collection.Path = ""

		collection_by_key[collection.Path] = modified_collection
	}

	return c.JSON(http.StatusOK, collection_by_key)
}

func (h *Handler) UpsertGalleryCollection(c echo.Context) error {
	db := h.DB

	var collection schema.GalleryCollectionDetail
	if err := c.Bind(&collection); err != nil {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "Unable to parse input.")
	}

	if err := db.Create(&collection).Error; errors.Is(err, gorm.ErrDuplicatedKey) {
		if err := db.Unscoped().Model(schema.GalleryCollectionDetail{}).Where("path = ?", collection.Path).Updates(&collection).Update("deleted_at", nil).Error; err != nil {
			return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
		}
	} else if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, &collection)
}

func (h *Handler) DeleteGalleryCollection(c echo.Context) error {
	path := c.Param("path")

	db := h.DB

	res := db.Where("path = ?", path).Delete(&schema.GalleryCollectionDetail{})
	row_count := res.RowsAffected
	err := res.Error

	if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	} else if row_count <= 0 {
		return h.ErrorResponseConstructor(c, http.StatusNotFound, "")
	}

	return c.JSON(http.StatusAccepted, nil)
}

func (h *Handler) GetAsset(c echo.Context) error {
	collection := c.Param("collection")

	asset_list, err := data.ReadAsset()

	if err != nil {
		return err
	}

	if collection == "all" {
		filtered := make(schema.AssetFileList)
		collections, err := data.ReadGalleryCollection()

		if err != nil {
			return err
		}

		for key, value := range asset_list {
			// Check if the key starts with the specified prefix
			if strings.HasPrefix(key, "/gallery") {

				// Check if the file is in a publicly listed collection
				for collection := range collections {
					if strings.HasPrefix(key, "/gallery/"+collection) {
						// Add the key-value pair to the filtered map
						filtered[key] = value
					}
				}
			}
		}

		combined := funk.FlatMap(filtered, func(key string, values []schema.Asset) []schema.Asset {
			with_collection := funk.Map(values, func(item schema.Asset) schema.Asset {
				item.Collection = key
				return item
			}).([]schema.Asset)

			return with_collection
		}).([]schema.Asset)

		sort.Slice(combined, func(i, j int) bool {
			return combined[i].Exif.Datetime > combined[j].Exif.Datetime
		})

		return c.JSON(http.StatusOK, combined)
	} else if collection != "" {
		formatted_collection := "/gallery/" + collection

		if values, exists := asset_list[formatted_collection]; exists {
			with_collection := funk.Map(values, func(item schema.Asset) schema.Asset {
				item.Collection = formatted_collection
				return item
			}).([]schema.Asset)

			sort.Slice(with_collection, func(i, j int) bool {
				return with_collection[i].Exif.Datetime < with_collection[j].Exif.Datetime
			})

			return c.JSON(http.StatusOK, with_collection)
		}

		return c.JSON(http.StatusNotFound, "")
	}

	return c.JSON(http.StatusBadRequest, "")
}
