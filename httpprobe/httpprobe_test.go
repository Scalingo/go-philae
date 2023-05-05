package httpprobe

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/go-philae/v5/tests"
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
				var errorMessage string

				if c.err == "" {
					errorMessage = "Error"
				} else {
					errorMessage = c.err
				}

				srv := tests.HTTPTestServer(map[tests.Route]tests.Response{
					{Method: "GET", Path: "/"}: {Status: 500, Body: errorMessage},
				})

				ctx := context.Background()
				p := NewHTTPProbe("http", srv.URL, HTTPOptions{
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

				srv := tests.HTTPTestServer(map[tests.Route]tests.Response{
					{Method: "GET", Path: "/"}: {Status: 200, Body: "OK"},
				})

				p := NewHTTPProbe("http", srv.URL, HTTPOptions{
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
