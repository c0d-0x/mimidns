package server

import (
	"encoding/binary"
	"log"

	"github.com/c0d-0x/mimidns/internals/globals"
)

func resolver(message *globals.Message, resourceRecords []globals.ResourceRecord) []globals.ResourceRecord {
	rr := []globals.ResourceRecord{}
	for _, query := range message.Question {
		for _, record := range resourceRecords {
			if query.NAME == record.Name && query.CLASS.String() == record.Class && query.TYPE.String() == record.Type {
				rr = append(rr, record)
			}
		}
	}
	return rr
}

func resourceRecordsToAnswers(rrs []globals.ResourceRecord) []globals.Answer {
	answers := []globals.Answer{}
	for _, rr := range rrs {
		aa := &globals.Answer{NAME: rr.Name, TTL: uint32(rr.TTL)}
		aa.TYPE.StrToMessageType(rr.Type)
		aa.CLASS.StrToMessageClass(rr.Class)

		aa.RDATA = append(aa.RDATA, rr.RData...)

		aa.RDLENGTH = uint16(len(aa.RDATA))
		answers = append(answers, *aa)
	}
	return answers
}

func prepareRespond(message globals.Message, resourceRecords []globals.ResourceRecord) *globals.Message {
	response := &globals.Message{}

	rrs := resolver(&message, resourceRecords)
	if rrs == nil {
		log.Printf("record not found: %v\n", message.Question)
		return nil
	}
	/* TODO: set flags accordingly */
	response.MHeader.ID = message.MHeader.ID
	binary.BigEndian.PutUint16(message.MHeader.FLAG[:], 0x8000)
	response.MHeader.QDCOUNT = message.MHeader.QDCOUNT
	response.MHeader.ANCOUNT = uint16(len(rrs))
	response.Question = message.Question
	response.Answer = append(response.Answer, resourceRecordsToAnswers(rrs)...)

	/* TODO: set authorities and additional data */
	return response
}
