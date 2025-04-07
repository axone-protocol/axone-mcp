package mcp

import (
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// requiredParam is a helper that retrieves a required parameter from the request.
func requiredParam[T comparable](r mcp.CallToolRequest, key string) (T, error) {
	var zero T

	val, ok := r.Params.Arguments[key]
	if !ok {
		return zero, fmt.Errorf("missing required parameter: %s", key)
	}

	typedVal, ok := val.(T)
	if !ok {
		return zero, fmt.Errorf("parameter %s is not of type %T", key, zero)
	}

	if typedVal == zero {
		return zero, fmt.Errorf("parameter %s must not be empty", key)
	}

	return typedVal, nil
}
