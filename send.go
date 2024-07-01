package main

import (
	"encoding/json"
	"io"
	"io/fs"
	"net/http"
	"os"
)

func sendMessage(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("db.json")
	FatalIfErr(err, "Can't read database")
	var db Database
	err = json.Unmarshal(data, &db)
	FatalIfErr(err, "Can't Unmarshal database")

	var send Send
	body, err := io.ReadAll(r.Body)
	FatalIfErr(err, "Can't read request body")
	err = json.Unmarshal(body, &send)
	FatalIfErr(err, "Can't unmarshal body")

	db[send.Receiver] = append(db[send.Receiver], Message{send.Sender, send.Content})

	newdb, err := json.Marshal(db)
	FatalIfErr(err, "Can't marshal new db")
	err = os.WriteFile("db.json", newdb, fs.ModePerm)
	FatalIfErr(err, "Can't write new db")

	w.WriteHeader(200)
}
