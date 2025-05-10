package mcp

import (
	"context"

	"github.com/axone-protocol/axone-mcp/internal/version"
	"google.golang.org/grpc"

	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
)

const (
	ServerName = "Axone MCP Server"
)

func NewServer(conn grpc.ClientConnInterface) (*server.MCPServer, error) {
	s := server.NewMCPServer(
		ServerName,
		version.Version,
		server.WithLogging(),
		server.WithToolCapabilities(false),
		WithHooksLogging(),
	)

	s.AddTool(getGovernanceCode(conn))

	return s, nil
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
