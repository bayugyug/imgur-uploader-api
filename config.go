package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"golang.org/x/oauth2"
)

//UserCredential data holder on client id/secrets
type UserCredential struct {
	ID     string        `json:"client_id,omitempty"`
	Secret string        `json:"client_secret,omitempty"`
	Code   string        `json:"code,omitempty"`
	Bearer *oauth2.Token `json:"bearer,omitempty"`
}

func formatConfig(s, t string) (bool, *UserCredential) {
	var conf []byte
	var err error

	//parse
	if s != "" {
		conf, err = ioutil.ReadFile(s)
		if err != nil {
			log.Println("formatConfig", err)
			return false, nil
		}
	} else if t != "" {
		log.Println(t)
		conf = []byte(t)
	}

	var row UserCredential
	json.Unmarshal(conf, &row)

	//sanity check
	if row.ID == "" || row.Secret == "" {
		log.Println("formatConfig", "Oops, fail to format the config credentials")
		return false, nil
	}

	return true, &row
}
