package server

import (
	"net/http"
	"sync"

	"github.com/charmbracelet/log"

	"github.com/NotBalds/cwe-server/internal/ws_api"
	"github.com/gorilla/websocket"
)

func handler(w http.ResponseWriter, r *http.Request) {
	var upg = websocket.Upgrader{ReadBufferSize: 4096, WriteBufferSize: 4096}
	ws, err := upg.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	ws_api.StartWsApi(ws)
	defer ws.Close()
}

func wsStart(port string, wg *sync.WaitGroup) {
	defer wg.Done()
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	http.ListenAndServe(":"+port, mux)
}
