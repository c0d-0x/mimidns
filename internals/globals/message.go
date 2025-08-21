package globals

const (
	HEADERlENGTH = 12
	ISRESPONSE   = uint16(0x8000)
)

type Header struct {
	ID      uint16
	FLAG    [2]byte
	QDCOUNT uint16
	ANCOUNT uint16
	NSCOUNT uint16
	ARCOUNT uint16
}

type Query struct {
	NAME  string
	TYPE  MessageType
	CLASS MessageClass
}

type Answer struct {
	NAME     string
	TYPE     MessageType
	CLASS    MessageClass
	TTL      uint32
	RDLENGTH uint16
	RDATA    []string
}

type Message struct {
	MHeader    Header
	Question   []Query
	Answer     []Answer
	Authority  []Answer
	Additional []Answer
}
