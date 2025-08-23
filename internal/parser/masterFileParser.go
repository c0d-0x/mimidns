package parser

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/c0d-0x/mimidns/internal/globals"
)

var Recglobals = []string{
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
}

var RecClasses = []string{
	"IN",
	"CS",
	"CH",
	"HS",
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

func splitStr(chunk *string, cc string) *[]string {
	if chunk == nil {
		return nil
	}

	reg, _ := regexp.Compile(`\s+`)
	subs := strings.Split(*chunk, cc)
	var tmpSubs []string
	for _, sub := range subs {
		if tmp := reg.ReplaceAllString(sub, ""); tmp != "" {
			tmpSubs = append(tmpSubs, tmp)
		}
	}
	return &tmpSubs
}

func isSingleLinedSOA(slice []string) bool {
	if slice == nil {
		return false
	}

	return slices.Contains(slice, "(") && slices.Contains(slice, ")")
}

func isValidRecClass(subStr *string) bool {
	return slices.Contains(RecClasses, *subStr)
}

func isValidRecType(subStr *string) bool {
	return slices.Contains(Recglobals, *subStr)
}

func newResourcesRecord(name *string, ttl *int, class *string, rData []string) *globals.ResourceRecord {
	if name == nil || ttl == nil || class == nil {
		return nil
	}

	if !isValidRecClass(class) {
		return nil
	}

	if !isValidRecType(&rData[0]) {
		return nil
	}

	hasSOA := slices.Contains(rData, "SOA")
	missingParens := !slices.Contains(rData, "(") || !slices.Contains(rData, ")")

	if hasSOA && missingParens {
		return nil
	}

	return &globals.ResourceRecord{
		Name:  *name,
		TTL:   *ttl,
		Class: *class,
		Type:  rData[0],
		RData: rData[1:],
	}
}

func ParseMasterFile(fileName string) ([]globals.ResourceRecord, error) {
	fd, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer fd.Close()
	stream := bufio.NewReader(fd)

	var tmpRecord *globals.ResourceRecord
	var rrlist []globals.ResourceRecord
	var defaultDomain string
	var baseDomain string
	defaultTTL := 0

	SoaFound := false
	for {
		buffer, err := stream.ReadString('\n')
		if err != nil {
			break
		}

		buffer = strings.ReplaceAll(buffer, "\t", " ")
		buffer = strings.TrimSpace(buffer)
		if buffer = *stripComment(&buffer, ";"); buffer == "" {
			continue
		}

		subs := *splitStr(&buffer, " ")

		switch subs[0] {
		case "$ORIGIN":
			defaultDomain = subs[1]
			baseDomain = subs[1]
			continue

		case "$TTL":
			defaultTTL, _ = strconv.Atoi(subs[1])
			continue
		case "$INCLUDE":
			/* TODO: To be handled properly with goroutines */
			continue

		}

		if !isValidRecClass(&subs[0]) {
			if subs[0] != "@" {
				if defaultDomain == "" {
					baseDomain = subs[0]
					defaultDomain = subs[0]
				} else {
					defaultDomain = fmt.Sprintf("%s.%s", subs[0], baseDomain)
				}
			} else {
				if defaultDomain == "" {
					errStr := fmt.Sprintf("No base damain found in: %s ", fileName)
					return nil, errors.New(errStr)
				}
			}

			if slices.Contains(subs, "SOA") && !SoaFound {
				SoaFound = true // only a single SOA can be found in a master file
				if !isSingleLinedSOA(subs) {
					/* parse muli-lined SOA  */
					for range 6 {
						opt, _ := stream.ReadString('\n')

						opt = strings.ReplaceAll(opt, "\t", " ")
						opt = *stripComment(&opt, ";")
						opt = strings.TrimSpace(opt)
						if opt != "" {
							subs = append(subs, opt)
						}
					}
				}
			}

			subs = subs[1:]
		}

		if tmpRecord = newResourcesRecord(&defaultDomain, &defaultTTL, &subs[0], subs[1:]); tmpRecord == nil {
			errStr := fmt.Sprintf("Invalid record in: %s\n", fileName)
			return nil, errors.New(errStr)
		}

		rrlist = append(rrlist, *tmpRecord)

	}

	return rrlist, nil
}
