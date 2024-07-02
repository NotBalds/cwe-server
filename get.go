package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
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

	data, err = os.ReadFile("register.json")
	FatalIfErr(err, "Can't read register")
	var register Register
	err = json.Unmarshal(data, &register)
	FatalIfErr(err, "Can't Unmarshal register")

	var usr User
	body, err := io.ReadAll(r.Body)
	FatalIfErr(err, "Can't read request body")
	err = json.Unmarshal(body, &usr)
	FatalIfErr(err, "Can't unmarshal body")

	btssig, err := base64.StdEncoding.DecodeString(usr.UuidSignature)
	btskey, err := base64.StdEncoding.DecodeString(register[usr.Uuid])
	key, err := x509.ParsePKCS1PublicKey(btskey)
	checksig := rsa.VerifyPKCS1v15(key, 0, []byte(usr.Uuid), btssig)
	if checksig != nil {
		w.WriteHeader(401)
		return
	}

	res, err := json.Marshal(db[usr.Uuid])
	FatalIfErr(err, "Can't marshal a response")
	w.Write(res)

	db[usr.Uuid] = []Message{}
	newdb, err := json.Marshal(db)
	FatalIfErr(err, "Can't marshal new db")
	err = os.WriteFile("db.json", newdb, fs.ModePerm)
	FatalIfErr(err, "Can't write new db")
}
