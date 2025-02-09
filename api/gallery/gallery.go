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

	return c.JSON(http.StatusBadRequest, "")
}

func GetGalleryCollection(c echo.Context) error {
	collection, err := data.ReadGalleryCollection()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, collection)
}

func GetAsset(c echo.Context) error {
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
