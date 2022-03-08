package channel_test

import (
	"testing"

	"github.com/catstyle/chatroom/pkg/channel"
	"github.com/stretchr/testify/assert"
)

type Api struct {
}

type ApiArgs struct {
}

type ApiReply struct {
}

func (api *Api) Exported(conn *channel.Conn, args *ApiArgs, reply *ApiReply) error {
	return nil
}

func (api *Api) inner() {
}

func TestRouter(t *testing.T) {
	routers, err := channel.NewRouters(&Api{}, "api")

	assert.Nil(t, err)
	assert.Equal(t, 1, len(routers))

	router := routers[0]
	assert.Equal(t, "api.Exported", router.GetName())
	assert.IsType(t, &ApiArgs{}, router.GetArgs().Interface())
	assert.IsType(t, &ApiReply{}, router.GetReply().Interface())
}
