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
	Auth Auth `json:"auth"`
}

type Auth struct {
	TenantName          string              `json:"tenantName"`
	PasswordCredentials PasswordCredentials `json:"passwordCredentials"`
}

type PasswordCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AccessContainer struct {
	Access struct {
		ServiceCatalog []struct {
			Endpoints []struct {
				AdminURL    string `json:"adminURL"`
				ID          string `json:"id"`
				InternalURL string `json:"internalURL"`
				PublicURL   string `json:"publicURL"`
				Region      string `json:"region"`
			} `json:"endpoints"`
			EndpointsLinks []interface{} `json:"endpoints_links"`
			Name           string        `json:"name"`
			Type           string        `json:"type"`
		} `json:"serviceCatalog"`
		Token struct {
			Expires  string `json:"expires"`
			ID       string `json:"id"`
			IssuedAt string `json:"issued_at"`
			Tenant   struct {
				Description string `json:"description"`
				Enabled     bool   `json:"enabled"`
				ID          string `json:"id"`
				Name        string `json:"name"`
			} `json:"tenant"`
		} `json:"token"`
		User struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"user"`
	} `json:"access"`
}

func getToken(ac AccessContainer) string {
	return ac.Access.Token.ID
}

func Authenticate() AccessContainer {

	url := "https://ident-r1nd1001.cnode.jp/v2.0/tokens"

	ac := AuthContainer{
		Auth: Auth{
			TenantName: os.Getenv("CONOHA_TENANT"),
			PasswordCredentials: PasswordCredentials{
				Username: os.Getenv("CONOHA_USER"),
				Password: os.Getenv("CONOHA_PASSWORD"),
			},
		},
	}

	json_data, err := json.Marshal(ac)
	if err != nil {
		log.Fatal(err)
	}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(json_data))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var a AccessContainer

	err = json.Unmarshal(body, &a)
	if err != nil {
		log.Fatal(err)
	}

	return a
}

func main() {
	var token string
	token = getToken(Authenticate())
	fmt.Println(token)
}
