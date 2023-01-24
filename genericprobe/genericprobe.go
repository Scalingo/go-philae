package genericprobe

import "context"

type Pinger interface {
	Ping() error
}

type GenericProbe struct {
	name   string
	pinger Pinger
}

func New(name string, pinger Pinger) *GenericProbe {
	return &GenericProbe{
		name: name, pinger: pinger,
	}
}

func (p *GenericProbe) Name() string {
	return p.name
}

func (p *GenericProbe) Check(_ context.Context) error {
	return p.pinger.Ping()
}
