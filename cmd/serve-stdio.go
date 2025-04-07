package cmd

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	MCPStdin  io.Reader = os.Stdin
	MCPStdout io.Writer = os.Stdout
	MCPStderr io.Writer = os.Stderr
)

var serveStdioCmd = &cobra.Command{
	Use:   "stdio",
	Short: "Serve the MCP over stdio (Standard Input/Output)",
	Long: `Start the MCP server using standard input and output streams.
This mode is typically used for local integrations and command-line tools that communicate via stdio.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		s, err := buildMCPServer()
		if err != nil {
			return err
		}

		zlog.Logger.Info().
			Str("transport", "stdio").
			Msg("ready")

		err = serveStdio(s, MCPStdin, MCPStdout, MCPStderr, WithZerolog())
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

// serveStdio creates and starts a StdioServer with the provided MCPServer and I/O streams.
// It sets up signal handling for graceful shutdown on SIGTERM and SIGINT.
// Returns an error if the server encounters any issues during operation.
func serveStdio(
	srv *server.MCPServer,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	opts ...server.StdioOption,
) error {
	s := server.NewStdioServer(srv)
	s.SetErrorLogger(log.New(stderr, "", log.LstdFlags))

	for _, opt := range opts {
		opt(s)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigChan
		cancel()
	}()

	return s.Listen(ctx, stdin, stdout)
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
