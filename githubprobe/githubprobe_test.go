package githubprobe

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGithubProbe(t *testing.T) {
	t.Run("When GitHub respond healthy", func(t *testing.T) {
		response := GithubStatusResponse{Status: GithubStatusResponseStatus{
			Indicator: "none",
		}}

		checker := GithubChecker{}

		buffer := new(bytes.Buffer)

		err := json.NewEncoder(buffer).Encode(&response)
		assert.NoError(t, err)

		err = checker.Check(buffer)
		assert.NoError(t, err)
	})

	t.Run("When GitHub respond not healthy", func(t *testing.T) {
		response := GithubStatusResponse{Status: GithubStatusResponseStatus{
			Indicator: "major",
		}}

		checker := GithubChecker{}
		buffer := new(bytes.Buffer)
		err := json.NewEncoder(buffer).Encode(&response)
		assert.NoError(t, err)

		err = checker.Check(buffer)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "GitHub is probably down")
	})
}
