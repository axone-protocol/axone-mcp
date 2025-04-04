package mcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	goctx "context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
	. "github.com/smartystreets/goconvey/convey"
)

func TestJSONRCPMessageHandling(t *testing.T) {
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

func TestOnRegisterSessionLog(t *testing.T) {
	Convey("Given a new MCP server", t, func() {
		s, err := NewServer()
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
