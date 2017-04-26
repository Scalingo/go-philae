package prober

import (
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
		So(res.Probes[0].Name, ShouldEqual, "a")
		So(res.Probes[0].Healthy, ShouldBeTrue)
		So(res.Probes[1].Name, ShouldEqual, "b")
		So(res.Probes[1].Healthy, ShouldBeTrue)
	})

	Convey("With unhealthy probes", t, func() {
		p := NewProber()
		p.AddProbe(sampleprobe.NewSampleProbe("a", false))
		p.AddProbe(sampleprobe.NewSampleProbe("b", false))

		res := p.Check()

		So(res.Healthy, ShouldBeFalse)
		So(len(res.Probes), ShouldEqual, 2)
		So(res.Probes[0].Name, ShouldEqual, "a")
		So(res.Probes[0].Healthy, ShouldBeFalse)
		So(res.Probes[1].Name, ShouldEqual, "b")
		So(res.Probes[1].Healthy, ShouldBeFalse)
	})

	Convey("With a healthy probe and a unhealthy probe", t, func() {
		p := NewProber()
		p.AddProbe(sampleprobe.NewSampleProbe("a", true))
		p.AddProbe(sampleprobe.NewSampleProbe("b", false))

		res := p.Check()

		So(res.Healthy, ShouldBeFalse)
		So(len(res.Probes), ShouldEqual, 2)
		So(res.Probes[0].Name, ShouldEqual, "a")
		So(res.Probes[0].Healthy, ShouldBeTrue)
		So(res.Probes[1].Name, ShouldEqual, "b")
		So(res.Probes[1].Healthy, ShouldBeFalse)
	})
}
