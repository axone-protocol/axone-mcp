package cmd

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"axone-protocol/axone-mcp/internal/mcp"

	"github.com/justinas/alice"

	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	FlagBaseURL    = "base-url"
	FlagListenAddr = "listen-addr"
)

const (
	ReadHeaderTimeout = 5 * time.Second
)

var (
	baseURL    string
	listenAddr string
)

var serveSseCmd = &cobra.Command{
	Use:   "serve-sse",
	Short: "Start the MCP server (SSE)",
	RunE: func(_ *cobra.Command, _ []string) error {
		log.Logger.Info().Msg("starting server")

		s, err := mcp.NewServer()
		if err != nil {
			return err
		}

		sseServer := server.NewSSEServer(s)
		chain := loggerChain().Then(sseServer)
		httpSrv := &http.Server{
			Addr:              listenAddr,
			ReadHeaderTimeout: ReadHeaderTimeout,
			Handler:           chain,
		}

		go func() {
			log.Logger.Info().
				Str("transport", "sse").
				Str("base_url", baseURL).
				Str("addr", listenAddr).
				Str("message_path", sseServer.CompleteMessagePath()).
				Str("sse_path", sseServer.CompleteSsePath()).
				Msg("listening for connections")
			if err := httpSrv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				log.Fatal().Err(err).Msg("failed to start server")
			}
		}()

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		<-ctx.Done()
		log.Info().Msg("shutdown signal received")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return sseServer.Shutdown(shutdownCtx)
	},
}

func loggerChain() alice.Chain {
	return alice.New(hlog.NewHandler(log.Logger),
		hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			zerolog.Ctx(r.Context()).
				Debug().
				Str("client", r.RemoteAddr).
				Str("method", r.Method).
				Stringer("path", r.URL).
				Str("user_agent", r.Header.Get("User-Agent")).
				Int("status", status).
				Int("size", size).
				Dur("duration", duration).
				Msg("")
		}),
	)
}

func init() {
	serveSseCmd.PersistentFlags().StringVar(&baseURL, FlagBaseURL, "", "The server's base URL")
	serveSseCmd.PersistentFlags().StringVar(&listenAddr, FlagListenAddr, "127.0.0.1:8081", "The server's listen address")

	rootCmd.AddCommand(serveSseCmd)
}
