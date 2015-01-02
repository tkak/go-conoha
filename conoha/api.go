package conoha

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type Client struct {
	Token   string
	URL     string
	Account string
	Http    *http.Client
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
		Token:   getToken(resp),
		URL:     getEndpoint("object-store", resp),
		Account: getAccount(resp),
		Http:    http.DefaultClient,
	}

	return &client, nil
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