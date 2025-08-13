package main

import (
	"fmt"
	"os"

	"github.com/c0d-0x/mimidns/internals/parser"
)

func main() {
	rrlist, err := parser.ParseMasterFile("zones/google.com.zone")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	rrlist2, err := parser.ParseMasterFile("zones/example.com.zone")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	rrlist = append(rrlist, rrlist2...)
	for _, rrecord := range rrlist {
		fmt.Println(rrecord)
	}
}
