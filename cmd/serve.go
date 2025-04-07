package cmd

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/axone-protocol/axone-mcp/internal/mcp"
	"github.com/axone-protocol/axone-sdk/dataverse"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	grpccreds "google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	FlagNodeGrpc          = "node-grpc"
	FlagGrpcNoTLS         = "grpc-no-tls"
	FlagGrpcTLSSkipVerify = "grpc-tls-skip-verify"
	FlagDataverseAddr     = "dataverse-addr"
	FlagGrpcTimeout       = "grpc-timeout"
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

	serveCmd.PersistentFlags().String(FlagNodeGrpc, "127.0.0.1:9090",
		"Address <host>:<port> of the gRPC endpoint exposed by the axone node")
	_ = viper.BindPFlag(FlagNodeGrpc, serveCmd.PersistentFlags().Lookup(FlagNodeGrpc))

	serveCmd.PersistentFlags().Bool(FlagGrpcNoTLS, false,
		"Disable TLS when connecting to the gRPC endpoint")
	_ = viper.BindPFlag(FlagGrpcNoTLS, serveCmd.PersistentFlags().Lookup(FlagGrpcNoTLS))

	serveCmd.PersistentFlags().Bool(FlagGrpcTLSSkipVerify, false,
		"Use TLS but skip certificate verification (insecure)")
	_ = viper.BindPFlag(FlagGrpcTLSSkipVerify, serveCmd.PersistentFlags().Lookup(FlagGrpcTLSSkipVerify))

	serveCmd.PersistentFlags().String(FlagDataverseAddr, "",
		"Address of the dataverse CosmWasm contract")
	_ = serveCmd.MarkPersistentFlagRequired(FlagDataverseAddr)
	_ = viper.BindPFlag(FlagDataverseAddr, serveCmd.PersistentFlags().Lookup(FlagDataverseAddr))

	serveCmd.PersistentFlags().Duration(FlagGrpcTimeout, 5*time.Second,
		"Timeout for establishing the gRPC connection to the axone node (e.g. 5s, 2m)")
	_ = viper.BindPFlag(FlagGrpcTimeout, serveCmd.PersistentFlags().Lookup(FlagGrpcTimeout))

	serveCmd.MarkFlagsMutuallyExclusive(FlagGrpcNoTLS, FlagGrpcTLSSkipVerify)
}

func buildMCPServer() (*server.MCPServer, error) {
	dvq, err := buildDataverseClient()
	if err != nil {
		return nil, err
	}

	return mcp.NewServer(dvq)
}

func buildDataverseClient() (dataverse.QueryClient, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), viper.GetDuration(FlagGrpcTimeout))
	defer cancelFn()

	return dataverse.NewQueryClient(
		ctx, viper.GetString(FlagNodeGrpc), viper.GetString(FlagDataverseAddr),
		grpc.WithTransportCredentials(getTransportCredentials()))
}

func getTransportCredentials() grpccreds.TransportCredentials {
	switch {
	case viper.GetBool(FlagGrpcNoTLS):
		return insecure.NewCredentials()
	case viper.GetBool(FlagGrpcTLSSkipVerify):
		return grpccreds.NewTLS(&tls.Config{InsecureSkipVerify: true}) //nolint:gosec
	default:
		return grpccreds.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12})
	}
}
