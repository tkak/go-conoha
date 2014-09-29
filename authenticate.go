package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Authentication data sample
// {
//   "auth": {
//     "tenantName": "1234567",
//     "passwordCredentials": {
//       "username": "1234567",
//       "password": "************"
//     }
//   }
// }
type passwordCredential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authData struct {
	TenantName         string              `json:"tenantName"`
	PasswordCredential *passwordCredential `json:"passwordCredentials"`
}

type Auth struct {
	Auth *authData `json:"auth"`
}

func main() {

	url := "https://ident-r1nd1001.cnode.jp/v2.0/tokens"

	a := &Auth{
		&authData{
			TenantName: os.Getenv("CONOHA_TENANT"),
			PasswordCredential: &passwordCredential{
				Username: os.Getenv("CONOHA_USER"),
				Password: os.Getenv("CONOHA_PASSWORD"),
			},
		},
	}

	b, _ := json.Marshal(a)
	fmt.Println(string(b))

	req, _ := http.NewRequest("POST", url, bytes.NewBufferString(string(b)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var i interface{}
	err = json.Unmarshal(body, &i)
	if err != nil {
		log.Fatal(err)
	}

	b, err = json.MarshalIndent(i, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}
