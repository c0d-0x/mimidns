package server

import (
	"log"
	"net"

	"github.com/c0d-0x/mimidns/internals/globals"
	"github.com/c0d-0x/mimidns/internals/parser"
)

type Server struct {
	Conn            net.UDPConn
	Addr            net.UDPAddr
	ResourseRecords []globals.ResourceRecord
}

func NewServer(addr string, rr []globals.ResourceRecord) (*Server, error) {
	resolvedAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	return &Server{
		Addr:            *resolvedAddr,
		ResourseRecords: rr,
	}, nil
}

func (s *Server) handleConn() {
	buf := make([]byte, 512)
	for {
		n, addr, err := s.Conn.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
		}

		data := make([]byte, n)
		copy(data, buf[:n])

		go func() {
			message, err := parser.ParseMessage(data)
			if err != nil {
				log.Println(err)
				return
			}

			log.Println("msg: ", message)
			response := prepareRespond(*message, s.ResourseRecords)

			log.Println("response: ", response)
			s.Conn.WriteToUDP(response.ToBytes(), addr)
		}()
		buf = buf[0:]

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
