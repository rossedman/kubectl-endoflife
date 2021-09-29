package endoflife

import (
	"fmt"
	"math"
	"net/http"
	"time"
)

// AmazonEKS represents the JSON data returned
// by the endoflife.date/api/amazon-eks/<version>.json
// endpoint
type AmazonEKS struct {
	EOL     string `json:"eol"`
	Release string `json:"release"`
	Latest  string `json:"latest"`
}

// GetDaysUntilEnd parses the EOL date and provides a count
// of days until that date is reached
func (a *AmazonEKS) GetDaysUntilEnd() (float64, error) {
	today := time.Now()
	endDate, err := time.Parse("2006-01-02", a.EOL)
	if err != nil {
		return 0, err
	}
	return math.Round(endDate.Sub(today).Hours() / 24), nil
}

// GetAmazonEKS returns the data for a single release of EKS
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
