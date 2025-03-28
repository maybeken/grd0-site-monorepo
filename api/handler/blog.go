package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/thoas/go-funk"

	"grd0.net/api/data"

	"grd0.net/api/schema"
)

const CONTENT_MAX = 1024

func (h *Handler) GetBlog(c echo.Context) error {
	uri := c.Param("uri")

	blog_posts, err := data.ReadBlogs()

	if err != nil {
		return err
	}

	if uri != "" {
		blog_posts = funk.Filter(blog_posts, func(post schema.Blog) bool {
			return post.Uri == uri
		}).([]schema.Blog)

		if len(blog_posts) > 0 {
			return c.JSON(http.StatusOK, blog_posts[0])
		}

		return c.JSON(http.StatusNotFound, "")
	}

	blog_posts = funk.Map(blog_posts, func(post schema.Blog) schema.Blog {
		if len(post.Content) > CONTENT_MAX {
			post.Content = post.Content[0:CONTENT_MAX]
			post.Content = post.Content + "..."
		}

		return post
	}).([]schema.Blog)

	blog_posts = funk.Filter(blog_posts, func(post schema.Blog) bool {
		publish_time, err := time.Parse(time.RFC3339, post.PublishedAt)

		if err != nil {
			return false
		}

		return publish_time.Before(time.Now())
	}).([]schema.Blog)

	return c.JSON(http.StatusOK, blog_posts)
}
