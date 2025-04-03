package mcp

import (
	"encoding/json"
	"fmt"
	"testing"

	goctx "context"

	"github.com/mark3labs/mcp-go/mcp"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewServer(t *testing.T) {
	Convey("Testing JSON-RPC message handling", t, func() {
		tests := []struct {
			name     string
			message  mcp.JSONRPCMessage
			validate func(response mcp.JSONRPCMessage)
		}{
			{
				name: "Ping",
				message: mcp.JSONRPCRequest{
					JSONRPC: mcp.JSONRPC_VERSION,
					ID:      42,
					Request: mcp.Request{
						Method: "ping",
					},
				},
				validate: func(response mcp.JSONRPCMessage) {
					So(response, ShouldNotBeNil)
					resp, ok := response.(mcp.JSONRPCResponse)
					So(ok, ShouldBeTrue)
					So(resp.ID, ShouldEqual, 42)
					So(resp.JSONRPC, ShouldEqual, mcp.JSONRPC_VERSION)
					_, ok = resp.Result.(mcp.EmptyResult)
					So(ok, ShouldBeTrue)
				},
			},
			{
				name: "hello_world tool (ko)",
				message: mcp.JSONRPCRequest{
					JSONRPC: mcp.JSONRPC_VERSION,
					ID:      42,
					Request: mcp.Request{
						Method: "tools/call",
					},
					Params: map[string]interface{}{
						"name": "hello_world",
						"arguments": map[string]interface{}{
							"name": 666,
						},
					},
				},
				validate: func(response mcp.JSONRPCMessage) {
					So(response, ShouldNotBeNil)
					resp, ok := response.(mcp.JSONRPCError)
					So(ok, ShouldBeTrue)
					So(resp.ID, ShouldEqual, 42)
					So(resp.JSONRPC, ShouldEqual, mcp.JSONRPC_VERSION)
					So(resp.Error.Message, ShouldEqual, "name must be a string")
				},
			},
			{
				name: "hello_world tool (ok)",
				message: mcp.JSONRPCRequest{
					JSONRPC: mcp.JSONRPC_VERSION,
					ID:      42,
					Request: mcp.Request{
						Method: "tools/call",
					},
					Params: map[string]interface{}{
						"name": "hello_world",
						"arguments": map[string]interface{}{
							"name": "test",
						},
					},
				},
				validate: func(response mcp.JSONRPCMessage) {
					So(response, ShouldNotBeNil)
					resp, ok := response.(mcp.JSONRPCResponse)
					So(ok, ShouldBeTrue)
					So(resp.ID, ShouldEqual, 42)
					So(resp.JSONRPC, ShouldEqual, mcp.JSONRPC_VERSION)
					ctr, ok := resp.Result.(mcp.CallToolResult)
					So(ok, ShouldBeTrue)
					So(ctr.IsError, ShouldBeFalse)
					So(ctr.Content, ShouldHaveLength, 1)
					content, ok := ctr.Content[0].(mcp.TextContent)
					So(ok, ShouldBeTrue)
					So(content.Text, ShouldEqual, "Hello, test!")
					So(content.Type, ShouldEqual, "text")
				},
			},
		}

		for _, tt := range tests {
			Convey(fmt.Sprintf("Given a new server for %s", tt.name), func() {
				s, err := NewServer()
				So(err, ShouldBeNil)

				messageBytes, err := json.Marshal(tt.message)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("When handling %s message", tt.name), func() {
					ctx := goctx.Background()
					got := s.HandleMessage(ctx, messageBytes)
					Convey("Then the response should be valid", func() {
						tt.validate(got)
					})
				})
			})
		}
	})
}
