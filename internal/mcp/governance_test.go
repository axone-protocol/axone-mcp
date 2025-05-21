package mcp

import (
	"encoding/base64"
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

func TestGovernanceJSONRCPMessageHandling(t *testing.T) {
	requestId := mcp.NewRequestId("42")

	Convey("Testing governance JSON-RPC message handling", t, func() {
		const selectQuery = `{"select":{"query":{"limit":1,"prefixes":[{"namespace":"https://w3id.org/axone/ontology/v4/schema/credential/governance/text/","prefix":"gov"}],"select":[{"variable":"code"}],"where":{"bgp":{"patterns":[{"object":{"node":{"named_node":{"full":"did:key:zQ3shTd79aJSfrNpMVpUVX1xrG9gabc6fmYJS4gFuwUnjKK3F"}}},"predicate":{"named_node":{"full":"dataverse:credential:body#subject"}},"subject":{"variable":"credId"}},{"object":{"node":{"named_node":{"prefixed":"gov:GovernanceTextCredential"}}},"predicate":{"named_node":{"full":"dataverse:credential:body#type"}},"subject":{"variable":"credId"}},{"object":{"variable":"claim"},"predicate":{"named_node":{"full":"dataverse:credential:body#claim"}},"subject":{"variable":"credId"}},{"object":{"variable":"gov"},"predicate":{"named_node":{"prefixed":"gov:isGovernedBy"}},"subject":{"variable":"claim"}},{"object":{"variable":"code"},"predicate":{"named_node":{"prefixed":"gov:fromGovernance"}},"subject":{"variable":"gov"}}]}}}}}`
		tests := []struct {
			name     string
			message  mcp.JSONRPCMessage
			fixture  func(connInterface *mocks.MockClientConnInterface)
			validate func(response mcp.JSONRPCMessage)
		}{
			{
				name: "get_resource_governance_code tool",
				message: mcp.JSONRPCRequest{
					JSONRPC: mcp.JSONRPC_VERSION,
					ID:      requestId,
					Request: mcp.Request{
						Method: "tools/call",
					},
					Params: map[string]interface{}{
						"name": "get_resource_governance_code",
						"arguments": map[string]interface{}{
							"dataverse": "axone1xt4ahzz2x8hpkc0tk6ekte9x6crw4w6u0r67cyt3kz9syh24pd7scvlt2w",
							"resource":  "did:key:zQ3shTd79aJSfrNpMVpUVX1xrG9gabc6fmYJS4gFuwUnjKK3F",
						},
					},
				},
				fixture: func(cc *mocks.MockClientConnInterface) {
					expectClientConn(cc, "axone1xt4ahzz2x8hpkc0tk6ekte9x6crw4w6u0r67cyt3kz9syh24pd7scvlt2w",
						`{"dataverse":{}}`,
						`{"triplestore_address":"axone1xa8wemfrzq03tkwqxnv9lun7rceec7wuhh8x3qjgxkaaj5fl50zsmj8u0n"}`,
						nil)

					expectClientConn(cc, "axone1xa8wemfrzq03tkwqxnv9lun7rceec7wuhh8x3qjgxkaaj5fl50zsmj8u0n",
						selectQuery,
						`{"head":{"vars":["code"]},"results":{"bindings":[{"code":{"type":"uri","value":{"full":"contract:law-stone:axone10tk8kmhhx49jahdyuxnn8d9luc9kxgc5m406k02s0y0ph59rdh7qstpynz"}}}]}}`,
						nil)

					expectClientConn(cc, "axone10tk8kmhhx49jahdyuxnn8d9luc9kxgc5m406k02s0y0ph59rdh7qstpynz",
						`{"program_code":{}}`,
						fmt.Sprintf(`"%s"`, base64.StdEncoding.EncodeToString([]byte(`hello(world).`))),
						nil)
				},
				validate: func(response mcp.JSONRPCMessage) {
					So(response, ShouldBeJSONRPCResponseSuccessWithText, "hello(world).")
				},
			},
			{
				name: "get_resource_governance_code tool - err1",
				message: mcp.JSONRPCRequest{
					JSONRPC: mcp.JSONRPC_VERSION,
					ID:      requestId,
					Request: mcp.Request{
						Method: "tools/call",
					},
					Params: map[string]interface{}{
						"name": "get_resource_governance_code",
						"arguments": map[string]interface{}{
							"dataverse": "axone1xt4ahzz2x8hpkc0tk6ekte9x6crw4w6u0r67cyt3kz9syh24pd7scvlt2w",
							"resource":  "did:key:zQ3shTd79aJSfrNpMVpUVX1xrG9gabc6fmYJS4gFuwUnjKK3F",
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
				name: "get_resource_governance_code tool - err2",
				message: mcp.JSONRPCRequest{
					JSONRPC: mcp.JSONRPC_VERSION,
					ID:      requestId,
					Request: mcp.Request{
						Method: "tools/call",
					},
					Params: map[string]interface{}{
						"name": "get_resource_governance_code",
						"arguments": map[string]interface{}{
							"dataverse": "axone1xt4ahzz2x8hpkc0tk6ekte9x6crw4w6u0r67cyt3kz9syh24pd7scvlt2w",
							"resource":  "did:key:zQ3shTd79aJSfrNpMVpUVX1xrG9gabc6fmYJS4gFuwUnjKK3F",
						},
					},
				},
				fixture: func(cc *mocks.MockClientConnInterface) {
					expectClientConn(cc, "axone1xt4ahzz2x8hpkc0tk6ekte9x6crw4w6u0r67cyt3kz9syh24pd7scvlt2w",
						`{"dataverse":{}}`,
						`{"triplestore_address":"axone1xa8wemfrzq03tkwqxnv9lun7rceec7wuhh8x3qjgxkaaj5fl50zsmj8u0n"}`,
						nil)

					expectClientConn(cc, "axone1xa8wemfrzq03tkwqxnv9lun7rceec7wuhh8x3qjgxkaaj5fl50zsmj8u0n",
						selectQuery,
						``,
						errors.New("err2"))
				},
				validate: func(response mcp.JSONRPCMessage) {
					So(response, ShouldBeJSONRPCResponseErrorWithText, "err2")
				},
			},
			{
				name: "get_resource_governance_code tool - err3",
				message: mcp.JSONRPCRequest{
					JSONRPC: mcp.JSONRPC_VERSION,
					ID:      requestId,
					Request: mcp.Request{
						Method: "tools/call",
					},
					Params: map[string]interface{}{
						"name": "get_resource_governance_code",
						"arguments": map[string]interface{}{
							"dataverse": "axone1xt4ahzz2x8hpkc0tk6ekte9x6crw4w6u0r67cyt3kz9syh24pd7scvlt2w",
							"resource":  "did:key:zQ3shTd79aJSfrNpMVpUVX1xrG9gabc6fmYJS4gFuwUnjKK3F",
						},
					},
				},
				fixture: func(cc *mocks.MockClientConnInterface) {
					expectClientConn(cc, "axone1xt4ahzz2x8hpkc0tk6ekte9x6crw4w6u0r67cyt3kz9syh24pd7scvlt2w",
						`{"dataverse":{}}`,
						`{"triplestore_address":"axone1xa8wemfrzq03tkwqxnv9lun7rceec7wuhh8x3qjgxkaaj5fl50zsmj8u0n"}`,
						nil)

					expectClientConn(cc, "axone1xa8wemfrzq03tkwqxnv9lun7rceec7wuhh8x3qjgxkaaj5fl50zsmj8u0n",
						selectQuery,
						`{"head":{"vars":["code"]},"results":{"bindings":[{"code":{"type":"uri","value":{"full":"contract:law-stone:axone10tk8kmhhx49jahdyuxnn8d9luc9kxgc5m406k02s0y0ph59rdh7qstpynz"}}}]}}`,
						nil)

					expectClientConn(cc, "axone10tk8kmhhx49jahdyuxnn8d9luc9kxgc5m406k02s0y0ph59rdh7qstpynz",
						`{"program_code":{}}`,
						``,
						errors.New("err3"))
				},
				validate: func(response mcp.JSONRPCMessage) {
					So(response, ShouldBeJSONRPCResponseErrorWithText, "err3")
				},
			},
			{
				name: "get_resource_governance_code tool - err4",
				message: mcp.JSONRPCRequest{
					JSONRPC: mcp.JSONRPC_VERSION,
					ID:      requestId,
					Request: mcp.Request{
						Method: "tools/call",
					},
					Params: map[string]interface{}{
						"name": "get_resource_governance_code",
						"arguments": map[string]interface{}{
							"dataverse": "axone1xt4ahzz2x8hpkc0tk6ekte9x6crw4w6u0r67cyt3kz9syh24pd7scvlt2w",
							"resource":  "did:key:zQ3shTd79aJSfrNpMVpUVX1xrG9gabc6fmYJS4gFuwUnjKK3F",
						},
					},
				},
				fixture: func(cc *mocks.MockClientConnInterface) {
					expectClientConn(cc, "axone1xt4ahzz2x8hpkc0tk6ekte9x6crw4w6u0r67cyt3kz9syh24pd7scvlt2w",
						`{"dataverse":{}}`,
						`{"triplestore_address":"axone1xa8wemfrzq03tkwqxnv9lun7rceec7wuhh8x3qjgxkaaj5fl50zsmj8u0n"}`,
						nil)

					expectClientConn(cc, "axone1xa8wemfrzq03tkwqxnv9lun7rceec7wuhh8x3qjgxkaaj5fl50zsmj8u0n",
						selectQuery,
						`{"head":{"vars":["code"]},"results":{"bindings":[{"code":{"type":"uri","value":{"full":"contract:law-stone:axone10tk8kmhhx49jahdyuxnn8d9luc9kxgc5m406k02s0y0ph59rdh7qstpynz"}}}]}}`,
						nil)

					expectClientConn(cc, "axone10tk8kmhhx49jahdyuxnn8d9luc9kxgc5m406k02s0y0ph59rdh7qstpynz",
						`{"program_code":{}}`,
						`"!!not_base64!!"`,
						nil)
				},
				validate: func(response mcp.JSONRPCMessage) {
					So(response, ShouldBeJSONRPCResponseErrorWithText, "failed to decode base64 code '!!not_base64!!': illegal base64 data at input byte 0")
				},
			},
			{
				name: "get_resource_governance_code tool - err5",
				message: mcp.JSONRPCRequest{
					JSONRPC: mcp.JSONRPC_VERSION,
					ID:      requestId,
					Request: mcp.Request{
						Method: "tools/call",
					},
					Params: map[string]interface{}{
						"name": "get_resource_governance_code",
						"arguments": map[string]interface{}{
							"dataverse": "axone1xt4ahzz2x8hpkc0tk6ekte9x6crw4w6u0r67cyt3kz9syh24pd7scvlt2w",
							"resource":  "did:key:zQ3shTd79aJSfrNpMVpUVX1xrG9gabc6fmYJS4gFuwUnjKK3F",
						},
					},
				},
				fixture: func(cc *mocks.MockClientConnInterface) {
					expectClientConn(cc, "axone1xt4ahzz2x8hpkc0tk6ekte9x6crw4w6u0r67cyt3kz9syh24pd7scvlt2w",
						`{"dataverse":{}}`,
						`{"triplestore_address":""}`,
						nil)
				},
				validate: func(response mcp.JSONRPCMessage) {
					So(response, ShouldBeJSONRPCResponseErrorWithText, "no triplestore address found")
				},
			},
			{
				name: "get_resource_governance_code tool - missing arg (1)",
				message: mcp.JSONRPCRequest{
					JSONRPC: mcp.JSONRPC_VERSION,
					ID:      requestId,
					Request: mcp.Request{
						Method: "tools/call",
					},
					Params: map[string]interface{}{
						"name": "get_resource_governance_code",
						"arguments": map[string]interface{}{
							"dataverse": "axone1xt4ahzz2x8hpkc0tk6ekte9x6crw4w6u0r67cyt3kz9syh24pd7scvlt2w",
						},
					},
				},
				validate: func(response mcp.JSONRPCMessage) {
					So(response, ShouldBeJSONRPCErrorWithText, `required argument "resource" not found`)
				},
			},
			{
				name: "get_resource_governance_code tool - missing arg (2)",
				message: mcp.JSONRPCRequest{
					JSONRPC: mcp.JSONRPC_VERSION,
					ID:      requestId,
					Request: mcp.Request{
						Method: "tools/call",
					},
					Params: map[string]interface{}{
						"name": "get_resource_governance_code",
						"arguments": map[string]interface{}{
							"resource": "did:key:zQ3shTd79aJSfrNpMVpUVX1xrG9gabc6fmYJS4gFuwUnjKK3F",
						},
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
