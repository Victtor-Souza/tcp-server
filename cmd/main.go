package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

const (
	TCP = "tcp"
)

type Server struct {
	listenAddr string
	ln         net.Listener
	quitchn    chan struct{}
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitchn:    make(chan struct{}),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen(TCP, s.listenAddr)

	if err != nil {
		return err
	}

	defer ln.Close()
	s.ln = ln

	go s.acceptLoop()

	<-s.quitchn

	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}

		go s.readLoop(conn)
	}
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {

			if io.EOF == err {
				return
			}

			fmt.Println("read error:", err)
			continue
		}

		msg := buf[:n]
		fmt.Println(string(msg))

	}
}

func main() {
	server := NewServer(":9999")
	log.Fatal(server.Start())
}
