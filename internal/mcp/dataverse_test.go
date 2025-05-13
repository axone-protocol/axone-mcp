package mcp

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	goctx "context"

	"github.com/axone-protocol/axone-mcp/internal/mocks"
	"github.com/mark3labs/mcp-go/mcp"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

func TestDataverseJSONRCPMessageHandling(t *testing.T) {
	Convey("Testing governance JSON-RPC message handling", t, func() {
		tests := []struct {
			name     string
			message  mcp.JSONRPCMessage
			fixture  func(connInterface *mocks.MockClientConnInterface)
			validate func(response mcp.JSONRPCMessage)
		}{
			{
				name: "get_dataverse_info tool",
				message: mcp.JSONRPCRequest{
					JSONRPC: mcp.JSONRPC_VERSION,
					ID:      42,
					Request: mcp.Request{
						Method: "tools/call",
					},
					Params: map[string]interface{}{
						"name": "get_dataverse_info",
						"arguments": map[string]interface{}{
							"dataverse": "axone1xt4ahzz2x8hpkc0tk6ekte9x6crw4w6u0r67cyt3kz9syh24pd7scvlt2w",
						},
					},
				},
				fixture: func(cc *mocks.MockClientConnInterface) {
					expectClientConn(cc, "axone1xt4ahzz2x8hpkc0tk6ekte9x6crw4w6u0r67cyt3kz9syh24pd7scvlt2w",
						`{"dataverse":{}}`,
						`{"name": "dataverse-42", "triplestore_address":"axone1xa8wemfrzq03tkwqxnv9lun7rceec7wuhh8x3qjgxkaaj5fl50zsmj8u0n"}`,
						nil)
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
					So(content.Text, ShouldEqual, `{"name":"dataverse-42","triplestore_address":"axone1xa8wemfrzq03tkwqxnv9lun7rceec7wuhh8x3qjgxkaaj5fl50zsmj8u0n"}`)
					So(content.Type, ShouldEqual, "text")
				},
			},
			{
				name: "get_dataverse_info tool - err1",
				message: mcp.JSONRPCRequest{
					JSONRPC: mcp.JSONRPC_VERSION,
					ID:      42,
					Request: mcp.Request{
						Method: "tools/call",
					},
					Params: map[string]interface{}{
						"name": "get_dataverse_info",
						"arguments": map[string]interface{}{
							"dataverse": "axone1xt4ahzz2x8hpkc0tk6ekte9x6crw4w6u0r67cyt3kz9syh24pd7scvlt2w",
						},
					},
				},
				fixture: func(cc *mocks.MockClientConnInterface) {
					expectClientConn(cc, "axone1xt4ahzz2x8hpkc0tk6ekte9x6crw4w6u0r67cyt3kz9syh24pd7scvlt2w",
						`{"dataverse":{}}`,
						"",
						errors.New("err1"))
				},
				validate: func(response mcp.JSONRPCMessage) {
					So(response, ShouldNotBeNil)
					resp, ok := response.(mcp.JSONRPCResponse)
					So(ok, ShouldBeTrue)
					So(resp.ID, ShouldEqual, 42)
					So(resp.JSONRPC, ShouldEqual, mcp.JSONRPC_VERSION)
					ctr, ok := resp.Result.(mcp.CallToolResult)
					So(ok, ShouldBeTrue)
					So(ctr.IsError, ShouldBeTrue)
					So(ctr.Content, ShouldHaveLength, 1)
					content, ok := ctr.Content[0].(mcp.TextContent)
					So(ok, ShouldBeTrue)
					So(content.Text, ShouldEqual, "err1")
					So(content.Type, ShouldEqual, "text")
				},
			},
			{
				name: "get_dataverse_info tool - missing arg",
				message: mcp.JSONRPCRequest{
					JSONRPC: mcp.JSONRPC_VERSION,
					ID:      42,
					Request: mcp.Request{
						Method: "tools/call",
					},
					Params: map[string]interface{}{
						"name": "get_dataverse_info",
					},
				},
				validate: func(response mcp.JSONRPCMessage) {
					So(response, ShouldNotBeNil)
					resp, ok := response.(mcp.JSONRPCError)
					So(ok, ShouldBeTrue)
					So(resp.ID, ShouldEqual, 42)
					So(resp.JSONRPC, ShouldEqual, mcp.JSONRPC_VERSION)
					So(resp.Error.Message, ShouldEqual, "missing required parameter: dataverse")
				},
			},
		}

		for _, tt := range tests {
			Convey(fmt.Sprintf("Given a new server for %s", tt.name), func() {
				ctrl := gomock.NewController(t)
				Reset(ctrl.Finish)

				cc := mocks.NewMockClientConnInterface(ctrl)
				if tt.fixture != nil {
					tt.fixture(cc)
				}
				s, err := NewServer(cc, ReadWrite)
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
