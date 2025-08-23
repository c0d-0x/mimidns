package server

import (
	"fmt"
	"log"
	"net"

	"github.com/c0d-0x/mimidns/internal/globals"
	"github.com/c0d-0x/mimidns/internal/parser"
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

		log.Printf("Request from: %v\n", addr.String())
		go func() {
			message, err := parser.ParseMessage(data)
			if err != nil {
				log.Println(err)
				return
			}

			response := prepareRespond(*message, s.ResourseRecords)
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
	fmt.Printf("Running @127.0.0.1%s\n\n", s.Addr.String())

	s.Conn = *conn
	s.handleConn()

	return nil
}
