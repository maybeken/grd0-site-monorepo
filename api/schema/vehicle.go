package schema

import (
	"time"
)

type Vehicle struct {
	BaseColumns

	Make  string `gorm:"notnull" json:"make"`
	Model string `gorm:"notnull" json:"model"`
	Year  int    `gorm:"notnull" json:"year"`

	BodyType     *string `json:"body_type"`
	Engine       *string `json:"engine"`
	Power        *string `json:"power"`
	Transmission *string `json:"transmission"`
	Drivetrain   *string `json:"drivetrain"`
	Fuel         *string `json:"fuel"`

	CountryCode *string `json:"country_code"`
	Country     *string `json:"country"`
	City        *string `json:"city"`

	DrivenAt *time.Time `json:"driven_at"`
	Notes    *string    `json:"notes"`
	Rating   *int       `json:"rating"`

	SortOrder *int `json:"sort_order"`

	CutoutKey     *string `json:"cutout_key"`
	SilhouetteKey *string `json:"silhouette_key"`
}
