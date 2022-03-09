package db

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Repo struct {
	*gorm.DB
}

var db *Repo

// Setup sql database by passing the DSN.
// Should be called before calling GetDB.
func Setup(dsn string) {
	gdb, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("error connecting database: %s", dsn))
	}
	db = &Repo{DB: gdb}
}

// Return the initialized global Repo object.
// Should call Setup before calling GetDB.
func GetDB() *Repo {
	return db
}
