package schema

type GalleryDetail struct {
	BaseColumns

	Path         string `json:"path"`
	Filename     string `json:"filename"`
	TzAdjustment int8   `json:"tz_adjustment,omitempty"`
	Description  string `json:"description,omitempty"`
}

type GalleryCollection map[string]GalleryCollectionDetail

type GalleryCollectionDetail struct {
	BaseColumns

	Path  string `gorm:"uniqueIndex" json:"path,omitempty"`
	Title string `json:"title,omitempty"`
	Cover string `json:"cover,omitempty"`
}
