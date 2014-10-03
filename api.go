package conoha

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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

type Client struct {
	Token string
	URL   string
	Http  *http.Client
}

type DoError struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

func NewClient(tenant, user, password string) (*Client, error) {

	URL := "https://ident-r1nd1001.cnode.jp/v2.0/tokens"

	if tenant == "" {
		tenant = os.Getenv("CONOHA_TENANT")
	}

	if user == "" {
		user = os.Getenv("CONOHA_USER")
	}

	if password == "" {
		password = os.Getenv("CONOHA_PASSWORD")
	}

	req := AuthRequest{
		Auth: Auth{
			TenantName: tenant,
			PasswordCredentials: PasswordCredentials{
				Username: user,
				Password: password,
			},
		},
	}

	resp, err := authenticate(&req, URL)

	if err != nil {
		log.Fatal(err)
	}

	client := Client{
		Token: resp.Access.Token.ID,
		URL:   getEndpoint("", resp),
		Http:  http.DefaultClient,
	}

	return &client, nil
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

func getToken(ar *AuthResponse) string {
	return ar.Access.Token.ID
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

// Creates a new request with the params
func (c *Client) NewRequest(params map[string]string, headerParams map[string]string, method string, endpoint string) (*http.Request, error) {

	p := url.Values{}
	u, err := url.Parse(c.URL + endpoint)

	if err != nil {
		return nil, fmt.Errorf("Error parsing base URL: %s", err)
	}

	for k, v := range params {
		p.Add(k, v)
	}

	u.RawQuery = p.Encode()

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %s", err)
	}

	req.Header.Add("X-Auth-Token", c.Token)
	for k, v := range headerParams {
		req.Header.Add(k, v)
	}

	return req, nil
}

// parseErr is used to take an error json resp
// and return a single string for use in error messages
func parseErr(resp *http.Response) error {
	errBody := new(DoError)

	err := decodeBody(resp, &errBody)

	// if there was an error decoding the body, just return that
	if err != nil {
		return fmt.Errorf("Error parsing error body for non-200 request: %s", err)
	}

	return fmt.Errorf("API Error: %s: %s", errBody.Id, errBody.Message)
}

// decodeBody is used to JSON decode a body
func decodeBody(resp *http.Response, out interface{}) error {
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, &out); err != nil {
		return err
	}

	return nil
}

// checkResp wraps http.Client.Do() and verifies that the
// request was successful. A non-200 request returns an error
// formatted to included any validation problems or otherwise
func checkResp(resp *http.Response, err error) (*http.Response, error) {
	// If the err is already there, there was an error higher
	// up the chain, so just return that
	if err != nil {
		return resp, err
	}
	fmt.Println(resp.StatusCode)

	switch i := resp.StatusCode; {
	case i == 200:
		return resp, nil
	case i == 201:
		return resp, nil
	case i == 202:
		return resp, nil
	case i == 204:
		return resp, nil
	case i == 422:
		return nil, parseErr(resp)
	case i == 400:
		return nil, parseErr(resp)
	default:
		return nil, fmt.Errorf("API Error: %s", resp.Status)
	}
}
