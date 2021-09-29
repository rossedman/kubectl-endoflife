package endoflife

import (
	"fmt"
	"math"
	"net/http"
	"time"
)

// Kubernetes represents the JSON data returned
// by the endoflife.date/api/kubernetes/<version>.json
// endpoint
type Kubernetes struct {
	EOL     string `json:"eol"`
	Release string `json:"release"`
	Latest  string `json:"latest"`
}

// GetDaysUntilEnd parses the EOL date and provides a count
// of days until that date is reached
func (k *Kubernetes) GetDaysUntilEnd() (float64, error) {
	today := time.Now()
	endDate, err := time.Parse("2006-01-02", k.EOL)
	if err != nil {
		return 0, err
	}
	return math.Round(endDate.Sub(today).Hours() / 24), nil
}

// GetKubernetes returns the data for a single release of Kubernetes
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