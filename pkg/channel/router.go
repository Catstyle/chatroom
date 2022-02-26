package channel

import (
	"log"
	"reflect"
	"strings"

	"github.com/catstyle/chatroom/pkg/protos"
	"github.com/catstyle/chatroom/utils"
)

type Router interface {
	GetName() string
	GetArgs() reflect.Value
	GetReply() reflect.Value
	Dispatch(*Conn, *protos.Message) (error, bool)
}

type Route struct {
	name      string
	receiver  reflect.Value
	method    reflect.Method
	argsType  reflect.Type
	replyType reflect.Type
}

func NewRouters(any interface{}, prefix string) ([]Router, error) {
	var err error
	if prefix != "" && !strings.Contains(prefix, ".") {
		prefix += "."
	}

	apiType := reflect.TypeOf(any)
	receiver := reflect.ValueOf(any)
	routers := []Router{}
	for i := 0; i < apiType.NumMethod(); i++ {
		method := apiType.Method(i)
		mType := method.Type
		routers = append(routers, &Route{
			name:     prefix + method.Name,
			receiver: receiver,
			method:   method,
			// assume the first args is conn
			argsType:  mType.In(2).Elem(),
			replyType: mType.In(3).Elem(),
		})
	}
	return routers, err

}

func (r *Route) GetName() string {
	return r.name
}

func (r *Route) GetArgs() reflect.Value {
	return reflect.New(r.argsType)
}

func (r *Route) GetReply() reflect.Value {
	return reflect.New(r.replyType)
}

// return value indicate disconnect
func (r *Route) Dispatch(
	conn *Conn, msg *protos.Message,
) (error, bool) {
	var err error

	args := reflect.New(r.argsType)
	if err = conn.Protocol.DecodeData(msg.Data, args.Interface()); err != nil {
		return err, true
	}

	reply := reflect.New(r.replyType)

	errValues := r.method.Func.Call([]reflect.Value{
		r.receiver, reflect.ValueOf(conn), args, reply,
	})

	errValue := errValues[0].Interface()
	if errValue != nil {
		err = errValue.(error)
		log.Println("call method error", r.name, args, reply, err)
	}

	var msgType protos.MsgType
	var replyData interface{}
	if err == nil {
		msgType = protos.RESP
		replyData = reply.Interface()
	}
	if err != nil {
		msgType = protos.ERROR
		replyData = utils.M{"error": err.Error()}
	}
	conn.SendMessage(msg.Convert(msgType), replyData)

	return err, false
}
