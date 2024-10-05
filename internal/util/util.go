package util

import (
	"log"
	"os"
)

func FatalIfErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg)
	}
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}
