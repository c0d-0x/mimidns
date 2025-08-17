package parser

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"

	"github.com/c0d-0x/mimidns/internals/globals"
)

func parseName(buf []byte) (*string, uint8) {
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

		_name = append(_name, string(tmp[1:labelLen+1]))
		nameLength += labelLen + 1
		tmp = tmp[labelLen+1:]

	}

	if _name == nil {
		return nil, 0
	}

	name := strings.Join(_name, ".")
	name = strings.TrimSpace(name)
	return &name, nameLength + 1
}

func ParseMessage(buf []byte) (*globals.Message, error) {
	message := globals.Message{}
	if len(buf) < globals.HEADERlENGTH {
		return nil, errors.New("invalid message")
	}

	reader := bytes.NewReader(buf)
	binary.Read(reader, binary.BigEndian, &message.MHeader)

	messageType := binary.BigEndian.Uint16(message.MHeader.FLAG[:2]) & globals.ISRESPONSE
	/* header is 12 bytes */
	buf = buf[12:]

	/* TODO: parse message accordingly */
	if globals.ISRESPONSE != messageType {
		/* parse query */
		for range message.MHeader.QDCOUNT {
			qname, len := parseName(buf)
			if qname == nil || len == 0 {
				return nil, errors.New("invalid qname")
			}
			buf = buf[len:]
			_query := globals.Query{
				NAME:  *qname,
				TYPE:  [2]byte(buf[:2]),
				CLASS: [2]byte(buf[2:4]),
			}
			buf = buf[4:]
			message.Question = append(message.Question, _query)
		}
	}

	return &message, nil
}
