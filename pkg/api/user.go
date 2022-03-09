package api

import (
	"github.com/catstyle/chatroom/pkg/channel"
	"github.com/catstyle/chatroom/pkg/services"
)

const (
	userErrorIndex = 1000
	UserNotFound   = userErrorIndex + iota
	UserNameUsed
)

type LoginArgs struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

type LoginReply struct {
	Uid  uint32 `json:"uid"`
	Name string `json:"name"`
}

// UserApi by interface, methods need to follow the Router signature.
type UserApi interface {
	Login(*channel.Conn, *LoginArgs, *LoginReply) error
}

type userApi struct {
}

func NewUserApi() UserApi {
	return &userApi{}
}

// Login will create an OnlineUser object.
// Create User object if needed.
func (api *userApi) Login(
	conn *channel.Conn, args *LoginArgs, reply *LoginReply,
) error {
	userSvc := services.GetUserService()
	user, err := userSvc.Login(conn, args.Name, args.Token)
	if err != nil {
		return err
	}
	reply.Uid = uint32(user.ID)
	reply.Name = user.Name
	return nil
}
