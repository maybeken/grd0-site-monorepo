package schema

import (
	"github.com/uptrace/bun"
)

type GalleryDetail struct {
	bun.BaseModel
	BaseColumns

	Path         string `json:"path"`
	Filename     string `json:"filename"`
	TzAdjustment int8   `json:"tz_adjustment,omitempty"`
	Description  string `json:"description,omitempty"`
}

type GalleryCollection map[string]GalleryCollectionDetail

type GalleryCollectionDetail struct {
	bun.BaseModel
	BaseColumns

	Title string `json:"title,omitempty"`
	Cover string `json:"cover,omitempty"`
}
