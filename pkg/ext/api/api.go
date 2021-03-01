package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type (
	// Client describes a client which makes requests to an external API
	Client struct {
		Hostname string
	}

	// HelloResponse describes a successful response to the hello endpoint
	HelloResponse struct {
		Message string `json:"message"`
	}
)

// GetHello calls the hello endpoint from the external API and returns the
// response if successful, or an error otherwise.
func (c Client) GetHello() (*HelloResponse, error) {
	res, err := http.Get(fmt.Sprintf("%s/hello", c.Hostname))
	if err != nil {
		return nil, fmt.Errorf("unable to get external API hello response, err = %s", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read hello response, err = %s", err)
	}

	decoded := &HelloResponse{}
	if err := json.Unmarshal(body, decoded); err != nil {
		return nil, fmt.Errorf("unable to decode hello response, err = %s", err)
	}

	return decoded, nil
}
