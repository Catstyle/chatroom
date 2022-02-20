package models

import "net"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"uniqueIndex"`
	TokenHash string
	CTime     int `gorm:"autoCreateTime"`
}

type OnlineUser struct {
	User *User
	// TODO: add custom Conn struct
	Conn net.Conn
}
