package models

import "github.com/catstyle/chatroom/pkg/db"

func Init() {
	db := db.GetDB()
	db.AutoMigrate(&User{})
}
