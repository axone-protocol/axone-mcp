package mcp

import (
	"context"

	"github.com/axone-protocol/axone-sdk/dataverse"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func getGovernanceCode(dqc dataverse.QueryClient) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	const resourceParam = "resource"
	return mcp.NewTool("get_resource_governance_code",
			mcp.WithDescription(`Get the governance code attached to the given resource (if any)`),
			mcp.WithString(resourceParam,
				mcp.Required(),
				mcp.Description("The DID URI of the resource"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			resourceDID, err := requiredParam[string](request, resourceParam)
			if err != nil {
				return nil, err
			}
			govDID, err := dqc.GetResourceGovAddr(ctx, resourceDID)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), err
			}
			code, err := dqc.GovCode(ctx, govDID)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), err
			}

			return mcp.NewToolResultText(code), nil
		}
}
