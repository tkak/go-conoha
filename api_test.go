package conoha

import (
	"fmt"
	"testing"
)

func TestNewClient(t *testing.T) {
	var token string
	var endpoint string
	var ar *AuthResponse

	ac, err := NewAuthClient("", "", "")

	if err != nil {
		fmt.Println(err)
	}

	ar, err = ac.Authenticate()

	if err != nil {
		fmt.Println(err)
	}

	token = GetToken(ar)
	endpoint = GetEndpoint("object-store", ar)

	c, err := NewClient(token, endpoint)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("client: %+v\n", c)
}
