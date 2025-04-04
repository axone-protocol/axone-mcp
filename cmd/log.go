package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	FlagLogLevel  = "log-level"
	FlagLogFormat = "log-format"
)

var (
	logLevel  string
	logFormat string
)

func InstallLogRunE(_ *cobra.Command, _ []string) error {
	var output io.Writer

	switch strings.ToLower(logFormat) {
	case "auto":
		if isatty.IsTerminal(os.Stdout.Fd()) {
			output = zerolog.ConsoleWriter{Out: os.Stderr}
		} else {
			output = os.Stderr
		}
	case "console":
		output = zerolog.ConsoleWriter{Out: os.Stderr}
	case "json":
		output = os.Stderr
	default:
		return fmt.Errorf("unknown log format: %s", logFormat)
	}

	l := zerolog.New(output).With().Timestamp().Logger()

	switch strings.ToLower(logLevel) {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		return fmt.Errorf("unknown log level: %s", logLevel)
	}

	log.Logger = l
	return nil
}

func init() {
	rootCmd.PersistentFlags().StringVar(&logLevel, FlagLogLevel, "info",
		`verbosity of logging ("trace", "debug", "info", "warn", "error")`)
	rootCmd.PersistentFlags().StringVar(&logFormat, FlagLogFormat, "auto",
		`format of logs ("auto", "console", "json")`)
}
