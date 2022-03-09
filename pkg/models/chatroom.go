package models

import (
	"time"

	"github.com/catstyle/chatroom/pkg/protos"
)

const (
	CRStatusNormal  = 1
	CRStatusDeleted = 100
)

type MsgType int

// Currently support text message only.
const (
	CMText MsgType = iota
)

// Holding message related data.
type ChatMessage struct {
	Sender  *User   `json:"sender"`
	CTime   int64   `json:"c_time"`
	MsgType MsgType `json:"msg_type"`
	MsgData string  `json:"msg_data"`
}

// Object that handle the chatroom data and provide some helper functions.
// Not save to db now.
type Chatroom struct {
	ID       uint32 `gorm:"primaryKey"`
	CTime    int64  `gorm:"autoCreateTime"`
	Status   int
	Users    map[uint32]*OnlineUser
	messages []*ChatMessage
}

func NewChatroom(roomId uint32) *Chatroom {
	return &Chatroom{
		ID:     roomId,
		CTime:  time.Now().UnixMilli(),
		Status: CRStatusNormal,
		Users:  make(map[uint32]*OnlineUser),
	}
}

func (room *Chatroom) NewChatMessage(
	msgType MsgType, sender *User, data string,
) *ChatMessage {
	msg := &ChatMessage{
		Sender:  sender,
		CTime:   time.Now().UnixMicro(),
		MsgType: msgType,
		MsgData: data,
	}
	room.messages = append(room.messages, msg)
	return msg
}

// Return the latest messages upto count size.
func (room *Chatroom) GetLatestMessages(count int) []*ChatMessage {
	start := len(room.messages) - count
	if start < 0 {
		start = 0
	}
	return room.messages[start:]
}

func (room *Chatroom) GetMessagesByType(msgType MsgType) []*ChatMessage {
	messages := []*ChatMessage{}
	for _, msg := range room.messages {
		if msg.MsgType == msgType {
			messages = append(messages, msg)
		}
	}
	return messages
}

// Broadcast message around attenders.
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

func (room *Chatroom) UserLeave(user *OnlineUser) error {
	delete(room.Users, user.User.ID)
	return nil
}
