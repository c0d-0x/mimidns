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

func (mc *MessageClass) StrToMessageClass(str string) {
	switch str {
	case "IN":
		*mc = IN
	case "CS":
		*mc = CS
	case "CH":
		*mc = CH
	case "HS":
		*mc = HS
	}
}

func (mc MessageClass) String() string {
	if clss, ok := MessageClasses[mc]; ok {
		return clss
	}
	return ""
}

func (mt MessageType) String() string {
	if tp, ok := MessageTypes[mt]; ok {
		return tp
	}
	return ""
}

func (mt *MessageType) StrToMessageType(str string) {
	switch str {
	case "A":
		*mt = A
	case "NS":
		*mt = NS
	case "MD":
		*mt = MD
	case "MF":
		*mt = MF
	case "CNAME":
		*mt = CNAME
	case "SOA":
		*mt = SOA
	case "MB":
		*mt = MB
	case "MG":
		*mt = MG
	case "MR":
		*mt = MR
	case "NULL":
		*mt = NULL
	case "WKS":
		*mt = WKS
	case "PTR":
		*mt = PTR
	case "HINFO":
		*mt = HINFO
	case "MINFO":
		*mt = MINFO
	case "MX":
		*mt = MX
	case "TXT":
		*mt = TXT

	}
}
