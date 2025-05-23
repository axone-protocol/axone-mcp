package cmd

import (
	"encoding/json"
	"strings"

	"github.com/axone-protocol/axone-mcp/internal/version"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const (
	flagLong   = "long"
	flagOutput = "output"
)

// NewVersionCommand returns a CLI command to interactively print the application binary version information.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the application binary version information",
	RunE: func(cmd *cobra.Command, _ []string) error {
		verInfo := version.NewInfo()

		if long, _ := cmd.Flags().GetBool(flagLong); !long {
			cmd.Println(verInfo.Version)
			return nil
		}

		var (
			bz  []byte
			err error
		)

		output, _ := cmd.Flags().GetString(flagOutput)
		switch strings.ToLower(output) {
		case "json":
			bz, err = json.Marshal(verInfo)

		default:
			bz, err = yaml.Marshal(&verInfo)
		}

		if err != nil {
			return err
		}

		cmd.Println(string(bz))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().Bool(flagLong, false, "Print long version information")
	_ = viper.BindPFlag(flagLong, versionCmd.Flags().Lookup(flagLong))

	versionCmd.Flags().StringP(flagOutput, "o", "text", "Output format (text|json)")
	_ = viper.BindPFlag(flagOutput, versionCmd.Flags().Lookup(flagOutput))
}
