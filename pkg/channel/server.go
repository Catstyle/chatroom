package channel

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/catstyle/chatroom/pkg/protos"
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
}

func NewTCPServer(config ServerConfig) *TCPServer {
	return &TCPServer{
		config: config,
		conns:  make(map[int]net.Conn),
		done:   make(chan bool),
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
	defer conn.Close()
	defer s.wg.Done()

	out := make(chan *protos.Message)
	go func() {
		buf := bufio.NewReader(conn)
		for {
			conn.SetReadDeadline(time.Now().Add(60 * time.Second))
			msg, err := s.config.Protocol.Decode(buf)
			if err == io.EOF {
				log.Printf("%s: read EOF\n", conn)
				break
			}
			if err != nil {
				log.Printf("%s: decode proto error %s\n", conn.RemoteAddr(), err)
				break
			}
			out <- msg
		}
	}()

	for {
		select {
		case <-done:
			log.Println("closing conn due to server shutdown")
			return

		case msg := <-out:
			log.Printf("%s receive message %#v", conn.RemoteAddr(), msg)
			if msg == nil {
				break
			}
			// TODO: handle the message
		}
	}
}
