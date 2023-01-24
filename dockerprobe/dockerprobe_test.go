package dockerprobe

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/goware/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDockerProbe(t *testing.T) {
	ctx := context.Background()
	headers := http.Header{}
	headers.Add("user-agent", "go-dockerclient")
	headers.Add("connection", "close")
	mockWorkingService := httpmock.NewMockHTTPServer("127.0.0.1:11000")

	requestUrl, _ := url.Parse("http://127.0.0.1:11000/containers/json?")
	mockWorkingService.AddResponses([]httpmock.MockResponse{
		{
			Request: http.Request{
				Method: "GET",
				URL:    requestUrl,
				Header: headers,
			},
			Response: httpmock.Response{
				StatusCode: 200,
				Body:       "[]",
			},
		},
	})

	mockNotWorkingService := httpmock.NewMockHTTPServer("127.0.0.1:11001")

	requestNotWorkingUrl, _ := url.Parse("http://127.0.0.1:11001/containers/json?")
	mockNotWorkingService.AddResponses([]httpmock.MockResponse{
		{
			Request: http.Request{
				Method: "GET",
				URL:    requestNotWorkingUrl,
				Header: headers,
			},
			Response: httpmock.Response{
				StatusCode: 500,
				Body:       "it's not alive!",
			},
		},
	})

	t.Run("With a working docker container", func(t *testing.T) {
		probe := NewDockerProbe("docker", "127.0.0.1:11000")
		err := probe.Check(ctx)

		assert.NoError(t, err)
	})
	t.Run("With a not working docker container", func(t *testing.T) {
		probe := NewDockerProbe("docker", "127.0.0.1:11001")
		err := probe.Check(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Unable to contact docker: API error (500): it's not alive!")
	})
}
