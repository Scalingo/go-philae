package elasticsearchprobe

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"

	"github.com/Scalingo/go-philae/v5/elasticsearchprobe/elasticsearchprobemock"

	"github.com/stretchr/testify/require"
)

func TestElasticsearchProbe_Check(t *testing.T) {
	runs := map[string]struct {
		pingerError   error
		expectedError string
	}{
		"A probe with a CA certificate should use it and return ok": {},
		"A probe without a required CA certificate and InsecureSkipVerify should return an error": {
			pingerError:   errors.New("failed to verify certificate: x509: certificate signed by unknown authority"),
			expectedError: "fail to get response",
		},
	}
	for name, run := range runs {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockPinger := elasticsearchprobemock.NewMockPinger(ctrl)
			mockPinger.EXPECT().Ping().Return(run.pingerError)

			probe := ElasticsearchProbe{
				name:     "testProbe",
				url:      "https://test.test:9200",
				insecure: false,
				caCert:   []byte("ThisShouldBeACACert"),
				pinger:   mockPinger,
			}

			err := probe.Check(context.Background())
			if run.expectedError != "" {
				require.Contains(t, err.Error(), run.expectedError)
				return
			}
			require.NoError(t, err)
		})
	}
}
