package dockerprobe

import (
	"context"
	"testing"

	"github.com/Scalingo/go-philae/v5/tests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDockerProbe(t *testing.T) {
	ctx := context.Background()

	goodServer := tests.HTTPTestServer(map[tests.Route]tests.Response{
		{Method: "GET", Path: "/containers/json"}: {Status: 200, Body: "[]"},
	})
	defer goodServer.Close()

	badServer := tests.HTTPTestServer(map[tests.Route]tests.Response{
		{Method: "GET", Path: "/containers/json"}: {Status: 500, Body: "it's not alive!"},
	})
	defer badServer.Close()

	t.Run("With a working docker container", func(t *testing.T) {
		probe := NewDockerProbe("docker", goodServer.URL)
		err := probe.Check(ctx)

		assert.NoError(t, err)
	})
	t.Run("With a not working docker container", func(t *testing.T) {
		probe := NewDockerProbe("docker", badServer.URL)
		err := probe.Check(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Unable to contact docker: API error (500): it's not alive!")
	})
}
