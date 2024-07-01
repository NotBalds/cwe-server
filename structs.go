package main

type Message struct {
	Sender  string // UUID
	Content string // Base64 of ciphered content
}

type Database map[string][]Message

type User struct {
	Uuid string // UUID
}

type Send struct {
	Receiver string
	Message
}
