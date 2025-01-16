package schema

type GalleryDetail struct {
	Path         string `json:"path"`
	Filename     string `json:"filename"`
	TzAdjustment int8   `json:"tz_adjustment,omitempty"`
	Description  string `json:"description,omitempty"`
}

type GalleryCategory map[string]GalleryCategoryDetail

type GalleryCategoryDetail struct {
	Title string `json:"title"`
	Cover string `json:"cover"`
}
