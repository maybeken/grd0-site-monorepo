package schema

type GalleryDetail struct {
	BaseColumns

	Path         string `gorm:"index:,unique,composite:filepath" json:"path"`
	Filename     string `gorm:"index:,unique,composite:filepath" json:"filename"`
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
