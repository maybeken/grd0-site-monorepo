package schema

type Music struct {
	BaseColumns

	Title       string `json:"title"`
	Artist      string `json:"artist,omitempty"`
	V           string `gorm:"uniqueIndex" json:"v"`
	Description string `json:"description"`
	Sorting     uint   `gorm:"default:0" json:"-"`
}
