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
)

func newResourcesRecord(name *string, ttl *int, class *string, rData []string) *RequestRecods {
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

func splitchunk(chunk *string, cc string) *[]string {
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

func isSingleLinedSOA(chunk []string) bool {
	if chunk == nil {
		return false
	}

	return slices.Contains(chunk, "(") && slices.Contains(chunk, ")")
}

func isValidRecClass(subStr *string) bool {
	return slices.Contains(RecClasses, *subStr)
}

func isValidRecType(subStr *string) bool {
	return slices.Contains(RecTypes, *subStr)
}

func ParseMasterFile(fileName string) ([]RequestRecods, error) {
	fd, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	stream := bufio.NewReader(fd)

	var tmpRecord *RequestRecods
	var rrlist []RequestRecods
	var defaultDomain string
	var baseDomain string
	defaultTTL := 0

	SoaFound := false
	for {
		line, err := stream.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.ReplaceAll(line, "\t", " ")
		line = strings.TrimSpace(line)
		if line = *stripComment(&line, ";"); line == "" {
			continue
		}

		subs := *splitchunk(&line, " ")

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

	defer fd.Close()
	return rrlist, nil
}
