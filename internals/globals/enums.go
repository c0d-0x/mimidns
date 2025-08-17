package globals

type (
	MessageClass uint16
	MessageType  uint16
)

const (
	A MessageType = iota + 1
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

const (
	IN MessageClass = iota + 1
	CS
	CH
	HS
)

var MessageTypes = map[MessageType]string{
	A:     "A",
	NS:    "NS",
	MD:    "MD",
	MF:    "MF",
	CNAME: "CNAME",
	SOA:   "SOA",
	MB:    "MB",
	MG:    "MG",
	MR:    "MR",
	NULL:  "NULL",
	WKS:   "WKS",
	PTR:   "PTR",
	HINFO: "HINFO",
	MINFO: "MINFO",
	MX:    "MX",
	TXT:   "TXT",
}

var MessageClasses = map[MessageClass]string{
	IN: "IN",
	CS: "CS",
	CH: "CH",
	HS: "HS",
}

func (MessageClass) String(cc MessageClass) string {
	if clss, ok := MessageClasses[cc]; ok {
		return clss
	}
	return ""
}

func (MessageType) String(tt MessageType) string {
	if tp, ok := MessageTypes[tt]; ok {
		return tp
	}
	return ""
}
