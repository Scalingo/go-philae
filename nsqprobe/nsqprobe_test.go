package nsqprobe

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Scalingo/go-philae/v5/tests"
)

func TestNSQProbe(t *testing.T) {
	ctx := context.Background()
	workingNsq := tests.HTTPTestServer(map[tests.Route]tests.Response{
		{Method: "GET", Path: "/ping"}: {Status: 200, Body: "it's alive!"},
	})

	notWorkingNsq := tests.HTTPTestServer(map[tests.Route]tests.Response{
		{Method: "GET", Path: "/ping"}: {Status: 500, Body: "it's not alive!"},
	})

	t.Run("With a working server", func(t *testing.T) {
		probe := NewNSQProbe("salut", "127.0.0.1", portFromURL(workingNsq.URL), nil)
		assert.NoError(t, probe.Check(ctx))
	})

	t.Run("With a non-working server", func(t *testing.T) {
		probe := NewNSQProbe("salut", "127.0.0.1", portFromURL(notWorkingNsq.URL), nil)
		assert.Error(t, probe.Check(ctx))
	})

}

func portFromURL(URL string) int {
	urlElements := strings.Split(URL, ":")
	port, _ := strconv.Atoi(urlElements[len(urlElements)-1])

	return port
}
