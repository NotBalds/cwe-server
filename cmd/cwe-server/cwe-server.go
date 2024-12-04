package main

import (
	"io/fs"
	"os"

	"github.com/NotBalds/cwe-server/internal/server"
	"github.com/NotBalds/cwe-server/internal/util"
	"github.com/joho/godotenv"
)

func main() {
	os.Mkdir("v1.0.0", os.ModePerm)
	os.Chdir("v1.0.0")
	godotenv.Load()
	if !util.Exists("db.json") {
		os.WriteFile("db.json", []byte("{}"), fs.ModePerm)
	}

	if !util.Exists("register.json") {
		os.WriteFile("register.json", []byte("{}"), fs.ModePerm)
	}

	server.Start()
}
