package api

import (
	"fmt"
	"log"

	"github.com/catstyle/chatroom/pkg/channel"
	"github.com/catstyle/chatroom/pkg/models"
	"github.com/catstyle/chatroom/pkg/services"
)

const (
	chatErrorIndex = 2000
	ChatroomFull   = chatErrorIndex + iota
)

type JoinArgs struct {
	RoomId uint32 `json:"room_id"`
}

type JoinReply struct {
	RoomId   uint32                        `json:"room_id"`
	Users    map[uint32]*models.OnlineUser `json:"users"`
	Messages []*models.ChatMessage         `json:"messages"`
}

type EmptyReply struct {
}

type TextMessage struct {
	Text string `json:"text"`
}

// ChatroomApi by interface, methods need to follow the Router signature.
type ChatroomApi interface {
	Join(*channel.Conn, *JoinArgs, *JoinReply) error
	SendText(*channel.Conn, *TextMessage, *EmptyReply) error
}

type chatroomApi struct {
}

func NewChatroomApi() ChatroomApi {
	return &chatroomApi{}
}

// Join will check if user is in another room.
// Broadcast UserJoin to other users in room.
// Return the 50 latest messages.
func (api *chatroomApi) Join(
	conn *channel.Conn, args *JoinArgs, reply *JoinReply,
) (err error) {
	userSvc := services.GetUserService()
	onlineUser, ok := userSvc.GetOnlineUser(conn)
	if !ok {
		// TODO: add Error Warning different level as return value
		log.Panic("please call Login first")
	}

	chatroomSvc := services.GetChatroomService()

	if onlineUser.RoomId != 0 {
		if onlineUser.RoomId != args.RoomId {
			chatroomSvc.Leave(onlineUser, args.RoomId)
			onlineUser.RoomId = 0
		} else {
			return fmt.Errorf("already in room")
		}
	}

	if room, err := chatroomSvc.Join(onlineUser, args.RoomId); err == nil {
		reply.RoomId = room.ID
		reply.Users = room.Users
		reply.Messages = room.GetMessages(50)
	}
	return err
}

func (api *chatroomApi) SendText(
	conn *channel.Conn, args *TextMessage, reply *EmptyReply,
) (err error) {
	onlineUser := services.GetUserService().MustGetOnlineUser(conn)
	return services.GetChatroomService().SendText(onlineUser, args.Text)
}
