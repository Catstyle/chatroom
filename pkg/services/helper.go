package services

import "github.com/catstyle/chatroom/pkg/channel"

func Init() {
	InitUserService()
	InitChatroomService()
	InitSensitiveService()
}

func OnConnClose(conn *channel.Conn) {
	userSvc := GetUserService()
	if ou, ok := userSvc.GetOnlineUser(conn); ok {
		if ou.RoomId != 0 {
			GetChatroomService().Leave(ou, ou.RoomId)
		}
		GetUserService().Logout(conn)
	}
}
