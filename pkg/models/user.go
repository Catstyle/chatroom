package models

import "github.com/catstyle/chatroom/pkg/channel"

type User struct {
	ID        uint32 `gorm:"primaryKey" json:"id"`
	Name      string `gorm:"uniqueIndex" json:"name"`
	TokenHash string `json:"-"`
	CTime     int64  `gorm:"autoCreateTime" json:"-"`
}

type OnlineUser struct {
	User      *User
	Conn      *channel.Conn `json:"-"`
	RoomId    uint32
	LoginTime int64
}
