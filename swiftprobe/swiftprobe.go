package swiftprobe

import (
	"errors"

	"github.com/ncw/swift/v2"
)

type SwiftProbe struct {
	name        string
	connBuilder func() (*swift.Connection, error)
	initErr     error
}

type SwiftProbeOpt func(p *SwiftProbe)

type AuthV2 struct {
	AuthURL  string
	Region   string
	Tenant   string
	Username string
	Password string
}

func WithAuthV2(cfg AuthV2) SwiftProbeOpt {
	return func(p *SwiftProbe) {
		p.connBuilder = func() (*swift.Connection, error) {
			return &swift.Connection{
				UserName:    cfg.Username,
				ApiKey:      cfg.Password,
				AuthUrl:     cfg.AuthURL,
				Tenant:      cfg.Tenant,
				Region:      cfg.Region,
				AuthVersion: 2,
			}, nil
		}
	}
}

type AuthV3 struct {
	AuthURL          string
	Region           string
	TenantID         string
	Username         string
	Password         string
	UserDomainName   string
	TenantDomainName string
}

func WithAuthV3(cfg AuthV3) SwiftProbeOpt {
	return func(p *SwiftProbe) {
		p.connBuilder = func() (*swift.Connection, error) {
			return &swift.Connection{
				UserName:     cfg.Username,
				ApiKey:       cfg.Password,
				AuthUrl:      cfg.AuthURL,
				TenantId:     cfg.TenantID,
				Region:       cfg.Region,
				Domain:       cfg.UserDomainName,
				TenantDomain: cfg.TenantDomainName,
				AuthVersion:  3,
			}, nil
		}
	}
}

func WithAuthFromEnv() SwiftProbeOpt {
	return func(p *SwiftProbe) {
		p.connBuilder = func() (*swift.Connection, error) {
			c := new(swift.Connection)
			err := c.ApplyEnvironment()
			if err != nil {
				return nil, err
			}
			return c, nil
		}
	}
}

func NewSwiftProbe(name string, opts ...SwiftProbeOpt) SwiftProbe {
	p := SwiftProbe{name: name}
	for _, opt := range opts {
		opt(&p)
	}
	if p.connBuilder == nil {
		p.initErr = errors.New("no swift connection configured")
	}
	return p
}

func (p SwiftProbe) Name() string {
	return p.name
}

func (p SwiftProbe) Check() error {
	if p.initErr != nil {
		return p.initErr
	}

	conn, err := p.connBuilder()
	if err != nil {
		return err
	}

	err = conn.Authenticate()
	if err != nil {
		return err
	}

	return nil
}
