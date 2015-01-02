package conoha

import (
	"fmt"
	"testing"
)

func TestContainer(t *testing.T) {
}

func TestCreateContainer(t *testing.T) {
	c, err := NewClient("", "", "")
	if err != nil {
		fmt.Println(err)
	}

	var cc *Container
	cc = &Container{
		Name: "test3",
	}

	//err = c.CreateContainer(cc)
	//err = c.ReadContainer(cc)
	err = c.DeleteContainer(cc)

	if err != nil {
		fmt.Println(err)
	}

}
