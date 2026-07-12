package schema

import (
	"time"

	"github.com/google/uuid"
)

type CoffeeBean struct {
	BaseColumns

	Name        string     `gorm:"notnull" json:"name"`
	Origin      *string    `json:"origin"`
	Roaster     *string    `json:"roaster"`
	RoastDate   *time.Time `json:"roast_date"`
	Variety     *string    `json:"variety"`
	Process     *string    `json:"process"`
	Altitude    *string    `json:"altitude"`
	Description *string    `json:"description"`
}

type BrewEquipment struct {
	BaseColumns

	Name        string  `gorm:"notnull" json:"name"`
	Type        *string `json:"type"`
	Description *string `json:"description"`
}

type TastingNote struct {
	BaseColumns

	BeanID      uuid.UUID     `json:"-"`
	Bean        CoffeeBean    `json:"bean"`
	EquipmentID uuid.UUID     `json:"-"`
	Equipment   BrewEquipment `json:"equipment"`

	GrinderID    uuid.UUID     `json:"-"`
	Grinder      BrewEquipment `json:"grinder"`
	GrindSize    *string       `json:"grind_size"`
	GrindSetting *float64      `json:"grind_setting"`

	CoffeeDose       *int     `json:"coffee_dose"`
	WaterIn          *int     `json:"water_in"`
	CoffeeOut        *int     `json:"coffee_out"`
	Ratio            *float64 `json:"ratio"`
	BrewTime         *int     `json:"brew_time"`
	WaterTemperature *int     `json:"water_temperature"`

	TasteFruity    *int `json:"taste_fruity"`
	TasteSour      *int `json:"taste_sour"`
	TasteFermented *int `json:"taste_fermented"`
	TasteNutty     *int `json:"taste_nutty"`
	TasteSpice     *int `json:"taste_spice"`
	TasteSweetness *int `json:"taste_sweetness"`
	TasteFloral    *int `json:"taste_floral"`
	TasteGreen     *int `json:"taste_green"`
	TasteTobacco   *int `json:"taste_tobacco"`
	TasteBitter    *int `json:"taste_bitter"`

	OverallNotes *string   `json:"overall_notes"`
	Rating       *int      `json:"rating"`
	TastedAt     time.Time `json:"tasted_at"`
	Pinned       bool      `gorm:"default:false" json:"pinned"`
}
