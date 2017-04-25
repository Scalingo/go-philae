package sampleprobe

// Used for tests only

type SampleProbe struct {
	name   string
	result bool
}

func NewSampleProbe(name string, result bool) SampleProbe {
	return SampleProbe{
		name:   name,
		result: result,
	}
}

func (s SampleProbe) Name() string {
	return s.name
}

func (s SampleProbe) Check() bool {
	return s.result
}
