package server

import (
	"net/http"
	"sync"

	"github.com/charmbracelet/log"

	"github.com/gorilla/websocket"
)

func handler(w http.ResponseWriter, r *http.Request) {
	var upg = websocket.Upgrader{ReadBufferSize: 2048, WriteBufferSize: 2048}
	ws, err := upg.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
}

func wsStart(port string, wg *sync.WaitGroup) {
	defer wg.Done()
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	http.ListenAndServe(":"+port, mux)
}
