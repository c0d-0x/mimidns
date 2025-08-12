package parser

type (
	CLASS string
)

type SoaType struct {
	name    string
	rname   string
	serial  int32
	refresh int
	entry   int
	expire  int
	minimum int
}

type RequestRecods struct {
	Name  string
	TTL   int
	class string
	rdata []string
}
