package services

import (
	"errors"

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
	repo *repo.UserRepo
}

func (svc *userService) Login(name, token string) (*models.User, error) {
	options := utils.GetOptions()
	token = utils.MD5Sum(token, options.GetString("SECRET_TOKEN"))
	user, err := svc.login(name, token)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		user, err = svc.createUser(name, token)
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
