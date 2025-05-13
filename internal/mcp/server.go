package mcp

import (
	"context"
	"fmt"

	"github.com/axone-protocol/axone-mcp/internal/version"
	"github.com/mark3labs/mcp-go/mcp"
	"google.golang.org/grpc"

	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"

	"github.com/samber/lo"
)

const (
	ServerName = "Axone MCP Server"
)

type AccessMode bool

const (
	ReadOnly  AccessMode = true
	ReadWrite AccessMode = false
)

type serverToolFactory func(grpc.ClientConnInterface) server.ServerTool

var serverToolFactories = []serverToolFactory{
	getDataverse,
	getGovernanceCode,
}

// NewServer creates a new MCP server instance.
// It takes a gRPC connection to the Axone node and a read-only flag which  restricts the server to read-only operations.
func NewServer(cc grpc.ClientConnInterface, mode AccessMode) (*server.MCPServer, error) {
	s := server.NewMCPServer(
		ServerName,
		version.Version,
		server.WithLogging(),
		server.WithToolCapabilities(false),
		server.WithToolFilter(func(_ context.Context, tools []mcp.Tool) []mcp.Tool {
			return lo.Filter(tools, func(tool mcp.Tool, _ int) bool {
				return mode != ReadOnly || lo.FromPtr(tool.Annotations.ReadOnlyHint)
			})
		}),
		WithHooksLogging(),
	)

	addServerTools(s, mode, cc, serverToolFactories...)

	return s, nil
}

func addServerTools(s *server.MCPServer, mode AccessMode, cc grpc.ClientConnInterface, factories ...serverToolFactory) {
	tools := lo.Map(factories, createTool(cc))
	addTools(s, mode, tools...)
}

func addTools(s *server.MCPServer, mode AccessMode, tools ...server.ServerTool) {
	guardedTools := lo.Map(tools, wrapToolWithAccessGuard(mode))
	s.AddTools(guardedTools...)
}

func createTool(cc grpc.ClientConnInterface) func(factory serverToolFactory, _ int) server.ServerTool {
	return func(factory serverToolFactory, _ int) server.ServerTool {
		return factory(cc)
	}
}

func wrapToolWithAccessGuard(mode AccessMode) func(srvTool server.ServerTool, _ int) server.ServerTool {
	return func(srvTool server.ServerTool, _ int) server.ServerTool {
		next := srvTool.Handler
		srvTool.Handler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			if mode == ReadOnly && !lo.FromPtr(srvTool.Tool.Annotations.ReadOnlyHint) {
				return mcp.NewToolResultError(
					fmt.Sprintf("The server is in read-only mode; tool %s cannot be invoked.", srvTool.Tool.Name),
				), nil
			}
			return next(ctx, request)
		}

		return srvTool
	}
}

func WithHooksLogging() server.ServerOption {
	hooks := &server.Hooks{}

	hooks.AddOnRegisterSession(func(_ context.Context, session server.ClientSession) {
		log.Logger.Info().
			Str("session_id", session.SessionID()).
			Msg("session created")
	})

	return server.WithHooks(hooks)
}
