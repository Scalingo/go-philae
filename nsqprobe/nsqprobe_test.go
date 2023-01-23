package nsqprobe

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/goware/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestNSQProbe(t *testing.T) {
	ctx := context.Background()
	mockWorkingService := httpmock.NewMockHTTPServer("127.0.0.1:10001")

	requestUrl, _ := url.Parse("http://127.0.0.1:10000/ping")
	mockWorkingService.AddResponses([]httpmock.MockResponse{
		{
			Request: http.Request{
				Method: "GET",
				URL:    requestUrl,
			},
			Response: httpmock.Response{
				StatusCode: 200,
				Body:       "it's alive!",
			},
		},
	})

	mockNotWorkingService := httpmock.NewMockHTTPServer("127.0.0.1:10000")

	requestNotWorkingUrl, _ := url.Parse("http://127.0.0.1:10000/ping")
	mockNotWorkingService.AddResponses([]httpmock.MockResponse{
		{
			Request: http.Request{
				Method: "GET",
				URL:    requestNotWorkingUrl,
			},
			Response: httpmock.Response{
				StatusCode: 500,
				Body:       "it's not alive!",
			},
		},
	})

	t.Run("With a working server", func(t *testing.T) {
		probe := NewNSQProbe("salut", "127.0.0.1", 10001, nil)
		assert.NoError(t, probe.Check(ctx))
	})

	t.Run("With a non-working server", func(t *testing.T) {
		probe := NewNSQProbe("salut", "127.0.0.1", 10000, nil)
		assert.Error(t, probe.Check(ctx))
	})

}
