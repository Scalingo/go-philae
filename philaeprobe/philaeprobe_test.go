package philaeprobe

import (
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPhilaeProbe(t *testing.T) {
	Convey("With a unaivalable server", t, func() {
		p := NewPhilaeProbe("http", "http://localhost:6666")
		err := p.Check()
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldStartWith, "Unable to send request")
	})

	Convey("With a server responding 5XX", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", "http://scalingo.com/",
			httpmock.NewStringResponder(500, `Error`))

		p := NewPhilaeProbe("http", "http://scalingo.com/")
		err := p.Check()
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldStartWith, "Invalid return code")
	})

	Convey("With a server responding 2XX but an invalid json", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", "http://scalingo.com/",
			httpmock.NewStringResponder(200, `Error`))
		p := NewPhilaeProbe("http", "http://scalingo.com/")
		err := p.Check()
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldStartWith, "Invalid json")
	})

	Convey("With a server responding 2XX but an unhealthy probe", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", "http://scalingo.com/",
			httpmock.NewStringResponder(200, `{"healthy": false, "probes": []}`))
		p := NewPhilaeProbe("http", "http://scalingo.com/")
		err := p.Check()
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldStartWith, "Node not healthy")
	})

	Convey("With a server responding 2XX and an healthy probe", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", "http://scalingo.com/",
			httpmock.NewStringResponder(200, `{"healthy": true, "probes": []}`))
		p := NewPhilaeProbe("http", "http://scalingo.com/")
		So(p.Check(), ShouldBeNil)
	})
}
