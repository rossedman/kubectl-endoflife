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

// InExpiryRange allows program to provide a timeframe that
// is considered "failed". Maybe the end of life deadline is 30
// days away but with a threshold of 60 days you would consider
// that failed
func (a *AmazonEKS) InExpiryRange(days int) (bool, error) {
	d, err := a.GetDaysUntilEnd()
	if err != nil {
		return false, err
	}
	return d < float64(days), nil
}

// IsExpired will return true if days left are less
// than one
func (a *AmazonEKS) IsExpired() (bool, error) {
	d, err := a.GetDaysUntilEnd()
	if err != nil {
		return false, err
	}
	return d < 1, nil
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
