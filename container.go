package conoha

import (
	"fmt"
)

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
