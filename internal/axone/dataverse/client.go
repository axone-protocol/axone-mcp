package dataverse

import (
	"context"
	"encoding/json"
	"fmt"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	schema "github.com/axone-protocol/axone-contract-schema/go/dataverse-schema/v6"
	"google.golang.org/grpc"
)

func Dataverse(ctx context.Context, cc grpc.ClientConnInterface,
	address string, req *schema.QueryMsg_Dataverse, opts ...grpc.CallOption,
) (*schema.DataverseResponse, error) {
	rawQueryData, err := json.Marshal(map[string]any{"dataverse": req})
	if err != nil {
		return nil, fmt.Errorf("encode dataverse query (%s): %w", address, err)
	}

	rawResponseData, err := queryContract(ctx, cc, address, rawQueryData, opts...)
	if err != nil {
		return nil, err
	}

	var response schema.DataverseResponse
	if err := json.Unmarshal(rawResponseData, &response); err != nil {
		return nil, fmt.Errorf("decode dataverse response (%s): %w", address, err)
	}

	return &response, nil
}

func queryContract(ctx context.Context, cc grpc.ClientConnInterface,
	address string, rawQueryData []byte, opts ...grpc.CallOption,
) ([]byte, error) {
	in := &wasmtypes.QuerySmartContractStateRequest{
		Address:   address,
		QueryData: rawQueryData,
	}
	out := &wasmtypes.QuerySmartContractStateResponse{}

	if err := cc.Invoke(ctx, "/cosmwasm.wasm.v1.Query/SmartContractState", in, out, opts...); err != nil {
		return nil, err
	}

	return out.Data, nil
}
