package models

import "github.com/catstyle/chatroom/pkg/channel"

type User struct {
	ID        uint32 `gorm:"primaryKey" json:"id"`
	Name      string `gorm:"uniqueIndex" json:"name"`
	TokenHash string `json:"-"`
	CTime     int64  `gorm:"autoCreateTime" json:"-"`
}

type OnlineUser struct {
	User *User
	// TODO: add custom Conn struct
	Conn *channel.Conn `json:"-"`
}
