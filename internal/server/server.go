package server

import (
	"sync"

	"github.com/charmbracelet/log"
)

func startHttpServer(wg *sync.WaitGroup) {
	var http_port = "1337"
	log.Info("Starting http server on port " + http_port + "...")
	go httpStart(http_port, wg)
}

func startWsServer(wg *sync.WaitGroup) {
	var ws_port = "2337"
	log.Info("Starting ws server on port " + ws_port + "...")
	go wsStart(ws_port, wg)
}

func Start() {
	var wg sync.WaitGroup
	wg.Add(2)
	startHttpServer(&wg)
	startWsServer(&wg)
	wg.Wait()
}
