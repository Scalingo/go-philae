package httpprobe

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/goware/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHttpProbe(t *testing.T) {
	t.Run("With a unavailable server", func(t *testing.T) {
		ctx := context.Background()
		p := NewHTTPProbe("http", "http://localhost:6666", HTTPOptions{testing: true})
		err := p.Check(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Unable to send request")
	})

	t.Run("With a server responding 5XX", func(t *testing.T) {
		cases := map[string]struct {
			status int
			err    string
		}{
			"With a normal node, it should refuse this node": {
				err: "Invalid return code",
			},
			"With a node accepting 500 responses, it should accept this node": {
				status: 500,
			},
			"With a node accepting 400 responses, it should not accept this node": {
				status: 400,
				err:    "Unexpected status code: 500 (expected: 400)",
			},
		}

		for title, c := range cases {
			t.Run(title, func(t *testing.T) {
				mockWorkingService := httpmock.NewMockHTTPServer("127.0.0.1:12345")
				defer mockWorkingService.Listener.Close()
				requestUrl, _ := url.Parse("http://127.0.0.1:12345/")
				mockWorkingService.AddResponses([]httpmock.MockResponse{
					{
						Request: http.Request{
							Method: "GET",
							URL:    requestUrl,
						},
						Response: httpmock.Response{
							StatusCode: 500,
							Body:       "Error",
						},
					},
				})

				ctx := context.Background()
				p := NewHTTPProbe("http", "http://127.0.0.1:12345/", HTTPOptions{
					testing:            true,
					ExpectedStatusCode: c.status,
				})
				err := p.Check(ctx)
				if c.err != "" {
					require.Error(t, err)
					assert.Contains(t, err.Error(), c.err)
					return
				}
				assert.NoError(t, err)
			})
		}
	})

	t.Run("With a server responding 2XX", func(t *testing.T) {
		cases := map[string]struct {
			checker HTTPChecker
			err     string
		}{
			"With a vanilla probe": {},
			"With a probe with custom successful checker": {
				checker: newTestChecker(nil),
			},
			"With a probe with custom failing checker": {
				checker: newTestChecker(errors.New("test error")),
				err:     "test error",
			},
		}

		for title, c := range cases {
			t.Run(title, func(t *testing.T) {
				ctx := context.Background()
				mockWorkingService := httpmock.NewMockHTTPServer("127.0.0.1:12345")
				defer mockWorkingService.Listener.Close()
				requestUrl, _ := url.Parse("http://127.0.0.1:12345/")
				mockWorkingService.AddResponses([]httpmock.MockResponse{
					{
						Request: http.Request{
							Method: "GET",
							URL:    requestUrl,
						},
						Response: httpmock.Response{
							StatusCode: 200,
							Body:       "OK",
						},
					},
				})

				p := NewHTTPProbe("http", "http://127.0.0.1:12345/", HTTPOptions{
					testing: true,
					Checker: c.checker,
				})

				err := p.Check(ctx)
				if c.err == "" {
					require.NoError(t, err)
					return
				}
				require.Error(t, err)
				assert.Contains(t, err.Error(), c.err)
			})
		}
	})
}
