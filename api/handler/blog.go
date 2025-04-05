package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"grd0.net/api/schema"
	"grd0.net/api/storage"
	"grd0.net/api/utils"
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
	db.Where("published_at NOT NULL AND published_at <= ?", time.Now()).Preload("Author").Order("published_at").Find(&blog_posts)

	blog_posts = funk.Map(blog_posts, func(post schema.Blog) schema.Blog {
		if len(post.Content) > CONTENT_MAX {
			post.Content = post.Content[0:CONTENT_MAX]
			post.Content = post.Content + "..."
		}

		return post
	}).([]schema.Blog)

	return c.JSON(http.StatusOK, blog_posts)
}

func (h *Handler) GetBlogByAdmin(c echo.Context) error {
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
	db.Preload("Author").Order("updated_at").Find(&blog_posts)

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
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "Unable to parse input.")
	}
	post.Uri = uri

	var author schema.Author
	if err := db.Where("email = ?", post.Author.Email).First(&author).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, fmt.Sprintf("Author %s not found", post.Author.Email))
	} else if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}
	post.AuthorID = author.ID
	post.Author = author

	if err := db.Omit(clause.Associations).Create(&post).Error; errors.Is(err, gorm.ErrDuplicatedKey) {
		if err := db.Omit(clause.Associations).Model(schema.Blog{}).Where("uri = ?", post.Uri).Updates(&post).Error; err != nil {
			return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
		}
	} else if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, &post)
}

func (h *Handler) DeleteBlog(c echo.Context) error {
	uri := c.Param("uri")

	db := h.DB

	res := db.Model(schema.Blog{}).Where("uri = ?", uri).Update("published_at", nil)
	err := res.Error
	row_count := res.RowsAffected

	if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	} else if row_count <= 0 {
		return h.ErrorResponseConstructor(c, http.StatusNotFound, "")
	}

	return c.JSON(http.StatusOK, uri)
}

type PresignedUrlResponseBody struct {
	Url       string `json:"url"`
	ObjectKey string `json:"key"`
}

func (h *Handler) GeneratePresignedUrl(c echo.Context) error {
	key := c.Param("key")

	s3 := h.S3
	year := time.Now().Year()
	month := time.Now().Month()
	path := fmt.Sprintf("blog/%d/%d/%s", year, month, key)

	url, err := s3.GeneratePresignedUrl(storage.S3PUT, utils.GetEnvWithFallback("S3_BUCKET", "assets-grd0-net"), path, 60)

	if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusAccepted, PresignedUrlResponseBody{
		Url:       url.String(),
		ObjectKey: path,
	})
}
