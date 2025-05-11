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

// NewServer creates a new MCP server instance.
// It takes a gRPC connection to the Axone node and a read-only flag which  restricts the server to read-only operations.
func NewServer(conn grpc.ClientConnInterface, mode AccessMode) (*server.MCPServer, error) {
	s := server.NewMCPServer(
		ServerName,
		version.Version,
		server.WithLogging(),
		server.WithToolCapabilities(false),
		server.WithToolFilter(func(_ context.Context, tools []mcp.Tool) []mcp.Tool {
			return lo.Filter(tools, func(tool mcp.Tool, _ int) bool {
				return mode != ReadOnly || tool.Annotations.ReadOnlyHint
			})
		}),
		WithHooksLogging(),
	)

	addTools(s, mode,
		getGovernanceCode(conn),
	)

	return s, nil
}

func addTools(s *server.MCPServer, mode AccessMode, tools ...server.ServerTool) {
	s.AddTools(
		wrapToolsWithAccessGuard(
			mode,
			tools,
		)...)
}

func wrapToolsWithAccessGuard(mode AccessMode, tools []server.ServerTool) []server.ServerTool {
	return lo.Map(tools, func(t server.ServerTool, _ int) server.ServerTool {
		return wrapToolWithAccessGuard(mode, t)
	})
}

func wrapToolWithAccessGuard(mode AccessMode, srvTool server.ServerTool) server.ServerTool {
	next := srvTool.Handler
	srvTool.Handler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if mode == ReadOnly && !srvTool.Tool.Annotations.ReadOnlyHint {
			return mcp.NewToolResultError(
				fmt.Sprintf("The server is in read-only mode; tool %s cannot be invoked.", srvTool.Tool.Name),
			), nil
		}
		return next(ctx, request)
	}

	return srvTool
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
