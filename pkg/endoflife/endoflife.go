package endoflife

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	BaseURL = "https://endoflife.date/api"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient() *Client {
	return &Client{
		BaseURL: BaseURL,
		HTTPClient: &http.Client{
			Timeout: time.Second * 2,
		},
	}
}

func (c *Client) send(req *http.Request, v interface{}) error {
	// req.Header.Set("Content-Type", "application/json; charset=utf-8")
	// req.Header.Set("Accept", "application/json; charset=utf-8")

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

func (c *Client) ListAmazonEKS() ([]AmazonEKS, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/amazon-eks.json", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}

	res := []AmazonEKS{}
	if err := c.send(req, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) GetAmazonEKS(version string) (AmazonEKS, error) {
	res := AmazonEKS{}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/amazon-eks/%s.json", c.BaseURL, version), nil)
	if err != nil {
		return res, err
	}

	if err := c.send(req, &res); err != nil {
		return res, err
	}

	return res, nil
}

func (c *Client) ListKubernetes() ([]Kubernetes, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/kubernetes.json", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}

	res := []Kubernetes{}
	if err := c.send(req, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) GetKubernetes(version string) (Kubernetes, error) {
	res := Kubernetes{}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/kubernetes/%s.json", c.BaseURL, version), nil)
	if err != nil {
		return res, err
	}

	if err := c.send(req, &res); err != nil {
		return res, err
	}

	return res, nil
}
