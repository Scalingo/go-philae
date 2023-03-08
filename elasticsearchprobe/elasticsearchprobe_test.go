package elasticsearchprobe

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/Scalingo/go-philae/v5/elasticsearchprobe/elasticsearchprobemock"
	"github.com/Scalingo/go-philae/v5/internal/tests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockSystemPoolWith(ctrl *gomock.Controller, probe *ElasticsearchProbe, ca []byte) {
	mock := elasticsearchprobemock.NewMockCertPoolGetter(ctrl)
	mock.EXPECT().SystemPool().Return(DefaultCertPoolGetter{}.FromCustomCA(ca), nil)

	probe.certPool = mock
}

func TestElasticsearchProbe_Check(t *testing.T) {
	ctx := context.Background()

	okHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("With an HTTP server", func(t *testing.T) {
		serv := httptest.NewServer(okHandler)
		defer serv.Close()
		assert.Contains(t, serv.URL, "http://")

		t.Run("It should succeed", func(t *testing.T) {
			probe := NewElasticsearchProbe("test", serv.URL)
			err := probe.Check(ctx)
			require.NoError(t, err)
		})
	})

	t.Run("With an HTTPs Server", func(t *testing.T) {
		ca, serv, err := tests.NewUnstartedServerWithTLSConfig(okHandler)
		require.NoError(t, err)
		serv.StartTLS()
		defer serv.Close()
		assert.Contains(t, serv.URL, "https://")

		t.Run("It should fail if it uses a custom certificate", func(t *testing.T) {
			probe := NewElasticsearchProbe("test", serv.URL)
			err := probe.Check(ctx)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "certificate signed by unknown authority")
		})

		t.Run("It should succeed if it uses a custom certificate with insecure verify", func(t *testing.T) {
			probe := NewElasticsearchProbe("test", serv.URL, WithInsecureSkipVerify())
			err := probe.Check(ctx)
			require.NoError(t, err)
		})

		t.Run("It should succeed if it uses a custom certificate and passes the CA", func(t *testing.T) {
			probe := NewElasticsearchProbe("test", serv.URL, WithCA(ca.CertificatePEM))
			err := probe.Check(ctx)
			require.NoError(t, err)
		})

		t.Run("Using the system CA", func(t *testing.T) {
			ctrl := gomock.NewController(t)

			probe := NewElasticsearchProbe("test", serv.URL)
			mockSystemPoolWith(ctrl, &probe, ca.CertificatePEM)
			err := probe.Check(ctx)
			require.NoError(t, err)
		})
	})
}
