package main

import (
	"context"
	"encoding/json"
	"io/fs"
	"os"
)

func registerUser(ctx context.Context, input *RegisterInput) (*StatusOutput, error) {
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

	if register[usr.Uuid] == "" || register[usr.Uuid] == usr.PublicKey {
		register[usr.Uuid] = usr.PublicKey
		newregister, err := json.Marshal(register)
		FatalIfErr(err, "Can't marshal new register")
		err = os.WriteFile("register.json", newregister, fs.ModePerm)
		FatalIfErr(err, "Can't write new register")
		return &StatusOutput{200}, nil
	}
	return &StatusOutput{401}, nil
}
