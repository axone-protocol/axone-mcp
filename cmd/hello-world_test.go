package cmd

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPrintHelloWorld(t *testing.T) {
	Convey("When ask for a welcome message", t, func() {
		result := getWelcomeMessage()
		Convey("It should return one", func() {
			So(result, ShouldEqual, "Hello world!")
		})
	})
}

func TestPrintHelloWorldCommand(t *testing.T) {
	Convey("When the hello-world command is executed", t, func() {
		helloWorld.SetArgs([]string{})
		err := helloWorld.Execute()
		Convey("It should not return an error", func() {
			So(err, ShouldBeNil)
		})
	})
}
