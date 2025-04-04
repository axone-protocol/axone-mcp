package cmd

import (
	"context"
	"errors"
	"log"
	"strings"

	"axone-protocol/axone-mcp/internal/mcp"

	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var serveStdioCmd = &cobra.Command{
	Use:   "stdio",
	Short: "Serve the MCP over stdio (Standard Input/Output)",
	Long: `Start the MCP server using standard input and output streams.
This mode is typically used for local integrations and command-line tools that communicate via stdio.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		s, err := mcp.NewServer()
		if err != nil {
			return err
		}

		zlog.Logger.Info().
			Str("transport", "stdio").
			Msg("ready")

		err = server.ServeStdio(s, WithZerolog())
		if err != nil && !errors.Is(err, context.Canceled) {
			return err
		}

		zlog.Info().Msg("shutdown signal received")

		return nil
	},
}

// WithZerolog configures the server to use zerolog for logging.
func WithZerolog() server.StdioOption {
	errorWriter := &logWriter{
		logger: zlog.With().Logger(),
	}

	// Create a standard library logger that writes to our custom writer
	stdLogger := log.New(errorWriter, "", 0)

	return server.WithErrorLogger(stdLogger)
}

// logWriter implements io.Writer by writing to a zerolog.Logger.
type logWriter struct {
	logger zerolog.Logger
}

// Write implements io.Writer.
func (w *logWriter) Write(p []byte) (n int, err error) {
	// Remove trailing newlines and spaces for cleaner log output
	msg := strings.TrimSpace(string(p))
	if msg != "" {
		w.logger.Error().Msg(msg)
	}
	return len(p), nil
}

func init() {
	serveCmd.AddCommand(serveStdioCmd)
}
