package main

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"io/fs"
	"math"
	"os"
	"strconv"
	"time"
)

func sendMessage(ctx context.Context, input *SendInput) (*StatusOutput, error) {
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

	var send = input.Body

	sendtime, _ := strconv.ParseInt(send.SendTime, 10, 64)

	if math.Abs(float64(time.Now().Unix()-sendtime)) > 10 {
		return &StatusOutput{400}, nil
	}

	keybl, _ := pem.Decode([]byte(register[send.Sender]))
	btskey := keybl.Bytes
	key, err := x509.ParsePKCS1PublicKey(btskey)

	btssig, _ := base64.StdEncoding.DecodeString(send.SendTimeSignature)

	hash := sha256.New()
	hash.Write([]byte(send.SendTime))
	checksig := rsa.VerifyPKCS1v15(key, crypto.SHA256, hash.Sum(nil), btssig)

	if checksig != nil {
		return &StatusOutput{401}, nil
	}

	db[send.Receiver] = append(db[send.Receiver], Message{send.Sender, send.Content})

	newdb, err := json.Marshal(db)
	FatalIfErr(err, "Can't marshal new db")
	err = os.WriteFile("db.json", newdb, fs.ModePerm)
	FatalIfErr(err, "Can't write new db")

	return &StatusOutput{200}, nil
}
