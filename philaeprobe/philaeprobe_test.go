package philaeprobe

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Scalingo/go-philae/prober"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPhilaeProbe(t *testing.T) {
	t.Run("With a unaivalable server", func(t *testing.T) {
		p := NewPhilaeProbe("http", "http://localhost:6666", 0, 0)
		err := p.Check()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Unable to send request")
	})

	t.Run("With a server responding 5XX", func(t *testing.T) {
		ts := launchTestServer(500, "Error")
		defer ts.Close()

		p := NewPhilaeProbe("http", ts.URL, 0, 0)
		err := p.Check()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid return code")
	})

	t.Run("With a server responding 2XX but an invalid json", func(t *testing.T) {
		ts := launchTestServer(200, "Salut salut")
		defer ts.Close()

		p := NewPhilaeProbe("http", ts.URL, 0, 0)
		err := p.Check()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid json")
	})

	t.Run("With a server responding 2XX but an unhealthy probe", func(t *testing.T) {
		result := &prober.Result{
			Healthy: false,
			Probes: []*prober.ProbeResult{
				&prober.ProbeResult{
					Name:    "node-1",
					Healthy: false,
					Comment: "pas bien",
				},
			},
		}
		ts := launchJSONTestServer(200, result)
		defer ts.Close()

		p := NewPhilaeProbe("http", ts.URL, 0, 0)
		err := p.Check()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "node-1 is down (pas bien),")
	})

	t.Run("With a server responding 2XX and an healthy probe", func(t *testing.T) {
		result := &prober.Result{
			Healthy: true,
			Probes:  []*prober.ProbeResult{},
		}
		ts := launchJSONTestServer(200, result)
		defer ts.Close()

		p := NewPhilaeProbe("http", ts.URL, 0, 0)
		assert.NoError(t, p.Check())
	})
}

func launchJSONTestServer(statusCode int, response interface{}) *httptest.Server {
	a, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	return launchTestServer(statusCode, string(a))
}

func launchTestServer(statusCode int, response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		fmt.Fprintln(w, response)
	}))
}
