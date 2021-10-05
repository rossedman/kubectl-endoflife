package endoflife

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
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
		eks, err := client.Get(AmazonEKSProduct, "1.19")
		assert.Nil(t, err)
		assert.Equal(t, "2022-04-01", eks.EOL)
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
		_, err := client.Get(AmazonEKSProduct, "1.666")
		assert.Error(t, err)
	})
}
