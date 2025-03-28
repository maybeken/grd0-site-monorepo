package schema

import (
	"github.com/uptrace/bun"
)

type GeoLocation struct {
}

type MapLocation struct {
	bun.BaseModel
	BaseColumns

	Title     string    `json:"title"`
	Subtitle  string    `json:"subtitle,omitempty"`
	Icon      string    `json:"icon"`
	Color     string    `json:"color"`
	Position  []float64 `json:"pos"`
	DisplayAt float64   `json:"display_at,omitempty"`
	HideAt    float64   `json:"hide_at,omitempty"`
	TextColor string    `json:"text_color,omitempty"`
}
