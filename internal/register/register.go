package register

import (
	"context"
	"encoding/json"
	"io/fs"
	"os"

	"github.com/NotBalds/cwe-server/internal/structs"
	"github.com/NotBalds/cwe-server/internal/util"
)

func RegisterUser(_ context.Context, input *structs.RegisterInput) (*structs.StatusOutput, error) {
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

	if register[usr.Uuid] == "" || register[usr.Uuid] == usr.PublicKey {
		register[usr.Uuid] = usr.PublicKey
		newregister, err := json.Marshal(register)
		util.FatalIfErr(err, "Can't marshal new register")
		err = os.WriteFile("register.json", newregister, fs.ModePerm)
		util.FatalIfErr(err, "Can't write new register")
		return &structs.StatusOutput{Status: 200}, nil
	}
	return &structs.StatusOutput{Status: 401}, nil
}
