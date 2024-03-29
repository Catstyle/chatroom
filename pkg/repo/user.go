package repo

import (
	"github.com/catstyle/chatroom/pkg/db"
	"github.com/catstyle/chatroom/pkg/models"
)

type UserRepo struct {
	DB *db.Repo
}

var userRepo *UserRepo

func GetUserRepo() *UserRepo {
	if userRepo == nil {
		userRepo = &UserRepo{DB: db.GetDB()}
	}
	return userRepo
}

func (repo *UserRepo) CreateUser(name, token string) (*models.User, error) {
	user := models.User{
		Name:      name,
		TokenHash: token,
	}
	if err := repo.DB.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepo) GetByName(name string) (*models.User, error) {
	var user models.User
	if err := repo.DB.Where(&models.User{Name: name}).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
