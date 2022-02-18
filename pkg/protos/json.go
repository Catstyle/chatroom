package protos

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	// "log"
)

type MsgType int

const (
	Req = iota
	Resp
	Broadcast
)

type Message struct {
	MsgID   uint32  `json:"msg_id"`
	MsgType MsgType `json:"msg_type"`
	Method  string  `json:"method"`
	Data    string  `json:"data"`
}

type Protocol interface {
	Encode(msg *Message) ([]byte, error)
	Decode(rd *bufio.Reader) (*Message, error)
}

type JSONProtocol struct {
}

func (p *JSONProtocol) Encode(msg *Message) ([]byte, error) {
	return nil, nil
}

func (p *JSONProtocol) Decode(rd *bufio.Reader) (*Message, error) {
	var message Message
	var err error

	_, err = rd.Peek(2)
	if err != nil {
		return nil, err
	}

	var size uint16
	err = binary.Read(rd, binary.BigEndian, &size)
	if err != nil {
		return nil, err
	}
	// log.Println("size", size)
	_, err = rd.Peek(int(size))
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
	return &message, err
}
