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
	huma.Register(api, huma.Operation{
		OperationID:  "register",
		Method:       http.MethodPost,
		Path:         "/register",
		Summary:      "Register a user",
		MaxBodyBytes: 100 * 1024 * 1024}, registerUser)
	huma.Register(api, huma.Operation{
		OperationID:  "get",
		Method:       http.MethodPost,
		Path:         "/get",
		Summary:      "Get a message",
		MaxBodyBytes: 100 * 1024 * 1024}, getMessages)
	huma.Register(api, huma.Operation{
		OperationID:  "send",
		Method:       http.MethodPost,
		Path:         "/send",
		Summary:      "Send a message",
		MaxBodyBytes: 100 * 1024 * 1024}, sendMessage)

	http.ListenAndServe(":1337", mux)
}
