package global

type ResourceRecord struct {
	Name  string
	TTL   int
	Class string
	Type  string
	RData []string
}
