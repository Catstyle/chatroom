package api

import (
	"net"

	"github.com/catstyle/chatroom/pkg/services"
)

const (
	userErrorIndex = 1000
	UserNotFound   = userErrorIndex + iota
	UserNameUsed
)

type LoginArgs struct {
	Nickname string
	Token    string
}

type LoginReply struct {
	Uid      uint32
	Nickname string
}

type UserApi interface {
	Login(net.Conn, *LoginArgs, *LoginReply) error
}

type userApi struct {
}

func NewUserApi() UserApi {
	return &userApi{}
}

func (api *userApi) Login(
	conn net.Conn, args *LoginArgs, reply *LoginReply,
) error {
	userSvc := services.GetUserService()
	user, err := userSvc.Login(conn, args.Nickname, args.Token)
	if err != nil {
		return err
	}
	reply.Uid = uint32(user.ID)
	return nil
}
