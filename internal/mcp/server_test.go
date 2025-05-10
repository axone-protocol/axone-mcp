package mcp

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	goctx "context"

	"github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
)

func TestJSONRCPMessageHandling(t *testing.T) {
	Convey("Testing JSON-RPC message handling", t, func() {
		tests := []struct {
			name     string
			message  mcp.JSONRPCMessage
			fixture  func(*MockClientConnInterface)
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
				name: "get_resource_governance_code tool",
				message: mcp.JSONRPCRequest{
					JSONRPC: mcp.JSONRPC_VERSION,
					ID:      42,
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
				fixture: func(cc *MockClientConnInterface) {
					expectClientConn(cc, "axone1xt4ahzz2x8hpkc0tk6ekte9x6crw4w6u0r67cyt3kz9syh24pd7scvlt2w",
						`{"dataverse":{}}`,
						`{"triplestore_address":"axone1xa8wemfrzq03tkwqxnv9lun7rceec7wuhh8x3qjgxkaaj5fl50zsmj8u0n"}`,
						nil)

					expectClientConn(cc, "axone1xa8wemfrzq03tkwqxnv9lun7rceec7wuhh8x3qjgxkaaj5fl50zsmj8u0n",
						`{"select":{"query":{"limit":1,"prefixes":[{"namespace":"https://w3id.org/axone/ontology/v4/schema/credential/governance/text/","prefix":"gov"}],"select":[{"variable":"code"}],"where":{"bgp":{"patterns":[{"object":{"node":{"named_node":{"full":"did:key:zQ3shTd79aJSfrNpMVpUVX1xrG9gabc6fmYJS4gFuwUnjKK3F"}}},"predicate":{"named_node":{"full":"dataverse:credential:body#subject"}},"subject":{"variable":"credId"}},{"object":{"node":{"named_node":{"prefixed":"gov:GovernanceTextCredential"}}},"predicate":{"named_node":{"full":"dataverse:credential:body#type"}},"subject":{"variable":"credId"}},{"object":{"variable":"claim"},"predicate":{"named_node":{"full":"dataverse:credential:body#claim"}},"subject":{"variable":"credId"}},{"object":{"variable":"gov"},"predicate":{"named_node":{"prefixed":"gov:isGovernedBy"}},"subject":{"variable":"claim"}},{"object":{"variable":"code"},"predicate":{"named_node":{"prefixed":"gov:fromGovernance"}},"subject":{"variable":"gov"}}]}}}}}`,
						`{"head":{"vars":["code"]},"results":{"bindings":[{"code":{"type":"uri","value":{"full":"contract:law-stone:axone10tk8kmhhx49jahdyuxnn8d9luc9kxgc5m406k02s0y0ph59rdh7qstpynz"}}}]}}`,
						nil)

					expectClientConn(cc, "axone10tk8kmhhx49jahdyuxnn8d9luc9kxgc5m406k02s0y0ph59rdh7qstpynz",
						`{"program_code":{}}`,
						fmt.Sprintf(`"%s"`, base64.StdEncoding.EncodeToString([]byte(`hello(world).`))),
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
					So(content.Text, ShouldEqual, "hello(world).")
					So(content.Type, ShouldEqual, "text")
				},
			},
			{
				name: "get_resource_governance_code tool - err1",
				message: mcp.JSONRPCRequest{
					JSONRPC: mcp.JSONRPC_VERSION,
					ID:      42,
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
				fixture: func(cc *MockClientConnInterface) {
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
				name: "get_resource_governance_code tool - err2",
				message: mcp.JSONRPCRequest{
					JSONRPC: mcp.JSONRPC_VERSION,
					ID:      42,
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
				fixture: func(cc *MockClientConnInterface) {
					expectClientConn(cc, "axone1xt4ahzz2x8hpkc0tk6ekte9x6crw4w6u0r67cyt3kz9syh24pd7scvlt2w",
						`{"dataverse":{}}`,
						`{"triplestore_address":"axone1xa8wemfrzq03tkwqxnv9lun7rceec7wuhh8x3qjgxkaaj5fl50zsmj8u0n"}`,
						nil)

					expectClientConn(cc, "axone1xa8wemfrzq03tkwqxnv9lun7rceec7wuhh8x3qjgxkaaj5fl50zsmj8u0n",
						`{"select":{"query":{"limit":1,"prefixes":[{"namespace":"https://w3id.org/axone/ontology/v4/schema/credential/governance/text/","prefix":"gov"}],"select":[{"variable":"code"}],"where":{"bgp":{"patterns":[{"object":{"node":{"named_node":{"full":"did:key:zQ3shTd79aJSfrNpMVpUVX1xrG9gabc6fmYJS4gFuwUnjKK3F"}}},"predicate":{"named_node":{"full":"dataverse:credential:body#subject"}},"subject":{"variable":"credId"}},{"object":{"node":{"named_node":{"prefixed":"gov:GovernanceTextCredential"}}},"predicate":{"named_node":{"full":"dataverse:credential:body#type"}},"subject":{"variable":"credId"}},{"object":{"variable":"claim"},"predicate":{"named_node":{"full":"dataverse:credential:body#claim"}},"subject":{"variable":"credId"}},{"object":{"variable":"gov"},"predicate":{"named_node":{"prefixed":"gov:isGovernedBy"}},"subject":{"variable":"claim"}},{"object":{"variable":"code"},"predicate":{"named_node":{"prefixed":"gov:fromGovernance"}},"subject":{"variable":"gov"}}]}}}}}`,
						``,
						errors.New("err2"))
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
					So(content.Text, ShouldEqual, "err2")
					So(content.Type, ShouldEqual, "text")
				},
			},
			{
				name: "get_resource_governance_code tool - missing arg",
				message: mcp.JSONRPCRequest{
					JSONRPC: mcp.JSONRPC_VERSION,
					ID:      42,
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
					So(response, ShouldNotBeNil)
					resp, ok := response.(mcp.JSONRPCError)
					So(ok, ShouldBeTrue)
					So(resp.ID, ShouldEqual, 42)
					So(resp.JSONRPC, ShouldEqual, mcp.JSONRPC_VERSION)
					So(resp.Error.Message, ShouldEqual, "missing required parameter: resource")
				},
			},
		}

		for _, tt := range tests {
			Convey(fmt.Sprintf("Given a new server for %s", tt.name), func() {
				ctrl := gomock.NewController(t)
				Reset(ctrl.Finish)

				cc := NewMockClientConnInterface(ctrl)
				if tt.fixture != nil {
					tt.fixture(cc)
				}
				s, err := NewServer(cc)
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

func TestOnRegisterSessionLog(t *testing.T) {
	Convey("Given a new MCP server", t, func() {
		ctrl := gomock.NewController(t)
		Reset(ctrl.Finish)

		s, err := NewServer(NewMockClientConnInterface(ctrl))
		So(err, ShouldBeNil)

		Convey("When RegisterSession is called with a new session", func() {
			session := fakeSession{
				sessionID:           "1234",
				notificationChannel: make(chan mcp.JSONRPCNotification),
				initialized:         false,
			}
			defer close(session.notificationChannel)

			output, err := captureLogOutput(func() error {
				return s.RegisterSession(goctx.Background(), session)
			})
			So(err, ShouldBeNil)

			Convey("Then the session ID and creation message should appear in the logs", func() {
				So(output, ShouldContainSubstring, "1234")
				So(output, ShouldContainSubstring, "session created")
			})
		})
	})
}

func expectClientConn(cc *MockClientConnInterface,
	address string,
	queryData string,
	respData string,
	err error,
) {
	cc.EXPECT().
		Invoke(gomock.Any(), gomock.Any(),
			&types.QuerySmartContractStateRequest{
				Address:   address,
				QueryData: []byte(queryData),
			},
			&types.QuerySmartContractStateResponse{},
			gomock.Any()).
		DoAndReturn(func(ctx goctx.Context, method string, req, reply any, opts ...grpc.CallOption) error {
			reply.(*types.QuerySmartContractStateResponse).Data = []byte(respData)
			return err
		}).AnyTimes()
}

func captureLogOutput(f func() error) (string, error) {
	var logBuffer bytes.Buffer
	originalLogger := log.Logger
	log.Logger = log.Logger.Output(&logBuffer)
	defer func() { log.Logger = originalLogger }()

	if err := f(); err != nil {
		return "", err
	}
	return logBuffer.String(), nil
}

type fakeSession struct {
	sessionID           string
	notificationChannel chan mcp.JSONRPCNotification
	initialized         bool
}

func (f fakeSession) SessionID() string {
	return f.sessionID
}

func (f fakeSession) NotificationChannel() chan<- mcp.JSONRPCNotification {
	return f.notificationChannel
}

func (f fakeSession) Initialize() {
}

func (f fakeSession) Initialized() bool {
	return f.initialized
}

var _ server.ClientSession = fakeSession{}
