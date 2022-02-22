package protos

import (
	"bufio"
	"encoding/binary"
	"encoding/json"

	"github.com/catstyle/chatroom/utils"
)

type MsgType int

const (
	REQ = iota
	RESP
	BROADCAST
	ERROR
)

type Message struct {
	MsgID         uint32  `json:"msg_id"`
	MsgType       MsgType `json:"msg_type"`
	Method        string  `json:"method"`
	ContentLength uint32  `json:"content_length"`
	Data          []byte  `json:"-"`
}

func NewMessage(msgId uint32, msgType MsgType, method string) *Message {
	return &Message{
		MsgID: msgId,
		MsgType: msgType,
		Method: method,
	}
}

func (m Message) Convert(msgType MsgType) *Message {
	return &Message{
		MsgID:   m.MsgID,
		MsgType: msgType,
	}
}

type Protocol interface {
	EncodeMessage(*Message) []byte
	EncodeMessageWithData(*Message, interface{}) ([]byte, error)
	EncodeData(interface{}) ([]byte, error)
	DecodeMessage(*bufio.Reader) (*Message, error)
	DecodeData([]byte, interface{}) error
}

type JSONProtocol struct {
}

func (p *JSONProtocol) EncodeMessage(msg *Message) []byte {
	data, _ := json.Marshal(msg)
	headerSize := make([]byte, 2)
	binary.BigEndian.PutUint16(headerSize, uint16(len(data)))
	return append(headerSize, data...)
}

func (p *JSONProtocol) EncodeData(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (p *JSONProtocol) EncodeMessageWithData(msg *Message, v interface{}) ([]byte, error) {
	var err error
	data, err := p.EncodeData(v)
	if err != nil {
		data, err = p.EncodeData(utils.M{"error": err})
	}
	msg.ContentLength = uint32(len(data))
	return append(p.EncodeMessage(msg), data...), err
}

func (p *JSONProtocol) DecodeMessage(rd *bufio.Reader) (*Message, error) {
	var message Message
	var err error

	var size uint16
	err = binary.Read(rd, binary.BigEndian, &size)
	if err != nil {
		return nil, err
	}

	body := make([]byte, size)
	err = binary.Read(rd, binary.BigEndian, body)
	if err != nil {
		return nil, err
	}
	// log.Println("body", string(body))
	if err = json.Unmarshal(body, &message); err != nil {
		return nil, err
	}

	data := make([]byte, message.ContentLength)
	err = binary.Read(rd, binary.BigEndian, data)
	if err != nil {
		return nil, err
	}
	message.Data = data

	return &message, err
}

func (p *JSONProtocol) DecodeData(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
