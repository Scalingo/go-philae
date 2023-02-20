package sampleprobe

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

// Used for tests only

type SampleProbe struct {
	name   string
	result bool
	time   time.Duration
	err    error
}

func NewSampleProbe(name string, result bool) SampleProbe {
	return SampleProbe{
		name:   name,
		result: result,
		time:   1 * time.Millisecond,
		err:    errors.New("error"),
	}
}

func NewTimedSampleProbe(name string, result bool, time time.Duration) SampleProbe {
	return SampleProbe{
		name:   name,
		result: result,
		time:   time,
		err:    errors.New("error"),
	}
}

func NewSampleProbeWithError(name string, result bool, err error) SampleProbe {
	return SampleProbe{
		name:   name,
		result: result,
		time:   1 * time.Millisecond,
		err:    err,
	}
}

func (s SampleProbe) Name() string {
	return s.name
}

func (s SampleProbe) Check(_ context.Context) error {
	time.Sleep(s.time)
	if s.result {
		return nil
	}
	return s.err
}
