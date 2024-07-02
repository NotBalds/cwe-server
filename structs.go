package main

type Message struct {
	Sender  string // UUID
	Content string // Base64 of ciphered content
}

type Database map[string][]Message

type User struct {
	Uuid             string // UUID
	GetTime          string
	GetTimeSignature string
}

type Register map[string]string

type Registration struct {
	Uuid      string
	PublicKey string
}

type Send struct {
	Receiver          string
	SendTime          string
	SendTimeSignature string
	Message
}
