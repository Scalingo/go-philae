package prober

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Scalingo/go-philae/v4/sampleprobe"
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

	t.Run("With a single probe that times out", func(t *testing.T) {
		p := NewProber(WithTimeout(200 * time.Millisecond))
		p.AddProbe(sampleprobe.NewTimedSampleProbe("test", true, 300*time.Millisecond))

		start := time.Now()
		res := p.CheckOneProbe(ctx, "test")
		duration := time.Since(start)

		assert.True(t, duration > 200*time.Millisecond)
		assert.False(t, res.Healthy)
		assert.Equal(t, "test", res.Name)
		assert.Equal(t, "error", res.Comment)
		assert.Equal(t, "probe check failed: prober: context deadline exceeded", res.Error.Error())
	})

	t.Run("With a single healthy probe", func(t *testing.T) {
		p := NewProber()
		p.AddProbe(sampleprobe.NewSampleProbe("test", true))

		res := p.CheckOneProbe(ctx, "test")

		assert.Equal(t, "test", res.Name)
		assert.True(t, res.Healthy)
		assert.NoError(t, res.Error)
	})
	t.Run("With a single not found probe", func(t *testing.T) {
		p := NewProber()

		res := p.CheckOneProbe(ctx, "test")

		assert.Equal(t, "", res.Name)
		assert.False(t, res.Healthy)

		// Checking that the error type is correct
		errNotFound := &NotFoundError{}
		assert.True(t, errors.As(res.Error, errNotFound))
		assert.Equal(t, "probe test is not present in prober", res.Error.Error())
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
