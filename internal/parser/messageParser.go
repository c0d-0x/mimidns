package parser

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"

	"github.com/c0d-0x/mimidns/internal/globals"
)

func parseName(buf []byte) (*string, uint16) {
	if buf == nil {
		return nil, 0
	}

	var _name []string

	nameLength := uint8(0)
	for {
		labelLen := buf[0]

		if len(buf) <= int(labelLen+1) || labelLen == 0 {
			break
		}

		_name = append(_name, string(buf[1:labelLen+1]))
		nameLength += labelLen + 1

		buf = buf[labelLen+1:]
	}

	if _name == nil {
		return nil, 0
	}

	name := strings.Join(_name, ".")
	name = strings.TrimSpace(name)
	name += "." // terminating '.'
	return &name, uint16(nameLength + 1)
}

func parseAnswer(buf []byte, count uint16) ([]globals.Answer, uint16) {
	nbSum := uint16(0)
	answers := []globals.Answer{}
	for range count {
		name, nb := parseName(buf)
		if nb == 0 {
			break
		}
		buf = buf[nb:]
		answer := &globals.Answer{
			NAME:     *name,
			TYPE:     globals.MessageType(binary.BigEndian.Uint16(buf[:2])),
			CLASS:    globals.MessageClass(binary.BigEndian.Uint16(buf[2:4])),
			TTL:      binary.BigEndian.Uint32(buf[4:8]),
			RDLENGTH: binary.BigEndian.Uint16(buf[8:10]),
		}

		answer.RDATA = append(answer.RDATA, string(buf[10:10+answer.RDLENGTH]))
		nbSum += 10 + answer.RDLENGTH

		answers = append(answers, *answer)

	}

	return answers, nbSum
}

func ParseMessage(buf []byte) (*globals.Message, error) {
	message := globals.Message{}
	if len(buf) <= globals.HEADERlENGTH {
		return nil, errors.New("invalid message")
	}

	reader := bytes.NewReader(buf)
	binary.Read(reader, binary.BigEndian, &message.MHeader)

	buf = buf[globals.HEADERlENGTH:]

	for range message.MHeader.QDCOUNT {
		qname, len := parseName(buf)
		if qname == nil {
			return nil, errors.New("invalid qname")
		}

		buf = buf[len:]
		_query := globals.Query{
			NAME:  *qname,
			TYPE:  globals.MessageType(binary.BigEndian.Uint16(buf[:2])),
			CLASS: globals.MessageClass(binary.BigEndian.Uint16(buf[2:4])),
		}
		buf = buf[4:]
		message.Question = append(message.Question, _query)
	}

	answers, n := parseAnswer(buf, message.MHeader.ANCOUNT)
	if answers != nil {
		message.Answer = append(message.Answer, answers...)
	}

	buf = buf[n:]
	authorities, n := parseAnswer(buf, message.MHeader.NSCOUNT)
	if authorities != nil {
		message.Authority = append(message.Authority, authorities...)
	}

	buf = buf[n:]
	additional, _ := parseAnswer(buf, message.MHeader.ARCOUNT)

	if additional != nil {
		message.Additional = append(message.Additional, additional...)
	}

	return &message, nil
}
