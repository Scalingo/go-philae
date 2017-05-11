package prober

import (
	"errors"
	"testing"

	"github.com/Scalingo/go-philae/sampleprobe"
	. "github.com/smartystreets/goconvey/convey"
)

func TestProber(t *testing.T) {
	Convey("With healthy probes", t, func() {
		p := NewProber()
		p.AddProbe(sampleprobe.NewSampleProbe("a", true))
		p.AddProbe(sampleprobe.NewSampleProbe("b", true))

		res := p.Check()

		So(res.Healthy, ShouldBeTrue)
		So(len(res.Probes), ShouldEqual, 2)
		So(validateProbe(res.Probes, "a", true), ShouldBeNil)
		So(validateProbe(res.Probes, "b", true), ShouldBeNil)
	})

	Convey("With unhealthy probes", t, func() {
		p := NewProber()
		p.AddProbe(sampleprobe.NewSampleProbe("a", false))
		p.AddProbe(sampleprobe.NewSampleProbe("b", false))

		res := p.Check()

		So(res.Healthy, ShouldBeFalse)
		So(len(res.Probes), ShouldEqual, 2)
		So(validateProbe(res.Probes, "a", false), ShouldBeNil)
		So(validateProbe(res.Probes, "b", false), ShouldBeNil)
	})

	Convey("With a healthy probe and a unhealthy probe", t, func() {
		p := NewProber()
		p.AddProbe(sampleprobe.NewSampleProbe("a", true))
		p.AddProbe(sampleprobe.NewSampleProbe("b", false))

		res := p.Check()

		So(res.Healthy, ShouldBeFalse)
		So(len(res.Probes), ShouldEqual, 2)
		So(validateProbe(res.Probes, "a", true), ShouldBeNil)
		So(validateProbe(res.Probes, "b", false), ShouldBeNil)
	})
}

func validateProbe(probes []*ProbeResult, name string, healthy bool) error {
	for _, probe := range probes {
		if probe.Name == name {
			So(probe.Healthy, ShouldEqual, healthy)
			return nil
		}
	}

	return errors.New("Unable to find node " + name)
}
