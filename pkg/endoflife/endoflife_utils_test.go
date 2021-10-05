package endoflife

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUtils(t *testing.T) {
	t.Run("returns days until end", func(t *testing.T) {
		days, err := GetDaysUntilEnd(time.Now().AddDate(0, 0, 7))
		assert.Nil(t, err)
		assert.Equal(t, float64(7), days)
	})

	t.Run("within expiry range returns true", func(t *testing.T) {
		threshold, err := InExpiryRange(time.Now().AddDate(0, 0, 25), 30)
		assert.Nil(t, err)
		assert.True(t, threshold)
	})

	t.Run("within expiry range returns false when outside of range", func(t *testing.T) {
		threshold, err := InExpiryRange(time.Now().AddDate(0, 0, 60), 30)
		assert.Nil(t, err)
		assert.False(t, threshold)
	})

	t.Run("if date is expired return true", func(t *testing.T) {
		expired, err := IsExpired(time.Now().AddDate(0, 0, -7))
		assert.Nil(t, err)
		assert.True(t, expired)
	})

	t.Run("if date is not expired return false", func(t *testing.T) {
		expired, err := IsExpired(time.Now().AddDate(0, 0, 30))
		assert.Nil(t, err)
		assert.False(t, expired)
	})
}
