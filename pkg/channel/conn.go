package channel

import (
	"bufio"
	"io"
	"log"
	"net"

	"github.com/catstyle/chatroom/pkg/protos"
)

type Conn struct {
	net.Conn
	ConnId    int
	Protocol  protos.Protocol
	sendQueue chan []byte
	RecvQueue chan *protos.Message
}

func NewConn(connId int, conn net.Conn, protocol protos.Protocol) *Conn {
	return &Conn{
		ConnId:    connId,
		Conn:      conn,
		Protocol:  protocol,
		sendQueue: make(chan []byte),
		// maybe some args to control this cap ?
		RecvQueue: make(chan *protos.Message, 10),
	}
}

func (conn *Conn) Close() error {
	close(conn.sendQueue)
	// close(conn.RecvQueue)
	return conn.Conn.Close()
}

func (conn *Conn) StartWriter() {
	for data := range conn.sendQueue {
		if _, err := conn.Conn.Write(data); err != nil {
			// TODO: what to do, close this conn ?
			log.Printf(
				"%s: write data error, %s", conn.Conn.RemoteAddr().String(), err,
			)
		}
	}
	log.Printf("%s: done writer", conn.Conn.RemoteAddr().String())
}

func (conn *Conn) StartReader() {
	addr := conn.RemoteAddr().String()
	buf := bufio.NewReader(conn)
	for {
		// need to add Heartbeat
		// conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		msg, err := conn.Protocol.DecodeMessage(buf)
		if err != nil {
			if err == io.EOF {
				log.Printf("%s: read EOF\n", addr)
			} else {
				log.Printf("%s: decode proto error %s\n", addr, err)
			}
			conn.RecvQueue <- nil
			close(conn.RecvQueue)
			break
		}
		conn.RecvQueue <- msg
	}
}

func (conn *Conn) SendMessage(msg *protos.Message, data interface{}) error {
	body, err := conn.Protocol.EncodeMessageWithData(msg, data)
	if err == nil {
		conn.sendQueue <- body
	} else {
		log.Println(conn.Conn.RemoteAddr().String(), "SendMessage error", err)
	}
	return err
}
