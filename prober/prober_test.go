package prober

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Scalingo/go-philae/v4/sampleprobe"
	"github.com/stretchr/testify/assert"
)

func TestProber(t *testing.T) {
	ctx := context.Background()
	t.Run("With healthy probes", func(t *testing.T) {
		p := NewProber()
		p.AddProbe(sampleprobe.NewSampleProbe("a", true))
		p.AddProbe(sampleprobe.NewSampleProbe("b", true))

		res := p.Check(ctx)

		assert.True(t, res.Healthy)
		assert.NoError(t, res.Error)
		assert.True(t, res.Error == nil)
		assert.Len(t, res.Probes, 2)
		assert.NoError(t, validateProbe(t, res.Probes, "a", true))
		assert.NoError(t, validateProbe(t, res.Probes, "b", true))
	})

	t.Run("With unhealthy probes", func(t *testing.T) {
		p := NewProber()
		p.AddProbe(sampleprobe.NewSampleProbe("a", false))
		// Postpone the failure of the second probe to ensure the order of errors
		p.AddProbe(sampleprobe.NewTimedSampleProbe("b", false, 10*time.Millisecond))

		res := p.Check(ctx)

		assert.Equal(t,
			"prober error: probe a: probe check failed: error, probe b: probe check failed: error",
			res.Error.Error(),
		)
		assert.False(t, res.Healthy)
		assert.Len(t, res.Probes, 2)
		assert.NoError(t, validateProbe(t, res.Probes, "a", false))
		assert.NoError(t, validateProbe(t, res.Probes, "b", false))
	})

	t.Run("With a healthy probe and a unhealthy probe", func(t *testing.T) {
		p := NewProber()
		p.AddProbe(sampleprobe.NewSampleProbe("a", true))
		p.AddProbe(sampleprobe.NewSampleProbe("b", false))

		res := p.Check(ctx)

		assert.False(t, res.Healthy)
		assert.Equal(t, res.Error.Error(), "prober error: probe b: probe check failed: error")
		assert.Len(t, res.Probes, 2)
		assert.NoError(t, validateProbe(t, res.Probes, "a", true))
		assert.NoError(t, validateProbe(t, res.Probes, "b", false))
	})

	t.Run("With a probe that times out", func(t *testing.T) {
		p := NewProber(WithTimeout(200 * time.Millisecond))
		p.AddProbe(sampleprobe.NewTimedSampleProbe("test1", true, 100*time.Millisecond))
		p.AddProbe(sampleprobe.NewTimedSampleProbe("test2", true, 300*time.Millisecond))
		start := time.Now()
		res := p.Check(ctx)
		duration := time.Now().Sub(start)

		assert.True(t, duration < 205*time.Millisecond)
		assert.True(t, duration > 200*time.Millisecond)

		assert.False(t, res.Healthy)
		assert.Equal(t, res.Error.Error(), "prober error: probe test2: probe check failed: prober: context deadline exceeded")
		assert.Len(t, res.Probes, 2)
		assert.True(t, strings.HasPrefix(res.Probes[0].Comment, "took "))
		assert.True(t, strings.HasSuffix(res.Probes[0].Comment, "ms"))
		assert.True(t, res.Probes[0].Duration < 105*time.Millisecond)
		assert.True(t, res.Probes[0].Duration > 100*time.Millisecond)
		assert.True(t, res.Probes[0].Healthy)
		assert.Equal(t, res.Probes[1].Comment, "error")
		assert.Equal(t, res.Probes[1].Error.Error(), "probe check failed: prober: context deadline exceeded")
		assert.True(t, res.Probes[1].Duration < 205*time.Millisecond)
		assert.True(t, res.Probes[1].Duration > 190*time.Millisecond)
		assert.False(t, res.Probes[1].Healthy)
	})
}

func validateProbe(t *testing.T, probes []*ProbeResult, name string, healthy bool) error {
	for _, probe := range probes {
		if probe.Name == name {
			assert.Equal(t, probe.Healthy, healthy)
			return nil
		}
	}

	return errors.New("Unable to find node " + name)
}
