package database

import (
	"gorm.io/gorm"

	"grd0.net/api/schema"
)

var tables = []interface{}{
	&schema.Author{},
	&schema.Article{},
	&schema.Blog{},
	&schema.BrewEquipment{},
	&schema.CoffeeBean{},
	&schema.GalleryDetail{},
	&schema.GalleryCollectionDetail{},
	&schema.MapLocation{},
	&schema.Music{},
	&schema.Playlist{},
	&schema.TastingNote{},
	&schema.Vehicle{},
}

func AutoMigrate(db *gorm.DB) {
	for _, table := range tables {
		db.AutoMigrate(table)
	}

	db.Model(&schema.MapLocation{}).Where("published IS NULL").UpdateColumn("published", true)
}
