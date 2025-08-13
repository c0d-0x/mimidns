package main

import (
	"fmt"
	"log"

	"github.com/c0d-0x/mimidns/internals/parser"
	"github.com/c0d-0x/mimidns/internals/server"
)

func main() {
	rrlist, err := parser.ParseMasterFile("zones/google.com.zone")
	if err != nil {
		log.Fatal(err)
	}

	rrlist2, err := parser.ParseMasterFile("zones/example.com.zone")
	if err != nil {
		log.Fatal(err)
	}

	rrlist = append(rrlist, rrlist2...)
	for _, rrecord := range rrlist {
		fmt.Println(rrecord)
	}

	serv, err := server.NewServer(":3000")
	if err != nil {
		log.Fatal(err)
	}

	if err := serv.Run(); err != nil {
		log.Fatal(err)
	}
}
