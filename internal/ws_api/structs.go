package ws_api

type WsRequest = struct {
	Method string `json:"action" doc:"register, send, or get string"`
	Data   string `json:"data" doc:"request body -> b64"`
}

type WsResponse = struct {
	Status int    `json:"status" doc:"integer http status"`
	Result string `json:"result" doc:"response body -> b64"`
}
