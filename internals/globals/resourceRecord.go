package globals

type SoaType struct {
	name    string
	rname   string
	serial  int32
	refresh int
	entry   int
	expire  int
	minimum int
}

type RequestRecord struct {
	Name  string
	TTL   int
	Class string
	Type  string
	RData []string
}
