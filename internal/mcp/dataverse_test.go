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
	requestId := mcp.NewRequestId("42")

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
					ID:      requestId,
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
					So(response, ShouldBeJSONRPCResponseSuccessWithText, `{"name":"dataverse-42","triplestore_address":"axone1xa8wemfrzq03tkwqxnv9lun7rceec7wuhh8x3qjgxkaaj5fl50zsmj8u0n"}`)
				},
			},
			{
				name: "get_dataverse_info tool - err1",
				message: mcp.JSONRPCRequest{
					JSONRPC: mcp.JSONRPC_VERSION,
					ID:      requestId,
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
					So(response, ShouldBeJSONRPCResponseErrorWithText, "err1")
				},
			},
			{
				name: "get_dataverse_info tool - missing arg",
				message: mcp.JSONRPCRequest{
					JSONRPC: mcp.JSONRPC_VERSION,
					ID:      requestId,
					Request: mcp.Request{
						Method: "tools/call",
					},
					Params: map[string]interface{}{
						"name": "get_dataverse_info",
					},
				},
				validate: func(response mcp.JSONRPCMessage) {
					So(response, ShouldBeJSONRPCErrorWithText, `required argument "dataverse" not found`)
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
