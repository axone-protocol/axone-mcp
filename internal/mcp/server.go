package mcp

import (
	"context"
	"errors"
	"fmt"

	"axone-protocol/axone-mcp/internal/version"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	MCPServerName = "Axone MCP Server"
)

func NewServer() (*server.MCPServer, error) {
	s := server.NewMCPServer(
		MCPServerName,
		version.Version,
	)

	tool := mcp.NewTool("hello_world",
		mcp.WithDescription("Say hello to someone"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the person to greet"),
		),
	)

	s.AddTool(tool, helloHandler)

	return s, nil
}

func helloHandler(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, ok := request.Params.Arguments["name"].(string)
	if !ok {
		return nil, errors.New("name must be a string")
	}

	return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
}
