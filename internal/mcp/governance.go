package mcp

import (
	"context"
	"encoding/base64"
	"fmt"

	dataverseschema "github.com/axone-protocol/axone-contract-schema/go/dataverse-schema/v6"
	lawstoneschema "github.com/axone-protocol/axone-contract-schema/go/law-stone-schema/v6"
	"github.com/axone-protocol/axone-mcp/internal/axone/cognitarium"
	"github.com/axone-protocol/axone-mcp/internal/axone/dataverse"
	"github.com/axone-protocol/axone-mcp/internal/axone/lawstone"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"google.golang.org/grpc"
)

func getGovernanceCode(cc grpc.ClientConnInterface) server.ServerTool {
	const dataverseAddressParam = "dataverse"
	const resourceParam = "resource"
	tool := mcp.NewTool("get_resource_governance_code",
		mcp.WithDescription(`Get the governance code attached to the given resource (if any) in the given dataverse`),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:         "Get the governance code for a resource",
			ReadOnlyHint:  mcp.ToBoolPtr(true),
			OpenWorldHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString(dataverseAddressParam,
			mcp.Required(),
			mcp.Description("The address of the dataverse contract")),
		mcp.WithString(resourceParam,
			mcp.Required(),
			mcp.Description("The DID URI of the resource")),
	)
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dataverseAddress, err := request.RequireString(dataverseAddressParam)
		if err != nil {
			return nil, err
		}

		resourceDID, err := request.RequireString(resourceParam)
		if err != nil {
			return nil, err
		}

		dataverseInfo, err := dataverse.Dataverse(ctx, cc, dataverseAddress, ref(dataverseschema.QueryMsg_Dataverse{}))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		cognitariumAddress := string(dataverseInfo.TriplestoreAddress)
		if cognitariumAddress == "" {
			return mcp.NewToolResultError("no triplestore address found"), nil
		}

		lawstoneAddress, err := cognitarium.GetGovernanceAddressForResource(ctx, cc, cognitariumAddress, resourceDID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		code, err := lawstone.ProgramCode(ctx, cc, lawstoneAddress, ref(lawstoneschema.QueryMsg_ProgramCode{}))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		decodedCode, err := base64.StdEncoding.DecodeString(*code)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to decode base64 code '%s': %v", *code, err)), nil
		}

		return mcp.NewToolResultText(string(decodedCode)), nil
	}

	return server.ServerTool{Tool: tool, Handler: handler}
}
