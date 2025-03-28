package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/oiime/logrusbun"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"

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

func OpenDatabase(path string) *bun.DB {
	log := utils.InitiateLogger()

	if err := ensureDirectoryExists(path); err != nil {
		panic(err)
	}

	sqldb, err := sql.Open(sqliteshim.ShimName, fmt.Sprintf("file:%s?cache=shared&mode=rwc", path))

	if err != nil {
		panic(err)
	}

	log.Infof("Successfully open database file at %s", path)

	db := bun.NewDB(sqldb, sqlitedialect.New())

	db.AddQueryHook(logrusbun.NewQueryHook(logrusbun.QueryHookOptions{QueryLevel: logrus.DebugLevel, ErrorLevel: logrus.ErrorLevel, Logger: log}))

	return db
}
