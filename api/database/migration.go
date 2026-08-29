package database

import (
	"gorm.io/gorm"

	"grd0.net/api/schema"
)

var tables = []interface{}{
	&schema.Author{},
	&schema.Blog{},
	&schema.BrewEquipment{},
	&schema.CoffeeBean{},
	&schema.GalleryDetail{},
	&schema.GalleryCollectionDetail{},
	&schema.MapLocation{},
	&schema.FeedEntry{},
	&schema.FeedBlock{},
	&schema.Music{},
	&schema.TastingNote{},
}

func AutoMigrate(db *gorm.DB) {
	for _, table := range tables {
		db.AutoMigrate(table)
	}
}
