package conoha

import (
	"fmt"
	"testing"
)

func TestContainer(t *testing.T) {
}

func TestCreateContainer(t *testing.T) {
	var token string
	var endpoint string
	var ar *AuthResponse

	ar = Authenticate()
	token = GetToken(ar)
	endpoint = GetEndpoint("object-store", ar)

	c, err := NewClient(token, endpoint)
	if err != nil {
		fmt.Println(err)
	}

	var cc *CreateContainer
	cc = &CreateContainer{
		Account:   "2dd5e62509b04ec5ab39d46944a443e8",
		Container: "test2",
	}

	out, err := c.CreateContainer(cc)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(out)
}
