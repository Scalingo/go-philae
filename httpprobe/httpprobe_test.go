package httpprobe

import (
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHttpProbe(t *testing.T) {
	Convey("With a unaivalable server", t, func() {
		p := NewHTTPProbe("http", "http://localhost:6666")
		err := p.Check()
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldStartWith, "Unable to send request")
	})

	Convey("With an invalid url", t, func() {
		p := NewHTTPProbe("http", "0xde:ad:be:ef")
		err := p.Check()
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldStartWith, "Unable to create request")
	})

	Convey("With a server responding 5XX", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", "http://scalingo.com/",
			httpmock.NewStringResponder(500, `Error`))

		p := NewHTTPProbe("http", "http://scalingo.com/")

		err := p.Check()
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldStartWith, "Invalid return code")

	})
	Convey("With a server responding 2XX", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", "http://scalingo.com/",
			httpmock.NewStringResponder(200, `Error`))

		p := NewHTTPProbe("http", "http://scalingo.com/")
		So(p.Check(), ShouldBeNil)

	})
}
