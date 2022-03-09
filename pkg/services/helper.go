package services

import "github.com/catstyle/chatroom/pkg/channel"

// Do some init when start, like create global singleton services.
func Init() {
	InitUserService()
	InitChatroomService()
	InitSensitiveService()
}

// Hook that will be called when a conn is closing.
func OnConnClose(conn *channel.Conn) {
	userSvc := GetUserService()
	if ou, ok := userSvc.GetOnlineUser(conn); ok {
		if ou.RoomId != 0 {
			GetChatroomService().Leave(ou, ou.RoomId)
		}
		GetUserService().Logout(conn)
	}
}
