package mcp

import (
	"context"

	"github.com/axone-protocol/axone-mcp/internal/version"
	"github.com/axone-protocol/axone-sdk/dataverse"

	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
)

const (
	ServerName = "Axone MCP Server"
)

func NewServer(dqc dataverse.QueryClient) (*server.MCPServer, error) {
	s := server.NewMCPServer(
		ServerName,
		version.Version,
		server.WithLogging(),
		server.WithToolCapabilities(false),
		WithHooksLogging(),
	)

	s.AddTool(getGovernanceCode(dqc))

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
