package mcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	goctx "context"

	"github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/axone-protocol/axone-mcp/internal/mocks"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
)

func TestJSONRCPMessageHandling(t *testing.T) {
	readWriteFooTool := server.ServerTool{
		Tool: mcp.NewTool("read_write_foo",
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        "ReadWriteFoo",
				ReadOnlyHint: ref(false),
			})),
		Handler: func(ctx goctx.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			t.Fatalf("read_write_foo tool shouldn't be called")

			return &mcp.CallToolResult{}, nil
		},
	}

	Convey("Testing JSON-RPC message handling", t, func() {
		tests := []struct {
			name     string
			message  mcp.JSONRPCMessage
			fixture  func(srv *server.MCPServer, cc *mocks.MockClientConnInterface)
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
				name: "Tools list",
				message: mcp.JSONRPCRequest{
					JSONRPC: mcp.JSONRPC_VERSION,
					ID:      42,
					Request: mcp.Request{
						Method: "tools/list",
					},
				},
				fixture: func(srv *server.MCPServer, cc *mocks.MockClientConnInterface) {
					addTools(srv, ReadOnly, readWriteFooTool)
				},
				validate: func(response mcp.JSONRPCMessage) {
					So(response, ShouldNotBeNil)
					resp, ok := response.(mcp.JSONRPCResponse)
					So(ok, ShouldBeTrue)
					So(resp.ID, ShouldEqual, 42)
					So(resp.JSONRPC, ShouldEqual, mcp.JSONRPC_VERSION)
					ctr, ok := resp.Result.(mcp.ListToolsResult)
					So(ok, ShouldBeTrue)
					tools := lo.Map(ctr.Tools, func(t mcp.Tool, _ int) string {
						return t.Name
					})
					So(tools, ShouldContain, "get_resource_governance_code")
					So(tools, ShouldNotContain, "read_write_foo")
				},
			},
			{
				name: "ReadWriteFoo",
				message: mcp.JSONRPCRequest{
					JSONRPC: mcp.JSONRPC_VERSION,
					ID:      42,
					Request: mcp.Request{
						Method: "tools/call",
					},
					Params: map[string]interface{}{
						"name": "read_write_foo",
					},
				},
				fixture: func(srv *server.MCPServer, cc *mocks.MockClientConnInterface) {
					addTools(srv, ReadOnly, readWriteFooTool)
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
					So(content.Text, ShouldEqual, "The server is in read-only mode; tool read_write_foo cannot be invoked.")
					So(content.Type, ShouldEqual, "text")
				},
			},
		}

		for _, tt := range tests {
			Convey(fmt.Sprintf("Given a new server for %s", tt.name), func() {
				ctrl := gomock.NewController(t)
				Reset(ctrl.Finish)

				cc := mocks.NewMockClientConnInterface(ctrl)
				s, err := NewServer(cc, ReadOnly)
				So(err, ShouldBeNil)

				if tt.fixture != nil {
					tt.fixture(s, cc)
				}

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

		s, err := NewServer(mocks.NewMockClientConnInterface(ctrl), ReadWrite)
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

func expectClientConn(cc *mocks.MockClientConnInterface,
	address string,
	queryData string,
	respData string,
	err error,
) {
	cc.EXPECT().
		Invoke(gomock.Any(), "/cosmwasm.wasm.v1.Query/SmartContractState",
			&types.QuerySmartContractStateRequest{
				Address:   address,
				QueryData: []byte(queryData),
			},
			&types.QuerySmartContractStateResponse{},
			gomock.Any()).
		DoAndReturn(func(ctx goctx.Context, method string, req, reply any, opts ...grpc.CallOption) error {
			reply.(*types.QuerySmartContractStateResponse).Data = []byte(respData)
			return err
		}).Times(1)
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
