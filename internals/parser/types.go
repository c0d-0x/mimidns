package parser

type (
	TYPE  uint16
	CLASS uint8
)

const (
	IN CLASS = iota + 1
	CS
	CH
	HS
)

const (
	A TYPE = iota + 1
	NS
	MD
	MF
	CNAME
	SOA
	MB
	MG
	MR
	NULL
	WKS
	PTR
	HINFO
	MINFO
	MX
	TXT
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
	domain string
	ttl    int64
	class  CLASS
	rdata  []byte
}
