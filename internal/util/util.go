package util

import (
	"os"

	"github.com/charmbracelet/log"
)

func FatalIfErr(err error, msg string) {
	if err != nil {
		log.Fatal(msg)
	}
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}
