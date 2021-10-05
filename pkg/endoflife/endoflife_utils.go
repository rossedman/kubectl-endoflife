package endoflife

import (
	"math"
	"time"
)

// GetDaysUntilEnd parses the EOL date and provides a count
// of days until that date is reached
func GetDaysUntilEnd(date time.Time) (float64, error) {
	today := time.Now()
	return math.Round(date.Sub(today).Hours() / 24), nil
}

// InExpiryRange allows program to provide a timeframe that
// is considered "failed". Maybe the end of life deadline is 30
// days away but with a threshold of 60 days you would consider
// that failed
func InExpiryRange(date time.Time, days int) (bool, error) {
	d, err := GetDaysUntilEnd(date)
	if err != nil {
		return false, err
	}
	return d < float64(days), nil
}

// IsExpired will return true if days left are less
// than one
func IsExpired(date time.Time) (bool, error) {
	d, err := GetDaysUntilEnd(date)
	if err != nil {
		return false, err
	}
	return d < 1, nil
}
