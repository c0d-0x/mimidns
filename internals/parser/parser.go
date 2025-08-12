package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var RecTypes []string = []string{
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

var RecClasses []string = []string{
	"IN",
	"CS",
	"CH",
	"HS",
}

func newResourcesRecord(name *string, ttl *int, class *string, rData []string) *RequestRecods {
	rr := RequestRecods{Name: *name, TTL: *ttl, class: *class, rdata: rData}
	return &rr
}

func stripComment(chunk *string, cc string) *string {
	if chunk == nil {
		return nil
	}

	before, _, found := strings.Cut(*chunk, cc)
	if !found {
		return chunk
	}

	return &before
}

func splitchunk(chunk *string, cc string) []string {
	if chunk == nil {
		return nil
	}

	reg, _ := regexp.Compile("s+")
	subs := strings.Split(*chunk, cc)
	for i, sub := range subs {
		subs[i] = reg.ReplaceAllString(sub, "")
	}
	return subs
}

func isSingleLinedSOA(chunk *string) bool {
	if chunk == nil {
		return false
	}

	return strings.Contains(*chunk, "(") && strings.Contains(*chunk, ")")
}

func isValidRecClass(subStr string) bool {
	return slices.Contains(RecClasses, subStr)
}

func ParseMaster(fileName string) ([]RequestRecods, error) {
	fd, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	stream := bufio.NewReader(fd)

	var tmpRecord *RequestRecods
	var rrlist []RequestRecods
	defaultDomain := ""
	baseDomain := ""
	defaultTTL := 0

	SoaFound := false
	for {
		line, err := stream.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		line = strings.ReplaceAll(line, "\t", " ")
		if line = *stripComment(&line, ";"); line == "" {
			continue
		}

		if strings.Contains(line, "SOA") && !SoaFound {
			SoaFound = true
			if !isSingleLinedSOA(&line) {
				/* parse muli-lined SOA  */
				continue
			} else {
				/* parse single-lined SOA  */
				/* fmt.Println("works") */
				continue
			}
		}

		subs := splitchunk(&line, " ")
		switch subs[0] {
		case "$ORIGIN":
			defaultDomain = subs[1]
			baseDomain = subs[1]
			continue

		case "$TTL":
			defaultTTL, _ = strconv.Atoi(subs[1])
			continue
		}

		if !isValidRecClass(subs[0]) {
			if subs[0] != "@" {
				defaultDomain = fmt.Sprintf("%s.%s", subs[0], baseDomain)
			}
			subs = subs[1:]
		}

		tmpRecord = newResourcesRecord(&defaultDomain, &defaultTTL, &subs[0], subs[1:])
		/* fmt.Println("name: ", (*tmpRecord).Name, "ttl: ", tmpRecord.TTL, "rData: ", tmpRecord.rdata) */
		rrlist = append(rrlist, *tmpRecord)

	}

	defer fd.Close()
	return rrlist, nil
}
