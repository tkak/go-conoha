//package conoha
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
type AuthContainer struct {
	Auth *Auth `json:"auth"`
}

type Auth struct {
	TenantName          string               `json:"tenantName"`
	PasswordCredentials *PasswordCredentials `json:"passwordCredentials"`
}

type PasswordCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func getAuthToken() string {

	url := "https://ident-r1nd1001.cnode.jp/v2.0/tokens"

	ac := &AuthContainer{
		&Auth{
			TenantName: os.Getenv("CONOHA_TENANT"),
			PasswordCredentials: &PasswordCredentials{
				Username: os.Getenv("CONOHA_USER"),
				Password: os.Getenv("CONOHA_PASSWORD"),
			},
		},
	}

	json_data, err := json.Marshal(ac)
	if err != nil {
		log.Fatal(err)
	}

	req, _ := http.NewRequest("POST", url, bytes.NewBufferString(string(json_data)))
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
	fmt.Println(i.(map[string]interface{}))

	m := i.(map[string]interface{})
	fmt.Println(m["access"])

	mm := m["access"].(map[string]interface{})
	fmt.Println(mm["token"])

	mmm := mm["token"].(map[string]interface{})
	fmt.Println(mmm["id"])

	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
	return "s"
}

func main() {
	getAuthToken()
}
