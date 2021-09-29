package endoflife

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	BaseURL = "https://endoflife.date/api"
)

// Client represents an HTTP client for
// accessing the endoflife.date API
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient provides an implementation of the endoflife.date
// HTTP Client for accessing the API
func NewClient(url string, client *http.Client) *Client {
	return &Client{
		BaseURL:    url,
		HTTPClient: client,
	}
}

// send is a generic method used by others to retrieve
// data from the endoflife.date API
func (c *Client) send(req *http.Request, v interface{}) error {
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("unknown error sending request")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &v)
	if err != nil {
		return err
	}

	return nil
}
