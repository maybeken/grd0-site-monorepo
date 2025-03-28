package database

import (
	"context"

	"github.com/uptrace/bun"

	"grd0.net/api/schema"
)

var tables = []interface{}{
	(*schema.Author)(nil),
	(*schema.Blog)(nil),
	(*schema.GalleryDetail)(nil),
	(*schema.GalleryCollectionDetail)(nil),
	(*schema.MapLocation)(nil),
	(*schema.Music)(nil),
}

func CreateTablesIfNotExists(db *bun.DB) {
	for _, table := range tables {
		db.NewCreateTable().
			Model(table).
			IfNotExists().
			Exec(context.Background())
	}
}

func DropTablesIfExists(db *bun.DB) {
	for _, table := range tables {
		db.NewDropTable().
			Model(table).
			IfExists().
			Exec(context.Background())
	}
}
