package main

import (
	"io/fs"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

func main() {
	if !exists("db.json") {
		os.WriteFile("db.json", []byte("{}"), fs.ModePerm)
	}

	if !exists("register.json") {
		os.WriteFile("register.json", []byte("{}"), fs.ModePerm)
	}

	mux := http.NewServeMux()
	api := humago.New(mux, huma.DefaultConfig("CWE API", "1.0.0"))
	huma.Post(api, "/register", registerUser)
	huma.Post(api, "/get", getMessages)
	huma.Post(api, "/send", sendMessage)

	http.ListenAndServe(":1337", mux)
}
