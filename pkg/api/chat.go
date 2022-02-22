package api

import (
	"errors"
	"net"

	"github.com/catstyle/chatroom/pkg/models"
	"github.com/catstyle/chatroom/pkg/services"
)

const (
	chatErrorIndex = 1000
	ChatroomFull   = chatErrorIndex + iota
)

type JoinArgs struct {
	RoomId uint32
}

type JoinReply struct {
	RoomId uint32
	Users  map[uint32]*models.OnlineUser
}

type ChatroomApi interface {
	Join(net.Conn, *JoinArgs, *JoinReply) error
}

type chatroomApi struct {
}

func NewChatroomApi() ChatroomApi {
	return &chatroomApi{}
}

func (api *chatroomApi) Join(
	conn net.Conn, args *JoinArgs, reply *JoinReply,
) (err error) {
	userSvc := services.GetUserService()
	onlineUser, ok := userSvc.GetOnlineUser(conn)
	if !ok {
		return errors.New("please call Login first")
	}

	chatroomSvc := services.GetChatroomService()
	if room, err := chatroomSvc.Join(conn, args.RoomId, onlineUser); err == nil {
		reply.RoomId = room.ID
		reply.Users = room.Users
	}
	return err
}
