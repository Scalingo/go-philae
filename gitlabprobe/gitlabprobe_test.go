package gitlabprobe

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/go-philae/statusioprobe"
)

func TestGitlabProbe(t *testing.T) {
	t.Run("When GitLab respond healthy", func(t *testing.T) {
		response := statusioprobe.StatusIOResponse{Result: &statusioprobe.StatusIOResult{
			Overall: &statusioprobe.StatusIOOverallResult{StatusCode: 100},
		}}

		checker := GitlabChecker{}

		buffer := new(bytes.Buffer)

		err := json.NewEncoder(buffer).Encode(&response)
		assert.NoError(t, err)

		err = checker.Check(buffer)
		assert.NoError(t, err)
	})

	t.Run("When GitHub respond not healthy", func(t *testing.T) {
		response := statusioprobe.StatusIOResponse{Result: &statusioprobe.StatusIOResult{
			Overall: &statusioprobe.StatusIOOverallResult{StatusCode: 400},
		}}

		checker := GitlabChecker{}
		buffer := new(bytes.Buffer)
		err := json.NewEncoder(buffer).Encode(&response)
		assert.NoError(t, err)

		err = checker.Check(buffer)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "One or more services from GitLab are down!")
	})
}
