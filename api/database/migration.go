package database

import (
	"gorm.io/gorm"

	"grd0.net/api/schema"
)

var tables = []interface{}{
	&schema.Author{},
	&schema.Blog{},
	&schema.GalleryDetail{},
	&schema.GalleryCollectionDetail{},
	&schema.MapLocation{},
	&schema.Music{},
}

func AutoMigrate(db *gorm.DB) {
	for _, table := range tables {
		db.AutoMigrate(table)
	}
}
