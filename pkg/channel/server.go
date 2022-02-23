package channel

import (
	"fmt"
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
	config         ServerConfig
	listener       net.Listener
	listening      bool
	conns          map[int]*Conn
	acceptingQueue chan *Conn
	closingQueue   chan *Conn
	done           chan bool
	wg             sync.WaitGroup

	routers map[string]Router
}

func NewTCPServer(config ServerConfig) *TCPServer {
	return &TCPServer{
		config:         config,
		conns:          make(map[int]*Conn),
		done:           make(chan bool),
		routers:        make(map[string]Router),
		acceptingQueue: make(chan *Conn),
		closingQueue:   make(chan *Conn),
	}
}

func (s *TCPServer) Start() {
	listener, err := net.Listen("tcp", s.config.Bind)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	s.listener = listener
	s.listening = true

	log.Printf("started tcp server... %s\n", s.config.Bind)

	go s.catchSignal()
	go s.accepter()

	ticker := time.Tick(10 * time.Second)
	for s.listening || len(s.conns) > 0 {
		log.Printf(
			"waiting for action, listening=%v, conns=%d",
			s.listening, len(s.conns),
		)
		select {
		case conn := <-s.acceptingQueue:
			s.conns[conn.ConnId] = conn
			s.wg.Add(1)
			go s.handler(conn)
		case conn := <-s.closingQueue:
			log.Printf("%s: closed", conn.Conn.RemoteAddr())
			conn.Close()
			delete(s.conns, conn.ConnId)
			s.wg.Done()
		case <-ticker:
			// do something or check if still listening
		}
	}

	s.wg.Wait()
	s.conns = nil
	log.Printf("bye")
}

func (s *TCPServer) accepter() {
	for s.listening {
		sock, err := s.listener.Accept()
		if err != nil {
			log.Println("error accepting conn", err)
			break
		}

		log.Println("accepted conn", sock.RemoteAddr())

		connId := len(s.conns)
		conn := NewConn(connId, sock, s.config.Protocol)
		s.acceptingQueue <- conn
	}
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
	s.listening = false
	log.Printf("%s: shutdown %d conns", s.config.Bind, len(s.conns))
	for idx := 0; idx < len(s.conns); idx++ {
		s.done <- true
	}
	log.Printf("%s: shutdown", s.config.Bind)
}

func (s *TCPServer) handler(conn *Conn) {
	addr := conn.RemoteAddr()

	defer func() {
		s.closingQueue <- conn
		if err := recover(); err != nil {
			log.Printf("%s: recover from panic error %+v\n", conn.RemoteAddr(), err)
		}
	}()

	go conn.StartWriter()
	go conn.StartReader()

	doneLogin := false
	for {
		select {
		case <-s.done:
			log.Println("closing conn due to server shutdown")
			return

		case msg := <-conn.RecvQueue:
			log.Printf("%s: receive message %d, %s", addr, msg.MsgID, msg.Method)
			if msg == nil {
				break
			}
			if !doneLogin && msg.Method != "User.Login" {
				log.Printf("%s: should call User.Login first, got %s", addr, msg.Method)
				conn.SendMessage(
					msg.Convert(protos.ERROR),
					utils.M{"error": "should call Login first" + msg.Method},
				)
				return
			}
			if router, ok := s.routers[msg.Method]; ok {
				err, exit := router.Dispatch(conn, msg)
				if err != nil {
					log.Printf(
						"%s: dispatch %s error %s, exit %v",
						addr, msg.Method, err, exit,
					)
				}
				if exit {
					return
				}
				doneLogin = true
			} else {
				log.Printf("%s: receive unknown message %s", addr, msg.Method)
				conn.SendMessage(
					msg.Convert(protos.ERROR),
					utils.M{"error": "unknown method" + msg.Method},
				)
			}
		}
	}
}
