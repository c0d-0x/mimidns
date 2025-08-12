package main

import (
	"fmt"

	"github.com/c0d-0x/mimidns/internals/parser"
)

func main() {
	rrlist, _ := parser.ParseMaster("zones/example.txt")

	for _, rrecord := range rrlist {
		fmt.Println(rrecord)
	}
}
