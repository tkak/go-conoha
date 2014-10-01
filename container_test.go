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

	var cc *Container
	cc = &Container{
		Account:       "2dd5e62509b04ec5ab39d46944a443e8",
		ContainerName: "test2",
	}

	//err = c.CreateContainer(cc)
	err = c.DeleteContainer(cc)

	if err != nil {
		fmt.Println(err)
	}

}
