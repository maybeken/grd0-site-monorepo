package schema

type Author struct {
	BaseColumns

	Email       string `gorm:"uniqueIndex" json:"email"`
	DisplayName string `json:"display_name"`
}

type Blog struct {
	BaseColumns

	Uri         string `gorm:"uniqueIndex" json:"uri"`
	AuthorID    string `gorm:",notnull" json:"-"`
	Author      Author `json:"author"`
	Title       string `json:"title"`
	SubTitle    string `json:"subtitle"`
	Content     string `json:"content"`
	PublishedAt string `json:"published_at"`
}
