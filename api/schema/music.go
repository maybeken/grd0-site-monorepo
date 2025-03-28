package schema

import (
	"github.com/uptrace/bun"
)

type Music struct {
	bun.BaseModel
	BaseColumns

	Title       string `json:"title"`
	Artist      string `json:"artist,omitempty"`
	V           string `json:"v"`
	Description string `json:"description"`
}
