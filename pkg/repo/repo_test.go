package repo_test

import (
	"errors"
	"testing"

	"github.com/catstyle/chatroom/pkg/db"
	"github.com/catstyle/chatroom/pkg/models"
	"github.com/catstyle/chatroom/pkg/repo"
	"github.com/catstyle/chatroom/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const SECRET_TOKEN = "12345"

func reset() {
	db.Setup("file::memory:")
	models.Init()
}

func TestCreateUser(t *testing.T) {
	reset()
	repo := repo.UserRepo{db.GetDB()}

	user, err := repo.CreateUser("test", utils.MD5Sum("test.token", SECRET_TOKEN))
	assert.Nil(t, err)
	assert.Equal(t, "test", user.Name)
}

func TestCreateUserDuplicatedName(t *testing.T) {
	reset()
	repo := repo.UserRepo{db.GetDB()}

	user, err := repo.CreateUser("test", utils.MD5Sum("test.token", SECRET_TOKEN))
	assert.Nil(t, err)
	assert.Equal(t, "test", user.Name)

	user2, err2 := repo.CreateUser("test", utils.MD5Sum("test.token", SECRET_TOKEN))
	assert.NotNil(t, err2)
	assert.Nil(t, user2)
}

func TestGetUserByName(t *testing.T) {
	TestCreateUser(t)
	repo := repo.UserRepo{db.GetDB()}

	user, err := repo.GetByName("test")
	assert.Nil(t, err)
	assert.Equal(t, "test", user.Name)
}

func TestGetUserByNameNotFound(t *testing.T) {
	TestCreateUser(t)
	repo := repo.UserRepo{db.GetDB()}

	user, err := repo.GetByName("test2")
	assert.NotNil(t, err)
	assert.Nil(t, user)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

