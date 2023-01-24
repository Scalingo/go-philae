package philaehandler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/go-philae/v5/prober"
	"github.com/Scalingo/go-philae/v5/sampleprobe"
)

func TestPhilaeHandler(t *testing.T) {
	TestWorkingEndpoint := func(resp *http.Response) {
		result := &prober.Result{}
		json.NewDecoder(resp.Body).Decode(&result)
		assert.True(t, result.Healthy)
		assert.Len(t, result.Probes, 1)
		assert.Equal(t, result.Probes[0].Name, "test")
		assert.True(t, result.Probes[0].Healthy)
	}

	t.Run("With an healthy prober", func(t *testing.T) {
		cases := map[string]struct {
			handler func(prober *prober.Prober) http.Handler
			status  int
			path    string
		}{
			"the ServeHTTP method should render a healthy node": {
				status: 200,
				path:   "/_health",
			},
			"the PhilaeRouter should route to the correct route": {
				status: 200,
				path:   "/_health",
				handler: func(prober *prober.Prober) http.Handler {
					return NewPhilaeRouter(http.NotFoundHandler(), prober)
				},
			},
			"the PhilarRouter should route all the other routes to the other router": {
				status: 301,
				path:   "/salut",
				handler: func(prober *prober.Prober) http.Handler {
					return NewPhilaeRouter(http.RedirectHandler("http://scalingo.com", 301), prober)
				},
			},
		}

		for title, c := range cases {
			t.Run(title, func(t *testing.T) {
				probe := prober.NewProber()
				probe.AddProbe(sampleprobe.NewSampleProbe("test", true))
				req := httptest.NewRequest("GET", "http://example.foo"+c.path, nil)
				w := httptest.NewRecorder()
				handler := NewHandler(probe)
				if c.handler != nil {
					handler = c.handler(probe)
				}
				handler.ServeHTTP(w, req)
				resp := w.Result()
				require.Equal(t, c.status, resp.StatusCode)
				if c.status == 200 {
					TestWorkingEndpoint(resp)
				}
			})
		}
	})
}
