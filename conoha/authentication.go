package conoha

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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
type AuthRequest struct {
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

type AuthResponse struct {
	Access struct {
		ServiceCatalogs []struct {
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

func authenticate(ar *AuthRequest, url string) (*AuthResponse, error) {
	json_data, err := json.Marshal(ar)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_data))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var a AuthResponse

	err = json.Unmarshal(body, &a)
	if err != nil {
		log.Fatal(err)
	}

	return &a, err
}

func getToken(ar *AuthResponse) string {
	return ar.Access.Token.ID
}

func getAccount(ar *AuthResponse) string {
	return ar.Access.Token.Tenant.ID
}

func getEndpoint(serviceType string, ac *AuthResponse) string {
	var endpoint string

	for _, element := range ac.Access.ServiceCatalogs {
		if element.Type == serviceType {
			endpoint = element.Endpoints[0].PublicURL
			break
		}
	}

	return endpoint
}
