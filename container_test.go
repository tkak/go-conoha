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
		Account:       "2dd5e62509b04ec5ab39d46944a443e8",
		ContainerName: "test",
	}

	//err = c.CreateContainer(cc)
	//err = c.DeleteContainer(cc)
	err = c.ReadContainer(cc)

	if err != nil {
		fmt.Println(err)
	}

}
