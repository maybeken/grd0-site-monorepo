package handler

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"grd0.net/api/schema"
)

var markdownImageRegex = regexp.MustCompile(`!\[.*?\]\((.*?)\)`)

const FEED_DEFAULT_LIMIT = 20
const FEED_MAX_LIMIT = 100
const FEED_SUMMARY_MAX = 1024

type FeedBlockInput struct {
	Type       string              `json:"type"`
	Order      int                 `json:"order"`
	CreateData *FeedBlockCreateData `json:"create_data,omitempty"`
	Ref        *string             `json:"ref,omitempty"`
}

type FeedBlockCreateData struct {
	Title    *string `json:"title,omitempty"`
	Subtitle *string `json:"subtitle,omitempty"`
	Content  *string `json:"content,omitempty"`
	V        *string `json:"v,omitempty"`
	Artist   *string `json:"artist,omitempty"`
}

type FeedUpsertInput struct {
	Title       string             `json:"title"`
	Subtitle    string             `json:"subtitle"`
	PublishedAt *time.Time         `json:"published_at"`
	Blocks      []FeedBlockInput   `json:"blocks"`
}

func (h *Handler) GetFeed(c echo.Context) error {
	db := h.DB

	limit := FEED_DEFAULT_LIMIT
	if l, err := strconv.Atoi(c.QueryParam("limit")); err == nil && l > 0 {
		limit = l
	}
	if limit > FEED_MAX_LIMIT {
		limit = FEED_MAX_LIMIT
	}

	query := db.Where("published_at IS NOT NULL AND published_at <= ?", time.Now())

	if before := c.QueryParam("before"); before != "" {
		if t, err := time.Parse(time.RFC3339, before); err == nil {
			query = query.Where("created_at < ?", t)
		}
	}

	var entries []schema.FeedEntry
	query.Preload("Author").
		Preload("Blocks", func(db *gorm.DB) *gorm.DB {
			return db.Order(clause.OrderByColumn{Column: clause.Column{Table: "feed_blocks", Name: "order"}, Desc: false})
		}).
		Preload("Blocks.Blog").
		Preload("Blocks.Blog.Author").
		Preload("Blocks.Music").
		Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true}).
		Limit(limit).
		Find(&entries)

	for i := range entries {
		if len(entries[i].Blocks) > 0 {
			firstBlock := entries[i].Blocks[0]
			switch firstBlock.Type {
			case "blog":
				if firstBlock.Blog != nil {
					content := firstBlock.Blog.Content
					if len(content) > FEED_SUMMARY_MAX {
						entries[i].Summary = content[:FEED_SUMMARY_MAX] + "..."
					} else {
						entries[i].Summary = content
					}
				}
			case "music":
				if firstBlock.Music != nil {
					if firstBlock.Music.Artist != "" {
						entries[i].Summary = firstBlock.Music.Title + " - " + firstBlock.Music.Artist
					} else {
						entries[i].Summary = firstBlock.Music.Title
					}
				}
			}
		}

		for _, block := range entries[i].Blocks {
			if block.Type == "blog" && block.Blog != nil {
				if matches := markdownImageRegex.FindStringSubmatch(block.Blog.Content); len(matches) > 1 {
					entries[i].CoverImage = matches[1]
					break
				}
			}
		}

		entries[i].Blocks = nil
	}

	return c.JSON(http.StatusOK, entries)
}

func (h *Handler) GetFeedByAdmin(c echo.Context) error {
	db := h.DB

	var entries []schema.FeedEntry
	db.Preload("Author").
		Preload("Blocks", func(db *gorm.DB) *gorm.DB {
			return db.Order(clause.OrderByColumn{Column: clause.Column{Table: "feed_blocks", Name: "order"}, Desc: false})
		}).
		Preload("Blocks.Blog").
		Preload("Blocks.Blog.Author").
		Preload("Blocks.Music").
		Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true}).
		Find(&entries)

	for i := range entries {
		if len(entries[i].Blocks) > 0 {
			firstBlock := entries[i].Blocks[0]
			switch firstBlock.Type {
			case "blog":
				if firstBlock.Blog != nil {
					content := firstBlock.Blog.Content
					if len(content) > FEED_SUMMARY_MAX {
						entries[i].Summary = content[:FEED_SUMMARY_MAX] + "..."
					} else {
						entries[i].Summary = content
					}
				}
			case "music":
				if firstBlock.Music != nil {
					if firstBlock.Music.Artist != "" {
						entries[i].Summary = firstBlock.Music.Title + " - " + firstBlock.Music.Artist
					} else {
						entries[i].Summary = firstBlock.Music.Title
					}
				}
			}
		}

		for _, block := range entries[i].Blocks {
			if block.Type == "blog" && block.Blog != nil {
				if matches := markdownImageRegex.FindStringSubmatch(block.Blog.Content); len(matches) > 1 {
					entries[i].CoverImage = matches[1]
					break
				}
			}
		}

		entries[i].Blocks = nil
	}

	return c.JSON(http.StatusOK, entries)
}

func (h *Handler) GetFeedBySlug(c echo.Context) error {
	slug := c.Param("slug")
	db := h.DB

	var entry schema.FeedEntry
	err := db.Where("slug = ? AND published_at IS NOT NULL AND published_at <= ?", slug, time.Now()).
		Preload("Author").
		Preload("Blocks", func(db *gorm.DB) *gorm.DB {
			return db.Order(clause.OrderByColumn{Column: clause.Column{Table: "feed_blocks", Name: "order"}, Desc: false})
		}).
		Preload("Blocks.Blog").
		Preload("Blocks.Blog.Author").
		Preload("Blocks.Music").
		First(&entry).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, "")
	}
	if err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, entry)
}

