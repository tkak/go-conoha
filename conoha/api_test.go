package conoha

import (
	"fmt"
	"testing"
)

func TestNewClient(t *testing.T) {
	c, err := NewClient("", "", "")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("client: %+v\n", c)
}
