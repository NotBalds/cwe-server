package main

import (
	"encoding/json"
	"io"
	"io/fs"
	"net/http"
	"os"
)

func registerUser(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("db.json")
	FatalIfErr(err, "Can't read database")
	var db Database
	err = json.Unmarshal(data, &db)
	FatalIfErr(err, "Can't Unmarshal database")

	data, err = os.ReadFile("register.json")
	FatalIfErr(err, "Can't read register")
	var register Register
	err = json.Unmarshal(data, &register)
	FatalIfErr(err, "Can't Unmarshal register")

	var usr Registration
	body, err := io.ReadAll(r.Body)
	FatalIfErr(err, "Can't read request body")
	err = json.Unmarshal(body, &usr)
	FatalIfErr(err, "Can't unmarshal body")

	if register[usr.Uuid] == "" || register[usr.Uuid] == usr.PublicKey {
		register[usr.Uuid] = usr.PublicKey
		newregister, err := json.Marshal(register)
		FatalIfErr(err, "Can't marshal new register")
		err = os.WriteFile("register.json", newregister, fs.ModePerm)
		FatalIfErr(err, "Can't write new register")
		w.WriteHeader(200)
		return
	}
	w.WriteHeader(401)
	return
}
