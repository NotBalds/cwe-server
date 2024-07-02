package main

import (
	"io/fs"
	"net/http"
	"os"
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

	mux.HandleFunc("POST /get", getMessages)
	mux.HandleFunc("POST /send", sendMessage)
	mux.HandleFunc("POST /register", registerUser)

	http.ListenAndServe(":1337", mux)
}
