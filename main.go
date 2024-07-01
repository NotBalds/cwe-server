package main

import "net/http"

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /get", getMessages)
	mux.HandleFunc("POST /send", sendMessage)

	http.ListenAndServe(":1337", mux)
}
