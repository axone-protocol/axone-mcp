package cognitarium

import (
	"context"
	"encoding/json"
	"fmt"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	schema "github.com/axone-protocol/axone-contract-schema/go/cognitarium-schema/v6"
	"google.golang.org/grpc"
)

func Select(ctx context.Context, cc grpc.ClientConnInterface,
	address string, req *schema.QueryMsg_Select,
	opts ...grpc.CallOption,
) (*schema.SelectResponse, error) {
	rawQueryData, err := json.Marshal(map[string]any{"select": req})
	if err != nil {
		return nil, fmt.Errorf("encode select query (%s): %w", address, err)
	}

	rawResponseData, err := queryContract(ctx, cc, address, rawQueryData, opts...)
	if err != nil {
		return nil, err
	}

	var response schema.SelectResponse
	if err := json.Unmarshal(rawResponseData, &response); err != nil {
		return nil, fmt.Errorf("decode select response (%s): %w", address, err)
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
