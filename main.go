package main

import (
	"github.com/catstyle/chatroom/pkg/api"
	"github.com/catstyle/chatroom/pkg/channel"
	"github.com/catstyle/chatroom/pkg/db"
	"github.com/catstyle/chatroom/pkg/models"
	"github.com/catstyle/chatroom/pkg/protos"
	"github.com/catstyle/chatroom/pkg/services"
	"github.com/catstyle/chatroom/utils"
)

const SECRET_TOKEN = "donot change me once used"

func main() {
	options := utils.LoadOptions("./conf.json")
	options.Set("SECRET_TOKEN", SECRET_TOKEN)
	options.SetDefault("SERVER_BIND", "localhost:5002")

	db.Setup(options.GetString("DB_DSN"))
	models.Init()
	services.Init()

	server := channel.NewTCPServer(
		channel.ServerConfig{
			Bind:     options.GetString("SERVER_BIND"),
			Protocol: &protos.JSONProtocol{},
		},
	)
	server.AddRouter(api.NewUserApi(), "User")
	server.Start()
}
