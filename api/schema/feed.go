package schema

import (
	"time"

	"github.com/google/uuid"
)

type FeedEntry struct {
	BaseColumns

	Slug        string       `gorm:"uniqueIndex" json:"slug"`
	Title       string       `json:"title"`
	SubTitle    string       `json:"subtitle"`
	AuthorID    uuid.UUID    `gorm:",notnull" json:"-"`
	Author      Author       `json:"author"`
	PublishedAt *time.Time   `gorm:"default:null" json:"published_at"`
	Blocks      []FeedBlock  `json:"blocks"`
	Summary     string       `gorm:"-" json:"summary"`
	CoverImage  string       `gorm:"-" json:"cover_image"`
}

type FeedBlock struct {
	BaseColumns

	FeedEntryID uuid.UUID `gorm:"index" json:"feed_entry_id"`
	FeedEntry   FeedEntry `json:"-"`
	Type        string    `json:"type"`
	Order       int       `json:"order"`
	BlogID      *uuid.UUID `gorm:"index" json:"-"`
	Blog        *Blog      `json:"blog,omitempty"`
	MusicID     *uuid.UUID `gorm:"index" json:"-"`
	Music       *Music     `json:"music,omitempty"`
}
