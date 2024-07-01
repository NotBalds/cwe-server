package main

import (
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
)

func FatalIfErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg)
	}
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("db.json")
	FatalIfErr(err, "Can't read database")
	var db Database
	err = json.Unmarshal(data, &db)
	FatalIfErr(err, "Can't Unmarshal database")

	var usr User
	body, err := io.ReadAll(r.Body)
	FatalIfErr(err, "Can't read request body")
	err = json.Unmarshal(body, &usr)
	FatalIfErr(err, "Can't unmarshal body")

	res, err := json.Marshal(db[usr.Uuid])
	FatalIfErr(err, "Can't marshal a response")
	w.Write(res)

	db[usr.Uuid] = []Message{}
	newdb, err := json.Marshal(db)
	FatalIfErr(err, "Can't marshal new db")
	err = os.WriteFile("db.json", newdb, fs.ModePerm)
	FatalIfErr(err, "Can't write new db")
}
