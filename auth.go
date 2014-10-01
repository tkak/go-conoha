package conoha

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type AuthClient struct {
	AuthRequest AuthRequest
	URL         string
}

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

func NewAuthClient(tenant, user, password string) (*AuthClient, error) {

	if tenant == "" {
		tenant = os.Getenv("CONOHA_TENANT")
	}

	if user == "" {
		user = os.Getenv("CONOHA_USER")
	}

	if password == "" {
		password = os.Getenv("CONOHA_PASSWORD")
	}

	ac := AuthClient{
		AuthRequest: AuthRequest{
			Auth: Auth{
				TenantName: tenant,
				PasswordCredentials: PasswordCredentials{
					Username: user,
					Password: password,
				},
			},
		},
		URL: "https://ident-r1nd1001.cnode.jp/v2.0/tokens",
	}

	return &ac, nil
}

func GetToken(ac *AuthResponse) string {
	return ac.Access.Token.ID
}

func GetEndpoint(serviceType string, ac *AuthResponse) string {

	var endpoint string

	for _, element := range ac.Access.ServiceCatalogs {
		if element.Type == serviceType {
			endpoint = element.Endpoints[0].PublicURL
			break
		}
	}

	return endpoint
}

func (ac *AuthClient) Authenticate() (*AuthResponse, error) {

	json_data, err := json.Marshal(ac.AuthRequest)

	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", ac.URL, bytes.NewBuffer(json_data))

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
	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var a AuthResponse

	err = json.Unmarshal(body, &a)

	if err != nil {
		log.Fatal(err)
	}
	//b, err := json.MarshalIndent(a, "", "  ")
	//fmt.Println("response Data:", string(b))

	return &a, err
}