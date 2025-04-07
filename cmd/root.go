package cmd

import (
	"os"
	"strings"

	"github.com/axone-protocol/axone-mcp/internal/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"resenje.org/casbab"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:               "axone-mcp",
	Short:             "Axone’s MCP server",
	Long:              "Gateway to the dataverse for AI-powered tools.",
	PersistentPreRunE: InstallLogRunE,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initConfig() {
	cobra.EnableTraverseRunHooks = true

	viper.SetEnvPrefix(casbab.ScreamingSnake(version.Name))
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

func init() {
	cobra.OnInitialize(initConfig)
}
