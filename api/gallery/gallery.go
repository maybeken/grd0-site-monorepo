package gallery

import (
	"net/http"

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

		return c.JSON(http.StatusOK, gallery_detail)
	}

	return c.JSON(http.StatusOK, gallery_detail)
}

func GetGalleryCategory(c echo.Context) error {
	blog_posts, err := data.ReadGalleryCategory()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, blog_posts)
}
