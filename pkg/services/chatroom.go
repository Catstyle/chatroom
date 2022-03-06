package services

import (
	"fmt"
	"strings"
	"sync"
	"time"

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
	err := room.UserJoin(user)
	if err == nil {
		user.RoomId = room.ID
		room.Broadcast("Chat.UserJoin", utils.M{"user": *user.User})
	}
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

	text = GetSensitiveService().Filter(text, '*')
	msg := room.NewChatMessage(models.CMText, user.User, text)
	room.Broadcast("Chat.TextMessage", msg)
	return nil
}

func (svc *chatroomService) PopularWord(roomId uint32) (word string, err error) {
	svc.lock.Lock()
	defer svc.lock.Unlock()
	if room, ok := svc.rooms[roomId]; ok {
		// message CTime is UnixMicro
		timeLimit := time.Now().Add(-10 * time.Minute).UnixMilli()
		maxCount := 0
		word = ""
		messages := room.GetMessagesByType(models.CMText)
		counter := make(map[string]int)
		for idx := len(messages) - 1; idx >= 0; idx-- {
			msg := messages[idx]
			if msg.CTime < timeLimit {
				break
			}
			for _, w := range strings.Split(msg.MsgData, " ") {
				counter[w] += 1
				if counter[w] > maxCount {
					maxCount = counter[w]
					word = w
				}
			}
		}
	} else {
		err = fmt.Errorf("room %d not found", roomId)
	}
	return word, err
}

func (svc *chatroomService) Leave(user *models.OnlineUser, roomId uint32) {
	svc.lock.Lock()
	defer svc.lock.Unlock()

	room, ok := svc.rooms[roomId]
	if ok {
		room.Broadcast("Chat.UserLeave", utils.M{"user": *user.User})
		delete(room.Users, user.User.ID)
	}
}
