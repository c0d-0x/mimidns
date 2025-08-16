package server

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	headerLength = 12
	isResponse   = uint16(0x8000)
)

func parseQname(buf []byte) (*string, uint8) {
	if buf == nil {
		return nil, 0
	}

	tmp := buf
	var _name []string

	nameLength := uint8(0)
	for {
		labelLen := tmp[0]
		if labelLen == 0 {
			break
		}

		_name = append(_name, string(tmp[1:labelLen]))
		fmt.Println("_name: ", _name)
		nameLength += labelLen + 1
		tmp = tmp[labelLen:]

	}
	name := strings.Join(_name, ".")
	name = strings.TrimSpace(name)
	return &name, nameLength + 1
}

func parseMessage(buf []byte) (*Message, error) {
	message := &Message{}
	if len(buf) < headerLength {
		return nil, errors.New("invalid message")
	}

	reader := bytes.NewReader(buf)
	binary.Read(reader, binary.BigEndian, &message.Header)

	messageType := binary.BigEndian.Uint16(message.Header.FLAG[:2]) & isResponse
	buf = buf[12:]

	/* TODO: parse message accordingly */
	if isResponse == messageType {
		log.Println("it's a response")
	} else {
		/* parse query */
		for range message.Header.QDCOUNT {
			qname, len := parseQname(buf)
			if len == 0 {
				break
			}
			buf = buf[len:]
			_query := query{
				NAME:  *qname,
				TYPE:  [2]byte(buf[:2]),
				CLASS: [2]byte(buf[2:4]),
			}
			buf = buf[4:]
			message.Question = append(message.Question, _query)
		}
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

		log.Println("msg: ", message.Question[len(message.Question)-1].NAME)

		/* TODO: decode Message and send a respond */
		s.Conn.WriteToUDP([]byte("example.com 300 A 127.0.0.1\r\n"), addr)
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
