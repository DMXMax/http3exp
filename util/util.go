package util

import (
	"log"
	"os"
	"path"
)

func GetCertFilePath(pathAdd string) string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	return path.Join(wd, pathAdd)
}