func (h *Handler) UpsertFeed(c echo.Context) error {
	slug := c.Param("slug")
	db := h.DB

	var input FeedUpsertInput
	if err := c.Bind(&input); err != nil {
		return h.ErrorResponseConstructor(c, http.StatusBadRequest, "Unable to parse input.")
	}

	var entry schema.FeedEntry
	err := db.Where("slug = ?", slug).First(&entry).Error
	is_new := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !is_new {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	tx := db.Begin()

	author := schema.Author{}
	if err := tx.Where("email = ?", "ken.lam@grd0.net").First(&author).Error; err != nil {
		tx.Rollback()
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, "Author not found.")
	}

	if is_new {
		entry = schema.FeedEntry{
			Slug:        slug,
			Title:       input.Title,
			SubTitle:    input.Subtitle,
			AuthorID:    author.ID,
			PublishedAt: input.PublishedAt,
		}
		if err := tx.Create(&entry).Error; err != nil {
			tx.Rollback()
			return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
		}
	} else {
		entry.Title = input.Title
		entry.SubTitle = input.Subtitle
		entry.AuthorID = author.ID
		entry.PublishedAt = input.PublishedAt
		if err := tx.Save(&entry).Error; err != nil {
			tx.Rollback()
			return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
		}

		if err := tx.Where("feed_entry_id = ?", entry.ID).Delete(&schema.FeedBlock{}).Error; err != nil {
			tx.Rollback()
			return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
		}
	}

	for _, blockInput := range input.Blocks {
		block := schema.FeedBlock{
			FeedEntryID: entry.ID,
			Type:        blockInput.Type,
			Order:       blockInput.Order,
		}

		if blockInput.Ref != nil && *blockInput.Ref != "" {
			ref := *blockInput.Ref
			switch blockInput.Type {
			case "blog":
				var blog schema.Blog
				if err := tx.Where("uri = ?", ref).First(&blog).Error; err != nil {
					tx.Rollback()
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return h.ErrorResponseConstructor(c, http.StatusNotFound, "Referenced blog not found.")
					}
					return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
				}
				block.BlogID = &blog.ID
			case "music":
				var music schema.Music
				if err := tx.Where("v = ?", ref).First(&music).Error; err != nil {
					tx.Rollback()
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return h.ErrorResponseConstructor(c, http.StatusNotFound, "Referenced music not found.")
					}
					return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
				}
				block.MusicID = &music.ID
			}
		} else if blockInput.CreateData != nil {
			cd := blockInput.CreateData
			switch blockInput.Type {
			case "blog":
				blog := schema.Blog{
					AuthorID: author.ID,
					Title:    stringOrDefault(cd.Title, ""),
					Content:  stringOrDefault(cd.Content, ""),
				}
				blog.Uri = slug + "-blog-" + uuid.New().String()[:8]

				if err := tx.Omit("Author").Create(&blog).Error; err != nil {
					tx.Rollback()
					return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
				}
				block.BlogID = &blog.ID
			case "music":
				music := schema.Music{
					Title:  stringOrDefault(cd.Title, ""),
					Artist: stringOrDefault(cd.Artist, ""),
					V:      stringOrDefault(cd.V, ""),
				}

				if music.V == "" {
					tx.Rollback()
					return h.ErrorResponseConstructor(c, http.StatusBadRequest, "YouTube video ID is required for music blocks.")
				}

				if err := tx.Create(&music).Error; err != nil {
					if errors.Is(err, gorm.ErrDuplicatedKey) {
						existing := schema.Music{}
						if err := tx.Where("v = ?", music.V).First(&existing).Error; err != nil {
							tx.Rollback()
							return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
						}
						block.MusicID = &existing.ID
					} else {
						tx.Rollback()
						return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
					}
				} else {
					block.MusicID = &music.ID
				}
			}
		}

		if err := tx.Create(&block).Error; err != nil {
			tx.Rollback()
			return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
		}
	}

	if err := tx.Commit().Error; err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	var result schema.FeedEntry
	db.Where("id = ?", entry.ID).
		Preload("Author").
		Preload("Blocks", func(db *gorm.DB) *gorm.DB {
			return db.Order(clause.OrderByColumn{Column: clause.Column{Table: "feed_blocks", Name: "order"}, Desc: false})
		}).
		Preload("Blocks.Blog").
		Preload("Blocks.Blog.Author").
		Preload("Blocks.Music").
		First(&result)

	return c.JSON(http.StatusOK, result)
}

func (h *Handler) DeleteFeed(c echo.Context) error {
	slug := c.Param("slug")
	db := h.DB

	var entry schema.FeedEntry
	if err := db.Where("slug = ?", slug).First(&entry).Error; err != nil {
		return h.ErrorResponseConstructor(c, http.StatusNotFound, "")
	}

	tx := db.Begin()

	if err := tx.Unscoped().Where("feed_entry_id = ?", entry.ID).Delete(&schema.FeedBlock{}).Error; err != nil {
		tx.Rollback()
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	if err := tx.Unscoped().Delete(&entry).Error; err != nil {
		tx.Rollback()
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return h.ErrorResponseConstructor(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusAccepted, nil)
}

func stringOrDefault(s *string, def string) string {
	if s != nil {
		return *s
	}
	return def
}
