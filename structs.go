package main

type Message struct {
	Sender  string `json:"sender" format:"uuid" doc:"Sender's UUID"`
	Content string `json:"content" doc:"Content, encoded with base64"`
}

type StatusOutput struct {
	Status int
}

type Database map[string][]Message

type GetInput struct {
	Body struct {
		Uuid             string `json:"uuid" format:"uuid" example:"045cd5a4-7d09-44fe-8140-51b61c7e9750"` // UUID
		GetTime          string `json:"gettime" doc:"Unixtime when request is sent"`
		GetTimeSignature string `json:"gettimesignature" doc:"RSA signature of gettime, in Base64, signed with publickey that was sent in /register"`
	}
}

type GetOutput struct {
	Body   []Message
	Status int
}

type Register map[string]string

type RegisterInput struct {
	Body struct {
		Uuid      string `json:"uuid" format:"uuid" example:"045cd5a4-7d09-44fe-8140-51b61c7e9750"`
		PublicKey string `json:"publickey" minLength:"360" maxLength:"360" doc:"PublicKey, base64-encoded"`
	}
}

type SendInput struct {
	Body struct {
		Receiver          string `json:"receiver" format:"uuid" doc:"UUID of receiver"`
		SendTime          string `json:"sendtime" doc:"Unixtime when request is sent"`
		SendTimeSignature string `json:"sendtimesignature" doc:"RSA signature of sendtime, in Base64, signed with publickey that was sent in /register"`
		Message
	}
}
