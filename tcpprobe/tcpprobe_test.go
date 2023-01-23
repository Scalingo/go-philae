package tcpprobe

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type dialerWithTimeout struct{}

func (dialerWithTimeout) DialContext(ctx context.Context, proto string, endpoint string) (net.Conn, error) {
	<-ctx.Done()
	return nil, ctx.Err()
}

type resolverWithTimeout struct{}

func (resolverWithTimeout) LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error) {
	<-ctx.Done()
	return nil, ctx.Err()
}

func TestTCPProbe_Check(t *testing.T) {
	cases := map[string]struct {
		dialer   Dialer
		resolver Resolver
		timeout  time.Duration
		err      string
	}{
		"it should work by default": {},
		"it should return an error if the resolver timeout": {
			resolver: resolverWithTimeout{},
			timeout:  50 * time.Millisecond,
			err:      "DNS resolution failed",
		},
		"it should return an error if the dial timeout": {
			dialer:  dialerWithTimeout{},
			timeout: 50 * time.Millisecond,
			err:     "fail to open TCP connection",
		},
	}
	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			ctx := context.Background()
			listener, err := net.Listen("tcp", ":")
			require.NoError(t, err)
			defer listener.Close()
			tcplistener, ok := listener.(*net.TCPListener)
			require.True(t, ok)

			opts := TCPOptions{}
			if c.dialer != nil {
				opts.Dialer = c.dialer
			}
			if c.resolver != nil {
				opts.Resolver = c.resolver
			}
			if c.timeout != 0 {
				opts.Timeout = c.timeout
			}

			probe := NewTCPProbe("test-probe", tcplistener.Addr().String(), opts)
			err = probe.Check(ctx)
			if c.err == "" {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)
			assert.Contains(t, err.Error(), c.err)
		})
	}
}
