package globals

import "encoding/binary"

const (
	HEADERlENGTH = 12
	ISRESPONSE   = uint16(0x8000)
)

type Header struct {
	ID      uint16
	FLAG    [2]byte
	QDCOUNT uint16
	ANCOUNT uint16
	NSCOUNT uint16
	ARCOUNT uint16
}

type Query struct {
	NAME  string
	TYPE  MessageType
	CLASS MessageClass
}

type Answer struct {
	NAME     string
	TYPE     MessageType
	CLASS    MessageClass
	TTL      uint32
	RDLENGTH uint16
	RDATA    []string
}

type Message struct {
	MHeader    Header
	Question   []Query
	Answer     []Answer
	Authority  []Answer
	Additional []Answer
}

func (msg *Message) ToBytes() []byte {
	var outBuffer []byte

	header := make([]byte, 12)
	binary.BigEndian.PutUint16(header[0:2], msg.MHeader.ID)
	copy(header[2:4], msg.MHeader.FLAG[:])
	binary.BigEndian.PutUint16(header[4:6], msg.MHeader.QDCOUNT)
	binary.BigEndian.PutUint16(header[6:8], msg.MHeader.ANCOUNT)
	binary.BigEndian.PutUint16(header[8:10], msg.MHeader.NSCOUNT)
	binary.BigEndian.PutUint16(header[10:12], msg.MHeader.ARCOUNT)
	outBuffer = append(outBuffer, header...)

	for _, query := range msg.Question {
		outBuffer = append(outBuffer, encodeQuery(&query)...)
	}

	for _, answer := range msg.Answer {
		outBuffer = append(outBuffer, encodeAnswer(&answer)...)
	}

	return outBuffer
}
