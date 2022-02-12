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

func Setup(dsn string) {
	gdb, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("error connecting database: %s", dsn))
	}
	db = &Repo{DB: gdb}
}

func GetDB() *Repo {
	return db
}
