package services

import (
	"errors"
	"net"
	"sync"

	"github.com/catstyle/chatroom/pkg/models"
	"github.com/catstyle/chatroom/pkg/repo"
	"github.com/catstyle/chatroom/utils"
	"gorm.io/gorm"
)

const (
	userErrorIndex = 1000
	UserNotFound   = userErrorIndex + iota
	UserNameUsed
)

type userService struct {
	repo        *repo.UserRepo
	onlineUsers map[string]*models.OnlineUser
	lock        sync.Mutex
}

var userSvc *userService

func InitUserService() {
	userSvc = &userService{
		repo: repo.GetUserRepo(),
		onlineUsers: make(map[string]*models.OnlineUser),
	}
}

func GetUserService() *userService {
	return userSvc
}

func (svc *userService) Login(
	conn net.Conn, name, token string,
) (*models.User, error) {
	options := utils.GetOptions()
	token = utils.MD5Sum(token, options.GetString("SECRET_TOKEN"))
	user, err := svc.login(name, token)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		user, err = svc.createUser(name, token)
	}
	if err == nil {
		svc.createOnlineUser(user, conn)
	}
	return user, err
}

func (svc *userService) login(name, token string) (*models.User, error) {
	user, err := svc.repo.GetByName(name)
	if err != nil {
		return nil, err
	}
	if user.TokenHash != token {
		return nil, errors.New("invalid token")
	}
	return user, nil
}

func (svc *userService) createUser(name, token string) (*models.User, error) {
	return svc.repo.CreateUser(name, token)
}

func (svc *userService) createOnlineUser(user *models.User, conn net.Conn) {
	svc.lock.Lock()
	defer svc.lock.Unlock()

	svc.onlineUsers[conn.RemoteAddr().String()] = &models.OnlineUser{
		User: user,
		Conn: conn,
	}
}
