package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	dataverseschema "github.com/axone-protocol/axone-contract-schema/go/dataverse-schema/v6"
	"github.com/axone-protocol/axone-mcp/internal/axone/dataverse"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"google.golang.org/grpc"
)

func getDataverse(cc grpc.ClientConnInterface) server.ServerTool {
	const dataverseAddressParam = "dataverse"
	tool := mcp.NewTool("get_dataverse_info",
		mcp.WithDescription(`Get information about the given dataverse`),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:         "Get the dataverse information",
			ReadOnlyHint:  mcp.ToBoolPtr(true),
			OpenWorldHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString(dataverseAddressParam,
			mcp.Required(),
			mcp.Description("The address of the dataverse contract")),
	)
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dataverseAddress, err := request.RequireString(dataverseAddressParam)
		if err != nil {
			return nil, err
		}

		dataverseInfo, err := dataverse.Dataverse(ctx, cc, dataverseAddress, ref(dataverseschema.QueryMsg_Dataverse{}))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		r, err := json.Marshal(dataverseInfo)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response: %w", err)
		}

		return mcp.NewToolResultText(string(r)), nil
	}

	return server.ServerTool{Tool: tool, Handler: handler}
}
