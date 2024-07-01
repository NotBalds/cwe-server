package main

import (
	"io/fs"
	"net/http"
	"os"
)

func main() {
	os.WriteFile("db.json", []byte("{}"), fs.ModePerm)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /get", getMessages)
	mux.HandleFunc("POST /send", sendMessage)

	http.ListenAndServe(":1337", mux)
}
