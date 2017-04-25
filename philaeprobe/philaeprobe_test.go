package philaeprobe

import (
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPhilaeProbe(t *testing.T) {
	Convey("With a unaivalable server", t, func() {
		p := NewPhilaeProbe("http", "http://localhost:6666")
		So(p.Check(), ShouldBeFalse)
	})

	Convey("With an invalid url", t, func() {
		p := NewPhilaeProbe("http", "0xde:ad:be:ef")
		So(p.Check(), ShouldBeFalse)
	})

	Convey("With a server responding 5XX", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", "http://scalingo.com/",
			httpmock.NewStringResponder(500, `Error`))

		p := NewPhilaeProbe("http", "http://scalingo.com/")
		So(p.Check(), ShouldBeFalse)
	})

	Convey("With a server responding 2XX but an invalid json", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", "http://scalingo.com/",
			httpmock.NewStringResponder(200, `Error`))
		p := NewPhilaeProbe("http", "http://scalingo.com/")
		So(p.Check(), ShouldBeFalse)
	})

	Convey("With a server responding 2XX but an unhealthy probe", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", "http://scalingo.com/",
			httpmock.NewStringResponder(200, `{"healthy": false, "probes": []}`))
		p := NewPhilaeProbe("http", "http://scalingo.com/")
		So(p.Check(), ShouldBeFalse)
	})

	Convey("With a server responding 2XX and an healthy probe", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", "http://scalingo.com/",
			httpmock.NewStringResponder(200, `{"healthy": true, "probes": []}`))
		p := NewPhilaeProbe("http", "http://scalingo.com/")
		So(p.Check(), ShouldBeTrue)
	})
}
