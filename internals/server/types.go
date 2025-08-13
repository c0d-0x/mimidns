package server

import "net"

type Server struct {
	Conn net.UDPConn
	Addr net.UDPAddr
	buf  [512]byte
}
