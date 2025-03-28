package schema

type Music struct {
	BaseColumns

	Title       string `json:"title"`
	Artist      string `json:"artist,omitempty"`
	V           string `json:"v"`
	Description string `json:"description"`
}
