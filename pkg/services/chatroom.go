package services

import (
	"fmt"
	"sync"

	"github.com/catstyle/chatroom/pkg/models"
	"github.com/catstyle/chatroom/utils"
)

const (
	chatroomErrorIndex = 2000
	ChatroomNotFound   = chatroomErrorIndex + iota
	ChatroomDeleted
)

type chatroomService struct {
	rooms map[uint32]*models.Chatroom
	lock  sync.Mutex
}

var chatroomSvc *chatroomService

func InitChatroomService() {
	chatroomSvc = &chatroomService{
		rooms: make(map[uint32]*models.Chatroom),
	}
}

func GetChatroomService() *chatroomService {
	return chatroomSvc
}

func (svc *chatroomService) Join(
	user *models.OnlineUser, roomId uint32,
) (*models.Chatroom, error) {
	svc.lock.Lock()
	defer svc.lock.Unlock()

	var room *models.Chatroom
	var err error
	room, ok := svc.rooms[roomId]
	if !ok {
		room, err = svc.createRoom(roomId)
		if err != nil {
			return nil, err
		}
		svc.rooms[roomId] = room
	}

	svc.join(room, user)

	return room, nil
}

func (svc *chatroomService) createRoom(
	roomId uint32,
) (*models.Chatroom, error) {
	room := models.NewChatroom(roomId)
	return room, nil
}

func (svc *chatroomService) join(
	room *models.Chatroom, user *models.OnlineUser,
) error {
	room.Broadcast("Chat.UserJoin", utils.M{"user": *user.User})
	err := room.UserJoin(user)
	return err
}

func (svc *chatroomService) GetRoom(roomId uint32) (*models.Chatroom, bool) {
	svc.lock.Lock()
	defer svc.lock.Unlock()
	room, ok := svc.rooms[roomId]
	return room, ok
}

func (svc *chatroomService) SendText(
	user *models.OnlineUser, text string,
) error {
	room, ok := svc.GetRoom(user.RoomId)
	if !ok {
		return fmt.Errorf("no such room: %d", user.RoomId)
	}
	// TODO: filter text
	room.Broadcast(
		"Chat.TextMessage", utils.M{"user": *user.User, "text": text},
	)
	return nil
}
