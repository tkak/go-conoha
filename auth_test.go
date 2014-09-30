package conoha

import (
	"fmt"
	"testing"
)

func TestAuth(t *testing.T) {
	var token string
	var endpoint string
	var ar *AuthResponse

	ar = Authenticate()
	token = GetToken(ar)
	endpoint = GetEndpoint("object-store", ar)

	fmt.Println(token)
	fmt.Println(endpoint)

	//if actual != expected {
	//	t.Errorf("got %v\nwant %v", actual, expected)
	//}
}
