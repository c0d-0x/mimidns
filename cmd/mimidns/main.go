package main

import (
	"flag"
	"log"

	"github.com/c0d-0x/mimidns/internal/globals"
	"github.com/c0d-0x/mimidns/internal/parser"
	"github.com/c0d-0x/mimidns/internal/server"
)

func main() {
	port := flag.String("p", "3000", "specify the port to run the server")
	zones := flag.String("zones", "zones", "<path> specify zones' directory")

	flag.Parse()

	rrlist := parser.LoadZoneFiles(*zones)
	if rrlist == nil {
		log.Fatal("no zone entries")
	}

	serv, err := server.NewServer(":"+*port, rrlist)
	if err != nil {
		log.Fatal(err)
	}

	globals.DrawASCIIArt()
	if err := serv.Run(); err != nil {
		log.Fatal(err)
	}
}
