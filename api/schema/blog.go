package schema

import (
	"github.com/uptrace/bun"
)

type Author struct {
	bun.BaseModel
	BaseColumns

	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

type Blog struct {
	bun.BaseModel
	BaseColumns

	Uri         string  `bun:",notnull" json:"uri"`
	AuthorID    string  `bun:",notnull" json:"author_id"`
	Author      *Author `bun:"rel:belongs-to" json:"author"`
	Title       string  `json:"title"`
	SubTitle    string  `json:"subtitle"`
	Content     string  `json:"content"`
	PublishedAt string  `json:"published_at"`
}
