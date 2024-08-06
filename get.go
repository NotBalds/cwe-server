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
	"log"
	"math"
	"os"
	"strconv"
	"time"
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

	var usr = input.Body

	btssig, _ := base64.StdEncoding.DecodeString(usr.GetTimeSignature)
	keybl, _ := pem.Decode([]byte(register[usr.Uuid]))
	btskey := keybl.Bytes
	key, err := x509.ParsePKCS1PublicKey(btskey)
	if err != nil {
		return &GetOutput{Status: 498}, nil
	}

	sendtime, _ := strconv.ParseInt(usr.GetTime, 10, 64)

	if math.Abs(float64(time.Now().Unix()-sendtime)) > 10 {
		return &GetOutput{Status: 400}, nil
	}

	hash := sha256.New()
	hash.Write([]byte(usr.GetTime))
	checksig := rsa.VerifyPKCS1v15(key, crypto.SHA256, hash.Sum(nil), btssig)

	if checksig != nil {
		return &GetOutput{Status: 401}, nil
	}

	var msgs = db[usr.Uuid]
	if msgs == nil {
		msgs = []Message{}
	}

	db[usr.Uuid] = []Message{}
	newdb, err := json.Marshal(db)
	FatalIfErr(err, "Can't marshal new db")
	err = os.WriteFile("db.json", newdb, fs.ModePerm)
	FatalIfErr(err, "Can't write new db")

	return &GetOutput{msgs, 200}, nil
}
