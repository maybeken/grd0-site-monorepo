package schema

import "github.com/google/uuid"

type Music struct {
	BaseColumns

	Title       string `json:"title"`
	Artist      string `json:"artist,omitempty"`
	V           string `gorm:"uniqueIndex" json:"v"`
	Description string `json:"description"`
	Sorting     uint   `gorm:"default:0" json:"sorting"`
}

type Playlist struct {
	BaseColumns

	Title       string  `gorm:"notnull" json:"title"`
	Description string  `json:"description,omitempty"`
	Sorting     uint    `gorm:"default:0" json:"sorting"`
	Songs       []Music `gorm:"many2many:playlist_songs" json:"songs"`
}

type PlaylistSong struct {
	PlaylistID uuid.UUID `gorm:"primaryKey" json:"playlist_id"`
	MusicID    uuid.UUID `gorm:"primaryKey" json:"music_id"`
}
