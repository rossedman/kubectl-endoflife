package endoflife

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetKubernetes(t *testing.T) {
	t.Run("retrieves response", func(t *testing.T) {
		// set mock testing server
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			content, err := ioutil.ReadFile("testdata/kubernetes-1.19.json")
			if err != nil {
				log.Fatal(err)
			}
			w.Write(content)
		}))
		defer svr.Close()

		// create client
		client := NewClient(svr.URL, &http.Client{})

		// create request
		eks, err := client.GetKubernetes("1.19")
		assert.Nil(t, err)
		assert.Equal(t, "2021-10-28", eks.EOL)
		assert.Equal(t, "2020-08-27", eks.Release)
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
		_, err := client.GetKubernetes("1.666")
		assert.Error(t, err)
	})
}
