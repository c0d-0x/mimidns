package server

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"net"
)

const (
	headerLength = 12
	isResponse   = uint16(0x8000)
)

func parseMessage(buf []byte) (*Message, error) {
	message := &Message{}
	if len(buf) < headerLength {
		return nil, errors.New("invalid message")
	}
	reader := bytes.NewReader(buf)

	binary.Read(reader, binary.BigEndian, &message.Header)
	/* message.Header.ID = binary.BigEndian.Uint16(buf[:2]) */
	/* message.Header.FLAG = [2]byte(buf[2:4]) */
	/* message.Header.QDCOUNT = binary.BigEndian.Uint16(buf[4:6]) */
	/* message.Header.ANCOUNT = binary.BigEndian.Uint16(buf[6:8]) */
	/* message.Header.ARCOUNT = binary.BigEndian.Uint16(buf[8:12]) */

	messageType := binary.BigEndian.Uint16(message.Header.FLAG[:2]) & isResponse

	/* TODO: parse message accordingly */
	if isResponse == messageType {
		log.Println("it's a response")
	} else {
		log.Println("it's a query")
	}

	return message, nil
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
		_, addr, err := s.Conn.ReadFromUDP(s.buf[:])
		if err != nil {
			log.Println(err)
		}

		message, err := parseMessage(s.buf)
		if err != nil {
			log.Println(err)
		}

		log.Println("msg: ", message)

		/* TODO: decode Message and send a respond */
		s.Conn.WriteToUDP([]byte("example.com 300 A 127.0.0.1\r\n"), addr)
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
