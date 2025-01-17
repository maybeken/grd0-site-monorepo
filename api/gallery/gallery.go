package gallery

import (
	"net/http"
	"sort"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/thoas/go-funk"

	"grd0.net/api/data"

	"grd0.net/api/schema"
)

func GetGalleryDetail(c echo.Context) error {
	path := c.Param("path")

	gallery_detail, err := data.ReadGalleryDetail()

	if err != nil {
		return err
	}

	if path != "" {
		gallery_detail = funk.Filter(gallery_detail, func(post schema.GalleryDetail) bool {
			return post.Path == path
		}).([]schema.GalleryDetail)

		if len(gallery_detail) > 0 {
			return c.JSON(http.StatusOK, gallery_detail)
		}

		return c.JSON(http.StatusNotFound, "")
	}

	return c.JSON(http.StatusOK, gallery_detail)
}

func GetGalleryCategory(c echo.Context) error {
	category, err := data.ReadGalleryCategory()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, category)
}

func GetAsset(c echo.Context) error {
	category := c.Param("category")

	asset_list, err := data.ReadAsset()

	if err != nil {
		return err
	}

	if category == "all" {
		filtered := make(schema.AssetFileList)

		for key, values := range asset_list {
			// Check if the key starts with the specified prefix
			if strings.HasPrefix(key, "/gallery") {
				// Add the key-value pair to the filtered map
				filtered[key] = values
			}
		}

		combined := funk.FlatMap(filtered, func(key string, values []schema.Asset) []schema.Asset {
			with_category := funk.Map(values, func(item schema.Asset) schema.Asset {
				item.Category = key
				return item
			}).([]schema.Asset)

			return with_category
		}).([]schema.Asset)

		sort.Slice(combined, func(i, j int) bool {
			return combined[i].Exif.Datetime > combined[j].Exif.Datetime
		})

		return c.JSON(http.StatusOK, combined)
	} else if category != "" {
		if values, exists := asset_list[category]; exists {
			with_category := funk.Map(values, func(item schema.Asset) schema.Asset {
				item.Category = category
				return item
			}).([]schema.Asset)

			sort.Slice(with_category, func(i, j int) bool {
				return with_category[i].Exif.Datetime < with_category[j].Exif.Datetime
			})

			return c.JSON(http.StatusOK, with_category)
		}

		return c.JSON(http.StatusNotFound, "")
	}

	return c.JSON(http.StatusBadRequest, "")
}
