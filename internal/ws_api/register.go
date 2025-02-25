package ws_api

import (
	"encoding/json"

	"github.com/NotBalds/cwe-server/internal/structs"
	"github.com/charmbracelet/log"
)

func register(data string) structs.StatusOutput {
	var req structs.RegisterInput
	err := json.Unmarshal([]byte(data), req)
	if err != nil {
		log.Error("Can't unmarshal data in WS message")
	}
}
