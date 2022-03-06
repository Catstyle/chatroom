package api

import (
	"github.com/catstyle/chatroom/pkg/channel"
	"github.com/catstyle/chatroom/pkg/models"
	"github.com/catstyle/chatroom/pkg/services"
)

const (
	statsErrorIndex = 3000
)

type OnlineStatsArgs struct {
	Username string `json:"username"`
}

type OnlineStatsReply struct {
	OnlineUser *models.OnlineUser
}

type PopularWordArgs struct {
	RoomId uint32 `json:"room_id"`
}

type PopularWordReply struct {
	Word string
}

type StatsApi interface {
	OnlineStats(*channel.Conn, *OnlineStatsArgs, *OnlineStatsReply) error
	PopularWord(*channel.Conn, *PopularWordArgs, *PopularWordReply) error
}

type statsApi struct {
}

func NewStatsApi() StatsApi {
	return &statsApi{}
}

func (api *statsApi) OnlineStats(
	conn *channel.Conn, args *OnlineStatsArgs, reply *OnlineStatsReply,
) (err error) {
	if onlineUser, err := services.GetUserService().GetOnlineUserByName(
		args.Username,
	); err == nil {
		reply.OnlineUser = onlineUser
	}
	return err
}

func (api *statsApi) PopularWord(
	conn *channel.Conn, args *PopularWordArgs, reply *PopularWordReply,
) (err error) {
	chatroomSvc := services.GetChatroomService()
	if word, err := chatroomSvc.PopularWord(args.RoomId); err == nil {
		reply.Word = word
	}
	return err
}
