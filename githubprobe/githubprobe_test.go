package githubprobe

import (
	"bytes"
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGithubProbe(t *testing.T) {
	Convey("When GitHub respond healthy", t, func() {
		response := GithubStatusResponse{Status: GithubStatusResponseStatus{
			Indicator: "none",
		}}

		checker := GithubChecker{}

		buffer := new(bytes.Buffer)

		err := json.NewEncoder(buffer).Encode(&response)
		So(err, ShouldBeNil)

		err = checker.Check(buffer)
		So(err, ShouldBeNil)
	})

	Convey("When GitHub respond not healthy", t, func() {
		response := GithubStatusResponse{Status: GithubStatusResponseStatus{
			Indicator: "mahor",
		}}

		checker := GithubChecker{}
		buffer := new(bytes.Buffer)
		err := json.NewEncoder(buffer).Encode(&response)
		So(err, ShouldBeNil)

		err = checker.Check(buffer)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "Github is probably down")
	})
}
