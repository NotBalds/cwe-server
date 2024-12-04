package server

import (
	"net/http"
	"sync"

	"github.com/NotBalds/cwe-server/internal/functions"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

func httpStart(port string, wg *sync.WaitGroup) {
	defer wg.Done()
	mux := http.NewServeMux()
	api := humago.New(mux, huma.DefaultConfig("CWE API", "1.0.0"))
	huma.Register(api, huma.Operation{
		OperationID:  "register",
		Method:       http.MethodPost,
		Path:         "/register",
		Summary:      "Register a user",
		MaxBodyBytes: 20 * 1024 * 1024}, functions.RegisterUser)
	huma.Register(api, huma.Operation{
		OperationID:  "get",
		Method:       http.MethodPost,
		Path:         "/get",
		Summary:      "Get a message",
		MaxBodyBytes: 20 * 1024 * 1024}, functions.GetMessages)
	huma.Register(api, huma.Operation{
		OperationID:  "send",
		Method:       http.MethodPost,
		Path:         "/send",
		Summary:      "Send a message",
		MaxBodyBytes: 20 * 1024 * 1024}, functions.SendMessage)

	http.ListenAndServe(":"+port, mux)
}
