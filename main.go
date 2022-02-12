package main

import (
	"github.com/catstyle/chatroom/pkg/db"
	"github.com/catstyle/chatroom/pkg/models"
	"github.com/catstyle/chatroom/utils"
)

const SECRET_TOKEN = "donot change me once used"

func main() {
	options := utils.LoadOptions("./conf.json")
	options.Set("SECRET_TOKEN", SECRET_TOKEN)

	db.Setup(options.GetString("DB_DSN"))
	models.Init()
}
