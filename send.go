package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"time"
)

func sendMessage(w http.ResponseWriter, r *http.Request) {
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

	var send Send
	body, err := io.ReadAll(r.Body)
	FatalIfErr(err, "Can't read request body")
	err = json.Unmarshal(body, &send)
	FatalIfErr(err, "Can't unmarshal body")

	sendtime, _ := strconv.ParseInt(send.SendTime, 10, 64)

	if time.Now().Unix()-sendtime > 10 {
		w.WriteHeader(400)
		fmt.Println("warning! spam!", strconv.FormatInt(time.Now().Unix(), 10), "and", strconv.FormatInt(sendtime, 10), "are not the same!")
		return
	}

	btssig, _ := base64.StdEncoding.DecodeString(send.SendTimeSignature)
	btskey, _ := base64.StdEncoding.DecodeString(register[send.Sender])
	key, _ := x509.ParsePKCS1PublicKey(btskey)
	checksig := rsa.VerifyPKCS1v15(key, 0, []byte(send.SendTime), btssig)

	if checksig != nil {
		w.WriteHeader(401)
		return
	}

	db[send.Receiver] = append(db[send.Receiver], Message{send.Sender, send.Content})

	newdb, err := json.Marshal(db)
	FatalIfErr(err, "Can't marshal new db")
	err = os.WriteFile("db.json", newdb, fs.ModePerm)
	FatalIfErr(err, "Can't write new db")

	w.WriteHeader(200)
}
