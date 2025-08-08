package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

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

func ParseMaster(fileName string) error {
	fd, err := os.Open(fileName)
	if err != nil {
		return err
	}

	stream := bufio.NewReader(fd)

	prevDomain := ""
	prevTTL := 0

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

		if strings.Contains(line, "SOA") {
			if !isSingleLinedSOA(&line) {
				/* parse muli-lined SOA  */
			} else {
				/* parse muli-lined SOA  */
			}
		}

		subs := splitchunk(&line, " ")
		switch subs[0] {
		case "$ORIGIN":
			prevDomain = subs[1]

		case "$TTL":
			prevTTL, _ = strconv.Atoi(subs[1])
		}

		fmt.Println("prevDomain: ", prevDomain, " prevTTL: ", prevTTL)
		fmt.Println(subs)
	}

	defer fd.Close()
	return nil
}
