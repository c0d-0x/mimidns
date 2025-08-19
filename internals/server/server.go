package server

import (
	"log"
	"net"

	"github.com/c0d-0x/mimidns/internals/parser"
)

type Server struct {
	Conn net.UDPConn
	Addr net.UDPAddr
	buf  []byte
}

func NewServer(addr string) (*Server, error) {
	resolvedAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	return &Server{
		Addr: *resolvedAddr,
		buf:  make([]byte, 512),
	}, nil
}

func (s *Server) handleConn() {
	for {
		n, addr, err := s.Conn.ReadFromUDP(s.buf)
		if err != nil {
			log.Println(err)
		}

		data := make([]byte, n)
		copy(data, s.buf[:n])

		go func() {
			message, err := parser.ParseMessage(data)

			if err != nil {
				log.Println(err)
			} else {
				log.Println("msg: ", message)
			}

			/*TODO: send a respond */
			s.Conn.WriteToUDP([]byte("example.com 300 A 127.0.0.1\r\n"), addr)
		}()
		s.buf = s.buf[0:]

	}
}

func (s *Server) Run() error {
	conn, err := net.ListenUDP("udp", &s.Addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	s.Conn = *conn
	s.handleConn()

	return nil
}
