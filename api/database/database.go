package database

import (
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"grd0.net/api/utils"
)

func ensureDirectoryExists(path string) error {
	// Get the directory path from the full file path
	dirPath := filepath.Dir(path)

	// Check if the directory exists
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		// Directory doesn't exist, create it
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err
		}
	} else if err != nil {
		// Some other error occurred
		return err
	}

	// Directory exists or was created successfully
	return nil
}

func OpenDatabase(path string) *gorm.DB {
	log := utils.InitiateLogger()

	if err := ensureDirectoryExists(path); err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	log.Infof("Successfully open database file at %s", path)

	return db
}
