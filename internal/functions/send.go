package send

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

func SendMessage(_ context.Context, input *structs.SendInput) (*structs.StatusOutput, error) {
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

	var send = input.Body

	sendtime, _ := strconv.ParseInt(send.SendTime, 10, 64)

	if math.Abs(float64(time.Now().Unix()-sendtime)) > 10 {
		return &structs.StatusOutput{Status: 400}, nil
	}

	keybl, _ := pem.Decode([]byte(register[send.Message.Sender]))
	btskey := keybl.Bytes
	key, err := x509.ParsePKCS1PublicKey(btskey)

	btssig, _ := base64.StdEncoding.DecodeString(send.SendTimeSignature)

	hash := sha256.New()
	hash.Write([]byte(send.SendTime))
	checksig := rsa.VerifyPKCS1v15(key, crypto.SHA256, hash.Sum(nil), btssig)

	if checksig != nil {
		return &structs.StatusOutput{Status: 401}, nil
	}

	db[send.Receiver] = append(db[send.Receiver], send.Message)

	newdb, err := json.Marshal(db)
	util.FatalIfErr(err, "Can't marshal new db")
	err = os.WriteFile("db.json", newdb, fs.ModePerm)
	util.FatalIfErr(err, "Can't write new db")

	return &structs.StatusOutput{Status: 200}, nil
}
