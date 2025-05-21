package mcp

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRef(t *testing.T) {
	Convey("Testing ref function", t, func() {
		Convey("with integer value", func() {
			val := 42
			ptr := ref(val)
			So(ptr, ShouldNotBeNil)
			So(*ptr, ShouldEqual, val)
		})

		Convey("with string value", func() {
			val := "test"
			ptr := ref(val)
			So(ptr, ShouldNotBeNil)
			So(*ptr, ShouldEqual, val)
		})

		Convey("with struct value", func() {
			type testStruct struct {
				Field string
			}
			val := testStruct{Field: "test"}
			ptr := ref(val)
			So(ptr, ShouldNotBeNil)
			So(ptr.Field, ShouldEqual, val.Field)
		})
	})
}
