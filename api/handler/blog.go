package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"

	"grd0.net/api/schema"
)

const CONTENT_MAX = 1024

func (h *Handler) GetBlog(c echo.Context) error {
	uri := c.Param("uri")

	db := h.DB

	if uri != "" {
		post := schema.Blog{
			Uri: uri,
		}
		error := db.Where(&post).Preload("Author").Take(&post).Error

		if errors.Is(error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "")
		}

		return c.JSON(http.StatusOK, post)
	}

	var blog_posts []schema.Blog
	db.Where("published_at <= ?", time.Now()).Preload("Author").Order("published_at").Find(&blog_posts)

	blog_posts = funk.Map(blog_posts, func(post schema.Blog) schema.Blog {
		if len(post.Content) > CONTENT_MAX {
			post.Content = post.Content[0:CONTENT_MAX]
			post.Content = post.Content + "..."
		}

		return post
	}).([]schema.Blog)

	return c.JSON(http.StatusOK, blog_posts)
}

func (h *Handler) UpsertBlog(c echo.Context) error {
	uri := c.Param("uri")

	db := h.DB

	var post schema.Blog
	if err := c.Bind(&post); err != nil {
		return ErrorResponseConstructor(c, http.StatusBadRequest, "Unable to parse input.")
	}
	post.Uri = uri

	existing_post := schema.Blog{
		Uri: uri,
	}
	error := db.Select("id").Where(&existing_post).Take(&existing_post).Error

	if errors.Is(error, gorm.ErrRecordNotFound) {
		db.Create(&post)
	} else {
		id := existing_post.ID
		db.Model(schema.Blog{}).Where("id = ?", id).Updates(&post)
	}

	return c.JSON(http.StatusOK, &post)
}
