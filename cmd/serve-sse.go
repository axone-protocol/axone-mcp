package cmd

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/justinas/alice"
	"github.com/spf13/viper"

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
	Use:   "sse",
	Short: "Serve the MCP over SSE (server-sent events)",
	Long: `Start the MCP server using Server-Sent Events (SSE) to enable streaming over HTTP.
Typically used for browser-based or reactive clients.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		s, err := buildMCPServer(cmd.Context())
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
				Msg("ready")
			if err := httpSrv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				log.Fatal().Err(err).Msg("failed to start server")
			}
		}()

		ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		<-ctx.Done()
		log.Info().Msg("shutdown signal received")

		shutdownCtx, cancel := context.WithTimeout(cmd.Context(), 10*time.Second)
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
	serveSseCmd.PersistentFlags().StringVar(&baseURL, FlagBaseURL, "",
		"The server's base URL")
	_ = viper.BindPFlag(FlagBaseURL, serveSseCmd.PersistentFlags().Lookup(FlagBaseURL))

	serveSseCmd.PersistentFlags().StringVar(&listenAddr, FlagListenAddr, "127.0.0.1:8081",
		"The server's listen address")
	_ = viper.BindPFlag(FlagListenAddr, serveSseCmd.PersistentFlags().Lookup(FlagListenAddr))

	serveCmd.AddCommand(serveSseCmd)
}
