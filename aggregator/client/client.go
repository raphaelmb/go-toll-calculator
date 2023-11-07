package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/raphaelmb/go-toll-calculator/types"
)

type Client struct {
	Endpoint string
}

func NewClient(endpoint string) *Client {
	return &Client{
		Endpoint: endpoint,
	}
}

func (c *Client) AggregateInvoice(dist types.Distance) error {
	b, err := json.Marshal(dist)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, c.Endpoint, bytes.NewReader(b))
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("the service responded with a non 200 status code: %d", resp.StatusCode)
	}
	return nil
}
