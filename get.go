package main

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/fs"
	"log"
	"os"
)

func FatalIfErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg)
	}
}

func getMessages(ctx context.Context, input *GetInput) (*GetOutput, error) {
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

	fmt.Println("uuid:", input.Body.Uuid)
	fmt.Println("time:", input.Body.GetTime)
	fmt.Println("time signature:", input.Body.GetTimeSignature)

	var usr = input.Body

	btssig, _ := base64.StdEncoding.DecodeString(usr.GetTimeSignature)
	keybl, _ := pem.Decode([]byte(register[usr.Uuid]))
	btskey := keybl.Bytes
	key, err := x509.ParsePKCS1PublicKey(btskey)
	if err != nil {
		return &GetOutput{Status: 498}, nil
	}

	checksig := rsa.VerifyPKCS1v15(key, crypto.SHA256, []byte(usr.GetTime), btssig)

	if checksig != nil {
		log.Println("Error verifying signature, ", "err", checksig)
		return &GetOutput{Status: 401}, nil
	}

	var msgs = db[usr.Uuid]

	db[usr.Uuid] = []Message{}
	newdb, err := json.Marshal(db)
	FatalIfErr(err, "Can't marshal new db")
	err = os.WriteFile("db.json", newdb, fs.ModePerm)
	FatalIfErr(err, "Can't write new db")

	return &GetOutput{msgs, 200}, nil
}
