package prober

import (
	"context"
	"fmt"
	"time"

	"github.com/Scalingo/go-utils/logger"
)

// Probe define a minimal set of methods that a probe should implement
type Probe interface {
	Name() string
	Check() error
}

// Prober entrypoint of the philae api. It will retain a set of probe and run
// checks when asked to
type Prober struct {
	timeout time.Duration
	probes  []Probe
}

// ProberOption is a function modifying some parameters of the Prober
type ProberOption func(p *Prober)

// WithTimeout is a ProberOption which defines a timeout the prober have to get
// executed into Recommandation: it should be higher than the timeout of the
// probes included in it, otherwise it would mask the real errors.
func WithTimeout(d time.Duration) ProberOption {
	return ProberOption(func(p *Prober) {
		p.timeout = d
	})
}

// Result is the data structure used to retain the data fetched from a single run of each probes
type Result struct {
	Healthy bool           `json:"healthy"`
	Probes  []*ProbeResult `json:"probes"`
}

// ProbeResult is the data structure used to retain the data fetched from a single probe
type ProbeResult struct {
	Name    string `json:"name"`
	Healthy bool   `json:"healthy"`
	Comment string `json:"comment"`
}

// NewProber is the default constructor of a Prober
func NewProber(opts ...ProberOption) *Prober {
	p := &Prober{
		timeout: 30 * time.Second,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (p *Prober) AddProbe(probe Probe) {
	p.probes = append(p.probes, probe)
}

// Check will run the check of each probes added and return the result in a Result struct
func (p *Prober) Check(ctx context.Context) *Result {
	probesResults := make([]*ProbeResult, len(p.probes))
	healthy := true
	resultChan := make(chan *ProbeResult, len(p.probes))

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	for _, probe := range p.probes {
		go p.CheckOneProbe(ctx, probe, resultChan)
	}

	for i := 0; i < len(p.probes); i++ {
		probeResult := <-resultChan
		if !probeResult.Healthy {
			healthy = false
		}
		probesResults[i] = probeResult
	}

	return &Result{
		Healthy: healthy,
		Probes:  probesResults,
	}
}

func (p *Prober) CheckOneProbe(ctx context.Context, probe Probe, res chan *ProbeResult) {
	log := logger.Get(ctx)
	probeRes := make(chan error)
	var err error

	go ProberWrapper(probe, probeRes)

	select {
	case e := <-probeRes:
		err = e
	case <-ctx.Done():
		err = fmt.Errorf("prober: %s", ctx.Err())
	}

	probe_healthy := true
	comment := ""
	if err != nil {
		comment = err.Error()
		probe_healthy = false
		log.Infof("[PHILAE] Probe %s failed, reason: %s\n", probe.Name(), err.Error())
	}
	probeResult := &ProbeResult{
		Name:    probe.Name(),
		Healthy: probe_healthy,
		Comment: comment,
	}

	res <- probeResult
}

func ProberWrapper(probe Probe, res chan error) {
	err := probe.Check()
	res <- err
}
