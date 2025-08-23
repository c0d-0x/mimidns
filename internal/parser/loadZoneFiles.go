package parser

import (
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/c0d-0x/mimidns/internal/globals"
)

func GetZonefiles(rootPath string) ([]string, error) {
	zoneFiles := []string{}

	entries, err := os.ReadDir(rootPath)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".zone" && !entry.IsDir() {
			zoneFiles = append(zoneFiles, entry.Name())
		}
	}

	return zoneFiles, nil
}

func LoadZoneFiles(rootPath string) []globals.ResourceRecord {
	files, err := GetZonefiles(rootPath)
	if err != nil {
		log.Println(err)
		return nil
	}

	rrlist := []globals.ResourceRecord{}
	for _, file := range files {
		tmplist, err := ParseMasterFile(path.Join(rootPath, file))
		if err != nil {
			log.Fatal(err)
		}

		rrlist = append(rrlist, tmplist...)

	}

	return rrlist
}
