package mcp

import (
	"fmt"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	. "github.com/smartystreets/goconvey/convey"
)

func makeRequest(args map[string]interface{}) mcp.CallToolRequest {
	var req mcp.CallToolRequest
	req.Params.Arguments = args
	return req
}

func TestRequiredParam_Int(t *testing.T) {
	Convey("Testing requiredParam", t, func() {
		tests := []struct {
			desc        string
			args        map[string]interface{}
			expectedErr string
			expectedVal string
		}{
			{
				desc: "success",
				args: map[string]interface{}{
					"foo": "bar",
					"arg": "42",
				},
				expectedVal: "42",
			},
			{
				desc: "missing parameter",
				args: map[string]interface{}{
					"foo": "bar",
				},
				expectedErr: "missing required parameter: arg",
			},
			{
				desc: "empty parameter",
				args: map[string]interface{}{
					"arg": "",
				},
				expectedErr: "parameter arg must not be empty",
			},
			{
				desc: "wrong type",
				args: map[string]interface{}{
					"arg": 42,
				},
				expectedErr: fmt.Sprintf("parameter %s is not of type %T", "arg", ""),
			},
		}

		for _, tt := range tests {
			Convey(tt.desc, func() {
				req := makeRequest(tt.args)
				got, err := requiredParam[string](req, "arg")
				if tt.expectedErr != "" {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, tt.expectedErr)
				} else {
					So(err, ShouldBeNil)
					So(got, ShouldEqual, tt.expectedVal)
				}
			})
		}
	})
}

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
