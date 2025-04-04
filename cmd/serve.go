package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// serveCmd represents the base serve command.
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the MCP using a specific transport",
	Long: `Start the Axone MCP server using the chosen transport:
SSE for web clients, stdio for command-line and local integrations.`,
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		log.Logger.Info().Msg("starting server...")
	},
	PersistentPostRun: func(_ *cobra.Command, _ []string) {
		log.Logger.Info().Msg("server stopped")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
