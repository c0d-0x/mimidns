package globals

import (
	"encoding/binary"
	"log"
	"net"
	"strconv"
	"strings"
)

func encodeName(name string) []byte {
	var buf []byte
	labels := strings.SplitSeq(name, ".")
	for label := range labels {
		if len(label) == 0 {
			continue
		}
		buf = append(buf, byte(len(label)))
		buf = append(buf, label...)
	}
	buf = append(buf, 0)
	return buf
}

func encodeQuery(query *Query) []byte {
	b := encodeName(query.NAME)
	tmp := make([]byte, 4)
	binary.BigEndian.PutUint16(tmp[0:2], uint16(query.TYPE))
	binary.BigEndian.PutUint16(tmp[2:4], uint16(query.CLASS))
	return append(b, tmp...)
}

func encodeAnswer(answer *Answer) []byte {
	buffer := encodeName(answer.NAME)

	var rdata []byte
	switch answer.TYPE {
	case A:
		ip := net.ParseIP(answer.RDATA[0]).To4()
		if ip != nil {
			rdata = ip
		}

	case 28: // AAAA record IPV6
		ip := net.ParseIP(answer.RDATA[0]).To16()
		if ip != nil {
			rdata = ip
		}

	case CNAME, NS, PTR:
		rdata = encodeName(answer.RDATA[0])

	case TXT:
		for _, txt := range answer.RDATA {
			if len(txt) > 255 {
				txt = txt[:255] // truncate per RFC
			}
			rdata = append(rdata, byte(len(txt)))
			rdata = append(rdata, txt...)
		}
	case MX:
		priority, err := strconv.Atoi((answer.RDATA[0]))
		if err != nil {
			log.Println(err)
			return nil
		}
		rdata = binary.BigEndian.AppendUint16(rdata, uint16(priority))
		rdata = append(rdata, encodeName(answer.RDATA[1])...)
	case SOA:
		rdata = encodeName(answer.RDATA[0])
		rdata = append(rdata, encodeName(answer.RDATA[1])...)

		for _, valStr := range answer.RDATA[3 : len(answer.RDATA)-1] {
			val, err := strconv.Atoi(valStr)
			if err != nil {
				log.Println(err)
				return nil
			}
			rdata = binary.BigEndian.AppendUint32(rdata, uint32(val))
		}
	default:
		for _, s := range answer.RDATA {
			rdata = append(rdata, []byte(s)...)
		}
	}

	tmp := make([]byte, 10)
	binary.BigEndian.PutUint16(tmp[0:2], uint16(answer.TYPE))
	binary.BigEndian.PutUint16(tmp[2:4], uint16(answer.CLASS))
	binary.BigEndian.PutUint32(tmp[4:8], answer.TTL)
	binary.BigEndian.PutUint16(tmp[8:10], uint16(len(rdata)))

	return append(buffer, append(tmp, rdata...)...)
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
