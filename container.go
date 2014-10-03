package conoha

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type ContainerResponse struct {
	Container Container `json:"container"`
}

type Container struct {
	Account       string
	ContainerName string
}

func (c *Client) CreateContainer(opts *Container) error {
	params := make(map[string]string)
	headerParams := make(map[string]string)

	params["account"] = opts.Account
	params["container"] = opts.ContainerName

	headerParams["Accept"] = "application/json"
	headerParams["Content-Length"] = "0"

	req, err := c.NewRequest(params, headerParams, "PUT", fmt.Sprintf("/%s", params["container"]))
	if err != nil {
		return err
	}

	_, err = checkResp(c.Http.Do(req))
	if err != nil {
		return fmt.Errorf("Error creating container: %s", err)
	}

	return nil
}

func (c *Client) ReadContainer(opts *Container) error {
	params := make(map[string]string)
	headerParams := make(map[string]string)

	params["account"] = opts.Account
	params["container"] = opts.ContainerName

	headerParams["Accept"] = "application/json"

	req, err := c.NewRequest(params, headerParams, "GET", fmt.Sprintf("/%s", params["container"]))

	if err != nil {
		return err
	}

	resp, err := checkResp(c.Http.Do(req))

	if err != nil {
		return fmt.Errorf("Error reading container: %s", err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var i interface{}

	err = json.Unmarshal(body, &i)

	if err != nil {
		log.Fatal(err)
	}

	b, err := json.MarshalIndent(i, "", "  ")

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))

	return nil
}

func (c *Client) DeleteContainer(opts *Container) error {
	params := make(map[string]string)
	headerParams := make(map[string]string)

	params["account"] = opts.Account
	params["container"] = opts.ContainerName

	req, err := c.NewRequest(params, headerParams, "DELETE", fmt.Sprintf("/%s", params["container"]))
	if err != nil {
		return err
	}

	_, err = checkResp(c.Http.Do(req))
	if err != nil {
		return fmt.Errorf("Error deleting container: %s", err)
	}

	return nil
}
