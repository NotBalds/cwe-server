package ws_api

import (
	"encoding/base64"
	"encoding/json"

	"github.com/NotBalds/cwe-server/internal/util"
	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
)

func response(resp WsResponse, ws *websocket.Conn) {
	res, err := json.Marshal(resp)
	util.FatalIfErr(err, "WS: Can't marshal WS response")
	ans := make([]byte, 1048576)
	base64.StdEncoding.Encode(ans, res)
	ws.WriteMessage(websocket.BinaryMessage, ans)
}

func StartWsApi(ws *websocket.Conn) {
	for {
		msgType, b64msg, err := ws.ReadMessage()
		if err != nil {
			log.Error("WS: Can't read message from WS", "err", err)
			ws.Close()
			continue
		}
		if msgType != websocket.BinaryMessage {
			log.Info("WS: Somebody tried to use non-binary message, closing")
			ws.WriteMessage(websocket.TextMessage, []byte("Go around, non-binary person!"))
			ws.Close()
			continue
		}

		msg := make([]byte, 1048576)
		_, err = base64.StdEncoding.Decode(msg, b64msg)
		if err != nil {
			log.Error("WS: Can't decode message from WS")
			ws.Close()
			continue
		}

		var req WsRequest
		err = json.Unmarshal(msg, &req)
		if err != nil {
			log.Error("WS: Can't unmarshal message from WS")
			ws.Close()
			continue
		}

		data, err := base64.StdEncoding.DecodeString(req.Data)
		if err != nil {
			log.Error("WS: Can't decode message{data} from WS")
			ws.Close()
			continue
		}
		sdata := string(data)

		switch req.Method {
		case "register":
			response(register(sdata), ws)
		case "send":
			response(send(sdata), ws)
		case "get":
			response(get(sdata), ws)
		}
	}
}
