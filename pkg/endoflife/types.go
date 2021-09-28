package endoflife

import (
	"math"
	"time"
)

type AmazonEKS struct {
	EOL     string `json:"eol"`
	Release string `json:"release"`
	Latest  string `json:"latest"`
}

func (a *AmazonEKS) GetDaysUntilEnd() (float64, error) {
	today := time.Now()
	endDate, err := time.Parse("2006-01-02", a.EOL)
	if err != nil {
		return 0, err
	}
	return math.Round(endDate.Sub(today).Hours() / 24), nil
}

type Kubernetes struct {
	EOL     string `json:"eol"`
	Release string `json:"release"`
	Latest  string `json:"latest"`
}

func (k *Kubernetes) GetDaysUntilEnd() (float64, error) {
	today := time.Now()
	endDate, err := time.Parse("2006-01-02", k.EOL)
	if err != nil {
		return 0, err
	}
	return math.Round(endDate.Sub(today).Hours() / 24), nil
}
