package models

import "github.com/catstyle/chatroom/pkg/db"

// Do some init when start, like AutoMigrate.
func Init() {
	db := db.GetDB()
	db.AutoMigrate(&User{})
}
