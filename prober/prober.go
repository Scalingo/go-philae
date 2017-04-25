package prober

// Probe define a minimal set of methods that a probe should implement
type Probe interface {
	Name() string
	Check() bool
}

// Prober entrypoint of the philae api. It will retain a set of probe and run
// checks when asked to
type Prober struct {
	probes []Probe
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
}

func NewProber() *Prober {
	return &Prober{}
}

func (p *Prober) AddProbe(probe Probe) {
	p.probes = append(p.probes, probe)
}

// Check will run the check of each probes added and return the result in a Result struct
func (p *Prober) Check() *Result {
	probesResult := make([]*ProbeResult, len(p.probes))
	healthy := true
	for i, probe := range p.probes {
		probesResult[i] = &ProbeResult{
			Name:    probe.Name(),
			Healthy: probe.Check(),
		}
		if !probesResult[i].Healthy {
			healthy = false
		}
	}

	return &Result{
		Healthy: healthy,
		Probes:  probesResult,
	}
}
