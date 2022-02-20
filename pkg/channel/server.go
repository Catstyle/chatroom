package channel

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/catstyle/chatroom/pkg/protos"
	"github.com/catstyle/chatroom/utils"
)

type ServerConfig struct {
	Bind     string
	Protocol protos.Protocol
}

type TCPServer struct {
	config   ServerConfig
	listener net.Listener
	conns    map[int]net.Conn
	done     chan bool
	wg       sync.WaitGroup

	routers map[string]Router
}

func NewTCPServer(config ServerConfig) *TCPServer {
	return &TCPServer{
		config:  config,
		conns:   make(map[int]net.Conn),
		done:    make(chan bool),
		routers: make(map[string]Router),
	}
}

func (s *TCPServer) Start() {
	listener, err := net.Listen("tcp", s.config.Bind)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	s.listener = listener

	log.Printf("started tcp server... %s\n", s.config.Bind)

	go s.catchSignal()

	s.done = make(chan bool)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("error accepting conn", err)
			break
		}
		log.Println("accepted conn", conn.RemoteAddr())

		s.conns[len(s.conns)] = conn
		s.wg.Add(1)
		go s.handler(conn, s.done)
	}

	s.wg.Wait()
}

func (s *TCPServer) AddRouter(any interface{}, prefix string) error {
	routers, err := NewRouters(any, prefix)
	if err != nil {
		return err
	}

	for _, router := range routers {
		name := router.GetName()
		if _, ok := s.routers[name]; ok {
			return fmt.Errorf("router name already used, %s", name)
		}
		s.routers[name] = router
	}
	return nil
}

func (s *TCPServer) catchSignal() {
	signal_chan := make(chan os.Signal, 1)
	signal.Notify(
		signal_chan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	sig := <-signal_chan
	s.shutdown(sig)
}

func (s *TCPServer) shutdown(sig os.Signal) {
	switch sig {
	// kill -SIGHUP XXXX
	case syscall.SIGHUP:
		log.Println("hungup")

	// kill -SIGINT XXXX or Ctrl+c
	case syscall.SIGINT:
		log.Println("Warikomi")

	// kill -SIGTERM XXXX
	case syscall.SIGTERM:
		log.Println("force stop")

	// kill -SIGQUIT XXXX
	case syscall.SIGQUIT:
		log.Println("stop and core dump")

	default:
		log.Println("Unknown signal.")
	}

	s.listener.Close()
	for idx := 0; idx < len(s.conns); idx++ {
		s.done <- true
	}
	s.conns = nil

}

func (s *TCPServer) handler(conn net.Conn, done chan bool) {
	addr := conn.RemoteAddr()

	defer func() {
		conn.Close()
		s.wg.Done()
		// if err := recover(); err != nil {
		// 	log.Printf("%s: recover from panic error %s\n", conn.RemoteAddr(), err)
		// }
		log.Printf("%s: closed\n", addr)
	}()

	out := make(chan *protos.Message)
	go func() {
		buf := bufio.NewReader(conn)
		for {
			conn.SetReadDeadline(time.Now().Add(60 * time.Second))
			msg, err := s.config.Protocol.DecodeMessage(buf)
			if err != nil {
				if err == io.EOF {
					log.Printf("%s: read EOF\n", conn)
				} else {
					log.Printf("%s: decode proto error %s\n", addr, err)
				}
				out <- nil
				close(out)
				break
			}
			out <- msg
		}
	}()

	// just do it here for convenience
	// should do it as hooks

	msg := <-out
	if msg.Method != "User.Login" {
		data, _ := s.config.Protocol.EncodeMessageWithData(
			msg.Convert(protos.ERROR),
			utils.M{"error": "should call Login first" + msg.Method},
		)
		conn.Write(data)
		return
	}
	err, exit := s.routers["User.Login"].Dispatch(conn, msg, s.config.Protocol)
	if err != nil {
		log.Printf(
			"%s: dispatch %s error %s, exit %v",
			addr, msg.Method, err, exit,
		)
	}
	if exit {
		return
	}

	for {
		select {
		case <-done:
			log.Println("closing conn due to server shutdown")
			return

		case msg := <-out:
			log.Printf("%s: receive message %d, %s", addr, msg.MsgID, msg.Method)
			if msg == nil {
				break
			}
			if router, ok := s.routers[msg.Method]; ok {
				err, exit := router.Dispatch(conn, msg, s.config.Protocol)
				if err != nil {
					log.Printf(
						"%s: dispatch %s error %s, exit %v",
						addr, msg.Method, err, exit,
					)
				}
				if exit {
					return
				}
			} else {
				log.Printf("%s: receive unknown message %#v", addr, msg)
				data, _ := s.config.Protocol.EncodeMessageWithData(
					msg.Convert(protos.ERROR),
					utils.M{"error": "unknown method" + msg.Method},
				)
				conn.Write(data)
			}
		}
	}
}
