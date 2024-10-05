package get

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

	"github.com/NotBalds/cwe-server/internal/structs"
	"github.com/NotBalds/cwe-server/internal/util"
)

func GetMessages(_ context.Context, input *structs.GetInput) (*structs.GetOutput, error) {
	data, err := os.ReadFile("db.json")
	util.FatalIfErr(err, "Can't read database")
	var db structs.Database
	err = json.Unmarshal(data, &db)
	util.FatalIfErr(err, "Can't Unmarshal database")

	data, err = os.ReadFile("register.json")
	util.FatalIfErr(err, "Can't read register")
	var register structs.Register
	err = json.Unmarshal(data, &register)
	util.FatalIfErr(err, "Can't Unmarshal register")

	var usr = input.Body

	btssig, _ := base64.StdEncoding.DecodeString(usr.GetTimeSignature)
	keybl, _ := pem.Decode([]byte(register[usr.Uuid]))
	btskey := keybl.Bytes
	key, err := x509.ParsePKCS1PublicKey(btskey)
	if err != nil {
		return &structs.GetOutput{Status: 498}, nil
	}

	sendtime, _ := strconv.ParseInt(usr.GetTime, 10, 64)

	if math.Abs(float64(time.Now().Unix()-sendtime)) > 10 {
		return &structs.GetOutput{Status: 400}, nil
	}

	hash := sha256.New()
	hash.Write([]byte(usr.GetTime))
	checksig := rsa.VerifyPKCS1v15(key, crypto.SHA256, hash.Sum(nil), btssig)

	if checksig != nil {
		return &structs.GetOutput{Status: 401}, nil
	}

	var msgs = db[usr.Uuid]
	if msgs == nil {
		msgs = []structs.Message{}
	}

	db[usr.Uuid] = []structs.Message{}
	newdb, err := json.Marshal(db)
	util.FatalIfErr(err, "Can't marshal new db")
	err = os.WriteFile("db.json", newdb, fs.ModePerm)
	util.FatalIfErr(err, "Can't write new db")

	return &structs.GetOutput{Body: msgs, Status: 200}, nil
}
