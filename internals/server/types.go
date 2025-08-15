package server

import (
	"net"
)

type Server struct {
	Conn net.UDPConn
	Addr net.UDPAddr
	buf  []byte
}

type header struct {
	ID      uint16
	FLAG    [2]byte
	QDCOUNT uint16
	ANCOUNT uint16
	NSCOUNT uint16
	ARCOUNT uint16
}

type query struct {
	NAME  string
	TYPE  int16
	CLASS int16
}

type answer struct {
	NAME     string
	TYPE     uint16
	CLASS    uint16
	TTL      uint32
	RDLENGTH uint16
	RDATA    []byte
}

type Message struct {
	Header     header
	Question   []query
	Answer     []answer
	Authority  []byte
	Additional []byte
}
