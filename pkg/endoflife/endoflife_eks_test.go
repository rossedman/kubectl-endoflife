package endoflife

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAmazonEKS(t *testing.T) {
	t.Run("returns days until end", func(t *testing.T) {
		now := time.Now()
		eks := AmazonEKS{
			EOL: now.AddDate(0, 0, 8).Format("2006-01-02"),
		}
		days, err := eks.GetDaysUntilEnd()
		assert.Nil(t, err)
		assert.Equal(t, float64(7), days)
	})

	t.Run("within expiry range returns true", func(t *testing.T) {
		now := time.Now()
		eks := AmazonEKS{
			EOL: now.AddDate(0, 0, 25).Format("2006-01-02"),
		}
		thres, err := eks.InExpiryRange(30)
		assert.Nil(t, err)
		assert.True(t, thres)
	})

	t.Run("within expiry range returns false when more than range", func(t *testing.T) {
		now := time.Now()
		eks := AmazonEKS{
			EOL: now.AddDate(0, 0, 60).Format("2006-01-02"),
		}
		thres, err := eks.InExpiryRange(30)
		assert.Nil(t, err)
		assert.False(t, thres)
	})

	t.Run("if date is expired return true", func(t *testing.T) {
		now := time.Now()
		eks := AmazonEKS{
			EOL: now.AddDate(0, 0, -7).Format("2006-01-02"),
		}
		expired, err := eks.IsExpired()
		assert.Nil(t, err)
		assert.True(t, expired)
	})

	t.Run("if date is not expired return false", func(t *testing.T) {
		now := time.Now()
		eks := AmazonEKS{
			EOL: now.AddDate(0, 0, 7).Format("2006-01-02"),
		}
		expired, err := eks.IsExpired()
		assert.Nil(t, err)
		assert.False(t, expired)
	})
}

func TestGetEKS(t *testing.T) {
	t.Run("retrieves response", func(t *testing.T) {
		// set mock testing server
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			content, err := ioutil.ReadFile("testdata/eks-1.19.json")
			if err != nil {
				log.Fatal(err)
			}
			w.Write(content)
		}))
		defer svr.Close()

		// create client
		client := NewClient(svr.URL, &http.Client{})

		// create request
		eks, err := client.GetAmazonEKS("1.19")
		assert.Nil(t, err)
		assert.Equal(t, "2022-04-01", eks.EOL)
		assert.Equal(t, "2021-02-16", eks.Release)
	})

	t.Run("fails when cannot be found", func(t *testing.T) {
		// set mock testing server
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
		}))
		defer svr.Close()

		// create client
		client := NewClient(svr.URL, &http.Client{})

		// create request
		_, err := client.GetAmazonEKS("1.666")
		assert.Error(t, err)
	})
}
