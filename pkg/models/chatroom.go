package models

import (
	"time"

	"github.com/catstyle/chatroom/pkg/protos"
)

const (
	CRStatusNormal  = 1
	CRStatusDeleted = 100
)

type Chatroom struct {
	ID     uint32 `gorm:"primaryKey"`
	CTime  int64  `gorm:"autoCreateTime"`
	Status int
	Users  map[uint32]*OnlineUser
}

func NewChatroom(roomId uint32) *Chatroom {
	return &Chatroom{
		ID:     roomId,
		CTime:  time.Now().UnixMilli(),
		Status: CRStatusNormal,
		Users:  make(map[uint32]*OnlineUser),
	}
}

func (room *Chatroom) Broadcast(method string, data interface{}) {
	msg := protos.NewMessage(0, protos.BROADCAST, method)
	for _, peer := range room.Users {
		peer.Conn.SendMessage(msg, data)
	}
}

func (room *Chatroom) UserJoin(user *OnlineUser) error {
	if _, ok := room.Users[user.User.ID]; !ok {
		room.Users[user.User.ID] = user
	}
	return nil
}
