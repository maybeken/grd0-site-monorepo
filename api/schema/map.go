package schema

type MapLocation struct {
	BaseColumns

	Slug      string    `gorm:"uniqueIndex" json:"slug"`
	Title     string    `json:"title"`
	Subtitle  string    `json:"subtitle,omitempty"`
	Icon      string    `json:"icon"`
	Color     string    `json:"color"`
	Position  []float64 `gorm:"-" json:"pos"`
	Longitude float64   `json:"-"`
	Latitude  float64   `json:"-"`
	DisplayAt float64   `json:"display_at,omitempty"`
	HideAt    float64   `json:"hide_at,omitempty"`
	TextColor string    `json:"text_color,omitempty"`
}
