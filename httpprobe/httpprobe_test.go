package httpprobe

import (
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHttpProbe(t *testing.T) {
	Convey("With a unaivalable server", t, func() {
		p := NewHTTPProbe("http", "http://localhost:6666", HTTPOptions{})
		err := p.Check()
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldStartWith, "Unable to send request")
	})

	Convey("With a server responding 5XX", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", "http://scalingo.com/",
			httpmock.NewStringResponder(500, `Error`))

		p := NewHTTPProbe("http", "http://scalingo.com/", HTTPOptions{})

		err := p.Check()
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldStartWith, "Invalid return code")

	})
	Convey("With a server responding 2XX", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", "http://scalingo.com/",
			httpmock.NewStringResponder(200, `Error`))

		p := NewHTTPProbe("http", "http://scalingo.com/", HTTPOptions{})
		So(p.Check(), ShouldBeNil)
	})
}
