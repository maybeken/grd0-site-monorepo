package schema

type Author struct {
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

type Blog struct {
	Uri       string `json:"uri"`
	Author    Author `json:"author"`
	Title     string `json:"title"`
	SubTitle  string `json:"subtitle"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
