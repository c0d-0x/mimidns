package server

import (
	"log"
	"net"
)

func NewServer(addr string) (*Server, error) {
	resolvedAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	return &Server{
		Addr: *resolvedAddr,
	}, nil
}

func (s *Server) handleConn() {
	for {
		n, addr, err := s.Conn.ReadFromUDP(s.buf[0:])
		if err != nil {
			log.Println(err)
		}

		log.Println("msg: ", string(s.buf[:n]))

		s.Conn.WriteToUDP([]byte("hello, there\r\r\n"), addr)
	}
}

func (s *Server) Run() error {
	conn, err := net.ListenUDP("udp", &s.Addr)
	if err != nil {
		return err
	}

	s.Conn = *conn
	s.handleConn()
	return nil
}
