package schema

type MapLocation struct {
	BaseColumns

	Title     string  `json:"title"`
	Subtitle  string  `json:"subtitle,omitempty"`
	Icon      string  `json:"icon"`
	Color     string  `json:"color"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	DisplayAt float64 `json:"display_at,omitempty"`
	HideAt    float64 `json:"hide_at,omitempty"`
	TextColor string  `json:"text_color,omitempty"`
}
