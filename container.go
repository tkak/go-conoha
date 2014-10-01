package conoha

import (
	"fmt"
)

type ContainerResponse struct {
	Out string
}

type CreateContainer struct {
	Account   string
	Container string
}

func (c *Client) CreateContainer(opts *CreateContainer) error {
	params := make(map[string]string)
	headerParams := make(map[string]string)

	params["account"] = opts.Account
	params["container"] = opts.Container

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
