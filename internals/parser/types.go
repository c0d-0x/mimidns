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

var RecTypes = []string{
	"A",
	"NS",
	"MD",
	"MF",
	"CNAME",
	"SOA",
	"MB",
	"MG",
	"MR",
	"NULL",
	"WKS",
	"PTR",
	"HINFO",
	"MINFO",
	"MX",
	"TXT",
	"AAAA",
}

var RecClasses = []string{
	"IN",
	"CS",
	"CH",
	"HS",
}
